import { useContext, useEffect } from "react";
import { CalendarDay } from "../components/calendar-day/calendar-day";
import { EventNote } from "../components/event-note/event-note";
import { getEventById } from "../http/get-event-by-id";
import { getEventsWithFilters } from "../http/get-events-with-filters";
import { Calendar } from "../components/calendar/calendar";
import { UserContext } from "../contexts/UserContext";

export const MainPage : React.FC = () => {
    const userContext = useContext(UserContext);
    if (!userContext) {
        throw new Error("Calendar must be used within a UserProvider");
    }
    const { user, setUser } = userContext;


    useEffect(
        ()=>{(async function fun() {
            try {
                const evts = await getEventsWithFilters({
                    before: new Date("2025-05-15T09:00:01Z"),
                    after: new Date("2025-05-15T08:00:00Z"),
                    visitorEmail: user?.email,
                })
                console.log(evts, user?.email)
            } catch (err: any) {
                console.log("caught exception", err)
            }
        })()}
    ,[])

    console.log("mmm", user?.email)

    return <div>
        Main
        {/* <EventNote isSignedUp={true} time="12:30" name="Название Мероприятия" category="Conference" /> */}
        {/* <CalendarDay day={1} events={[
            {
                isSignedUp: false,
                name: "Название Мероприятия",
                category: "Conference",
                time: "12:30"
            },
            {
                isSignedUp: false,
                name: "Название Мероприятия",
                category: "Conference",
                time: "12:30"
            },
            {
                isSignedUp: true,
                name: "Название Мероприятия",
                category: "Conference",
                time: "12:30"
            },
            {
                isSignedUp: false,
                name: "Название Мероприятия",
                category: "Conference",
                time: "12:30"
            },
            {
                isSignedUp: false,
                name: "Название Мероприятия",
                category: "Conference",
                time: "12:30"
            },
            {
                isSignedUp: false,
                name: "Название Мероприятия",
                category: "Conference",
                time: "12:30"
            }
        ]}/> */}
        <Calendar before={new Date("2025-05-15T09:00:01Z")} after={new Date("2025-05-15T08:00:00Z")} />
    </div>;
}