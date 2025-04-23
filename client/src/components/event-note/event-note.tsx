import { FC, useState } from "react";
import { Category } from "../../types";
import styles from "./event-note.module.scss"
import { Typography } from "@mui/material";
import { getCategoryIcon } from "../../utils/get-category-icon";
import { EventModal } from "../../pages/Event";

export interface EventProps {
    id: string;
    title: string;
    date: string;
    time : string;
    name : string;
    category : Category;
    description?: string;
    location: [number, number];
    seats: number;
}

export interface EventNoteProps {
    isSignedUp : boolean;
    event: EventProps;
    createdByUser?: boolean;
}

export const EventNote : FC<EventNoteProps> = ({isSignedUp, event, createdByUser}) => {
    const [open, setOpen] = useState<boolean>(false);
    const handleOpen = () => setOpen(true);
    const handleClose = () => setOpen(false);

    return (<>
        <div className={`${styles.eventNote} ${isSignedUp ? styles.background : ""}`} onClick={handleOpen}>
            <div className={styles.top}>
                <Typography variant="subtitle2" className={styles.timeCaption}>{event.time}</Typography>
                {getCategoryIcon(event.category)}
            </div>
            <Typography  variant="subtitle2" className={styles.titleCaption}>{event.name}</Typography>
        </div>
        <EventModal event={event} open={open} handleClose={handleClose} isSignedUp={isSignedUp} createdByUser={createdByUser}/>
    </>);
}  