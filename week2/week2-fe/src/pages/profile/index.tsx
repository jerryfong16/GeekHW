import { FaUserCircle } from "react-icons/fa";
import { useState } from "react";
import { Account, getDefaultAccount } from "@/dto/account";
import { useRouter } from "next/router";
import { Button } from "@mui/material";
import LoadingModal from "@/pages/profile/LoadingModal";
import EditModal from "@/pages/profile/EditModal";

export default function ProfilePage() {
    const router = useRouter();

    const [loading, setLoading] = useState<boolean>(false);
    const [account, setAccount] = useState<Account>(getDefaultAccount());

    const [isEditModalOpen, setIsEditModalOpen] = useState<boolean>(false);
    const [editing, setEditing] = useState<boolean>(false);

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
                        <div className="flex justify-between">
                            <span className="font-bold">Email</span>
                            <span>{account.email}</span>
                        </div>
                        <div className="flex justify-between">
                            <span className="font-bold">Name</span>
                            <span>{account.name}</span>
                        </div>
                        <div className="flex justify-between">
                            <span className="font-bold">Birth</span>
                            <span>{account.birth}</span>
                        </div>
                        <div className="flex flex-col space-y-2">
                            <span className="font-bold">About</span>
                            <span>{account.about}</span>
                        </div>
                    </div>
                    <Button fullWidth variant="contained" disabled={editing} onClick={() => setIsEditModalOpen(true)}>
                        Edit
                    </Button>
                </div>
            </div>
            <LoadingModal open={loading} />
            <EditModal
                open={isEditModalOpen}
                accountInfo={account}
                onClose={() => setIsEditModalOpen(false)}
                onFinish={(name, birth, about) => {}}
            />
        </>
    );
}
