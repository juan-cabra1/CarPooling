/**
 * Switch Component
 * Toggle switch for boolean values
 */

import React from 'react'
import { Switch as RNSwitch, View, Text } from 'react-native'
import { styles } from './Switch.styles'
import { colors } from '@/styles/colors'

interface SwitchProps {
  value: boolean
  onValueChange: (value: boolean) => void
  label?: string
  disabled?: boolean
}

export function Switch({ value, onValueChange, label, disabled = false }: SwitchProps) {
  return (
    <View style={styles.container}>
      {label && <Text style={styles.label}>{label}</Text>}
      <RNSwitch
        value={value}
        onValueChange={onValueChange}
        disabled={disabled}
        trackColor={{ false: colors.border, true: colors.primary[300] }}
        thumbColor={value ? colors.primary[500] : colors.muted}
        ios_backgroundColor={colors.border}
      />
    </View>
  )
}
