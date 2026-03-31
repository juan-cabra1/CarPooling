/**
 * Badge Component
 * Small status indicator
 */

import React from 'react'
import { View, Text, type ViewProps } from 'react-native'
import { styles, textStyles } from './Badge.styles'

export interface BadgeProps extends ViewProps {
  /** Badge variant */
  variant?: 'default' | 'success' | 'warning' | 'destructive' | 'secondary' | 'outline'
  /** Badge text */
  children: React.ReactNode
}

export default function Badge({ variant = 'default', children, style, ...props }: BadgeProps) {
  return (
    <View style={[styles[variant], style]} {...props}>
      <Text style={textStyles[variant]}>{children}</Text>
    </View>
  )
}
