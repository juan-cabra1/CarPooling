/**
 * Badge Component Styles
 */

import { StyleSheet } from 'react-native'
import { colors } from '@/styles/colors'

export const styles = StyleSheet.create({
  // Badge variants
  default: {
    backgroundColor: colors.primary[100],
    paddingHorizontal: 12,
    paddingVertical: 4,
    borderRadius: 9999,
  },
  success: {
    backgroundColor: colors.green[100],
    paddingHorizontal: 12,
    paddingVertical: 4,
    borderRadius: 9999,
  },
  warning: {
    backgroundColor: colors.yellow[100],
    paddingHorizontal: 12,
    paddingVertical: 4,
    borderRadius: 9999,
  },
  destructive: {
    backgroundColor: colors.red[100],
    paddingHorizontal: 12,
    paddingVertical: 4,
    borderRadius: 9999,
  },
  secondary: {
    backgroundColor: colors.secondary[100],
    paddingHorizontal: 12,
    paddingVertical: 4,
    borderRadius: 9999,
  },
  outline: {
    borderWidth: 1,
    borderColor: colors.border,
    paddingHorizontal: 12,
    paddingVertical: 4,
    borderRadius: 9999,
  },
})

// Text styles as separate StyleSheet
export const textStyles = StyleSheet.create({
  default: {
    color: colors.primary[700],
    fontSize: 12,
    fontWeight: '600',
  },
  success: {
    color: colors.green[700],
    fontSize: 12,
    fontWeight: '600',
  },
  warning: {
    color: colors.yellow[700],
    fontSize: 12,
    fontWeight: '600',
  },
  destructive: {
    color: colors.red[700],
    fontSize: 12,
    fontWeight: '600',
  },
  secondary: {
    color: colors.secondary[700],
    fontSize: 12,
    fontWeight: '600',
  },
  outline: {
    color: colors.foreground,
    fontSize: 12,
    fontWeight: '600',
  },
})
