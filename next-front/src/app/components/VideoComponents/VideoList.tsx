import { VideoModel } from "@/model";
import Link from "next/link";
import { VideoCard } from "./VideoCard";


export async function getVideos(search: string): Promise<VideoModel[]> {
    const url = search
        ? `${process.env.DJANGO_API_URL}/videos?q=${search}`
        : `${process.env.DJANGO_API_URL}/videos`;

    const response = await fetch(url, {
        cache: "no-cache"
    });
    return response.json();
}

export type VideoListProps = {
    search: string;
};

export async function VideosList(props: VideoListProps) {
    const { search } = props;
    const videos = await getVideos(search);

    return videos.length ? (
        videos.map((video) => (
            <Link key={video.id} href={`/${video.slug}/play`}>
                <VideoCard
                    title={video.title}
                    thumbnail={video.thumbnail}
                    views={video.views}
                    />
            </Link>
        ))
    ) : (
        <div className="flex items-center justify-center w-full col-span-full">
            <p className="text-gray-600 text-xl font-semibold">
                Nenhum vídeo encontrado
            </p>
        </div>
    );
}