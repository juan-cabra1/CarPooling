/**
 * Button Component
 * Reusable button with multiple variants matching web frontend design
 */

import React from 'react'
import { TouchableOpacity, Text, ActivityIndicator, type TouchableOpacityProps, type ViewStyle, type TextStyle } from 'react-native'
import { styles, textStyles } from './Button.styles'

export interface ButtonProps extends TouchableOpacityProps {
  /** Button text */
  title?: string
  /** Button variant */
  variant?: 'primary' | 'secondary' | 'outline' | 'destructive' | 'ghost' | 'link'
  /** Button size */
  size?: 'default' | 'sm' | 'lg' | 'icon'
  /** Loading state */
  loading?: boolean
  /** Disabled state */
  disabled?: boolean
  /** Children (overrides title) */
  children?: React.ReactNode
}

export default function Button({
  title,
  variant = 'primary',
  size = 'default',
  loading = false,
  disabled = false,
  children,
  style,
  ...props
}: ButtonProps) {
  const isDisabled = disabled || loading

  // Combine button styles
  const buttonStyle: ViewStyle[] = [
    styles.base,
    styles[variant],
    size !== 'default' && styles[size],
    isDisabled && styles.disabled,
    style as ViewStyle,
  ].filter(Boolean)

  // Combine text styles
  const textStyle: TextStyle[] = [
    textStyles[variant],
    size !== 'default' && size !== 'icon' && textStyles[size],
    isDisabled && textStyles.disabled,
  ].filter(Boolean)

  return (
    <TouchableOpacity
      style={buttonStyle}
      disabled={isDisabled}
      activeOpacity={0.7}
      {...props}
    >
      {loading ? (
        <ActivityIndicator
          color={variant === 'outline' || variant === 'ghost' || variant === 'link' ? '#0ea5e9' : '#ffffff'}
          size="small"
        />
      ) : (
        <>
          {children || <Text style={textStyle}>{title}</Text>}
        </>
      )}
    </TouchableOpacity>
  )
}
