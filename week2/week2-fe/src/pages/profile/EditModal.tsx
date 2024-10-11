import { Account } from "@/dto/account";
import { useEffect, useState } from "react";
import { Backdrop, Fade, Modal } from "@mui/material";
import { Box } from "@mui/system";

interface EditModalProps {
    open: boolean;
    accountInfo: Account;
    onClose: boolean;
    onFinish: (account: Account) => void;
}

export default function EditModal({ open, accountInfo, onClose, onFinish }: EditModalProps) {
    const [account, setAccount] = useState<Account>(accountInfo);
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);

    useEffect(() => {
        setIsModalOpen(open);
    }, [open]);

    useEffect(() => {
        setAccount(accountInfo);
    }, [accountInfo]);

    return (
        <Modal
            aria-labelledby="transition-modal-title"
            aria-describedby="transition-modal-description"
            open={isModalOpen}
            closeAfterTransition
            slots={{ backdrop: Backdrop }}
            slotProps={{
                backdrop: {
                    timeout: 500,
                },
            }}
        >
            <Fade in={isModalOpen}>
                <Box
                    sx={{
                        position: "absolute",
                        top: "50%",
                        left: "50%",
                        transform: "translate(-50%, -50%)",
                        width: 400,
                        outline: "none",
                    }}
                >
                    <div className="w-full h-full p-4 bg-white rounded-2xl shadow-2xl flex flex-col justify-center items-center"></div>
                </Box>
            </Fade>
        </Modal>
    );
}
