/**
 * Avatar Component
 * Circular user avatar with initials fallback
 */

import React from 'react'
import { View, Text, StyleSheet } from 'react-native'
import { colors } from '@/styles/colors'

export interface AvatarProps {
  name: string
  size?: number
  backgroundColor?: string
  textColor?: string
}

export default function Avatar({ name, size = 40, backgroundColor, textColor }: AvatarProps) {
  const initials = name
    .split(' ')
    .slice(0, 2)
    .map(word => word.charAt(0).toUpperCase())
    .join('')

  const fontSize = size * 0.38

  return (
    <View
      style={[
        styles.container,
        {
          width: size,
          height: size,
          borderRadius: size / 2,
          backgroundColor: backgroundColor ?? colors.primary[500],
        },
      ]}
    >
      <Text style={[styles.text, { fontSize, color: textColor ?? colors.white }]}>
        {initials}
      </Text>
    </View>
  )
}

const styles = StyleSheet.create({
  container: {
    alignItems: 'center',
    justifyContent: 'center',
  },
  text: {
    fontWeight: '700',
  },
})
