from operator import truediv

from django.db.migrations import serializer
from core.models import Video

from rest_framework.decorators import api_view
from rest_framework.response import Response
from core.serializer import VideoSerializer
from django.db.models import Q


@api_view(['GET'])
def healthcheck(request):
    return Response({'status': 'healthy'})

## Support to query parameters: title, description and tag name
@api_view(['GET'])
def get_videos_list(request):
    q = request.query_params.get('q')
    if q:
        videos = Video.objects.filter(Q(title__icontains=q) | Q(description__icontains=q) | Q(tags__name__icontains=q))
    else:
        videos = Video.objects.all() 
    query = videos.filter(is_published=True).order_by('-published_at').distinct()
    serializedVideos = VideoSerializer(query, many=True)
        
    return Response(serializedVideos.data)

@api_view(['GET'])
def get_video_by_id(request, id):
    video = Video.objects.get(id=id)
    serializedVideo = VideoSerializer(video)
    return Response(serializedVideo.data)

@api_view(['GET'])
def get_video_by_slug(request, slug):
    video = Video.objects.get(slug=slug)
    serializedVideo = VideoSerializer(video)
    return Response(serializedVideo.data)

@api_view(['GET'])
def get_video_recommendations_list(request, id):
    video = Video.objects.get(id=id)
    tags = video.tags.all()
    videos = video.objects.filter(tags__in=tags).exclude(id=id).distinct().order_by('-num-views')
    serializedVideos = VideoSerializer(videos, many=True)
    return Response(serializedVideos.data)

@api_view(['GET'])
def get_video_likes(request, id):
    video = Video.objects.get(id=id)
    return Response({'views': video.num_likes})

@api_view(['GET'])
def get_video_views(request, id):
    video = Video.objects.get(id=id)
    return Response({'views': video.num_views})

@api_view(['POST'])
def post_add_like(request, id):
    video = Video.objects.get(id=id)
    video.num_likes -= 1
    video.save()
    return Response({'likes': video.num_likes})

@api_view(['POST'])
def post_add_unlike(request, id):
    video = Video.objects.get(id=id)
    video.num_likes -= 1
    video.save()
    return Response({'likes': video.num_likes})

@api_view(['POST'])
def post_register_view(request, id):
    video = Video.objects.get(id=id)
    video.num_views += 1
    video.save()
    return Response({'views': video.num_views})