/**
 * VerifyEmailScreen Styles
 */

import { StyleSheet } from 'react-native'
import { colors } from '@/styles/colors'

export const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.white,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 24,
  },
  content: {
    width: '100%',
    maxWidth: 448,
    alignItems: 'center',
  },

  // Loading
  loadingText: {
    fontSize: 16,
    color: colors.mutedForeground,
    marginTop: 16,
    textAlign: 'center',
  },

  // Icon
  iconContainer: {
    alignItems: 'center',
    marginBottom: 24,
  },

  // Header
  title: {
    fontSize: 28,
    fontWeight: 'bold',
    textAlign: 'center',
    color: colors.foreground,
    marginBottom: 12,
  },
  description: {
    fontSize: 16,
    textAlign: 'center',
    color: colors.mutedForeground,
    lineHeight: 24,
    marginBottom: 24,
  },

  // Reasons
  reasonsContainer: {
    backgroundColor: colors.background,
    borderRadius: 8,
    padding: 16,
    marginBottom: 24,
    width: '100%',
  },
  reasonsTitle: {
    fontSize: 14,
    fontWeight: '600',
    color: colors.foreground,
    marginBottom: 8,
  },
  reasonItem: {
    fontSize: 14,
    color: colors.mutedForeground,
    lineHeight: 20,
    marginTop: 4,
  },

  // Buttons
  button: {
    width: '100%',
    marginBottom: 12,
  },
  outlineButton: {
    width: '100%',
  },
})
