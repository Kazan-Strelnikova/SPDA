import axios from "axios";
import { UUID } from "crypto";

export const deleteEventEnrollmentById = async (id: UUID): Promise<boolean> => {
    try {
      const response = await axios.delete(`/events/${id}/enrollment`, {
        headers: {
          "Content-Type": "application/json",
        },
        withCredentials: true,
      });
  
      return response.status === 204;
    } catch (error: any) {
        throw new Error(
            `Failed to fetch event with ID ${id}: ${error.response?.status || error.message}`
        );
    }
};