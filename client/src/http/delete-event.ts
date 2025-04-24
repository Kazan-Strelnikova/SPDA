import axios from "axios";
import { UUID } from "crypto";

export const deleteEventById = async (id: UUID): Promise<boolean> => {
    try {
      const response = await axios.delete(`/events/${id}`, {
        headers: {
          "Content-Type": "application/json",
        },
      });
  
      return response.status === 204;
    } catch (error: any) {
        throw new Error(
            `Failed to fetch event with ID ${id}: ${error.response?.status || error.message}`
        );
    }
};