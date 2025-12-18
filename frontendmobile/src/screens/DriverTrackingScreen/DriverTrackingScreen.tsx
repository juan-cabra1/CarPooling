
import React, { useEffect, useState } from 'react';
import { View, Text, ActivityIndicator, Alert } from 'react-native';
import { useRoute, useNavigation, type RouteProp } from '@react-navigation/native';
import LiveTrackingMap from '../../components/map/LiveTrackingMap';
import NavigationBanner from '../../components/map/NavigationBanner';
import StepList from '../../components/map/StepList';
import { useTrackingContext } from '../../context/TrackingContext';
import Button from '../../components/ui/Button/Button';
import { styles } from './DriverTrackingScreen.styles';
import { tripsService, bookingsService, getErrorMessage } from '../../services';
import type { Trip, Booking } from '../../types';
import type { RootStackParamList } from '../../navigation/types';

type DriverTrackingRouteProp = RouteProp<RootStackParamList, 'DriverTracking'>;

const DriverTrackingScreen: React.FC = () => {
  const route = useRoute<DriverTrackingRouteProp>();
  const navigation = useNavigation();
  const { tripId } = route.params;
  const {
    activeTrip,
    eta,
    distanceRemaining,
    startTracking,
    stopTracking,
    // Estado de navegación
    navigationSteps,
    currentStepIndex,
    distanceToNextStep,
    nextStepPreview,
    isApproachingTurn,
  } = useTrackingContext();

  const [loading, setLoading] = useState(true);
  const [trip, setTrip] = useState<Trip | null>(null);
  const [bookings, setBookings] = useState<Booking[]>([]);

  useEffect(() => {
    console.log('DriverTrackingScreen mounted with tripId:', tripId);
    console.log('Route params:', route.params);

    if (!tripId) {
      Alert.alert('Error', 'No se proporcionó el ID del viaje');
      navigation.goBack();
      return;
    }

    loadTripAndBookings();
  }, [tripId]);

  const loadTripAndBookings = async () => {
    try {
      setLoading(true);

      const tripData = await tripsService.getTripById(tripId);
      setTrip(tripData);

      const tripBookings = await bookingsService.getBookingsByTrip(tripId);
      setBookings(tripBookings);

      console.log(`Loaded ${tripBookings.length} confirmed bookings for trip ${tripId}`);

      // Iniciar tracking, now returns a result object
      const trackingResult = await startTracking(tripData, true);

      if (!trackingResult.success) {
        // Show an alert but stay on the screen
        Alert.alert('Error de Rastreo', trackingResult.message);
      }
    } catch (err) {
      console.error('Error loading trip data:', err);
      // Only show alert, do not navigate away
      Alert.alert('Error Crítico', getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  const handleFinalize = async () => {
    Alert.alert(
      'Finalizar Viaje',
      '¿Estás seguro de que quieres finalizar este viaje?',
      [
        { text: 'Cancelar', style: 'cancel' },
        {
          text: 'Finalizar',
          onPress: async () => {
            try {
              await tripsService.completeTrip(tripId);
              await stopTracking();
              navigation.goBack();
            } catch (err) {
              Alert.alert('Error', getErrorMessage(err));
            }
          },
        },
      ]
    );
  };

  if (loading) {
    return (
      <View style={[styles.container, { justifyContent: 'center', alignItems: 'center' }]}>
        <ActivityIndicator size="large" color="#0ea5e9" />
        <Text style={{ marginTop: 16, fontSize: 16 }}>Cargando viaje...</Text>
      </View>
    );
  }

  const currentStep = navigationSteps[currentStepIndex] || null;

  return (
    <View style={styles.container}>
      {/* Navigation Banner (Top) */}
      <NavigationBanner
        currentStep={currentStep}
        distanceToNextStep={distanceToNextStep}
        nextStepPreview={nextStepPreview}
        isApproachingTurn={isApproachingTurn}
      />

      <LiveTrackingMap />
      <View style={styles.bottomContainer}>
        <Text style={styles.etaText}>
          {eta ? `ETA: ${eta} min` : 'Calculating ETA...'}
        </Text>
        <Text style={styles.distanceText}>
          {distanceRemaining ? `${distanceRemaining.toFixed(1)} km remaining` : ''}
        </Text>

        {/* NUEVO: Lista de pasos de navegación */}
        {navigationSteps.length > 0 && (
          <StepList
            steps={navigationSteps}
            currentStepIndex={currentStepIndex}
          />
        )}

        <View style={styles.passengersContainer}>
          <Text style={styles.passengersTitle}>
            Pasajeros confirmados: {bookings.length}
          </Text>
          {bookings.length === 0 ? (
            <Text style={{ fontSize: 14, color: '#6b7280', marginTop: 4 }}>
              No hay pasajeros confirmados todavía
            </Text>
          ) : (
            <>
              <Text style={{ fontSize: 12, color: '#6b7280', marginTop: 4, marginBottom: 8 }}>
                Total de asientos reservados: {bookings.reduce((sum, b) => sum + b.seats_requested, 0)}
              </Text>
              {bookings.map((booking) => (
                <Text key={booking.id} style={{ fontSize: 14, marginTop: 4, color: '#374151' }}>
                  • Pasajero ID: {booking.passenger_id} - {booking.seats_requested} asiento{booking.seats_requested > 1 ? 's' : ''}
                </Text>
              ))}
            </>
          )}
        </View>
        <Button onPress={handleFinalize}>
          <Text>Finalizar Viaje</Text>
        </Button>
      </View>
    </View>
  );
};

export default DriverTrackingScreen;
