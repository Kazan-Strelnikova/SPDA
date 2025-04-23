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

  useEffect(() => {
    if (position){
      onLocationChange(position);
    }
  }, [position]);

  function LocationMarker() {
    useMapEvents({
      click(e: any) {
        const { lat, lng } = e.latlng;
        setPosition({ lat, lng });
      },
    });

    return position ? <Marker position={position} /> : null;
  }

  return position ? (
    <div style={{ height: '360px', width: '480px', zIndex: 0 }}>
        <MapContainer
        center={position as LatLngExpression}
        zoom={13}
        scrollWheelZoom={false}
        style={{ height: "100%", width: "100%" }}
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
