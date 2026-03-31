/**
 * LocationInput Component Styles (Mobile)
 */

import { StyleSheet } from 'react-native'
import { colors } from '@/styles/colors'

export const styles = StyleSheet.create({
  container: {
    marginBottom: 16,
  },

  inputContainer: {
    position: 'relative',
    flexDirection: 'row',
    alignItems: 'center',
  },
  input: {
    flex: 1,
    height: 48,
    backgroundColor: colors.white,
    borderWidth: 1,
    borderColor: '#d1d5db',
    borderRadius: 8,
    paddingHorizontal: 16,
    fontSize: 16,
  },

  clearButton: {
    position: 'absolute',
    right: 12,
  },
  required: {
    color: colors.red[700],
  },

  loadingContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
    marginTop: 8,
  },
  loadingText: {
    fontSize: 12,
    color: '#6b7280',
  },

  errorText: {
    fontSize: 14,
    color: colors.red[700],
    marginTop: 4,
  },

  selectedLocationCard: {
    marginTop: 8,
    padding: 12,
    backgroundColor: colors.muted,
    borderRadius: 8,
    borderWidth: 1,
    borderColor: colors.border,
  },
  selectedLocationRow: {
    flexDirection: 'row',
    marginBottom: 4,
  },
  selectedLocationLabel: {
    fontWeight: '600',
    fontSize: 14,
    color: '#374151',
    width: 96,
  },
  selectedLocationValue: {
    fontSize: 14,
    color: '#111827',
    flex: 1,
  },
  coordinatesText: {
    fontSize: 12,
    color: '#6b7280',
    marginTop: 4,
  },

  // Modal styles
  modalOverlay: {
    flex: 1,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    justifyContent: 'center',
    alignItems: 'center',
    paddingHorizontal: 16,
  },
  modalContent: {
    backgroundColor: colors.white,
    borderRadius: 16,
    width: '100%',
    maxHeight: 384,
    overflow: 'hidden',
  },

  modalHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 16,
    borderBottomWidth: 1,
    borderBottomColor: colors.border,
  },
  modalTitle: {
    fontSize: 18,
    fontWeight: '600',
    color: '#111827',
  },

  suggestionsList: {
    maxHeight: 320,
  },
  suggestionItem: {
    padding: 16,
    borderBottomWidth: 1,
    borderBottomColor: '#f3f4f6',
  },
  suggestionCity: {
    fontWeight: '600',
    fontSize: 16,
    color: '#111827',
  },
  suggestionProvince: {
    fontSize: 14,
    color: '#6b7280',
    marginTop: 4,
  },
})
