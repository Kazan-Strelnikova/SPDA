import axios from "axios";
import { Categories, Event } from "../types";
import { UUID } from "crypto";

export const postEventEnrollment = async function(id: UUID): Promise<number> {
    try {
        const response = await axios.post(`/events/${id}/enrollment`,
            {
                headers: {
                'Content-Type': 'application/json',
                },
                withCredentials: true,
            }
        );

        return response.status
    } catch (err: any) {
        throw new Error(
            `failed to post enrollment`
        )
    }
} 