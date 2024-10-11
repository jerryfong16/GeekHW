import { useRouter } from "next/router";
import { useEffect, useRef, useState } from "react";
import { FaUser, FaUserCircle } from "react-icons/fa";
import { Button, IconButton, InputAdornment, TextField } from "@mui/material";
import { IoMdCloseCircle, IoMdEye } from "react-icons/io";
import { RiLockPasswordFill } from "react-icons/ri";
import { IoEyeOff } from "react-icons/io5";

export default function SignupPage() {
    const router = useRouter();

    const [username, setUsername] = useState<string>("");
    const usernameTextFieldRef = useRef<HTMLInputElement | null>(null);

    const [password, setPassword] = useState<string>("");
    const [showPassword, setShowPassword] = useState<boolean>(false);
    const passwordTextFieldRef = useRef<HTMLInputElement | null>(null);

    const [confirmPassword, setConfirmPassword] = useState<string>("");
    const [showConfirmPassword, setShowConfirmPassword] = useState<boolean>(false);
    const confirmPasswordTextFieldRef = useRef<HTMLInputElement | null>(null);

    const [waiting, setWaiting] = useState<boolean>(false);

    useEffect(() => {
        usernameTextFieldGetFocused();
    }, []);

    const usernameTextFieldGetFocused = () => {
        if (usernameTextFieldRef.current) {
            usernameTextFieldRef.current.focus();
        }
    };

    const passwordTextFieldGetFocused = () => {
        if (passwordTextFieldRef.current) {
            passwordTextFieldRef.current.focus();
            setTimeout(() => {
                const position = passwordTextFieldRef.current?.value.length ?? 0;
                passwordTextFieldRef.current?.setSelectionRange(position, position);
            }, 0);
        }
    };

    const confirmPasswordTextFieldGetFocused = () => {
        if (confirmPasswordTextFieldRef.current) {
            confirmPasswordTextFieldRef.current.focus();
            setTimeout(() => {
                const position = confirmPasswordTextFieldRef.current?.value.length ?? 0;
                confirmPasswordTextFieldRef.current?.setSelectionRange(position, position);
            }, 0);
        }
    };

    const onSignUp = () => {
        console.log(username);
        console.log(password);
    };

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
                        w-[20vw]
                        p-4 rounded shadow-lg
                        bg-white
                        flex flex-col justify-center items-center space-y-2
                    "
                >
                    <FaUserCircle size={64} color={"#0ea5e9"} />
                    <TextField
                        inputRef={usernameTextFieldRef}
                        fullWidth
                        label="Username"
                        variant="outlined"
                        slotProps={{
                            input: {
                                startAdornment: (
                                    <InputAdornment position="start">
                                        <FaUser />
                                    </InputAdornment>
                                ),
                                endAdornment: (
                                    <InputAdornment position="end">
                                        {username.length > 0 ? (
                                            <IconButton
                                                tabIndex={-1}
                                                onClick={() => {
                                                    setUsername("");
                                                    usernameTextFieldGetFocused();
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
                        value={username}
                        onChange={event => {
                            if (event.target.value.length >= 16) {
                                return;
                            }
                            setUsername(event.target.value.trim());
                        }}
                    />
                    <TextField
                        inputRef={passwordTextFieldRef}
                        fullWidth
                        label="Password"
                        variant="outlined"
                        type={showPassword ? "text" : "password"}
                        slotProps={{
                            input: {
                                startAdornment: (
                                    <InputAdornment position="start">
                                        <RiLockPasswordFill />
                                    </InputAdornment>
                                ),
                                endAdornment: (
                                    <InputAdornment position="end">
                                        <div className="flex">
                                            {password.length > 0 ? (
                                                <IconButton
                                                    tabIndex={-1}
                                                    onClick={() => {
                                                        setPassword("");
                                                        passwordTextFieldGetFocused();
                                                    }}
                                                >
                                                    <IoMdCloseCircle />
                                                </IconButton>
                                            ) : (
                                                <></>
                                            )}
                                            <IconButton
                                                tabIndex={-1}
                                                onClick={() => {
                                                    setShowPassword(!showPassword);
                                                    passwordTextFieldGetFocused();
                                                }}
                                            >
                                                {showPassword ? <IoMdEye /> : <IoEyeOff />}
                                            </IconButton>
                                        </div>
                                    </InputAdornment>
                                ),
                            },
                        }}
                        value={password}
                        onChange={event => {
                            if (event.target.value.length >= 16) {
                                return;
                            }
                            setPassword(event.target.value.trim());
                        }}
                    />
                    <TextField
                        inputRef={confirmPasswordTextFieldRef}
                        fullWidth
                        label="Confirm Password"
                        variant="outlined"
                        type={showConfirmPassword ? "text" : "password"}
                        slotProps={{
                            input: {
                                startAdornment: (
                                    <InputAdornment position="start">
                                        <RiLockPasswordFill />
                                    </InputAdornment>
                                ),
                                endAdornment: (
                                    <InputAdornment position="end">
                                        <div className="flex">
                                            {password.length > 0 ? (
                                                <IconButton
                                                    tabIndex={-1}
                                                    onClick={() => {
                                                        setConfirmPassword("");
                                                        confirmPasswordTextFieldGetFocused();
                                                    }}
                                                >
                                                    <IoMdCloseCircle />
                                                </IconButton>
                                            ) : (
                                                <></>
                                            )}
                                            <IconButton
                                                tabIndex={-1}
                                                onClick={() => {
                                                    setShowConfirmPassword(!showConfirmPassword);
                                                    confirmPasswordTextFieldGetFocused();
                                                }}
                                            >
                                                {showConfirmPassword ? <IoMdEye /> : <IoEyeOff />}
                                            </IconButton>
                                        </div>
                                    </InputAdornment>
                                ),
                            },
                        }}
                        value={confirmPassword}
                        onChange={event => {
                            if (event.target.value.length >= 16) {
                                return;
                            }
                            setConfirmPassword(event.target.value.trim());
                        }}
                    />
                    <Button fullWidth variant="contained" onClick={() => onSignUp()}>
                        Sign Up
                    </Button>
                    <div className="w-full px-2 flex justify-end">
                        <span
                            className="
                                text-[14px] text-[#3b82f6]
                                hover:cursor-pointer hover:text-[#93c5fd]
                            "
                            onClick={() => router.replace("/login").catch(console.log)}
                        >
                            Login
                        </span>
                    </div>
                </div>
            </div>
        </>
    );
}
