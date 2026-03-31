import React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { useSafeAreaInsets } from 'react-native-safe-area-context';
import { NavigationStep, getManeuverIcon, formatDistance } from '../../services/locationService';

interface NavigationBannerProps {
  currentStep: NavigationStep | null;
  distanceToNextStep: number | null; // meters
  nextStepPreview?: NavigationStep | null;
  isApproachingTurn: boolean;
}

export const NavigationBanner: React.FC<NavigationBannerProps> = ({
  currentStep,
  distanceToNextStep,
  nextStepPreview,
  isApproachingTurn,
}) => {
  const insets = useSafeAreaInsets();

  if (!currentStep) {
    return null;
  }

  const maneuverIcon = getManeuverIcon(currentStep.maneuver);

  return (
    <View style={[styles.container, { paddingTop: insets.top + 16 }]}>
      {/* Card Principal */}
      <View style={[styles.card, isApproachingTurn && styles.cardApproaching]}>
        {/* Ícono de maniobra (grande, izquierda) */}
        <View style={styles.iconContainer}>
          <Text style={styles.icon}>{maneuverIcon.unicode}</Text>
        </View>

        {/* Detalles */}
        <View style={styles.instructionContainer}>
          <View style={styles.topRow}>
            <Text style={styles.distance}>
              {formatDistance(distanceToNextStep !== null ? distanceToNextStep : currentStep.distance)}
            </Text>
            {currentStep.durationText && (
              <Text style={styles.duration}>en {currentStep.durationText}</Text>
            )}
          </View>

          <Text style={styles.instruction} numberOfLines={2}>
            {currentStep.instruction}
          </Text>

          {currentStep.streetName && (
            <Text style={styles.streetName} numberOfLines={1}>
              {currentStep.streetName}
            </Text>
          )}
        </View>
      </View>

      {/* Preview del Siguiente Paso */}
      {nextStepPreview && (
        <View style={styles.previewCard}>
          <Text style={styles.previewIcon}>
            {getManeuverIcon(nextStepPreview.maneuver).unicode}
          </Text>
          <Text style={styles.previewText} numberOfLines={1}>
            Luego: {nextStepPreview.instruction}
          </Text>
        </View>
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    paddingHorizontal: 16,
    zIndex: 10,
  },
  card: {
    backgroundColor: '#FFFFFF',
    borderRadius: 16,
    padding: 16,
    flexDirection: 'row',
    alignItems: 'center',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.15,
    shadowRadius: 8,
    elevation: 8,
  },
  cardApproaching: {
    backgroundColor: '#FFF9E6', // Subtle yellow highlight
    borderWidth: 2,
    borderColor: '#FFB800',
  },
  iconContainer: {
    width: 64,
    height: 64,
    borderRadius: 32,
    backgroundColor: '#0EA5E9',
    alignItems: 'center',
    justifyContent: 'center',
    marginRight: 16,
  },
  icon: {
    fontSize: 36,
    color: '#FFFFFF',
  },
  instructionContainer: {
    flex: 1,
  },
  topRow: {
    flexDirection: 'row',
    alignItems: 'baseline',
    marginBottom: 4,
  },
  distance: {
    fontSize: 28,
    fontWeight: 'bold',
    color: '#1F2937',
    marginRight: 8,
  },
  duration: {
    fontSize: 14,
    color: '#6B7280',
  },
  instruction: {
    fontSize: 16,
    fontWeight: '600',
    color: '#374151',
    marginBottom: 4,
  },
  streetName: {
    fontSize: 14,
    color: '#6B7280',
    fontStyle: 'italic',
  },
  previewCard: {
    backgroundColor: '#F3F4F6',
    borderRadius: 12,
    padding: 12,
    marginTop: 8,
    flexDirection: 'row',
    alignItems: 'center',
  },
  previewIcon: {
    fontSize: 20,
    marginRight: 8,
    color: '#6B7280',
  },
  previewText: {
    fontSize: 14,
    color: '#4B5563',
    flex: 1,
  },
});

export default NavigationBanner;
