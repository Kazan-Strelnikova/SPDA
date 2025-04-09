import axios from "axios"
import { UUID } from "crypto";
import { Event, EventResponse } from "../types";
import { getEventFromResponse } from "../utils/get-event-from-response";

export const getEventById = async (id: UUID): Promise<Event> => {
    try {
      const response = await axios.get<EventResponse>(`https://178.236.23.92/team-1/events/${id}`, {
        headers: {
          "Content-Type": "application/json",
        },
        // We may need to disable certificate verification for local testing (self-signed certs)
        // httpsAgent: new (require("https").Agent)({
        //   rejectUnauthorized: false,
        // }),
      });
  
      return getEventFromResponse(response.data);
    } catch (error: any) {
        throw new Error(
            `Failed to fetch event with ID ${id}: ${error.response?.status || error.message}`
        );
    }
};