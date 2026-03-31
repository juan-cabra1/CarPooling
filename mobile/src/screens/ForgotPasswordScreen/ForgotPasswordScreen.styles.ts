/**
 * ForgotPasswordScreen Styles
 */

import { StyleSheet } from 'react-native'
import { colors } from '@/styles/colors'

export const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.white,
    alignItems: 'center',
    justifyContent: 'center',
    padding: 24,
  },
  card: {
    width: '100%',
  },

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
  successIconContainer: {
    width: 64,
    height: 64,
    backgroundColor: '#22c55e',
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

  errorContainer: {
    padding: 12,
    borderRadius: 8,
    backgroundColor: '#fee2e2',
    borderWidth: 1,
    borderColor: '#fecaca',
    flexDirection: 'row',
    alignItems: 'flex-start',
    gap: 8,
    marginBottom: 16,
  },
  errorText: {
    fontSize: 14,
    color: colors.destructive,
    flex: 1,
  },

  infoBox: {
    padding: 16,
    borderRadius: 8,
    backgroundColor: '#dbeafe',
    borderWidth: 1,
    borderColor: '#bfdbfe',
    marginBottom: 16,
  },
  infoText: {
    fontSize: 14,
    color: '#1e3a8a',
  },

  warningBox: {
    padding: 16,
    borderRadius: 8,
    backgroundColor: '#fef9c3',
    borderWidth: 1,
    borderColor: '#fef08a',
    marginBottom: 16,
  },
  warningText: {
    fontSize: 14,
    color: '#854d0e',
  },

  button: {
    width: '100%',
    marginBottom: 12,
  },
  registerContainer: {
    width: '100%',
    alignItems: 'center',
    marginTop: 16,
  },
  registerText: {
    textAlign: 'center',
    fontSize: 14,
    color: colors.mutedForeground,
    marginTop: 16,
  },
  registerLink: {
    color: colors.primary[500],
    fontWeight: '600',
  },
})
