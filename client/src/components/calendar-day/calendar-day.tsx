import { FC } from "react"
import { EventNote, EventNoteProps } from "../event-note/event-note";
import styles from "./calendar-day.module.scss"
import { Chip, Divider, Typography } from "@mui/material";

interface CalendarDayProps {
    day : number;
    events : EventNoteProps[];
}

export const CalendarDay : FC<CalendarDayProps> = ({day, events}) => {
    return (
    <div className={styles.calendarDay}>
        
        <Chip className={styles.dayIndex} label={day}/>
        
        <div className={styles.eventsContainer}>
            {events.length <= 0 
            ? 
            <>
                <Divider className={styles.divider}/>
                <Typography variant="body2" className={styles.noEventsTypography}>События отсутствуют</Typography>
            </>
            : events.map((event, idx) => 
                <>
                    <Divider className={styles.divider}/>
                    <EventNote isSignedUp={event.isSignedUp} name={event.name} time={event.time} category={event.category}></EventNote>
                    {idx === events.length - 1 && <Divider className={styles.divider}/>}
                </>
            )}
        </div>
    </div>);
}