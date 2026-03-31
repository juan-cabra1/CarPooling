/**
 * TextArea Component
 * Multiline text input
 */

import React from 'react'
import { TextInput, View, Text, type TextInputProps } from 'react-native'
import { styles } from './TextArea.styles'

export interface TextAreaProps extends Omit<TextInputProps, 'multiline'> {
  label?: string
  error?: string
  helperText?: string
}

export default function TextArea({ label, error, helperText, ...props }: TextAreaProps) {
  return (
    <View style={styles.container}>
      {label && <Text style={styles.label}>{label}</Text>}
      <TextInput
        {...props}
        multiline
        textAlignVertical="top"
        style={error ? styles.inputError : styles.input}
        placeholderTextColor="#94a3b8"
      />
      {helperText && !error && <Text style={styles.helperText}>{helperText}</Text>}
      {error && <Text style={styles.errorText}>{error}</Text>}
    </View>
  )
}
