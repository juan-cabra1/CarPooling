import React, { useEffect, useRef } from 'react';
import { StyleSheet, View, Text } from 'react-native';
import MapView, { Marker, Polyline } from 'react-native-maps';
import { useTrackingContext } from '../../context/TrackingContext';

const LiveTrackingMap: React.FC = () => {
  const {
    driverLocation,
    driverHeading,
    activeTrip,
    routeCoordinates,
    routeToOrigin,
    navigationSteps,
    distanceToNextStep,
    activeRouteType,
  } = useTrackingContext();
  const mapRef = useRef<MapView>(null);

  // Debug logs
  useEffect(() => {
    console.log('🗺️ LiveTrackingMap - driverLocation:', driverLocation);
    console.log('🗺️ LiveTrackingMap - activeTrip:', activeTrip ? 'Yes' : 'No');
    console.log('🗺️ LiveTrackingMap - activeRouteType:', activeRouteType);
    console.log('🗺️ LiveTrackingMap - routeCoordinates (origin → destination):', routeCoordinates ? `${routeCoordinates.length} points` : 'none');
    console.log('🗺️ LiveTrackingMap - routeToOrigin (current → origin):', routeToOrigin ? `${routeToOrigin.length} points` : 'none');
  }, [driverLocation, activeTrip, routeCoordinates, routeToOrigin, activeRouteType]);


  // Nota: animateMarkerToCoordinate no está disponible en Expo Go
  // El marcador se actualizará automáticamente cuando cambien las coordenadas

  // Cámara inteligente con zoom ajustado según navegación
  useEffect(() => {
    if (!mapRef.current || !driverLocation) return;

    // Si hay pasos de navegación, centrar en conductor con zoom cercano
    if (navigationSteps.length > 0) {
      const baseZoom = 17;
      let adjustedZoom = baseZoom;

      // Ajustar zoom según distancia al giro
      if (distanceToNextStep !== null) {
        if (distanceToNextStep < 50) adjustedZoom = 18;
        else if (distanceToNextStep < 200) adjustedZoom = 17.5;
        else if (distanceToNextStep > 1000) adjustedZoom = 16;
      }

      mapRef.current.animateCamera(
        {
          center: { latitude: driverLocation.lat, longitude: driverLocation.lng },
          zoom: adjustedZoom,
          heading: driverHeading || 0, // Rotar mapa hacia adelante
        },
        { duration: 1000 }
      );
    } else {
      // Sin navegación activa, mostrar todas las rutas
      const allCoordinates: Array<{ latitude: number; longitude: number }> = [];

      if (driverLocation) {
        allCoordinates.push({ latitude: driverLocation.lat, longitude: driverLocation.lng });
      }

      if (routeToOrigin && routeToOrigin.length > 0) {
        allCoordinates.push(...routeToOrigin.map(c => ({ latitude: c.lat, longitude: c.lng })));
      }

      if (routeCoordinates && routeCoordinates.length > 0) {
        allCoordinates.push(...routeCoordinates.map(c => ({ latitude: c.lat, longitude: c.lng })));
      }

      if (allCoordinates.length > 1) {
        mapRef.current.fitToCoordinates(allCoordinates, {
          edgePadding: { top: 200, right: 80, bottom: 400, left: 80 }, // Account for banner & panel
          animated: true,
        });
      } else if (driverLocation) {
        mapRef.current.animateCamera({
          center: {
            latitude: driverLocation.lat,
            longitude: driverLocation.lng,
          },
          zoom: 15,
        });
      }
    }
  }, [driverLocation, driverHeading, distanceToNextStep, navigationSteps, routeCoordinates, routeToOrigin]);


  if (!activeTrip) {
    return null;
  }

  return (
      <MapView
        ref={mapRef}
        style={StyleSheet.absoluteFill}
        initialRegion={driverLocation ? {
          latitude: driverLocation.lat,
          longitude: driverLocation.lng,
          latitudeDelta: 0.0922,
          longitudeDelta: 0.0421,
        } : undefined}
        showsUserLocation={false} // Use custom marker instead
        followsUserLocation={false}
        rotateEnabled={true}
        pitchEnabled={false}
      >
        {/* Custom Driver Marker with Rotation */}
        {driverLocation && (
          <Marker
            coordinate={{
              latitude: driverLocation.lat,
              longitude: driverLocation.lng,
            }}
            anchor={{ x: 0.5, y: 0.5 }}
            flat={true}
            rotation={driverHeading || 0}
          >
            <View style={styles.vehicleMarkerContainer}>
              <Text style={styles.vehicleArrow}>▲</Text>
            </View>
          </Marker>
        )}
        {activeTrip.origin?.coordinates && (
          <Marker
            coordinate={{
              latitude: activeTrip.origin.coordinates.lat,
              longitude: activeTrip.origin.coordinates.lng,
            }}
            title="Origen"
            pinColor="green"
          />
        )}
        {activeTrip.destination?.coordinates && (
          <Marker
            coordinate={{
              latitude: activeTrip.destination.coordinates.lat,
              longitude: activeTrip.destination.coordinates.lng,
            }}
            title="Destino"
            pinColor="red"
          />
        )}
        {/* Ruta origen → destino: borde blanco + línea azul (estilo Google Maps) */}
        {routeCoordinates && routeCoordinates.length > 0 && activeRouteType === 'toDestination' && (
            <>
                <Polyline
                    coordinates={routeCoordinates.map(p => ({ latitude: p.lat, longitude: p.lng }))}
                    strokeColor="#FFFFFF"
                    strokeWidth={10}
                />
                <Polyline
                    coordinates={routeCoordinates.map(p => ({ latitude: p.lat, longitude: p.lng }))}
                    strokeColor="#1A73E8"
                    strokeWidth={7}
                />
            </>
        )}

        {/* Ruta actual → origen: borde blanco + línea naranja (estilo Google Maps) */}
        {routeToOrigin && routeToOrigin.length > 0 && activeRouteType === 'toOrigin' && (
            <>
                <Polyline
                    coordinates={routeToOrigin.map(p => ({ latitude: p.lat, longitude: p.lng }))}
                    strokeColor="#FFFFFF"
                    strokeWidth={10}
                />
                <Polyline
                    coordinates={routeToOrigin.map(p => ({ latitude: p.lat, longitude: p.lng }))}
                    strokeColor="#FF6B00"
                    strokeWidth={7}
                />
            </>
        )}
      </MapView>
  );
};

const styles = StyleSheet.create({
  vehicleMarkerContainer: {
    width: 48,
    height: 48,
    borderRadius: 24,
    backgroundColor: '#0EA5E9',
    alignItems: 'center',
    justifyContent: 'center',
    borderWidth: 3,
    borderColor: '#FFFFFF',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.3,
    shadowRadius: 4,
    elevation: 5,
  },
  vehicleArrow: {
    fontSize: 28,
    color: '#FFFFFF',
    fontWeight: 'bold',
  },
});

export default LiveTrackingMap;
