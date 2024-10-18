import { Button, IconButton, InputAdornment, TextField } from "@mui/material";
import { useRouter } from "next/router";
import { FaUser, FaUserCircle } from "react-icons/fa";
import { RiLockPasswordFill } from "react-icons/ri";
import { useEffect, useRef, useState } from "react";
import { IoEyeOff } from "react-icons/io5";
import { IoMdCloseCircle, IoMdEye } from "react-icons/io";
import { API_URL } from "@/constants";

export default function LoginPage() {
    const router = useRouter();

    const [username, setUsername] = useState<string>("");
    const usernameTextFieldRef = useRef<HTMLInputElement | null>(null);
    const [password, setPassword] = useState<string>("");
    const [showPassword, setShowPassword] = useState<boolean>(false);
    const passwordTextFieldRef = useRef<HTMLInputElement | null>(null);

    const [waiting, setWaiting] = useState<boolean>(false);

    const [hasError, setHasError] = useState<boolean>(false);
    const [error, setError] = useState<string>("");

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

    const onLogin = () => {
        setWaiting(true);
        fetch(`${API_URL}/account/login`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                email: username,
                password,
            }),
        })
            .then(resp => {
                if (resp.ok) {
                    const jwt = resp.headers.get("x-jwt");
                    if (jwt === null) {
                        displayError("Fail to Login, Token Not Found");
                        return;
                    }
                    localStorage.setItem("token", jwt);
                    router.replace("/profile").catch(console.log);
                } else {
                    displayError("Fail to Login");
                }
            })
            .catch(error => displayError("Fail to send request"))
            .finally(() => {
                setWaiting(false);
            });
    };

    const displayError = (info: string) => {
        setHasError(true);
        setError(info);
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
                            if (event.target.value.length >= 32) {
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
                    <Button fullWidth variant="contained" disabled={waiting} onClick={() => onLogin()}>
                        Login
                    </Button>
                    {hasError ? (
                        <div className="w-full flex justify-items-start">
                            <span className="text-red-500">{error}</span>
                        </div>
                    ) : (
                        <></>
                    )}
                    <div className="w-full px-2 flex justify-end">
                        <span
                            className="
                                text-[14px] text-[#3b82f6]
                                hover:cursor-pointer hover:text-[#93c5fd]
                            "
                            onClick={() => router.replace("/signup").catch(console.log)}
                        >
                            Sign Up
                        </span>
                    </div>
                </div>
            </div>
        </>
    );
}