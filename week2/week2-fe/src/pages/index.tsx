import { FiLoader } from "react-icons/fi";
import { useRouter } from "next/router";
import { useEffect } from "react";

export default function Home() {
    const router = useRouter();

    useEffect(() => {
        setTimeout(() => router.replace("/login").catch(console.log), 2000);
    }, []);

    return (
        <>
            <div
                className="
                    w-screen h-screen
                    bg-gradient-to-r from-cyan-500 to-blue-500
                    flex justify-center items-center
                "
            >
                <FiLoader size={128} className="animate-spin" />
            </div>
        </>
    );
}
