import { VideoModel } from "@/model";
import Link from "next/link";
import { VideoCard } from "./VideoCard";

const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms))

export async function getVideosRecommendations(
    videoId: number
): Promise<VideoModel[]> {
    await sleep(2000);
    const response = await fetch(`http://localhost:8000/api/videos/${videoId}/recommended`, {
        cache: "no-cache",
    });

    return response.json();
}

export type VideoRecommendationListProps = {
    videoId: number;
};

export async function VideoRecommendationsList(props: VideoRecommendationListProps) {
    const { videoId } = props;
    const videos = await getVideosRecommendations(videoId);

    return videos.length ? (
        videos.map((video) => (
            <Link key={video.id} href={`/${video.slug}/play`} >
                <VideoCard
                    title={video.title}
                    thumbnail={video.thumbnail}
                    views={video.views}
                    orientation="horizontal"
                    />
            </Link>
        ))
    ) : (
        <div className="flex items-center justify-center w-full col-span-full">
            <p className="text-gray-600 text-xl font-semibold">
                Nenhum vídeo encontrado.
            </p>
        </div>
    );
}