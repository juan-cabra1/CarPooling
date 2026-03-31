import React, { useRef, useEffect } from 'react';
import { View, Text, StyleSheet, ScrollView, TouchableOpacity } from 'react-native';
import { NavigationStep, getManeuverIcon } from '../../services/locationService';

interface StepListProps {
  steps: NavigationStep[];
  currentStepIndex: number;
  onStepPress?: (step: NavigationStep, index: number) => void;
}

export const StepList: React.FC<StepListProps> = ({
  steps,
  currentStepIndex,
  onStepPress,
}) => {
  const scrollViewRef = useRef<ScrollView>(null);

  // Auto-scroll to current step
  useEffect(() => {
    if (scrollViewRef.current && currentStepIndex >= 0) {
      // Scroll to current step (approximate position)
      const yOffset = currentStepIndex * 80; // Approximate height per item
      scrollViewRef.current.scrollTo({ y: yOffset, animated: true });
    }
  }, [currentStepIndex]);

  if (steps.length === 0) {
    return (
      <View style={styles.emptyContainer}>
        <Text style={styles.emptyText}>No hay instrucciones de navegación</Text>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.headerText}>
          Instrucciones ({currentStepIndex + 1} de {steps.length})
        </Text>
      </View>

      <ScrollView
        ref={scrollViewRef}
        style={styles.scrollView}
        showsVerticalScrollIndicator={false}
      >
        {steps.map((step, index) => {
          const isCurrent = index === currentStepIndex;
          const isCompleted = index < currentStepIndex;
          const maneuverIcon = getManeuverIcon(step.maneuver);

          return (
            <TouchableOpacity
              key={step.stepIndex}
              style={[
                styles.stepItem,
                isCurrent && styles.stepItemCurrent,
                isCompleted && styles.stepItemCompleted,
              ]}
              onPress={() => onStepPress?.(step, index)}
              disabled={!onStepPress}
            >
              {/* Step Number & Icon */}
              <View style={styles.stepIconContainer}>
                {isCompleted ? (
                  <Text style={styles.stepCheckmark}>✓</Text>
                ) : (
                  <Text style={styles.stepIcon}>{maneuverIcon.unicode}</Text>
                )}
              </View>

              {/* Step Details */}
              <View style={styles.stepContent}>
                <Text
                  style={[
                    styles.stepInstruction,
                    isCurrent && styles.stepInstructionCurrent,
                    isCompleted && styles.stepInstructionCompleted,
                  ]}
                  numberOfLines={2}
                >
                  {step.instruction}
                </Text>

                {step.streetName && (
                  <Text style={styles.stepStreetName} numberOfLines={1}>
                    {step.streetName}
                  </Text>
                )}

                <Text style={styles.stepDistance}>
                  {step.distanceText} • {step.durationText}
                </Text>
              </View>
            </TouchableOpacity>
          );
        })}
      </ScrollView>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    maxHeight: 300,
    marginTop: 16,
  },
  header: {
    paddingVertical: 8,
    borderBottomWidth: 1,
    borderBottomColor: '#E5E7EB',
    marginBottom: 8,
  },
  headerText: {
    fontSize: 14,
    fontWeight: '600',
    color: '#6B7280',
  },
  scrollView: {
    flex: 1,
  },
  emptyContainer: {
    padding: 20,
    alignItems: 'center',
  },
  emptyText: {
    fontSize: 14,
    color: '#9CA3AF',
  },
  stepItem: {
    flexDirection: 'row',
    alignItems: 'flex-start',
    paddingVertical: 12,
    paddingHorizontal: 8,
    borderRadius: 8,
    marginBottom: 4,
  },
  stepItemCurrent: {
    backgroundColor: '#EBF5FF',
    borderLeftWidth: 4,
    borderLeftColor: '#0EA5E9',
  },
  stepItemCompleted: {
    opacity: 0.5,
  },
  stepIconContainer: {
    width: 40,
    height: 40,
    borderRadius: 20,
    backgroundColor: '#F3F4F6',
    alignItems: 'center',
    justifyContent: 'center',
    marginRight: 12,
  },
  stepIcon: {
    fontSize: 20,
    color: '#374151',
  },
  stepCheckmark: {
    fontSize: 24,
    color: '#10B981',
  },
  stepContent: {
    flex: 1,
  },
  stepInstruction: {
    fontSize: 15,
    fontWeight: '500',
    color: '#374151',
    marginBottom: 4,
  },
  stepInstructionCurrent: {
    fontWeight: '700',
    color: '#0EA5E9',
  },
  stepInstructionCompleted: {
    color: '#9CA3AF',
  },
  stepStreetName: {
    fontSize: 13,
    color: '#6B7280',
    fontStyle: 'italic',
    marginBottom: 4,
  },
  stepDistance: {
    fontSize: 12,
    color: '#9CA3AF',
  },
});

export default StepList;
