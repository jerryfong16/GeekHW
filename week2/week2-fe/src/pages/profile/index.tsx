import { FaUserCircle } from "react-icons/fa";
import { useEffect, useState } from "react";
import { Account, getDefaultAccount } from "@/dto/account";
import { useRouter } from "next/router";
import { Button } from "@mui/material";
import LoadingModal from "@/pages/profile/LoadingModal";
import EditModal from "@/pages/profile/EditModal";
import { API_URL } from "@/constants";

export default function ProfilePage() {
    const router = useRouter();

    const [loading, setLoading] = useState<boolean>(false);
    const [account, setAccount] = useState<Account>(getDefaultAccount());

    const [isEditModalOpen, setIsEditModalOpen] = useState<boolean>(false);
    const [editing, setEditing] = useState<boolean>(false);

    const [hasError, setHasError] = useState<boolean>(false);
    const [error, setError] = useState<string>("");

    useEffect(() => {
        getAccountInfo();
    }, []);

    const getAccountInfo = () => {
        setLoading(true);
        fetch(`${API_URL}/account/profile`, {
            method: "GET",
            headers: {
                Authorization: `Bearer ${localStorage.getItem("token")}`,
            },
        })
            .then(resp => {
                if (resp.ok) {
                    resp.json().then(content => {
                        const account: Account = content;
                        setAccount(account);
                    });
                } else {
                    localStorage.removeItem("token");
                    router.replace("/login").catch(console.log);
                }
            })
            .catch(error => {})
            .finally(() => setLoading(false));
    };

    const onEdit = (name: string, birth: string, about: string) => {
        setEditing(true);
        fetch(`${API_URL}/account/edit`, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${localStorage.getItem("token")}`,
            },
            body: JSON.stringify({
                name,
                birth,
                about,
            }),
        })
            .then(resp => {
                if (resp.ok) {
                    const jwt = resp.headers.get("x-jwt");
                    if (jwt !== null) {
                        localStorage.setItem("token", jwt);
                    }
                    setAccount({ ...account, name, birth, about });
                } else {
                    displayError("Fail to edit");
                }
            })
            .catch(error => displayError("Fail to send request"))
            .finally(() => {
                setEditing(false);
            });
    };

    const displayError = (info: string) => {
        setHasError(true);
        setError(info);
    };

    return (
        <>
            <div className="w-screen h-screen bg-gradient-to-r from-cyan-500 to-blue-500 flex justify-center items-center">
                <div className="w-[32vw] p-4 rounded shadow-lg bg-white flex flex-col justify-center items-center space-y-2">
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
                onFinish={(name, birth, about) => onEdit(name, birth, about)}
            />
        </>
    );
}
