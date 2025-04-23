import axios from "axios";
import { User } from "../contexts/UserContext";
import { getUserFromLoginResponse } from "../utils/get-user-from-login-response";

export const getUserByToken = async function (): Promise<User> {
    try {
        const user = await axios.get("/team-1/users/signin/cookie", 
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
            `failed to fetch user`
        )
    } 
}