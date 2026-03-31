/**
 * LoginScreen Styles
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

  // Card
  card: {
    width: '100%',
    maxWidth: 448,
  },

  // Header
  iconContainer: {
    width: 64,
    height: 64,
    backgroundColor: colors.primary[500],
    borderRadius: 32,
    alignItems: 'center',
    justifyContent: 'center',
    marginBottom: 16,
    alignSelf: 'center',
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    textAlign: 'center',
    color: colors.foreground,
  },
  description: {
    textAlign: 'center',
    color: colors.mutedForeground,
    marginTop: 8,
  },

  // Messages
  successContainer: {
    padding: 12,
    borderRadius: 8,
    backgroundColor: '#dcfce7',
    borderWidth: 1,
    borderColor: '#bbf7d0',
    flexDirection: 'row',
    alignItems: 'flex-start',
    gap: 8,
    marginBottom: 16,
  },
  successText: {
    fontSize: 14,
    color: '#16a34a',
    flex: 1,
  },

  errorContainer: {
    padding: 12,
    borderRadius: 8,
    backgroundColor: '#fee2e2',
    borderWidth: 1,
    borderColor: '#fecaca',
    marginBottom: 16,
  },
  errorInner: {
    flexDirection: 'row',
    alignItems: 'flex-start',
    gap: 8,
  },
  errorText: {
    fontSize: 14,
    color: colors.destructive,
    flex: 1,
  },
  resendLink: {
    marginLeft: 28,
    marginTop: 8,
  },
  resendLinkText: {
    fontSize: 14,
    color: colors.primary[500],
    fontWeight: '500',
  },

  // Form
  formContainer: {
    gap: 16,
  },
  inputGroup: {
    gap: 8,
  },

  // Label with icon
  labelContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
    marginBottom: 8,
  },
  label: {
    fontSize: 14,
    fontWeight: '500',
    color: colors.foreground,
  },

  // Footer
  forgotPassword: {
    fontSize: 14,
    color: colors.primary[500],
    textAlign: 'right',
    marginBottom: 16,
  },

  // Button
  button: {
    width: '100%',
    backgroundColor: colors.primary[500],
    paddingVertical: 16,
    borderRadius: 8,
    alignItems: 'center',
    justifyContent: 'center',
    flexDirection: 'row',
  },
  buttonText: {
    color: colors.white,
    fontSize: 16,
    fontWeight: 'bold',
    marginLeft: 8,
  },

  // Register link
  registerContainer: {
    textAlign: 'center',
    marginTop: 16,
  },
  registerText: {
    fontSize: 14,
    color: colors.mutedForeground,
  },
  registerLink: {
    color: colors.primary[500],
    fontWeight: '600',
  },
})
