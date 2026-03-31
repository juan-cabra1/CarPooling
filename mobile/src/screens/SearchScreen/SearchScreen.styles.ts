/**
 * SearchScreen Styles
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
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: colors.white,
  },
  loadingText: {
    marginTop: 12,
    fontSize: 16,
    color: '#6b7280',
  },

  header: {
    marginBottom: 24,
  },
  headerTitle: {
    fontSize: 30,
    fontWeight: 'bold',
    color: '#111827',
    marginBottom: 4,
  },
  headerSubtitle: {
    fontSize: 16,
    color: '#6b7280',
  },

  searchCard: {
    marginBottom: 16,
  },
  searchCardHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
    marginBottom: 16,
  },
  searchCardTitle: {
    fontSize: 18,
    fontWeight: '600',
    color: '#111827',
  },

  searchForm: {
    gap: 16,
  },

  inputGroup: {
    marginBottom: 16,
  },
  dateInputContainer: {
    position: 'relative',
    flexDirection: 'row',
    alignItems: 'center',
  },
  dateInput: {
    flex: 1,
  },

  filtersToggle: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
    paddingVertical: 12,
    borderTopWidth: 1,
    borderTopColor: colors.border,
  },
  filtersToggleText: {
    fontSize: 16,
    fontWeight: '500',
    color: colors.primary[500],
  },

  advancedFilters: {
    gap: 16,
    paddingVertical: 16,
    borderTopWidth: 1,
    borderTopColor: colors.border,
  },
  filterItem: {
    marginBottom: 12,
  },
  filterInputContainer: {
    position: 'relative',
    flexDirection: 'row',
    alignItems: 'center',
  },

  actionButtons: {
    flexDirection: 'row',
    gap: 12,
    marginTop: 16,
  },
  clearButton: {
    flex: 1,
  },
  searchButton: {
    flex: 1,
  },

  errorContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
    padding: 16,
    marginBottom: 16,
    backgroundColor: '#fee2e2',
    borderWidth: 1,
    borderColor: '#fecaca',
    borderRadius: 8,
  },
  errorText: {
    fontSize: 14,
    color: '#b91c1c',
    flex: 1,
  },

  resultsHeader: {
    marginBottom: 16,
  },
  resultsTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    color: '#111827',
  },

  emptyContainer: {
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: 48,
  },
  emptyTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    color: '#111827',
    marginBottom: 8,
    marginTop: 16,
  },
  emptyText: {
    fontSize: 16,
    color: '#6b7280',
    textAlign: 'center',
    paddingHorizontal: 32,
    marginBottom: 24,
  },
  emptyButton: {
    marginTop: 16,
  },

  pagination: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingVertical: 24,
    borderTopWidth: 1,
    borderTopColor: colors.border,
    marginTop: 16,
  },
  paginationButton: {
    flex: 1,
    marginHorizontal: 4,
  },
  paginationText: {
    fontSize: 14,
    color: '#374151',
    paddingHorizontal: 16,
  },

  // TripCard styles
  tripCard: {
    backgroundColor: colors.white,
    borderRadius: 8,
    borderWidth: 1,
    borderColor: colors.border,
    padding: 16,
    marginBottom: 12,
    shadowColor: colors.black,
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.05,
    shadowRadius: 2,
    elevation: 1,
  },
  tripCardHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    marginBottom: 8,
  },
  tripCardContent: {
    flex: 1,
  },
  tripCardRoute: {
    fontWeight: '600',
    fontSize: 18,
    color: '#111827',
  },
  tripCardDate: {
    color: '#6b7280',
    fontSize: 14,
    marginTop: 4,
  },
  tripCardPrice: {
    fontWeight: 'bold',
    fontSize: 18,
    color: colors.primary[600],
  },
  tripCardFooter: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginTop: 8,
  },
  tripCardInfo: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  tripCardInfoText: {
    color: '#6b7280',
    fontSize: 14,
    marginLeft: 4,
  },
  tripCardDriver: {
    flexDirection: 'row',
    alignItems: 'center',
    marginTop: 8,
    paddingTop: 8,
    borderTopWidth: 1,
    borderTopColor: '#f3f4f6',
  },
  tripCardDriverText: {
    color: '#6b7280',
    fontSize: 14,
    marginLeft: 4,
  },

  // CustomPicker styles
  pickerContainer: {
    borderWidth: 1,
    borderColor: colors.input,
    borderRadius: 8,
    backgroundColor: colors.white,
  },
  pickerOption: {
    paddingVertical: 12,
    paddingHorizontal: 16,
  },
  pickerOptionBorder: {
    borderBottomWidth: 1,
    borderBottomColor: '#f3f4f6',
  },
  pickerOptionSelected: {
    backgroundColor: colors.primary[50],
  },
  pickerOptionText: {
    color: '#111827',
  },
  pickerOptionTextSelected: {
    color: colors.primary[600],
    fontWeight: '500',
  },
})
