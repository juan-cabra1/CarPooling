/**
 * BookingModal Component Styles (Mobile)
 */

import { StyleSheet } from 'react-native'
import { colors } from '@/styles/colors'

export const styles = StyleSheet.create({
  overlay: {
    flex: 1,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    justifyContent: 'flex-end',
  },
  modalContainer: {
    backgroundColor: colors.white,
    borderTopLeftRadius: 24,
    borderTopRightRadius: 24,
    maxHeight: '83.33%',
  },

  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 24,
    borderBottomWidth: 1,
    borderBottomColor: '#e5e7eb',
  },
  headerTitle: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#111827',
  },

  content: {
    padding: 24,
  },

  successContainer: {
    alignItems: 'center',
    paddingVertical: 32,
  },
  successTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    color: '#111827',
    marginTop: 16,
    marginBottom: 8,
  },
  successText: {
    fontSize: 16,
    color: '#6b7280',
  },

  form: {
    gap: 24,
  },

  tripInfoCard: {
    backgroundColor: '#f9fafb',
    borderRadius: 8,
    padding: 16,
    gap: 8,
  },
  tripRoute: {
    fontWeight: '600',
    fontSize: 18,
    color: '#111827',
  },
  tripDate: {
    fontSize: 14,
    color: '#6b7280',
  },
  tripSeatsInfo: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
  },
  tripSeatsText: {
    fontSize: 14,
    color: '#6b7280',
  },

  seatsSection: {
    gap: 12,
  },
  seatsLabel: {
    fontSize: 16,
    fontWeight: '600',
  },
  seatsControl: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: 16,
  },
  seatsButton: {
    width: 48,
    height: 48,
    backgroundColor: '#e5e7eb',
    borderRadius: 24,
    alignItems: 'center',
    justifyContent: 'center',
  },
  seatsButtonText: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#374151',
  },
  seatsValue: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#111827',
    width: 64,
    textAlign: 'center',
  },
  seatsHint: {
    fontSize: 14,
    color: '#6b7280',
    textAlign: 'center',
  },

  priceSummary: {
    backgroundColor: colors.primary[50],
    borderRadius: 8,
    padding: 16,
    gap: 8,
  },
  priceRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  priceLabel: {
    fontSize: 14,
    color: '#6b7280',
  },
  priceValue: {
    fontWeight: '500',
    color: '#111827',
  },
  priceDivider: {
    height: 1,
    backgroundColor: colors.primary[200],
    marginVertical: 8,
  },
  priceTotalRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  priceTotalLabel: {
    fontWeight: '600',
    color: '#111827',
  },
  priceTotalValue: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
  },
  priceTotalAmount: {
    fontSize: 24,
    fontWeight: 'bold',
    color: colors.primary[500],
  },

  errorContainer: {
    flexDirection: 'row',
    alignItems: 'flex-start',
    gap: 8,
    padding: 16,
    backgroundColor: '#fef2f2',
    borderWidth: 1,
    borderColor: '#fecaca',
    borderRadius: 8,
  },
  errorText: {
    fontSize: 14,
    color: '#991b1b',
    flex: 1,
  },

  actionButtons: {
    flexDirection: 'row',
    gap: 12,
    paddingTop: 16,
  },
  cancelButton: {
    flex: 1,
  },
  confirmButton: {
    flex: 1,
  },
})
