from operator import truediv
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
def videos_list(request):
    q = request.query_params.get('q')
    if q:
        videos = Video.objects.filter(Q(title__icontains=q) | Q(description__icontains=q) | Q(tags__name__icontains=q))
    else:
        videos = Video.objects.all()
    query = videos.filter(is_published=True).order_by('-published_at').distinct()
    serializedVideos = VideoSerializer(query, many=True)
        
    return Response(serializedVideos.data)