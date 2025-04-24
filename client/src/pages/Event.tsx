import * as React from 'react';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import Modal from '@mui/material/Modal';
import { getCategoryIcon } from "../utils/get-category-icon";
import { EventProps } from "../components/event-note/event-note";
import variables from '../variables.module.scss';
import CalendarIcon from '../assets/calendar.svg';
import TimeIcon from '../assets/time.svg';
import SeatsIcon from '../assets/seats.svg';
import PlaceIcon from '../assets/place.svg';

import {
  MapContainer,
  Marker,
  TileLayer,
} from "react-leaflet";
import { ButtonAKAM } from '../components/button/button';
import { Button, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle } from '@mui/material';
import { deleteEventById } from '../http/delete-event';
import { UUID } from 'crypto';
import { useNavigate } from 'react-router-dom';

interface EventModalProps {
    event: EventProps;
    open: boolean;
    handleClose: () => void;
    isSignedUp?: boolean;
    createdByUser?: boolean;
}

function Delete({id}: {id: UUID}) {
  const navigate = useNavigate();
  const [open, setOpen] = React.useState(false);
  const handleOpen = () => {
    setOpen(true);
  };
  const handleClose = () => {
    setOpen(false);
  };

  const handleAgreeClose = () => {
    deleteEventById(id);
    setOpen(false);
    window.location.reload();
  };

  return (
    <React.Fragment>
      <ButtonAKAM outlined onClick={handleOpen}>Удалить</ButtonAKAM>
      <Dialog
        open={open}
        onClose={handleClose}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
      >
        <DialogTitle id="alert-dialog-title">
          {"Удалить событие?"}
        </DialogTitle>
        <DialogContent>
          <DialogContentText id="alert-dialog-description">
            При удалении события его не получится восстановить автоматически
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose}>Отменить</Button>
          <Button onClick={handleAgreeClose} autoFocus>
            Подтвердить
          </Button>
        </DialogActions>
      </Dialog>
    </React.Fragment>
  );
}

export const EventModal : React.FC<EventModalProps> = ({event, open, handleClose, isSignedUp=false, createdByUser=false} : EventModalProps) => {
    return (
        <>
        <Modal
        open={open}
        onClose={handleClose}
        aria-labelledby="modal-modal-title"
        aria-describedby="modal-modal-description"
      >
        <Box sx={style}>
          <Box sx={{display: 'flex', flexDirection: 'row', justifyContent: 'space-between', width: '100%'}}>
          <Typography id="modal-modal-title" variant="subtitle1" fontWeight={600}>
            {event?.title}
          </Typography>
          {getCategoryIcon(event.category)}
          </Box>
          <Box sx={{display: 'flex', flexDirection: 'column', width: '100%', gap: '10px'}}>
            <Typography variant='body2'>
              {event.description}
            </Typography>
            <Box sx={{display: 'flex', flexDirection: 'row', width: '100%', gap: '5px', alignItems: 'center'}}>
              <img src={CalendarIcon} alt="Calendar" />
              <Typography variant='body2'>
                {event.date}
              </Typography>
            </Box>
            <Box sx={{display: 'flex', flexDirection: 'row', width: '100%', gap: '5px', alignItems: 'center'}}>
              <img src={TimeIcon} alt="Time" />
              <Typography variant='body2'>
                {event.time}
              </Typography>
            </Box>
            <Box sx={{display: 'flex', flexDirection: 'row', width: '100%', gap: '5px', alignItems: 'center'}}>
              <img src={SeatsIcon} alt="Seats" />
              <Typography variant='body2'>
                {event.seats === -1 ? 'Количество мест не ограничено' : `Осталось свободных мест: ${event.seats}`}
              </Typography>
            </Box>
            {event.location && <Box sx={{display: 'flex', flexDirection: 'column', width: '100%', gap: '5px', alignItems: 'center'}}>
            <Box sx={{display: 'flex', flexDirection: 'row', width: '100%', gap: '5px', alignItems: 'center'}}>
              <img src={PlaceIcon} alt="Place" />
              <Typography variant='body2'>
                Локация
              </Typography>
              </Box>
              <MapContainer
                    center={event.location}
                    zoom={15}
                    style={{ height: "300px", width: "100%" }}
                  >
                    <TileLayer
                      attribution='&copy; OpenStreetMap contributors'
                      url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                    />
                    <Marker
                      key={event.id}
                      position={[event.location[0], event.location[1]]}
                    />
                  </MapContainer>
            </Box>}

          </Box>
          <Box sx={{display: 'flex', flexDirection: 'row', gap: '15px', alignItems: 'center'}}>
            {createdByUser ? <>
              <Delete id={event.id as UUID}/>
              <ButtonAKAM filled>Изменить</ButtonAKAM>
            </> :
            (isSignedUp ? 
              <ButtonAKAM outlined>Отменить запись</ButtonAKAM>
              :
              <ButtonAKAM filled>Зарегистрироваться</ButtonAKAM>
            )}
          </Box>
        </Box>
      </Modal>
        </>
    );
}

const style = {
  position: 'absolute',
  top: '50%',
  left: '50%',
  transform: 'translate(-50%, -50%)',
  width: 400,
  bgcolor: variables.background,
  borderRadius: '12px',
  p: '25px',
  BorderInner: 'none',
  outline: 'none',
  display: 'flex',
  flexDirection: 'column',
  gap: '10px',
  boxShadow: '0 4px 16px rgba(29,51,113, 0.4)',
  alignItems: 'end',
};
