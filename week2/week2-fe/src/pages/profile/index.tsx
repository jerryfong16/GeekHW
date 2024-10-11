import { FaUserCircle } from "react-icons/fa";
import { useState } from "react";
import { Account } from "@/dto/account";
import { useRouter } from "next/router";
import { Button } from "@mui/material";
import LoadingModal from "@/pages/profile/LoadingModal";

export default function ProfilePage() {
    const router = useRouter();

    const [loading, setLoading] = useState<boolean>(false);
    const [account, setAccount] = useState<Account>();

    return (
        <>
            <div
                className="
                    w-screen h-screen
                    bg-gradient-to-r from-cyan-500 to-blue-500
                    flex justify-center items-center
                "
            >
                <div
                    className="
                        w-[32vw]
                        p-4 rounded shadow-lg
                        bg-white
                        flex flex-col justify-center items-center space-y-2
                    "
                >
                    <FaUserCircle size={64} color={"#0ea5e9"} />
                    <div className="w-full h-full flex flex-col space-y-2">
                        <div className="flex space-x-2">
                            <span className="font-bold">Email</span>
                            <span>{account?.email}</span>
                        </div>
                        <div className="flex space-x-2">
                            <span className="font-bold">Name</span>
                            <span>{account?.name}</span>
                        </div>
                        <div className="flex space-x-2">
                            <span className="font-bold">Birth</span>
                            <span>{account?.birth}</span>
                        </div>
                        <div className="flex flex-col space-x-2">
                            <span className="font-bold">About</span>
                            <span>{account?.about}</span>
                        </div>
                    </div>
                    <Button fullWidth variant="contained">
                        Edit
                    </Button>
                </div>
            </div>
            <LoadingModal open={loading} />
        </>
    );
}
