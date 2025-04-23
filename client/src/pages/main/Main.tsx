import { Button, Typography } from "@mui/material";
import { Calendar } from "../../components/calendar/calendar";
import styles from "./Main.module.scss"
import AddIcon from "../../assets/add-icon.svg";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { EventMap } from "../../components/events-map/events-map";


function getMonthString(date: Date): string {
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
    };
};

const oneDayMs = 24 * 60 * 60 * 1000;

// useEffect(
//     ()=>{(async function fun() {
//         try {
//             setVisitedEventIds((await getEventsWithFilters({
//                 before: before,
//                 after: after,
//                 visitorEmail: user == undefined ? "." : user.email,
//               })).map(evt => evt.id))
//             const copy = new Date(after);
//             copy.setDate(copy.getDate() + 3);
//             setFrom(copy)
//         } catch (err: any) {
//             console.log("caught exception", err)
//         }
//     })()}
// ,[user, before, after])
function getMonday(d = new Date()) {
    const date = new Date(d);
    const day = date.getDay();
    const diff = date.getDate() - day + (day === 0 ? -6 : 1);
    date.setDate(diff);
    return date;
  }

  function getMonthOfWeek(mondayDate: Date) {
    const thursday = new Date(mondayDate);
    thursday.setDate(mondayDate.getDate() + 3);
    switch (thursday.getMonth()) {
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
    };
  }

export const MainPage : React.FC = () => {
    const [monday, setMonday] = useState<Date>(getMonday(new Date()));
    console.log(monday)
    const [month, setMonth] = useState<string>("");
    
    useEffect(()=>{
        setMonth(getMonthOfWeek(monday));
        console.log("hehe", monday, getMonthOfWeek(monday))
    }, [monday]);

    const navigate = useNavigate();


    return <>
        <div className={styles.eventBar}>
            <Typography variant="h5">События</Typography>
            <div className={styles.monthCaption} >{month}</div>
            <Button className={styles.addButton} onClick={() => navigate("/create")}>
                <img src={AddIcon} alt="AddIcon" />
                <Typography variant="body2">
                    Добавить
                </Typography>
            </Button>
        </div>
        <Calendar from={monday} setFrom={setMonday}/>
        <EventMap after={monday} before={new Date(monday.getTime() + 7 * oneDayMs)}/>
    </>;
}