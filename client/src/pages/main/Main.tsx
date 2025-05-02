import { Button, Typography } from "@mui/material";
import styles from "./Main.module.scss"
import AddIcon from "../../assets/add-icon.svg";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { WeekContent } from "../../components/week-content/week-content";

function getMonday(d = new Date()) {
    const date = new Date(d);
    const day = date.getDay();
    const diff = date.getDate() - day + (day === 0 ? -6 : 1);
    date.setDate(diff);
    return date;
  }

export const MainPage : React.FC = () => {
    const [monday, setMonday] = useState<Date>(getMonday(new Date()));

    const navigate = useNavigate();


    return <>
        <div className={styles.eventBar}>
            <Typography variant="h5">События</Typography>
            <Button className={styles.addButton} onClick={() => navigate("/create")}>
                <img src={AddIcon} alt="AddIcon" />
                <Typography variant="body2">
                    Добавить
                </Typography>
            </Button>
        </div>
        <WeekContent from={monday} setFrom={setMonday}></WeekContent>
    </>;
}