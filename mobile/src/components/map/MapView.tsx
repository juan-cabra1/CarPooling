
import React from 'react';
import MapView, { Region, Marker, Polyline, MapMarkerProps, MapPolylineProps } from 'react-native-maps';
import { StyleSheet } from 'react-native';

export interface MapViewProps {
  region?: Region;
  markers?: MapMarkerProps[];
  polylines?: MapPolylineProps[];
  onRegionChange?: (region: Region) => void;
  showsUserLocation?: boolean;
  followsUserLocation?: boolean;
}

const CustomMapView: React.FC<MapViewProps> = ({
  region,
  markers,
  polylines,
  onRegionChange,
  showsUserLocation,
  followsUserLocation,
}) => {
  console.log('🗺️ MapView rendering - polylines:', polylines?.length || 0, 'markers:', markers?.length || 0);
  if (polylines && polylines.length > 0) {
    console.log('🗺️ First polyline props:', {
      coordinatesCount: polylines[0].coordinates?.length,
      strokeColor: polylines[0].strokeColor,
      strokeWidth: polylines[0].strokeWidth,
    });
  }

  return (
    <MapView
      style={styles.map}
      region={region}
      onRegionChangeComplete={onRegionChange}
      showsUserLocation={showsUserLocation}
      followsUserLocation={followsUserLocation}
    >
      {markers && markers.map((markerProps, index) => (
        <Marker key={index} {...markerProps} />
      ))}
      {polylines && polylines.map((polylineProps, index) => (
        <Polyline key={index} {...polylineProps} />
      ))}
    </MapView>
  );
};

const styles = StyleSheet.create({
  map: {
    ...StyleSheet.absoluteFillObject,
  },
});

export default CustomMapView;
