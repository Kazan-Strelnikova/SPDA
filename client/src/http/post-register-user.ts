import axios from "axios";
import { User } from "../contexts/UserContext";
import { getUserFromLoginResponse } from "../utils/get-user-from-login-response";

export const postRegisterUser = async function (name: string, last_name: string, email: string, password: string): Promise<User> {
    try {
        const user = await axios.post("/team-1/users/signup", 
            {   
                name: name,
                last_name: last_name,
                email: email,
                password: password,
            }, 
            {
                headers: {
                'Content-Type': 'application/json',
                },
                withCredentials: true,
            },
        );

        return getUserFromLoginResponse(user.data?.user)
    } catch (err: any) {
        throw new Error(
            `failed to fetch user ${email}`
        )
    } 
}