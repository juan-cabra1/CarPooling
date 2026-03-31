
import React from 'react';
import { View, Text } from 'react-native';
import LiveTrackingMap from '../../components/map/LiveTrackingMap';
import { useTrackingContext } from '../../context/TrackingContext';
import Button from '../../components/ui/Button/Button';
import { styles } from './TripTrackingScreen.styles';

const TripTrackingScreen: React.FC = () => {
  const { activeTrip, eta, distanceRemaining } = useTrackingContext();

  return (
    <View style={styles.container}>
      <LiveTrackingMap />
      <View style={styles.bottomContainer}>
        <Text style={styles.etaText}>
          {eta ? `ETA: ${eta} min` : 'Calculating ETA...'}
        </Text>
        <Text style={styles.distanceText}>
          {distanceRemaining ? `${distanceRemaining} km remaining` : ''}
        </Text>
        {activeTrip?.driver && (
          <View style={styles.driverContainer}>
            <Text style={styles.driverName}>{activeTrip.driver.name}</Text>
            <Text>{activeTrip.car.brand} {activeTrip.car.model}</Text>
            <Text>{activeTrip.car.plate}</Text>
          </View>
        )}
        <View style={styles.buttonsContainer}>
          <Button style={styles.button}>
            <Text>Contact</Text>
          </Button>
          <Button style={styles.button} variant="danger">
            <Text>Cancel Booking</Text>
          </Button>
        </View>
      </View>
    </View>
  );
};

export default TripTrackingScreen;
