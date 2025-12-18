/**
 * TextArea Component Styles
 */

import { StyleSheet } from 'react-native'
import { colors } from '@/styles/colors'

export const styles = StyleSheet.create({
  container: {
    marginBottom: 16,
  },
  label: {
    fontSize: 14,
    fontWeight: '500',
    color: colors.foreground,
    marginBottom: 8,
  },
  input: {
    borderWidth: 1,
    borderColor: colors.border,
    borderRadius: 8,
    padding: 12,
    fontSize: 16,
    backgroundColor: colors.white,
    minHeight: 96,
    color: colors.foreground,
  },
  inputError: {
    borderWidth: 1,
    borderColor: colors.destructive,
    borderRadius: 8,
    padding: 12,
    fontSize: 16,
    backgroundColor: colors.white,
    minHeight: 96,
    color: colors.foreground,
  },
  helperText: {
    fontSize: 12,
    color: colors.mutedForeground,
    marginTop: 4,
  },
  errorText: {
    fontSize: 14,
    color: colors.destructive,
    marginTop: 4,
  },
})
