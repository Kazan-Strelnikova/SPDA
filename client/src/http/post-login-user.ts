import axios from "axios";
import { User } from "../contexts/UserContext";
import { getUserFromLoginResponse } from "../utils/get-user-from-login-response";

export const postLoginUser = async function (email: string, password: string): Promise<User> {
    try {
        const user = await axios.post("/team-1/users/signin", 
            {
                email: email,
                password: password,
            }, 
            {
                headers: {
                'Content-Type': 'application/json',
                },
                withCredentials: true,
            }

        );

        return getUserFromLoginResponse(user.data?.user)
    } catch (err: any) {
        throw new Error(
            `failed to fetch user ${email}`
        )
    } 
}