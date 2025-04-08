import { Event, EventResponse } from "../types";
import { getCategoryFromNumber } from "./get-category-from-number";

export const getEventFromResponse = (data : EventResponse): Event => {
    return {
        id: data.id,
        title: data.title,
        type: getCategoryFromNumber(data.type),
        date: new Date(data.date),
        total_seats: data.total_seats,
        available_seats: data.available_seats,
        creator_email: data.creator_email,
        location: data.location,
        has_unlimited_seats: data.has_unlimited_seats,
        description: data.description == "" ? undefined : data.description,
    }
}

export {}