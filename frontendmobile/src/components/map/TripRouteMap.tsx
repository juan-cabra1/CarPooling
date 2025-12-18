
import React, { useState, useEffect } from 'react';
import { View, StyleSheet } from 'react-native';
import Constants from 'expo-constants';
import MapView from './MapView';
import { Coordinates } from '../../types';
import { PolylineProps } from 'react-native-maps';

interface TripRouteMapProps {
  origin: Coordinates;
  destination: Coordinates;
}

const GOOGLE_MAPS_API_KEY = Constants.expoConfig?.extra?.googleMapsApiKey;

const TripRouteMap: React.FC<TripRouteMapProps> = ({ origin, destination }) => {
  const [route, setRoute] = useState<Coordinates[] | null>(null);

  useEffect(() => {
    const getDirections = async () => {
      if (!GOOGLE_MAPS_API_KEY) {
        console.error('Google Maps API key is not set');
        return;
      }

      const url = `https://maps.googleapis.com/maps/api/directions/json?origin=${origin.lat},${origin.lng}&destination=${destination.lat},${destination.lng}&key=${GOOGLE_MAPS_API_KEY}`;

      try {
        const response = await fetch(url);
        const json = await response.json();

        if (json.routes.length > 0) {
          const points = decode(json.routes[0].overview_polyline.points);
          setRoute(points);
        }
      } catch (error) {
        console.error('Error fetching directions:', error);
      }
    };

    getDirections();
  }, [origin, destination]);

  // Decode polyline from Google Directions API
  const decode = (t: string): Coordinates[] => {
    let points: Coordinates[] = [];
    for (let step = 0; step < t.length; ) {
      let lat = 0;
      let lng = 0;
      let shifter = 0;
      let result = 0;
      let byte = 0;
      do {
        byte = t.charCodeAt(step++) - 63;
        result |= (byte & 0x1f) << shifter;
        shifter += 5;
      } while (byte >= 0x20);
      let dlat = (result & 1) != 0 ? ~(result >> 1) : (result >> 1);
      lat += dlat;

      shifter = 0;
      result = 0;
      do {
        byte = t.charCodeAt(step++) - 63;
        result |= (byte & 0x1f) << shifter;
        shifter += 5;
      } while (byte >= 0x20);
      let dlng = (result & 1) != 0 ? ~(result >> 1) : (result >> 1);
      lng += dlng;

      points.push({ lat: lat / 1e5, lng: lng / 1e5 });
    }
    return points;
  };

  const markers = [
    { coordinate: { latitude: origin.lat, longitude: origin.lng }, title: 'Origin' },
    { coordinate: { latitude: destination.lat, longitude: destination.lng }, title: 'Destination' },
  ];

  const polylines: PolylineProps[] = route ? [{
    coordinates: route.map(p => ({ latitude: p.lat, longitude: p.lng })),
    strokeColor: '#000',
    strokeWidth: 3,
  }] : [];

  return (
    <View style={styles.container}>
      <MapView
        markers={markers}
        polylines={polylines}
        region={{
            latitude: (origin.lat + destination.lat) / 2,
            longitude: (origin.lng + destination.lng) / 2,
            latitudeDelta: Math.abs(origin.lat - destination.lat) + 0.05,
            longitudeDelta: Math.abs(origin.lng - destination.lng) + 0.05,
        }}
      />
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    ...StyleSheet.absoluteFillObject,
    justifyContent: 'flex-end',
    alignItems: 'center',
  },
});

export default TripRouteMap;
