/**
 * RegisterScreen Styles
 */

import { StyleSheet } from 'react-native'
import { colors } from '@/styles/colors'

export const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.white,
  },
  scrollContent: {
    paddingBottom: 32,
  },

  // Card
  card: {
    marginHorizontal: 16,
    marginTop: 16,
    marginBottom: 32,
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

  // Success screen
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
  successInfo: {
    padding: 16,
    borderRadius: 8,
    backgroundColor: '#dbeafe',
    borderWidth: 1,
    borderColor: '#bfdbfe',
    marginBottom: 16,
  },
  successInfoText: {
    fontSize: 14,
    color: '#1e3a8a',
    textAlign: 'center',
  },
  redirectText: {
    textAlign: 'center',
    fontSize: 14,
    color: colors.mutedForeground,
  },

  // Error
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

  // Form sections
  section: {
    padding: 16,
    borderRadius: 8,
    backgroundColor: colors.muted,
    marginBottom: 16,
    opacity: 0.5,
  },
  sectionTitle: {
    fontWeight: '600',
    fontSize: 14,
    color: colors.mutedForeground,
    marginBottom: 16,
  },

  // Inputs
  inputGroup: {
    marginBottom: 16,
  },
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
  button: {
    width: '100%',
    marginBottom: 16,
  },
  loginContainer: {
    textAlign: 'center',
  },
  loginText: {
    fontSize: 14,
    color: colors.mutedForeground,
  },
  loginLink: {
    color: colors.primary[500],
    fontWeight: '600',
  },
})
