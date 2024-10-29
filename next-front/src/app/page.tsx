import Image from "next/image";
import { VideoCard } from "./components/VideoCard./VideoCard";
import VideoCardSkeleton from "./components/VideoCard./VideoCardSkeleton";

export default function Home() {
  return (
    <div className="container mx-auto px-4 py-6">
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-6">       
        <VideoCard title="Video Title" thumbnail="/thumbnail.jpg" views={142} likes={42}/>
        <VideoCard title="Video Title" thumbnail="/thumbnail.jpg" views={142} likes={42}/>
        <VideoCard title="Video Title" thumbnail="/thumbnail.jpg" views={142} likes={42}/>
        <VideoCard title="Video Title" thumbnail="/thumbnail.jpg" views={142} likes={42}/>
        <VideoCard title="Video Title" thumbnail="/thumbnail.jpg" views={142} likes={42}/>
        <VideoCard title="Video Title" thumbnail="/thumbnail.jpg" views={142} likes={42}/>
        <VideoCard title="Video Title" thumbnail="/thumbnail.jpg" views={142} likes={42}/>
        <VideoCard title="Video Title" thumbnail="/thumbnail.jpg" views={142} likes={42}/>
      </div>
    </div>
  );
}
