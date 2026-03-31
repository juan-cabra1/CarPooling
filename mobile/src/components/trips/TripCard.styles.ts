/**
 * TripCard Component Styles
 */

import { StyleSheet } from 'react-native'
import { colors } from '@/styles/colors'

export const styles = StyleSheet.create({
  card: {
    borderWidth: 2,
    borderColor: colors.border,
    marginBottom: 12,
  },

  header: {
    backgroundColor: colors.primary[50],
    padding: 16,
  },

  headerContent: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
  },

  routeContainer: {
    flex: 1,
  },

  routeRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
    marginBottom: 4,
  },

  routeLine: {
    width: 1,
    height: 16,
    backgroundColor: colors.border,
    marginLeft: 8,
  },

  cityText: {
    fontWeight: '600',
    color: '#1f2937',
  },

  priceContainer: {
    alignItems: 'flex-end',
  },

  price: {
    fontSize: 24,
    fontWeight: 'bold',
    color: colors.primary[500],
  },

  priceLabel: {
    fontSize: 12,
    color: colors.mutedForeground,
  },

  content: {
    padding: 16,
  },

  driverSection: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 12,
    marginBottom: 16,
    paddingBottom: 16,
    borderBottomWidth: 1,
    borderBottomColor: colors.border,
  },

  driverAvatar: {
    width: 48,
    height: 48,
    borderRadius: 24,
    backgroundColor: colors.primary[100],
    alignItems: 'center',
    justifyContent: 'center',
  },

  driverAvatarText: {
    color: colors.primary[500],
    fontWeight: 'bold',
    fontSize: 18,
  },

  driverInfo: {
    flex: 1,
  },

  driverName: {
    fontWeight: '600',
    color: '#1f2937',
  },

  driverRating: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
  },

  driverRatingText: {
    fontSize: 14,
    fontWeight: '500',
  },

  driverTrips: {
    fontSize: 14,
    color: colors.mutedForeground,
  },

  detailsGrid: {
    gap: 12,
  },

  detailRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
  },

  detailText: {
    fontSize: 14,
    color: '#374151',
    flex: 1,
  },

  preferencesSection: {
    marginTop: 12,
    paddingTop: 12,
    borderTopWidth: 1,
    borderTopColor: colors.border,
  },

  preferencesTitle: {
    fontSize: 12,
    fontWeight: '600',
    color: colors.mutedForeground,
    marginBottom: 8,
  },

  preferencesBadges: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 8,
  },

  footer: {
    padding: 12,
    backgroundColor: '#f9fafb',
    borderTopWidth: 1,
    borderTopColor: colors.border,
  },

  button: {
    backgroundColor: colors.primary[500],
    paddingVertical: 12,
    borderRadius: 8,
  },

  buttonText: {
    color: colors.white,
    textAlign: 'center',
    fontWeight: 'bold',
  },
})
