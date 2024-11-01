import Link from "next/link";
import { Suspense } from "react";
import { SearchBar } from "./SearchBar";


export function NavBar() {
    return (
        <header className="bg-primary border-secondary border-b-2">
            <div className="container mx-auto px-4 py-2 flex items-center justify-between">
                <div className="flex-items-center space-x-4">
                    <Link href="/" className="text-2xl font-bold text-[#ffcd00]">
                    Tube
                    </Link>
                </div>

                <div className="w-1/2 relative">
                    {/* https://next.js.or/docs/messages/missing-suspense-ith-csr-bailout */}
                    <Suspense>
                        <SearchBar />
                    </Suspense>
                </div>

                <div className="flex items-center space-x-4">
                    <a href="#" className="text-primary">
                        Login
                    </a>
                </div>
            </div>
        </header>
    )
}