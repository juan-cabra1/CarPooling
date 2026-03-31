
import * as TaskManager from 'expo-task-manager';
import * as Location from 'expo-location';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { socketService, LocationUpdateEvent } from './socketService';
import { User } from '../types';

const LOCATION_TASK_NAME = 'CARPOOLING_LOCATION_TRACKING';
const ACTIVE_TRIP_STORAGE_KEY = '@carpooling:active_trip';
const USER_STORAGE_KEY = 'user';
const TOKEN_STORAGE_KEY = 'token';

TaskManager.defineTask(LOCATION_TASK_NAME, async ({ data, error }) => {
  if (error) {
    console.error('Error in background location task:', error);
    return;
  }

  if (data) {
    const { locations } = data as { locations: Location.LocationObject[] };

    try {
      const activeTripData = await AsyncStorage.getItem(ACTIVE_TRIP_STORAGE_KEY);
      if (!activeTripData) {
        // No active trip, so we can stop the task
        await Location.stopLocationUpdatesAsync(LOCATION_TASK_NAME);
        return;
      }

      const { tripId, isDriver } = JSON.parse(activeTripData);

      if (!isDriver) {
        // Not a driver, should not be running this task
        return;
      }

      const userStr = await AsyncStorage.getItem(USER_STORAGE_KEY);
      const token = await AsyncStorage.getItem(TOKEN_STORAGE_KEY);

      if (!userStr || !token) {
        // Not authenticated, cannot send updates
        return;
      }

      if (!socketService.isConnected()) {
        const user: User = JSON.parse(userStr);
        // Attempt to reconnect, but don't wait for it.
        // The service has its own reconnection logic.
        socketService.connect(user.id, token);
      }

      // Give it a moment to connect if it wasn't already
      if (!socketService.isConnected()) {
        await new Promise(resolve => setTimeout(resolve, 2000));
      }

      if (socketService.isConnected()) {
        const location = locations[0];
        if (location) {
          const locationUpdate: LocationUpdateEvent = {
            tripId: tripId,
            coordinates: {
              lat: location.coords.latitude,
              lng: location.coords.longitude,
            },
            speed: location.coords.speed ? location.coords.speed * 3.6 : 0, // m/s to km/h
            heading: location.coords.heading,
            timestamp: location.timestamp,
          };
          socketService.emitLocationUpdate(tripId, locationUpdate);
          console.log('Sent location update from background task');
        }
      } else {
          console.log('Socket not connected, cannot send location update from background');
      }

    } catch (err) {
      console.error('Error processing background location update:', err);
    }
  }
});

export const startBackgroundTask = async () => {
    try {
        const isRegistered = await TaskManager.isTaskRegisteredAsync(LOCATION_TASK_NAME);
        if(!isRegistered){
            console.log("Task not registered");
        }
    } catch(err) {
        console.log(err);
    }
}
