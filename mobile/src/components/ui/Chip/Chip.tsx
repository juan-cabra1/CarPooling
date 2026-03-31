/**
 * Chip Component
 * Selectable/dismissible filter tag
 */

import React from 'react'
import { TouchableOpacity, Text, View, StyleSheet } from 'react-native'
import { Ionicons } from '@expo/vector-icons'
import { colors } from '@/styles/colors'

export interface ChipProps {
  label: string
  selected?: boolean
  onPress?: () => void
  onDismiss?: () => void
  disabled?: boolean
}

export default function Chip({ label, selected = false, onPress, onDismiss, disabled = false }: ChipProps) {
  return (
    <TouchableOpacity
      onPress={onPress}
      disabled={disabled || (!onPress && !onDismiss)}
      activeOpacity={0.7}
      style={[styles.chip, selected ? styles.selected : styles.unselected, disabled && styles.disabled]}
    >
      <Text style={[styles.label, selected ? styles.labelSelected : styles.labelUnselected]}>
        {label}
      </Text>
      {onDismiss && (
        <View style={styles.dismissButton}>
          <Ionicons
            name="close"
            size={12}
            color={selected ? colors.primary[600] : colors.mutedForeground}
            onPress={onDismiss}
          />
        </View>
      )}
    </TouchableOpacity>
  )
}

const styles = StyleSheet.create({
  chip: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: 12,
    paddingVertical: 6,
    borderRadius: 9999,
    borderWidth: 1,
  },
  selected: {
    backgroundColor: colors.primary[100],
    borderColor: colors.primary[400],
  },
  unselected: {
    backgroundColor: colors.white,
    borderColor: colors.border,
  },
  disabled: {
    opacity: 0.5,
  },
  label: {
    fontSize: 13,
    fontWeight: '500',
  },
  labelSelected: {
    color: colors.primary[700],
  },
  labelUnselected: {
    color: colors.mutedForeground,
  },
  dismissButton: {
    marginLeft: 4,
  },
})
