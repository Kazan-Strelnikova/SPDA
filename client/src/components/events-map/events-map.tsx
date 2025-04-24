import { useEffect, useState } from "react";
import {
  MapContainer,
  Marker,
  Popup,
  TileLayer,
  useMapEvents,
} from "react-leaflet";
import { Event } from "../../types";
import { getEventsWithFilters } from "../../http/get-events-with-filters";
import { EventCont, EventModal } from "../../pages/Event";
import { getDateString } from "../../utils/get-date-string";
import { getCategoryIcon } from "../../utils/get-category-icon";
import { Box, Typography } from "@mui/material";
import CalendarIcon from '../../assets/calendar.svg';
import TimeIcon from '../../assets/time.svg';
import SeatsIcon from '../../assets/seats.svg';
import styles from './events-map.module.scss';

interface EventMapProps {
  before: Date;
  after: Date;
}

const ZOOM_RADIUS_MAP: Record<number, number> = {
  18: 500, 17: 1000, 16: 2000, 15: 4000,
  14: 4 * 1500, 13: 4 * 3000, 12: 4 * 6000, 11: 60000,
  10: 4 * 25000, 9: 4 * 50000, 8: 4 * 100000, 7: 4 * 200000,
  6: 4 * 400000, 5: 4 * 800000, 4: 4 * 1600000, 3: 4 * 3200000,
  2: 4 * 6400000, 1: 4 * 12800000,
};

// Hook component to update geo filters on map move
function MapLocationUpdater({
  onChange,
}: {
  onChange: (filters: { lat: number; lon: number; radius: number }) => void;
}) {
  useMapEvents({
    moveend(e) {
      const center = e.target.getCenter();
      const zoom = e.target.getZoom();
      const radius = ZOOM_RADIUS_MAP[zoom] || 1000;
      onChange({ lat: center.lat, lon: center.lng, radius });
    },
  });

  return null;
}

export const EventMap = ({ before, after }: EventMapProps) => {
  const [userLocation, setUserLocation] = useState<[number, number] | null>(null);
  const [events, setEvents] = useState<Event[]>([]);
  const [geoFilters, setGeoFilters] = useState<{ lat: number; lon: number; radius: number } | null>(null);

  useEffect(() => {
    navigator.geolocation.getCurrentPosition(
      (position) => {
        setUserLocation([
          position.coords.latitude,
          position.coords.longitude,
        ]);
        setGeoFilters({ lat: position.coords.latitude, lon: position.coords.longitude, radius: 1000 });

      },
      (err) => {
        console.error("Geolocation error:", err);
        setUserLocation([55.751244, 37.618423]); // Moscow
      }
    );
  }, []);

  useEffect(() => {
    if (!geoFilters) return;

    const fetchEvents = async () => {
      try {
        const fetched = await getEventsWithFilters({
          latitude: geoFilters.lon,
          longtitude: geoFilters.lat,
          radius: geoFilters.radius,
          after,
          before,
        });
        setEvents(fetched);
        console.log(fetched)
      } catch (e) {
        console.error("Failed to fetch events:", e);
      }
    };

    fetchEvents();
  }, [geoFilters, after, before]);


  if (!userLocation) return <p>Loading map...</p>;

  return (
    <MapContainer
      center={userLocation}
      zoom={13}
      style={{ height: "100%", width: "100%" }}
    >
      <TileLayer
        attribution='&copy; OpenStreetMap contributors'
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
      />
      <MapLocationUpdater onChange={setGeoFilters} />
      {events.map((event) => (
        <Marker
          key={event.id}
          position={[event.location[0], event.location[1]]}
        >
          <Popup className={styles.popup}>
            <div>
              <Box sx={{ display: 'flex', flexDirection: 'row', justifyContent: 'space-between', width: '100%' }}>
                <Typography id="modal-modal-title" variant="subtitle1" fontWeight={600}>
                  {event?.title}
                </Typography>
                {getCategoryIcon(event.type)}
              </Box>
              <Box sx={{ display: 'flex', flexDirection: 'column', width: '100%', gap: '10px' }}>
                <Typography variant='body2'>
                  {event.description}
                </Typography>
                <Box sx={{ display: 'flex', flexDirection: 'row', width: '100%', gap: '5px', alignItems: 'center' }}>
                  <img src={CalendarIcon} alt="Calendar" />
                  <Typography variant='body2'>
                    {getDateString(event.date)}
                  </Typography>
                </Box>
                <Box sx={{ display: 'flex', flexDirection: 'row', width: '100%', gap: '5px', alignItems: 'center' }}>
                  <img src={TimeIcon} alt="Time" />
                  <Typography variant='body2'>
                    {event.date.toTimeString().slice(0, 5)}
                  </Typography>
                </Box>
                <Box sx={{ display: 'flex', flexDirection: 'row', width: '100%', gap: '5px', alignItems: 'center' }}>
                  <img src={SeatsIcon} alt="Seats" />
                  <Typography variant='body2'>
                    {event.has_unlimited_seats ? 'Количество мест не ограничено' : `Осталось свободных мест: ${event.available_seats}`}
                  </Typography>
                </Box>
              </Box>
            </div>

          </Popup>
        </Marker>
      ))}
    </MapContainer>
  );
};
