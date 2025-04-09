import axios from "axios";
import { Event, EventResponse } from "../types";
import { getEventFromResponse } from "../utils/get-event-from-response";


interface Filters {
    type?: number;
    limit?: number;
    offset?: number;
    creatorEmail?: string;
    visitorEmail?: string;
    longtitude?: number;
    latitude?: number;
    radius?: number;
    before?: Date;
    after?: Date;
}

export const getEventsWithFilters = async (filters: Filters): Promise<Event[]> => {
    try {

        const params = new URLSearchParams();

        if (filters.type !== undefined) params.append("type", filters.type.toString());
        if (filters.limit !== undefined) params.append("limit", filters.limit.toString());
        if (filters.offset !== undefined) params.append("offset", filters.offset.toString());
        if (filters.creatorEmail) params.append("creator_email", filters.creatorEmail);
        if (filters.visitorEmail) params.append("visitor_email", filters.visitorEmail);
        if (filters.longtitude !== undefined) params.append("lon", filters.longtitude.toString());
        if (filters.latitude !== undefined) params.append("lat", filters.latitude.toString());
        if (filters.radius !== undefined) params.append("radius", filters.radius.toString());
        if (filters.before) params.append("before", filters.before.toISOString());
        if (filters.after) params.append("after", filters.after.toISOString());

        console.log(params, filters.visitorEmail)

        const response = await axios.get<EventResponse[]>(`https://178.236.23.92/team-1/events`, {
            headers: {
            "Content-Type": "application/json",
            },
            params,
        });
    
        return response.data.map(getEventFromResponse);
    } catch (error: any) {
        throw new Error(
            `Failed to fetch events with filters ${filters}: ${error.response?.status || error.message}`
        );
    }
}