import { Button, Typography } from "@mui/material";
import { Calendar } from "../../components/calendar/calendar";
import styles from "./Main.module.scss"
import AddIcon from "../../assets/add-icon.svg";
import { useState } from "react";

function getMonth(date: Date): string {
    switch (date.getMonth()) {
        default:
            return ""
        case 0: return "Январь"
        case 1: return "Февраль"
        case 2: return "Март"
        case 3: return "Апрель"
        case 4: return "Май"
        case 5: return "Июнь"
        case 6: return "Июль"
        case 7: return "Август"
        case 8: return "Сентябрь"
        case 9: return "Октябрь"
        case 10: return "Ноябрь"
        case 11: return "Декабрь"
    }
}

export const MainPage : React.FC = () => {

    const [today, setToday] = useState<Date>(new Date())

    return <div>
        <div className={styles.eventBar}>
            <Typography variant="h5">События</Typography>
            <div className={styles.monthCaption} >{getMonth(today)}</div>
            <Button className={styles.addButton}>
                <img src={AddIcon} alt="AddIcon" />
                <Typography variant="body2">
                    Добавить
                </Typography>
            </Button>
        </div>
        <Calendar from={today} setFrom={setToday}/>
    </div>;
}