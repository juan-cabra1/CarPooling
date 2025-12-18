/**
 * CreateTripScreen Styles
 */

import { StyleSheet } from 'react-native'
import { colors } from '@/styles/colors'

export const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f9fafb',
  },

  header: {
    backgroundColor: colors.white,
    paddingHorizontal: 16,
    paddingTop: 16,
    paddingBottom: 24,
    borderBottomWidth: 1,
    borderBottomColor: '#e5e7eb',
  },

  backButton: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
    marginBottom: 16,
  },

  backButtonText: {
    fontSize: 16,
    fontWeight: '500',
    color: colors.primary[500],
  },

  headerTitle: {
    fontSize: 30,
    fontWeight: 'bold',
    color: '#111827',
    marginBottom: 8,
  },

  headerSubtitle: {
    fontSize: 16,
    color: '#4b5563',
  },

  content: {
    padding: 16,
  },

  errorContainer: {
    flexDirection: 'row',
    alignItems: 'flex-start',
    gap: 12,
    padding: 16,
    marginBottom: 16,
    backgroundColor: '#fef2f2',
    borderWidth: 1,
    borderColor: '#fecaca',
    borderRadius: 8,
  },

  errorText: {
    fontSize: 14,
    color: '#b91c1c',
    flex: 1,
  },

  card: {
    marginBottom: 16,
  },

  sectionHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
    marginBottom: 16,
  },

  sectionTitle: {
    fontSize: 18,
    fontWeight: '600',
    color: '#111827',
  },

  inputGroup: {
    marginBottom: 16,
  },

  dateButton: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 12,
    padding: 16,
    backgroundColor: colors.white,
    borderWidth: 1,
    borderColor: '#d1d5db',
    borderRadius: 8,
  },

  dateButtonText: {
    fontSize: 16,
    color: '#111827',
  },

  required: {
    color: colors.destructive,
  },

  preferenceItem: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingVertical: 12,
    borderBottomWidth: 1,
    borderBottomColor: '#f3f4f6',
  },

  preferenceLabel: {
    fontSize: 16,
    color: '#111827',
  },

  actionButtons: {
    flexDirection: 'row',
    gap: 12,
    marginBottom: 32,
  },

  cancelButton: {
    flex: 1,
  },

  submitButton: {
    flex: 1,
  },
})
