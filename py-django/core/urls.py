from django.urls import path
from core.api import videos_list, healthcheck

urlpatterns = [
    path("healthcheck", healthcheck, name="APIhealthcheck"),
    path("api/videos", videos_list, name="APIvideo_list")
]
