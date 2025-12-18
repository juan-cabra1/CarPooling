
import React from 'react';
import { View, StyleSheet } from 'react-native';
import { useRoute } from '@react-navigation/native';
import TripRouteMap from '../../components/map/TripRouteMap';
import { Coordinates } from '../../types';

const TripRouteMapScreen: React.FC = () => {
  const route = useRoute();
  const { origin, destination } = route.params as { origin: Coordinates, destination: Coordinates };

  return (
    <View style={styles.container}>
      <TripRouteMap origin={origin} destination={destination} />
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
});

export default TripRouteMapScreen;
