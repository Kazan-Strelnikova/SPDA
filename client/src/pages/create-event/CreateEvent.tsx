import { styled, TextField, Typography } from "@mui/material";
import React, { ChangeEvent, FC, useState } from "react";
import BackIcon from "../../assets/back.svg";
import style from "./CreateEvent.module.scss"
import { Event } from "../../types";
import { DatePicker, DateValidationError, LocalizationProvider, PickerChangeHandlerContext, TimePicker, TimeValidationError } from "@mui/x-date-pickers";
import { AdapterDateFns } from "@mui/x-date-pickers/AdapterDateFns";
import { Dropdown, Menu, MenuButton as BaseMenuButton, MenuItem as BaseMenuItem, menuItemClasses, MenuListboxSlotProps, PopupContext, CssTransition} from "@mui/base";
import { PickerValue } from "@mui/x-date-pickers/internals";

const Categories = [
    ["Конференция", "Conference"],
    ["Митап", "Meetup"],
    ["Концерт", "Concert"],
    ["Выставка", "Exhibition"],
    ["Вечеринка", "Party"],
    ["Спорт", "Sport"],
    ["Образование", "Education"],
    ["Соревнование", "Competition"],
    ["Другое", "Other"]
]


export const CreateEventPage : FC = () => {
    
    const [menuOpen, setMenuOpen] = useState<boolean>()
    const [evt, setEvt] = useState<Event>()
    const createHandleMenuClick = (menuItem: string) => {

        return () => setEvt(prev => ({
            ...prev,
            type: menuItem
        } as Event));
    };
    
    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setEvt({
            ...evt,
            [name]: value
        } as Event);
        console.log(evt);
    };

    const handleDateChange = (
      value: PickerValue,
      _context: PickerChangeHandlerContext<DateValidationError>
    ) => {
      if (value !== null) {
        const newDate = new Date(value); 
        if (evt?.date) {
          
          newDate.setHours(evt.date.getHours(), evt.date.getMinutes(), 0, 0);
        }
        setEvt({
          ...evt,
          date: value,
        } as Event);
      }
    };
    
    const handleTimeChange = (
      value: PickerValue,
      _context: PickerChangeHandlerContext<TimeValidationError>
    ) => {
      if (value !== null) {
        const newTime = new Date(value); 
        const baseDate = evt?.date ?? new Date();
        
        newTime.setFullYear(baseDate.getFullYear(), baseDate.getMonth(), baseDate.getDate());
    
        setEvt({
          ...evt,
          date: value,
        } as Event);
      }
    };
    
    
    return (
        <div className={style.page}>
        <div className={style.label}>
            <img src={BackIcon} alt="back" />
            <Typography variant="h5" > Новое Событие </Typography>
        </div>

        <div className={style.box}>
            <div className={style.inputs}>
                <TextField
                    className={menuOpen ? style.clipped : ""}
                    margin="normal"                     
                    id="title"
                    label="Название"
                    name="title"
                    autoFocus
                    value={evt?.title}
                    onChange={handleChange}
                    variant="outlined"
                />
                <TextField
                    className={menuOpen ? style.clipped : ""}
                    margin="normal"                     
                    id="description"
                    label="Описание"
                    name="description"
                    autoFocus
                    value={evt?.description}
                    onChange={handleChange}
                    variant="outlined"
                    multiline
                    minRows={5}
                />
                <TextField
                    className={menuOpen ? style.clipped : ""}
                    margin="normal"                     
                    id="seats"
                    label="Количество мест"
                    name="seats"
                    autoFocus
                    value={evt?.available_seats}
                    onChange={(e) => {
                        const newValue = e.target.value;
                        if (/^\d*$/.test(newValue)) {
                            handleChange(e as ChangeEvent<HTMLInputElement>);
                        }
                    }}                  
                    variant="outlined"
                    helperText="от 1 до 150 000"
                />
                <Dropdown onOpenChange={() => setMenuOpen((prev) => (!prev))}>
                  <MenuButton>
                      {evt?.type
                      ? Categories.find((value) => value[1] === evt.type)?.[0] 
                      : "Категория"}
                  </MenuButton>
                        <Menu slots={{ listbox: AnimatedListbox }}>
                            {Categories.map((value, _idx) => {
                                return <MenuItem onClick={createHandleMenuClick(value[1])}>{value[0]}</MenuItem>
                            })}
                        </Menu>
                </Dropdown>

                <LocalizationProvider dateAdapter={AdapterDateFns}>
                    <DatePicker label="Дата" onChange={handleDateChange}/>
                </LocalizationProvider >

                <LocalizationProvider dateAdapter={AdapterDateFns}>
                    <TimePicker label="Время" onChange={handleTimeChange}/>
                </LocalizationProvider >
            </div>

        </div>
    </div>
)}


const blue = {
    50: '#F0F7FF',
    100: '#C2E0FF',
    200: '#99CCF3',
    300: '#66B2FF',
    400: '#3399FF',
    500: '#007FFF',
    600: '#0072E6',
    700: '#0059B3',
    800: '#004C99',
    900: '#003A75',
  };
  
  const grey = {
    50: '#F3F6F9',
    100: '#E5EAF2',
    200: '#DAE2ED',
    300: '#C7D0DD',
    400: '#B0B8C4',
    500: '#9DA8B7',
    600: '#6B7A90',
    700: '#434D5B',
    800: '#303740',
    900: '#1C2025',
  };
  
  const Listbox = styled('ul')(
    ({ theme }) => `
    font-family: 'IBM Plex Sans', sans-serif;
    font-size: 0.875rem;
    box-sizing: border-box;
    padding: 6px;
    margin: 12px 0;
    min-width: 200px;
    border-radius: 12px;
    overflow: auto;
    outline: 0;
    background: ${theme.palette.mode === 'dark' ? grey[900] : '#fff'};
    border: 1px solid ${theme.palette.mode === 'dark' ? grey[700] : grey[200]};
    color: ${theme.palette.mode === 'dark' ? grey[300] : grey[900]};
    box-shadow: 0 4px 30px ${theme.palette.mode === 'dark' ? grey[900] : grey[200]};
    z-index: 1;
  
    .closed & {
      opacity: 0;
      transform: scale(0.95, 0.8);
      transition: opacity 200ms ease-in, transform 200ms ease-in;
    }
    
    .open & {
      opacity: 1;
      transform: scale(1, 1);
      transition: opacity 100ms ease-out, transform 100ms cubic-bezier(0.43, 0.29, 0.37, 1.48);
    }
  
    .placement-top & {
      transform-origin: bottom;
    }
  
    .placement-bottom & {
      transform-origin: top;
    }
    `,
  );
  
  const AnimatedListbox = React.forwardRef(function AnimatedListbox(
    props: MenuListboxSlotProps,
    ref: React.ForwardedRef<HTMLUListElement>,
  ) {
    const { ownerState, ...other } = props;
    const popupContext = React.useContext(PopupContext);
  
    if (popupContext == null) {
      throw new Error(
        'The `AnimatedListbox` component cannot be rendered outside a `Popup` component',
      );
    }
  
    const verticalPlacement = popupContext.placement.split('-')[0];
  
    return (
      <CssTransition
        className={`placement-${verticalPlacement} ${style.dropdown}`}
        enterClassName="open"
        exitClassName="closed"
      >
        <Listbox {...other} ref={ref} />
      </CssTransition>
    );
  });
  
  const MenuItem = styled(BaseMenuItem)(
    ({ theme }) => `
    list-style: none;
    padding: 8px;
    border-radius: 8px;
    cursor: default;
    user-select: none;
  
    &:last-of-type {
      border-bottom: none;
    }
  
    &:focus {
      outline: 3px solid ${theme.palette.mode === 'dark' ? blue[600] : blue[200]};
      background-color: ${theme.palette.mode === 'dark' ? grey[800] : grey[100]};
      color: ${theme.palette.mode === 'dark' ? grey[300] : grey[900]};
    }
  
    &.${menuItemClasses.disabled} {
      color: ${theme.palette.mode === 'dark' ? grey[700] : grey[400]};
    }
    `,
  );
  
  const MenuButton = styled(BaseMenuButton)(
    ({ theme }) => `
    // font-weight: 600;
    font-size: 1rem;
    line-height: 1.5;
    text-align: start;
    padding: 16.5px 14px;
    border-radius: 8px;
    transition: all 150ms ease;
    cursor: pointer;
    background: ${theme.palette.mode === 'dark' ? grey[900] : '#fff'};
    border: 1px solid ${theme.palette.mode === 'dark' ? grey[700] : grey[400]};
    color: ${theme.palette.mode === 'dark' ? grey[200] : grey[700]};
    box-shadow: 0 1px 2px 0 rgb(0 0 0 / 0.05);
  
    &:hover {
      background: ${theme.palette.mode === 'dark' ? grey[800] : grey[50]};
      border-color: ${theme.palette.mode === 'dark' ? grey[600] : grey[300]};
    }
  
    &:active {
      background: ${theme.palette.mode === 'dark' ? grey[700] : grey[100]};
    }
  
    &:focus-visible {
      box-shadow: 0 0 0 4px ${theme.palette.mode === 'dark' ? blue[300] : blue[200]};
      outline: none;
    }
    `,
  );