import { useEffect, useState } from "react";
import {
  MapContainer,
  TileLayer,
  Marker,
  useMapEvents,
} from "react-leaflet";

import "leaflet/dist/leaflet.css";
import "leaflet-defaulticon-compatibility";
import "leaflet-defaulticon-compatibility/dist/leaflet-defaulticon-compatibility.css";
import { LatLngExpression } from "leaflet";

type LatLng = {
  lat: number;
  lng: number;
};

export default function LocationPicker({
  onLocationChange,
}: {
  onLocationChange: (coords: LatLng) => void;
}) {
  const [position, setPosition] = useState<LatLng | null>(null);

  useEffect(() => {
    navigator.geolocation.getCurrentPosition(
      (pos) => {
        setPosition({
          lat: pos.coords.latitude,
          lng: pos.coords.longitude,
        });
      },
      () => {
        setPosition({ lat: 40.7128, lng: -74.006 });
      }
    );
  }, []);

  function LocationMarker() {
    useMapEvents({
      click(e: any) {
        const { lat, lng } = e.latlng;
        setPosition({ lat, lng });
        onLocationChange({ lat, lng });
      },
    });

    return position ? <Marker position={position} /> : null;
  }

  return position ? ( // я починил, не хватало пакета с их типапи @types/react-leaflet
    <div style={{ height: '400px', width: '100%' }}>
        <MapContainer
        center={position as LatLngExpression}   // Так странно, у них в доке есть все эти поля в примерах, а уменя ругается
        zoom={13} //версия другая мож  https://react-leaflet.js.org/docs/api-map/
        scrollWheelZoom={false} // v5.x, у нас 5.0.0 
        style={{ height: "400px", width: "100%" }}
        >
        <TileLayer
            attribution='&copy; <a href="https://www.openstreetmap.org/">OpenStreetMap</a>'
            url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        />
        <LocationMarker />
        </MapContainer>
    </div>
  ) : (
    <p>Loading map...</p>
  );
}
