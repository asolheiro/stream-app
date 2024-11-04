import { getVideos } from "@/app/components/VideoComponents/VideoList";
import { VideoPlayer } from "@/app/components/VideoComponents/VideoPlayer";
import { after } from "node:test";
import { Suspense } from "react";
import { VideoViews } from "./VideoViews";
import { formatDistance } from "date-fns";
import { ptBR } from "date-fns/locale";
import { getVideo } from "./getVideo";
import { VideoLikesCounter } from "./VideoLike";
import VideoCardSkeleton from "@/app/components/VideoComponents/VideoCardSkeleton";
import { VideosRecomendationList } from "@/app/components/VideoComponents/VideoRecommendations";


export default async function VideoPlayPage({
    params, 
}: {
    params: { slug: string };
}) {
    const video = await getVideo(params.slug);
    after(async () => {
        await fetch(`http://localhost:8000/videos/${video.id}/register-view`, {
            method: "POST"
        })
    });

    return (
        <main className="container mx-auto px-4 py-6">
            <div className="flex flex-col md:flex-row gap-4">
                <div className="flex-1">
                    <div className="bg-black rounded-lg overflow-x-hidden">
                        <VideoPlayer url={video.video_url} poster={video.thumbnail} />
                    </div>

                    <h2 className="mt-4 text-2xl font-bold text-primary">
                        {video.title}
                    </h2>

                    <div className="mt-2 flex items-center justify-between text-secondary">
                        <div className="flex flex-row items-center">
                            <Suspense
                                fallback={
                                    <div className="w-48 h-8 bg-secondary animate-pulse rounded mr-2"></div>
                                }
                            >
                                <VideoViews videoId={video.id} />
                                <span>
                                    &nbsp;há&nbsp;
                                    {formatDistance(video.published_at, new Date(), {
                                        locale: ptBR
                                    })}
                                </span>
                            </Suspense>
                        </div>
                        <Suspense
                            fallback={
                                <div className="w-24 h-8 bg-secondary animate-pulse rounded mr-2"></div>
                            }
                            >
                            <VideoLikesCounter videoId={video.id} />
                        </Suspense>
                    </div>
                    
                    <div className="bg-primary rounded-lg mt-2">
                        <p className="text-primary">
                            {video.description}
                        </p>
                    </div>
                </div>

                <div className="w-full md:2-1/3">
                    <h3 className="text-xl font-semibold text-primary mb-4">
                        Vídeos recomendados
                    </h3>

                    <div className="flex flex-col gap-4">
                        <Suspense
                            fallback={ new Array(10).fill(0).map((_, i) => (
                                <VideoCardSkeleton orientation="horizontal" key={i} />
                            ))}
                        >
                            <VideosRecomendationList videoId={video.id} />
                        </Suspense>
                    </div>
                </div>
            </div>
        </main>
    );
}