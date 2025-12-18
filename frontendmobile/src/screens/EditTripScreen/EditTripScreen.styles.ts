/**
 * EditTripScreen Styles
 */

import { StyleSheet } from 'react-native'
import { colors } from '@/styles/colors'

export const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.white,
  },

  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: colors.white,
  },

  loadingText: {
    marginTop: 16,
    color: '#4b5563',
    fontSize: 16,
  },

  errorContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: colors.white,
    paddingHorizontal: 24,
  },

  errorTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    color: '#111827',
    marginTop: 16,
  },

  errorText: {
    color: '#4b5563',
    textAlign: 'center',
    marginTop: 8,
    marginBottom: 24,
  },

  header: {
    paddingHorizontal: 24,
    paddingTop: 24,
    paddingBottom: 16,
    borderBottomWidth: 1,
    borderBottomColor: '#e5e7eb',
  },

  backButton: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 16,
  },

  backButtonText: {
    marginLeft: 8,
    color: colors.primary[500],
    fontSize: 16,
    fontWeight: '600',
  },

  headerTitle: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#111827',
  },

  headerSubtitle: {
    color: '#4b5563',
    fontSize: 14,
    marginTop: 4,
  },

  content: {
    paddingHorizontal: 24,
    paddingVertical: 16,
  },

  errorBox: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#fef2f2',
    borderWidth: 1,
    borderColor: '#fecaca',
    borderRadius: 8,
    padding: 16,
    marginBottom: 16,
  },

  errorBoxText: {
    marginLeft: 12,
    color: '#b91c1c',
    fontSize: 14,
    flex: 1,
  },

  card: {
    marginBottom: 16,
    padding: 16,
    backgroundColor: colors.white,
    borderWidth: 1,
    borderColor: '#e5e7eb',
    borderRadius: 8,
  },

  sectionHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 16,
    paddingBottom: 12,
    borderBottomWidth: 1,
    borderBottomColor: '#f3f4f6',
  },

  sectionTitle: {
    marginLeft: 8,
    fontSize: 18,
    fontWeight: '600',
    color: '#111827',
  },

  inputGroup: {
    marginBottom: 16,
  },

  dateButtonText: {
    paddingHorizontal: 16,
    paddingVertical: 12,
    backgroundColor: '#f0f9ff',
    borderWidth: 1,
    borderColor: '#bae6fd',
    borderRadius: 8,
    color: '#0369a1',
    fontWeight: '500',
  },

  preferenceItem: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingVertical: 12,
    borderBottomWidth: 1,
    borderBottomColor: '#f3f4f6',
  },

  preferenceLabel: {
    fontSize: 16,
    color: '#374151',
    fontWeight: '500',
  },

  actionButtons: {
    flexDirection: 'row',
    gap: 12,
    marginTop: 24,
    marginBottom: 32,
  },

  cancelButton: {
    flex: 1,
  },

  submitButton: {
    flex: 1,
  },
})
