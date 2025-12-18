/**
 * ResendVerificationScreen Styles
 */

import { StyleSheet } from 'react-native'
import { colors } from '@/styles/colors'

export const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.white,
  },
  scrollContent: {
    flexGrow: 1,
    justifyContent: 'center',
    padding: 24,
  },
  content: {
    width: '100%',
    maxWidth: 448,
    alignSelf: 'center',
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
    marginBottom: 32,
  },

  // Form
  inputContainer: {
    marginBottom: 24,
  },

  // Buttons
  resendButton: {
    marginBottom: 16,
  },
  backButton: {
    marginBottom: 24,
  },

  // Help Text
  helpText: {
    fontSize: 14,
    textAlign: 'center',
    color: colors.mutedForeground,
    lineHeight: 20,
  },
})
