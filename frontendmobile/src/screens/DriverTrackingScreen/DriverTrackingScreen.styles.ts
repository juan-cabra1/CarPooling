
import { StyleSheet } from 'react-native';

export const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  bottomContainer: {
    position: 'absolute',
    bottom: 0,
    left: 0,
    right: 0,
    backgroundColor: 'white',
    padding: 20,
    borderTopLeftRadius: 20,
    borderTopRightRadius: 20,
  },
  etaText: {
    fontSize: 20,
    fontWeight: 'bold',
  },
  distanceText: {
    fontSize: 16,
    color: 'gray',
  },
  passengersContainer: {
    marginVertical: 20,
  },
  passengersTitle: {
    fontSize: 18,
    fontWeight: 'bold',
  },
});
