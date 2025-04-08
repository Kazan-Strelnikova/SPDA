import { NavigateBefore, NavigateNext } from "@mui/icons-material";
import { FC, useContext, useEffect, useState } from "react"
import { getEventsWithFilters } from "../../http/get-events-with-filters";
import { Event } from "../../types";
import { CalendarDay } from "../calendar-day/calendar-day";
import { EventNoteProps } from "../event-note/event-note";
import styles from "./calendar.module.scss"
import { UUID } from "crypto";
import { UserContext } from "../../contexts/UserContext";

interface CalendarProps {
    before : Date;
    after : Date;
}

export const Calendar : FC<CalendarProps> = ({before, after}) => {
    const [events, setEvents] = useState<Event[]>();
    const [dayBuckets, setDayBuckets] = useState<Event[][]>([]);
    const [visitedEventIds, setVisitedEventIds] = useState<UUID[]>();
    const userContext = useContext(UserContext);
    if (!userContext) {
        throw new Error("Calendar must be used within a UserProvider");
    }
    const { user, setUser } = userContext;
  
    useEffect(() => {
      (async function fetchEvents() {
        try {
          const evts = await getEventsWithFilters({
            before: before,
            after: after,
          });
          setEvents(evts);
  
          const buckets: Event[][] = Array.from({ length: 7 }, () => []);
  
          for (const event of evts) {
            const eventDate = new Date(event.date);
            const dayIndex = Math.floor(
              (eventDate.getTime() - after.getTime()) / (1000 * 60 * 60 * 24)
            );
  
            if (dayIndex >= 0 && dayIndex <= 6) {
              buckets[dayIndex].push(event);
            }
          }
  
          setDayBuckets(buckets);
        } catch (err: any) {
          console.log("caught exception", err);
        }
      })();
    }, [before, after]);

    useEffect(
        ()=>{(async function fun() {
            try {
                setVisitedEventIds((await getEventsWithFilters({
                    before: before,
                    after: after,
                    visitorEmail: user?.email,
                  })).map(evt => evt.id))
                  console.log(visitedEventIds)
            } catch (err: any) {
                console.log("caught exception", err)
            }
        })()}
    ,[user])

    return (
    <div className={styles.calendar}>
        <NavigateBefore />
        <div className={styles.calendarBox}>
            {dayBuckets.map((dayEvents, idx) => <CalendarDay day={idx + 1} events={dayEvents.map<EventNoteProps>(
                function(evt, _idx, _arr): EventNoteProps {
                    const incl = visitedEventIds?.includes(evt.id)
                    console.log(incl, visitedEventIds, user?.email)
                    return {
                        isSignedUp: incl == undefined ? false : incl,
                        name: evt.title,
                        time: evt.date.toTimeString().slice(0, 5),
                        category: evt.type,
                    }
                }
            )} />
            )}
        </div>
        <NavigateNext />
    </div>)
}