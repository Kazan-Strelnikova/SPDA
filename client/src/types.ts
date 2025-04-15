import { UUID } from "crypto";

export type Category = "Conference" | "Meetup" | "Concert" | "Exhibition" | "Party" | "Sport" | "Education" | "Competition" | "Other";

export interface EventResponse {
    id: UUID;
    title: string;
    type: number;
    date: string; // ISO 8601 datetime string (e.g., "2025-05-15T09:00:00Z")
    total_seats: number;
    available_seats: number;
    creator_email: string;
    location: [number, number]; // [latitude, longitude]
    has_unlimited_seats: string; // "true" or "false" as string
    description: string;
}
  
export interface Event {
    id: UUID;
    title: string;
    type: Category;
    date: Date;
    total_seats: number;
    available_seats: number;
    creator_email: string;
    location: [number, number];
    has_unlimited_seats: string;
    description?: string;
}

export interface LoginResponse {
    name: string;
    last_name: string;
    email: string;
}