/**
 * MyTripsScreen Styles
 */

import { StyleSheet } from 'react-native'
import { colors } from '@/styles/colors'

export const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  list: {
    paddingHorizontal: 16,
    paddingTop: 16,
  },

  emptyContainer: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    padding: 32,
  },
  emptyIcon: {
    marginBottom: 16,
  },
  emptyTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    color: colors.foreground,
    marginBottom: 8,
  },
  emptyText: {
    color: colors.mutedForeground,
    textAlign: 'center',
    marginBottom: 24,
  },
  emptyButton: {
    paddingHorizontal: 32,
  },

  loadingContainer: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
  },
  loadingText: {
    color: colors.mutedForeground,
    marginTop: 12,
  },

  errorContainer: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    padding: 24,
  },
  errorText: {
    color: colors.destructive,
    textAlign: 'center',
    marginBottom: 16,
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
    elevation: 2,
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
  tripRoute: {
    fontWeight: '600',
    fontSize: 18,
    color: '#111827',
  },
  tripDate: {
    color: '#6b7280',
    fontSize: 14,
    marginTop: 4,
  },
  tripPrice: {
    fontWeight: 'bold',
    fontSize: 18,
    color: colors.primary[600],
  },
  tripMetaRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginTop: 8,
  },
  tripSeatsContainer: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  tripSeatsText: {
    color: '#6b7280',
    fontSize: 14,
    marginLeft: 4,
  },
  tripStatus: {
    color: '#6b7280',
    fontSize: 14,
  },
  tripActionsRow: {
    flexDirection: 'row',
    gap: 8,
    marginTop: 12,
    paddingTop: 12,
    borderTopWidth: 1,
    borderTopColor: '#f3f4f6',
  },
  tripActionButton: {
    flex: 1,
  },
})
