/**
 * Button Component Styles
 * Separated styles for Button component
 */

import { StyleSheet } from 'react-native'
import { colors } from '@/styles/colors'

export const styles = StyleSheet.create({
  // Base button styles
  base: {
    borderRadius: 8,
    alignItems: 'center',
    justifyContent: 'center',
  },

  // Variants
  primary: {
    backgroundColor: colors.primary[500],
    paddingHorizontal: 24,
    paddingVertical: 12,
  },
  secondary: {
    backgroundColor: colors.secondary[500],
    paddingHorizontal: 24,
    paddingVertical: 12,
  },
  outline: {
    borderWidth: 2,
    borderColor: colors.primary[500],
    backgroundColor: colors.transparent,
    paddingHorizontal: 24,
    paddingVertical: 12,
  },
  destructive: {
    backgroundColor: colors.destructive,
    paddingHorizontal: 24,
    paddingVertical: 12,
  },
  ghost: {
    backgroundColor: colors.transparent,
    paddingHorizontal: 24,
    paddingVertical: 12,
  },
  link: {
    backgroundColor: colors.transparent,
    paddingHorizontal: 0,
    paddingVertical: 0,
  },

  // Sizes
  default: {
    paddingHorizontal: 24,
    paddingVertical: 12,
  },
  sm: {
    paddingHorizontal: 16,
    paddingVertical: 8,
  },
  lg: {
    paddingHorizontal: 32,
    paddingVertical: 16,
  },
  icon: {
    width: 40,
    height: 40,
    padding: 0,
  },

  // Disabled
  disabled: {
    backgroundColor: colors.muted,
    opacity: 0.5,
  },
})

// Text styles as separate StyleSheet
export const textStyles = StyleSheet.create({
  // Variant text styles
  primary: {
    color: colors.white,
    textAlign: 'center',
    fontWeight: 'bold',
    fontSize: 16,
  },
  secondary: {
    color: colors.white,
    textAlign: 'center',
    fontWeight: 'bold',
    fontSize: 16,
  },
  outline: {
    color: colors.primary[500],
    textAlign: 'center',
    fontWeight: 'bold',
    fontSize: 16,
  },
  destructive: {
    color: colors.white,
    textAlign: 'center',
    fontWeight: 'bold',
    fontSize: 16,
  },
  ghost: {
    color: colors.foreground,
    textAlign: 'center',
    fontWeight: '600',
    fontSize: 16,
  },
  link: {
    color: colors.primary[500],
    fontWeight: '600',
    fontSize: 16,
    textDecorationLine: 'underline',
  },
  disabled: {
    color: colors.mutedForeground,
  },

  // Size text styles
  sm: {
    fontSize: 14,
  },
  default: {
    fontSize: 16,
  },
  lg: {
    fontSize: 18,
  },
})
