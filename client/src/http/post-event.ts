import axios from "axios";
import { Categories, Event } from "../types";

export const postEvent = async function(evt: Event, email: string): Promise<number> {
    try {
        console.log(evt.total_seats.valueOf())
        const response = await axios.post("https://178.236.23.92/team-1/events", 
            {
                title: evt.title,
                type: evt?.type ? Categories.findIndex((value) => value[1] === evt.type) : 8,
                date: evt.date,
                total_seats: Number(evt?.total_seats),
                creator_email: email,
                location: evt?.location ?? {
                    latitude: "40.7128",
                    longitude: "-74.0060"
                },
                has_unlimited_seats: evt?.has_unlimited_seats ?? "false",
                description: evt?.description ?? ""
            }, 
            {
                headers: {
                'Content-Type': 'application/json',
                },
            }
        );

        return response.status
    } catch (err: any) {
        throw new Error(
            `failed to fetch user ${email}`
        )
    }
} 