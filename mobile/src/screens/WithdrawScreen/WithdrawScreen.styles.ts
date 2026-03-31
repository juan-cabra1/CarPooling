/**
 * WithdrawScreen Styles
 */

import { StyleSheet } from 'react-native'
import { colors } from '@/styles/colors'

export const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.muted,
  },

  loadingContainer: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: colors.white,
    gap: 16,
  },
  loadingText: {
    fontSize: 16,
    color: '#6b7280',
  },

  // Balance Card
  balanceCard: {
    margin: 16,
    marginBottom: 8,
    alignItems: 'center',
    padding: 24,
  },
  balanceLabel: {
    fontSize: 14,
    color: '#6b7280',
    marginBottom: 8,
  },
  balanceAmount: {
    fontSize: 32,
    fontWeight: 'bold',
    color: colors.primary[500],
  },

  // Amount Card
  amountCard: {
    marginHorizontal: 16,
    marginBottom: 16,
  },
  inputLabel: {
    fontSize: 14,
    fontWeight: '500',
    color: '#6b7280',
    marginBottom: 12,
  },
  amountInputContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#f9fafb',
    borderRadius: 12,
    paddingHorizontal: 16,
    marginBottom: 16,
  },
  currencySymbol: {
    fontSize: 32,
    fontWeight: 'bold',
    color: '#111827',
  },
  amountInput: {
    flex: 1,
    fontSize: 32,
    fontWeight: 'bold',
    color: '#111827',
    paddingVertical: 16,
    paddingHorizontal: 8,
  },
  withdrawAllButton: {
    alignSelf: 'flex-start',
  },

  // Summary Card
  summaryCard: {
    marginHorizontal: 16,
    marginBottom: 16,
  },
  summaryTitle: {
    fontSize: 16,
    fontWeight: '600',
    color: '#111827',
    marginBottom: 16,
  },
  summaryRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 12,
  },
  summaryLabel: {
    fontSize: 14,
    color: '#6b7280',
  },
  summaryValue: {
    fontSize: 14,
    color: '#111827',
  },
  summaryLabelBold: {
    fontSize: 16,
    fontWeight: '600',
    color: '#111827',
  },
  summaryValueBold: {
    fontSize: 18,
    fontWeight: 'bold',
    color: colors.primary[500],
  },
  summaryDivider: {
    height: 1,
    backgroundColor: '#e5e7eb',
    marginVertical: 12,
  },

  // Info Card
  infoCard: {
    marginHorizontal: 16,
    marginBottom: 16,
    backgroundColor: '#f0f9ff',
    borderWidth: 1,
    borderColor: '#bae6fd',
  },
  infoItem: {
    flexDirection: 'row',
    alignItems: 'flex-start',
    gap: 12,
    marginBottom: 12,
  },
  infoText: {
    flex: 1,
    fontSize: 14,
    color: '#0369a1',
  },

  // Error
  errorContainer: {
    flexDirection: 'row',
    alignItems: 'flex-start',
    gap: 8,
    marginHorizontal: 16,
    marginBottom: 16,
    padding: 16,
    backgroundColor: '#fef2f2',
    borderWidth: 1,
    borderColor: '#fecaca',
    borderRadius: 8,
  },
  errorText: {
    flex: 1,
    fontSize: 14,
    color: '#991b1b',
  },

  // Submit
  submitContainer: {
    padding: 16,
    paddingBottom: 32,
  },
  submitButton: {
    width: '100%',
  },
})
