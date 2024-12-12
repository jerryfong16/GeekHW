import { Account } from "@/dto/account";
import { useEffect, useState } from "react";
import { Backdrop, Button, Fade, IconButton, InputAdornment, Modal, TextField } from "@mui/material";
import { Box } from "@mui/system";
import { IoMdCloseCircle } from "react-icons/io";

interface EditModalProps {
    open: boolean;
    accountInfo: Account;
    onClose: () => void;
    onFinish: (name: string, birth: string, about: string) => void;
}

export default function EditModal({ open, accountInfo, onClose, onFinish }: EditModalProps) {
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);

    const [name, setName] = useState<string>("");
    const [birth, setBirth] = useState<string>("");
    const [about, setAbout] = useState<string>("");

    useEffect(() => {
        setIsModalOpen(open);
    }, [open]);

    useEffect(() => {
        initInfo(accountInfo);
    }, [accountInfo]);

    const initInfo = (account: Account) => {
        setName(account.name);
        setBirth(account.birth);
        setAbout(account.about);
    };

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
                    <div className="w-full h-full p-4 bg-white rounded-2xl shadow-2xl flex flex-col justify-center items-center space-y-2">
                        <TextField
                            autoComplete="off"
                            fullWidth
                            label="Name"
                            variant="outlined"
                            slotProps={{
                                input: {
                                    endAdornment: (
                                        <InputAdornment position="end">
                                            {name.length > 0 ? (
                                                <IconButton
                                                    tabIndex={-1}
                                                    onClick={() => {
                                                        setName("");
                                                    }}
                                                >
                                                    <IoMdCloseCircle />
                                                </IconButton>
                                            ) : (
                                                <></>
                                            )}
                                        </InputAdornment>
                                    ),
                                },
                            }}
                            value={name}
                            onChange={event => {
                                setName(event.target.value.trim());
                            }}
                        />
                        <TextField
                            autoComplete="off"
                            fullWidth
                            label="Birth"
                            variant="outlined"
                            value={birth}
                            onChange={event => {
                                setBirth(event.target.value.trim());
                            }}
                        />
                        <TextField
                            autoComplete="off"
                            fullWidth
                            label="About"
                            variant="outlined"
                            value={about}
                            onChange={event => {
                                setAbout(event.target.value.trim());
                            }}
                        />
                        <Button
                            fullWidth
                            variant="contained"
                            onClick={() => {
                                onFinish(name, birth, about);
                                onClose();
                            }}
                        >
                            Confirm
                        </Button>
                        <Button
                            fullWidth
                            variant="contained"
                            color="warning"
                            onClick={() => {
                                onClose();
                                initInfo(accountInfo);
                            }}
                        >
                            Cancel
                        </Button>
                    </div>
                </Box>
            </Fade>
        </Modal>
    );
}
