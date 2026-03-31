/**
 * Label Component
 * Simple text label for forms
 */

import React from 'react'
import { Text, type TextProps } from 'react-native'

export interface LabelProps extends TextProps {
  children: React.ReactNode
}

export default function Label({ children, className, ...props }: LabelProps) {
  return (
    <Text className={`text-sm font-medium text-foreground mb-2 ${className || ''}`} {...props}>
      {children}
    </Text>
  )
}
