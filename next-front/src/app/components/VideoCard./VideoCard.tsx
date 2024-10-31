import Image from "next/image";
import { space } from "postcss/lib/list";

export type VideoCardProps = {
    title: string;
    thumbnail: string;
    views: number;
    likes?: number;
    orientation?: "vertical" | "horizontal"
}

export function VideoCard(props: VideoCardProps) {
    const {
        title,
        thumbnail,
        views,
        likes,
        orientation = "vertical",
    } = props;
    const flexDirection = orientation === "vertical" ? "flex-col": "flex-row";
    const imageSize = 
        orientation === "horizontal" ? "w-5/12 h-24" : "w-full h-40";
    return (
        <div className={`flex ${flexDirection} overflow-hidden gap-2`}>
            <div className={`${imageSize} relative`}>
                <Image 
                    src={thumbnail}
                    priority={true}
                    alt="Video thumbnail"
                    fill={true}
                    sizes="100%"
                    className="object-cover rounded-lg"
                    />
            <div className="flex items-center justify-between mt-1 text-sm text-secondary">
                <span>{views} visualizações</span>
                <div className="flex items-center space-x-1">
                    {likes !== undefined && <span>{likes}</span>}
                </div>
            </div>
            </div>
        </div>
    );
}