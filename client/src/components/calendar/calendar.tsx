import { FC, useContext, useEffect, useState } from "react"
import { getEventsWithFilters } from "../../http/get-events-with-filters";
import { Event } from "../../types";
import { CalendarDay } from "../calendar-day/calendar-day";
import { EventNoteProps } from "../event-note/event-note";
import styles from "./calendar.module.scss"
import { UUID } from "crypto";
import { UserContext } from "../../contexts/UserContext";
import { getDateString } from "../../utils/get-date-string";

interface CalendarProps {
    before : Date;
    after : Date;
}

function getDaysInMonth(date: Date): number {
    const year = date.getFullYear();
    const month = date.getMonth();
    return new Date(year, month + 1, 0).getDate();
}


export const Calendar : FC<CalendarProps> = ({before, after}) => {
    const [dayBuckets, setDayBuckets] = useState<Event[][]>([]);
    const [visitedEventIds, setVisitedEventIds] = useState<UUID[]>();
      

    const userContext = useContext(UserContext);
    if (!userContext) {
        throw new Error("Calendar must be used within a UserProvider");
    }
    const { user } = userContext;
  
    useEffect(() => {
      (async function fetchEvents() {
        try {
          const evts = await getEventsWithFilters({
            before: before,
            after: after,
          });
  
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
                    visitorEmail: user === undefined ? "." : user.email,
                  })).map(evt => evt.id))
            } catch (err: any) {
                console.log("caught exception", err)
            }
        })()}
    ,[user, before, after])

    return (
        <div className={styles.calendarBox}>
            {dayBuckets.map((dayEvents, idx) => <CalendarDay 
            key={(idx + after.getDate())}
            
            day={
                (idx + after.getDate()) > getDaysInMonth(after) 
                ? (idx + after.getDate()) % getDaysInMonth(after) 
                : (idx + after.getDate())
            } 
            
            events={dayEvents
              .sort((a, b) => a.date.getTime() - b.date.getTime())
              .map<EventNoteProps>(evt => {
                    const incl = visitedEventIds?.includes(evt.id)
                    return {
                      isSignedUp: incl === undefined ? false : incl,
                      event: {
                          id: evt.id,
                          title: evt.title,
                          name: evt.title,
                          date: getDateString(evt.date),
                          time: evt.date.toTimeString().slice(0, 5),
                          category: evt.type,
                          description: evt.description,
                          location: evt.location,
                          seats: evt.has_unlimited_seats ? -1 : evt.available_seats,
                        },
                        createdByUser: evt.creator_email === user?.email,
                        eventObj: evt,
                    }
                }
            )} />
            )}
        </div>)
}