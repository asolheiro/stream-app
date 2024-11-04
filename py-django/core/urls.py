from os import name
from django.urls import path
from core.api import (
    healthcheck, get_videos_list, get_video_by_id, get_video_by_slug,
    get_video_recommendations_list, get_video_likes, get_video_views,
    post_add_like, post_add_unlike, post_register_view
    )

urlpatterns = [
    path("healthcheck", healthcheck, name="APIhealthcheck"),
    path("api/videos", get_videos_list, name="APIget_video_list"),
    path("api/videos/<int:id>", get_video_by_id, name="APIget_video_by_id"),
    path("api/videos/<slug:slug>", get_video_by_slug, name="APIget_video_by_slug"),
    path("api/videos/<int:id>/recommendations", get_video_recommendations_list, name="APIget_video_recommendations_list"),
    path("api/videos/<int:id>/likes", get_video_likes, name="APIget_video_likes"),
    path("api/videos/<int:id>/views", get_video_views, name="APIget_video_views"),
    path("api/videos/<int:id>/like", post_add_like, name="APIpost_add_like"),
    path("api/videos/<int:id>/unlike", post_add_unlike, name="APIpost_add_unlike"),    
    path("api/videos/<int:id>/register-view", post_register_view, name="APIpost_register_view"),    
]