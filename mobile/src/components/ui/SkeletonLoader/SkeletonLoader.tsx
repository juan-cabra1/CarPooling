/**
 * SkeletonLoader Component
 * Animated placeholder for loading states
 */

import React, { useEffect, useRef } from 'react'
import { View, Animated, StyleSheet, type ViewStyle } from 'react-native'
import { colors } from '@/styles/colors'

export interface SkeletonLoaderProps {
  width?: number | string
  height?: number
  borderRadius?: number
  style?: ViewStyle
}

export default function SkeletonLoader({ width = '100%', height = 16, borderRadius = 8, style }: SkeletonLoaderProps) {
  const opacity = useRef(new Animated.Value(0.3)).current

  useEffect(() => {
    const animation = Animated.loop(
      Animated.sequence([
        Animated.timing(opacity, { toValue: 1, duration: 700, useNativeDriver: true }),
        Animated.timing(opacity, { toValue: 0.3, duration: 700, useNativeDriver: true }),
      ])
    )
    animation.start()
    return () => animation.stop()
  }, [opacity])

  return (
    <Animated.View
      style={[
        styles.skeleton,
        { width: width as any, height, borderRadius, opacity },
        style,
      ]}
    />
  )
}

/**
 * TripCardSkeleton — pre-built skeleton matching TripCard layout
 */
export function TripCardSkeleton() {
  return (
    <View style={styles.card}>
      <View style={styles.cardHeader}>
        <View style={{ flex: 1, gap: 8 }}>
          <SkeletonLoader height={18} width="70%" />
          <SkeletonLoader height={14} width="50%" />
        </View>
        <View style={{ alignItems: 'flex-end', gap: 6 }}>
          <SkeletonLoader height={22} width={80} />
          <SkeletonLoader height={12} width={60} />
        </View>
      </View>
      <View style={styles.cardBody}>
        <SkeletonLoader height={14} width="90%" />
        <SkeletonLoader height={14} width="75%" />
      </View>
      <View style={styles.cardFooter}>
        <SkeletonLoader width={36} height={36} borderRadius={18} />
        <View style={{ flex: 1, marginLeft: 12, gap: 6 }}>
          <SkeletonLoader height={14} width="60%" />
          <SkeletonLoader height={12} width="40%" />
        </View>
      </View>
    </View>
  )
}

const styles = StyleSheet.create({
  skeleton: {
    backgroundColor: colors.border,
  },
  card: {
    backgroundColor: colors.card,
    borderRadius: 12,
    borderWidth: 1,
    borderColor: colors.border,
    padding: 16,
    marginBottom: 12,
  },
  cardHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: 12,
  },
  cardBody: {
    gap: 8,
    marginBottom: 12,
  },
  cardFooter: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingTop: 12,
    borderTopWidth: 1,
    borderTopColor: colors.border,
  },
})
