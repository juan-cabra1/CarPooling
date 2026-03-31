/**
 * Input Component Styles
 */

import { StyleSheet } from 'react-native'
import { colors } from '@/styles/colors'

export const styles = StyleSheet.create({
  // Container
  container: {
    marginBottom: 16,
  },

  // Input states
  base: {
    backgroundColor: colors.white,
    borderWidth: 1,
    borderColor: colors.input,
    borderRadius: 8,
    paddingHorizontal: 16,
    paddingVertical: 12,
    fontSize: 16,
    color: colors.foreground,
  },
  focused: {
    borderColor: colors.primary[500],
    borderWidth: 2,
  },
  error: {
    borderColor: colors.destructive,
    borderWidth: 2,
  },
  disabled: {
    backgroundColor: colors.muted,
    color: colors.mutedForeground,
  },

  // Label
  label: {
    fontSize: 14,
    fontWeight: '500',
    color: colors.foreground,
    marginBottom: 8,
  },

  // Error text
  errorText: {
    color: colors.destructive,
    fontSize: 14,
    marginTop: 4,
  },

  // Helper text
  helperText: {
    color: colors.mutedForeground,
    fontSize: 14,
    marginTop: 4,
  },
})
