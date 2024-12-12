export interface Account {
    id: number;
    email: string;
    password: string;
    name: string;
    birth: string;
    about: string;
    createdTime: number;
    updatedTime: number;
}

export function getDefaultAccount(): Account {
    return {
        id: 0,
        email: "",
        password: "",
        name: "",
        birth: "",
        about: "",
        createdTime: 0,
        updatedTime: 0,
    };
}
