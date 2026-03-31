/**
 * RatingStars Component
 * Displays a star rating (read-only)
 */

import React from 'react'
import { View, Text, StyleSheet } from 'react-native'
import { Ionicons } from '@expo/vector-icons'

export interface RatingStarsProps {
  rating: number
  maxStars?: number
  size?: number
  showValue?: boolean
  color?: string
}

export default function RatingStars({
  rating,
  maxStars = 5,
  size = 16,
  showValue = false,
  color = '#facc15',
}: RatingStarsProps) {
  return (
    <View style={styles.container}>
      {Array.from({ length: maxStars }, (_, i) => {
        const filled = i < Math.floor(rating)
        const half = !filled && i < rating
        return (
          <Ionicons
            key={i}
            name={filled ? 'star' : half ? 'star-half' : 'star-outline'}
            size={size}
            color={filled || half ? color : '#d1d5db'}
            style={styles.star}
          />
        )
      })}
      {showValue && (
        <Text style={[styles.value, { fontSize: size * 0.85 }]}>{rating.toFixed(1)}</Text>
      )}
    </View>
  )
}

const styles = StyleSheet.create({
  container: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  star: {
    marginRight: 2,
  },
  value: {
    marginLeft: 4,
    color: '#6b7280',
    fontWeight: '600',
  },
})
