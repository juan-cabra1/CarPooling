/**
 * Input Component
 * Reusable text input with label and error states
 */

import React, { useState } from 'react'
import { View, TextInput, Text, type TextInputProps, type ViewStyle, type TextStyle } from 'react-native'
import { styles } from './Input.styles'

export interface InputProps extends TextInputProps {
  /** Input label */
  label?: string
  /** Error message */
  error?: string
  /** Helper text */
  helper?: string
  /** Container style */
  containerStyle?: ViewStyle
}

export default function Input({
  label,
  error,
  helper,
  containerStyle,
  style,
  editable = true,
  ...props
}: InputProps) {
  const [isFocused, setIsFocused] = useState(false)

  const inputStyle: TextStyle[] = [
    styles.base,
    isFocused && styles.focused,
    error && styles.error,
    !editable && styles.disabled,
    style as TextStyle,
  ].filter(Boolean)

  return (
    <View style={containerStyle || styles.container}>
      {label && <Text style={styles.label}>{label}</Text>}

      <TextInput
        style={inputStyle}
        placeholderTextColor="#94a3b8"
        editable={editable}
        onFocus={(e) => {
          setIsFocused(true)
          props.onFocus?.(e)
        }}
        onBlur={(e) => {
          setIsFocused(false)
          props.onBlur?.(e)
        }}
        {...props}
      />

      {error && <Text style={styles.errorText}>{error}</Text>}
      {!error && helper && <Text style={styles.helperText}>{helper}</Text>}
    </View>
  )
}
