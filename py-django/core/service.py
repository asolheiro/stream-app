from dataclasses import dataclass
import os
import shutil

from django.conf import settings
from django.contrib.gis.gdal.raster import source
from django.db import IntegrityError, transaction

from core.models import Video, VideoMedia
from .rabbitmq import create_rabbitmq_connection

@dataclass
class VideoService:
    storage: 'Storage'
    
    def get_chunk_directory(self, video_id: int) -> str:
        return f'/media/tmp/videos/{video_id}'
    
    def find_video(self, video_id: int) -> Video:
        return Video.objects.get(id=video_id)
    
    def process_upload(self, video_id: int, chunk_index: int, chunk: bytes) -> None:
        video = self.find_video(video_id)
        directory = self.get_chunk_directory(video_id)
        
        with transaction.atomic():
            video_media = self.__prepare_video_media(video)
            
            if video_media.status == VideoMedia.Status.PROCESSING_STARTED:
                raise VideoMediaInvalidStatusException("An upload is already in progress.")
            
            if video_media.status == VideoMedia.Status.UPLOAD_STARTED:
                self.storage.storage_chunk(str(video_media.video_path), chunk_index, chunk)
                return
            
            if video_media.status == VideoMedia.Status.PROCESSING_FINISHED:
                video_media.video_path = directory
                video.is_published = False
                video.save()
            
            video_media.status = VideoMedia.Status.UPLOAD_STARTED
            video_media.save()
            self.storage.storage_chunk(str(video_media.video_path), chunk_index, chunk).status

    def __prepare_video_media(self, video: Video) -> VideoMedia:
        try:
            video_media = video.video_media
        except Video.video_media.RelatedObjectDoesNotExist:
            try: video_media = VideoMedia.objects.create(
                video=video,
                status=VideoMedia.Status.UPLOAD_STARTED,
                video_path=self.get_chunk_directory(video.id)
            )
            except IntegrityError:
                video_media= VideoMedia.objects.get(video=video)
                
        return video_media
    
    def finalize_upload(self, video_id: int, total_chunks) -> None:
        video = self.find_video(video_id)
        
        try:
            video_media = video.video_media
            if video_media.status != VideoMedia.Status.UPLOAD_STARTED:
                raise VideoMediaInvalidStatusException('Upload must be started to finish it.')
            is_chunk_valid = self.__validate_chunks(str(video.video_media.video_path), total_chunks)
            if not is_chunk_valid:
                raise VideoChunkUploadException("Chunks are not valid.")
            
            video_media.status = VideoMedia.Status.PROCESSING_STARTED
            video_media.save()
            
            
            self.__produce_message(video_id, video_media.video_path, 'chunks')
            
        except Video.video_media.RelatedObjectDoesNotExist:
            raise VideoMediaNotExistsException('Upload not started.')    
          
    
    def __validate_chunks(self, video_path: str, total_chunks: int) -> bool:
        if not os.path.exists(video_path):
            return False
        
        for i in range(total_chunks):
            chunk_path = os.path.join(video_path, f'{i}.chunk')
            if not os.path.exists(chunk_path):
                return False
        
        return True
    
    
    def upload_chunks_to_external_storage(self, video_id: int) -> None:
        self.find_video(video_id)
        
        source_path = self.get_chunk_directory(video_id)
        dest_path = f'/media/uploads/{video_id}'
        
        self.storage.move_chunks(source_path, dest_path)
        
        conversion_key = settings.CONVERSION_KEY
        
        self.__produce_message(video_id, dest_path, conversion_key)
        
    def register_processed_video_path(self, video_id: int, video_path) -> None:
        video = self.find_video(video_id)
        video_media = video.video_media
        if video_media.status != VideoMedia.Status.PROCESSING_STARTED:
            raise VideoMediaInvalidStatusException('Processing must be started to finish it.')
        video_media.video_path = video_path
        video_media.status = VideoMedia.Status.PROCESSING_FINISHED
        video_media.save()
        
    def __produce_message(self, video_id: int, path: str, routing_key: str):
        with create_rabbitmq_connection() as conn:
            producer = conn.Producer(serializer='json')
            producer.publish(
                {
                'video_id': video_id,
                'path': path
                },
                exchange='conversionExchange',
                routing_key=routing_key
            )
        
    
def create_video_service_factory() -> VideoService:
    return VideoService(Storage())


class VideoMediaInvalidStatusException(Exception):
    pass

class VideoMediaNotExistsException(Exception):
    pass

class VideoChunkUploadException(Exception):
    pass

class Storage:
    
    def storage_chunk(self, directory: str, chunk_index: int, chunk: bytes) -> None:
        if not os.path.exists(directory):
            os.makedirs(directory)
            
        chunk_path = os.path.join(directory, f'{chunk_index}.chunk')
        
        with open(chunk_path, 'wb') as chunk_file:
            chunk_file.write(chunk)
    
    def move_chunks(self, source_path: str, dest_path: str) -> None:
        if not os.path.exists(dest_path):
            os.makedirs(dest_path, exist_ok=True)
        
        for filename in os.listdir(source_path):
            file_path = os.path.join(source_path, filename)
            
            if os.path.isfile(file_path):
                try:
                    shutil.move(file_path, os.path.join(dest_path, filename))
                    print(f"Arquivo {source_path}/{filename} movido para {dest_path}.")
                except Exception as e:
                    print(f"Erro ao mover o arquivo {source_path}/{filename}: {e}")
            else:
                print(f"{file_path} não é um arquivo")