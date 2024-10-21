from django.db import models

class Video(models.Model):
    title = models.CharField(max_length=100, unique=True, verbose_name='Título')
    description = models.TextField(verbose_name='Descrição')
    thumbnail = models.ImageField(upload_to='thumbnails/', verbose_name='Thumbnail')
    slug = models.SlugField(unique=True)
    published_at = models.DateTimeField(
        verbose_name='Publicado em', 
        null=True, 
        editable=False
        )
    is_published = models.BooleanField(default=False, verbose_name='Publicado')
    num_likes = models.IntegerField(default=0, verbose_name='Likes', editable=False)
    num_views = models.IntegerField(default=0, verbose_name='Visualizações', editable=False)
    tags = models.ManyToManyField('Tag', verbose_name='Tags', related_name='videos')
    author = models.ForeignKey(
        'auth.User', 
        on_delete=models.PROTECT, 
        verbose_name='Autor', 
        related_name='videos',
        editable=False 
        )
    
    def get_video_status_display(self):
        if not hasattr(self, 'video_media'):
            return 'Pendente'
        return self.video_media.get_status_display()
    
    def __str__(self):
        return self.title
        
    class Meta:
        verbose_name = 'Vídeo'
        verbose_name_plural = "Vídeos"
        

class VideoMedia(models.Model):
    
    class Status(models.TextChoices):
        UPLOAD_STARTED = "UPLOAD_STARTED", "Upload iniciado"
        PROCESS_STARTED = "PROCESS_STARTED", "Processamento iniciado"
        PROCESS_FINISHED = "PROCESS_FINISHED", "Processamento finalizado"
        PROCESS_ERROR = "PROCESS_ERROR", "Erro no processamento"
        
    video_path = models.CharField(max_length=255, verbose_name='Video')    
    status = models.CharField(
        max_length=20, 
        choices=Status.choices, 
        default=Status.UPLOAD_STARTED, 
        verbose_name='status'
        )
    video = models.OneToOneField("Video", on_delete=models.PROTECT, verbose_name="Vídeo", related_name="video_media")
    
    def get_status_display(self):
        return VideoMedia.Status(self.status).label
    
    class Meta:
        verbose_name = 'Mídia'
        verbose_name_plural = 'Mídias'
        
        
class Tag(models.Model):
    name = models.CharField(max_length=50, unique=True, verbose_name='Nome')
    
    def __str__(self):
        return self.name
    
    class Meta:
        verbose_name = 'Tag'
        verbose_name_plural = 'Tags'

    