import axios from "axios";
import { User } from "../contexts/UserContext";
import { LoginResponse } from "../types";
import { getUserFromLoginResponse } from "../utils/get-user-from-login-response";

export const postLoginUser = async function (email: string, password: string): Promise<User> {
    try {
        const user = await axios.post("https://178.236.23.92/team-1/users/signin", 
            {
                email: email,
                password: password,
            }, 
            {
                headers: {
                'Content-Type': 'application/json',
                },
            }
        );

        return getUserFromLoginResponse(user.data?.user)
    } catch (err: any) {
        throw new Error(
            `failed to fetch user ${email}`
        )
    } 
}