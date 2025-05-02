import { NavigateBefore, NavigateNext } from "@mui/icons-material";
import { Dispatch, FC, SetStateAction, useEffect, useState } from "react";
import styles from './week-content.module.scss';
import { Calendar } from "../calendar/calendar";
import { EventMap } from "../events-map/events-map";
import { ToggleButtonGroup, ToggleButton } from "@mui/material";

interface WeekContentProps {
  from: Date;
  setFrom : Dispatch<SetStateAction<Date>>;
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

  function getMonthDate(startDate: Date){
        const MONTH_NAMES = [
          "", "ЯНВАРЯ", "ФЕВРАЛЯ", "МАРТА", "АПРЕЛЯ", "МАЯ", "ИЮНЯ",
          "ИЮЛЯ", "АВГУСТА", "СЕНТЯБРЯ", "ОКТЯБРЯ", "НОЯБРЯ", "ДЕКАБРЯ"
        ];
      
        const isMonday = startDate.getDay() === 1 || (startDate.getDay() === 0 && startDate.toLocaleDateString('ru-RU', { weekday: 'long' }) === 'понедельник');
      
        if (!isMonday) {
          throw new Error("Дата должна быть понедельником");
        }
      
        const endDate = new Date(startDate);
        endDate.setDate(startDate.getDate() + 6);
      
        const startDay = startDate.getDate();
        const endDay = endDate.getDate();
        const startMonth = MONTH_NAMES[startDate.getMonth() + 1];
        const endMonth = MONTH_NAMES[endDate.getMonth() + 1];
      
        if (startDate.getMonth() === endDate.getMonth()) {
          return `${startDay} - ${endDay} ${endMonth}`;
        } else {
          return `${startDay} ${startMonth} - ${endDay} ${endMonth}`;
        }
  }

export const WeekContent : FC<WeekContentProps> = ({from, setFrom}) =>{
    const [after, setAfter] = useState<Date>(from);
    const [before, setBefore] = useState<Date>(() => {
        const copy = new Date(from);
        copy.setDate(copy.getDate() + 7);
        return copy;
    });

    const [view, setView] = useState('calendar');
    const [month, setMonth] = useState<string>("");
    const [monthdate, setMonthDate] = useState<string>("");
    
    useEffect(()=>{
        setMonth(getMonthOfWeek(after));
        setMonthDate(getMonthDate(after));
    }, [after]);

    const handleChange = (
        event: React.MouseEvent<HTMLElement>,
        newView: string,
    ) => {
        setView(newView);
    };
    
    return (
        <div className={styles.weekCont}>
        <div className={styles.weekHeader}>
            <ToggleButtonGroup
                size="small" 
                color="primary"
                value={view}
                exclusive
                onChange={handleChange}
                aria-label="view"
                >
                <ToggleButton value="calendar" sx={{lineHeight: 1}}>Календарь</ToggleButton>
                <ToggleButton value="map" sx={{lineHeight: 1}}>Карта</ToggleButton>
            </ToggleButtonGroup>
            <div className={styles.monthCaption}>{view === 'calendar' ? month : monthdate}</div>
            <div style={{width: '154.52px'}}/>
            </div>
        
    <div className={`${styles.week} ${styles[view]}`}>
        <NavigateBefore onClick={() => {
            const newAfter = new Date(after);
            newAfter.setDate(newAfter.getDate() - 7);
            
            const newBefore = new Date(before);
            newBefore.setDate(newBefore.getDate() - 7);
            
            setAfter(newAfter);
            setBefore(newBefore);
            setFrom(newAfter);
        }} />
            {view === 'calendar' && <Calendar before={before} after={after}/>}
            {view === 'map' && <EventMap before={before} after={after}/>}
        <NavigateNext onClick={() => {
            const newAfter = new Date(after);
            newAfter.setDate(newAfter.getDate() + 7);
            
            const newBefore = new Date(before);
            newBefore.setDate(newBefore.getDate() + 7);
            
            setAfter(newAfter);
            setBefore(newBefore);
            setFrom(newAfter);
        }} />
        </div>
        </div>
        );
      };