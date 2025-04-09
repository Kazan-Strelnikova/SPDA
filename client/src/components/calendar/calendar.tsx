import { NavigateBefore, NavigateNext } from "@mui/icons-material";
import { Dispatch, FC, SetStateAction, useContext, useEffect, useState } from "react"
import { getEventsWithFilters } from "../../http/get-events-with-filters";
import { Event } from "../../types";
import { CalendarDay } from "../calendar-day/calendar-day";
import { EventNoteProps } from "../event-note/event-note";
import styles from "./calendar.module.scss"
import { UUID } from "crypto";
import { UserContext } from "../../contexts/UserContext";

interface CalendarProps {
    from : Date;
    setFrom : Dispatch<SetStateAction<Date>>
}

function getDaysInMonth(date: Date): number {
    const year = date.getFullYear();
    const month = date.getMonth();
    return new Date(year, month + 1, 0).getDate();
}


export const Calendar : FC<CalendarProps> = ({from, setFrom}) => {
    const [dayBuckets, setDayBuckets] = useState<Event[][]>([]);
    const [visitedEventIds, setVisitedEventIds] = useState<UUID[]>();
    const [after, setAfter] = useState<Date>(from)
    const [before, setBefore] = useState<Date>(() => {
        const copy = new Date(from);
        copy.setDate(copy.getDate() + 7);
        return copy;
    });
      

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
                    visitorEmail: user == undefined ? "." : user.email,
                  })).map(evt => evt.id))
                const copy = new Date(after);
                copy.setDate(copy.getDate() + 3);
                setFrom(copy)
            } catch (err: any) {
                console.log("caught exception", err)
            }
        })()}
    ,[user, before, after])

    return (
    <div className={styles.calendar}>
        <NavigateBefore onClick={() => {
            const newAfter = new Date(after);
            newAfter.setDate(newAfter.getDate() - 7);

            const newBefore = new Date(before);
            newBefore.setDate(newBefore.getDate() - 7);

            setAfter(newAfter);
            setBefore(newBefore);
        }} />
        <div className={styles.calendarBox}>
            {dayBuckets.map((dayEvents, idx) => <CalendarDay 
            
            day={
                (idx + after.getDate()) > getDaysInMonth(after) 
                ? (idx + after.getDate()) % getDaysInMonth(after) 
                : (idx + after.getDate())
            } 
            
            events={dayEvents
              .sort((a, b) => a.date.getTime() - b.date.getTime())
              .map<EventNoteProps>(
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
        <NavigateNext onClick={() => {
            const newAfter = new Date(after);
            newAfter.setDate(newAfter.getDate() + 7);

            const newBefore = new Date(before);
            newBefore.setDate(newBefore.getDate() + 7);

            setAfter(newAfter);
            setBefore(newBefore);
        }} />
    </div>)
}