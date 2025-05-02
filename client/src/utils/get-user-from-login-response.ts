import { User } from "../contexts/UserContext";
import { LoginResponse } from "../types";

export const getUserFromLoginResponse = function(response: LoginResponse): User {
    return {
        name: response.name,
        lastName: response.last_name,
        email: response.email,
    }
}