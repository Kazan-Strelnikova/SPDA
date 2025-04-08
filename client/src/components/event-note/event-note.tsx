import { FC } from "react";
import { Category } from "../../types";
import styles from "./event-note.module.scss"
import { Typography } from "@mui/material";
import { getCategoryIcon } from "../../utils/get-category-icon";

export interface EventNoteProps {
    isSignedUp : boolean;
    time : string;
    name : string;
    category : Category;
}

export const EventNote : FC<EventNoteProps> = (props) => {
    return (<>
        <div className={`${styles.eventNote} ${props.isSignedUp ? styles.background : ""}`}>
            <div className={styles.top}>
                <Typography variant="subtitle2" className={styles.timeCaption}>{props.time}</Typography>
                {getCategoryIcon(props.category)}
            </div>
            <Typography  variant="subtitle2" className={styles.titleCaption}>{props.name}</Typography>
        </div>
    </>);
}  