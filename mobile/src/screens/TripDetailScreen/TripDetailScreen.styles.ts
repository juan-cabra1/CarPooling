/**
 * TripDetailScreen Styles
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

  errorContainer: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    paddingHorizontal: 32,
  },
  errorTitle: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#111827',
    marginTop: 16,
    marginBottom: 8,
  },
  errorText: {
    fontSize: 16,
    color: '#6b7280',
    textAlign: 'center',
    marginBottom: 24,
  },
  errorButton: {
    marginTop: 16,
  },

  backButtonContainer: {
    paddingHorizontal: 16,
    paddingTop: 16,
    paddingBottom: 8,
  },
  backButton: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
  },
  backButtonText: {
    fontSize: 16,
    fontWeight: '500',
    color: colors.primary[500],
  },

  mainCard: {
    marginHorizontal: 16,
    marginTop: 8,
  },

  // Header
  header: {
    backgroundColor: colors.primary[50],
    padding: 16,
  },
  headerContent: {
    flex: 1,
    marginBottom: 16,
  },
  headerTitleRow: {
    flexDirection: 'row',
    alignItems: 'flex-start',
    gap: 12,
    marginBottom: 8,
  },
  headerTitle: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#111827',
    flex: 1,
  },
  headerSubtitle: {
    fontSize: 16,
    color: '#6b7280',
  },
  priceContainer: {
    alignItems: 'flex-end',
  },
  price: {
    fontSize: 28,
    fontWeight: 'bold',
    color: colors.primary[500],
  },
  priceLabel: {
    fontSize: 12,
    color: '#6b7280',
  },

  // Route Details
  routeDetails: {
    gap: 24,
    paddingVertical: 24,
    borderTopWidth: 1,
    borderTopColor: '#e5e7eb',
  },
  routeItem: {
    gap: 8,
  },
  routeIconRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
  },
  routeLabel: {
    fontWeight: '600',
    fontSize: 16,
  },
  routeInfo: {
    marginLeft: 28,
  },
  routeCity: {
    fontWeight: '500',
    color: '#111827',
  },
  routeAddress: {
    fontSize: 14,
    color: '#6b7280',
  },

  // Date Section
  dateSection: {
    gap: 16,
    paddingVertical: 24,
    borderTopWidth: 1,
    borderTopColor: '#e5e7eb',
  },
  dateItem: {
    gap: 8,
  },
  dateIconRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
  },
  dateLabel: {
    fontWeight: '600',
    color: '#374151',
  },
  dateValue: {
    marginLeft: 28,
    color: '#111827',
  },

  // Trip Info
  tripInfoSection: {
    gap: 24,
    paddingVertical: 24,
    borderTopWidth: 1,
    borderTopColor: '#e5e7eb',
  },
  infoItem: {
    gap: 4,
  },
  infoIconRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
    marginBottom: 8,
  },
  infoLabel: {
    fontWeight: '600',
    color: '#374151',
  },
  infoValue: {
    marginLeft: 28,
    color: '#111827',
  },
  infoSubValue: {
    marginLeft: 28,
    fontSize: 14,
    color: '#6b7280',
  },

  // Driver Section
  driverSection: {
    paddingVertical: 24,
    borderTopWidth: 1,
    borderTopColor: '#e5e7eb',
  },
  driverIconRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
    marginBottom: 16,
  },
  driverSectionTitle: {
    fontSize: 18,
    fontWeight: '600',
    color: '#111827',
  },
  driverCard: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 16,
    padding: 16,
    backgroundColor: colors.primary[50],
    borderRadius: 8,
  },
  driverAvatar: {
    width: 64,
    height: 64,
    borderRadius: 32,
    backgroundColor: colors.primary[100],
    alignItems: 'center',
    justifyContent: 'center',
  },
  driverAvatarText: {
    color: colors.primary[500],
    fontWeight: 'bold',
    fontSize: 24,
  },
  driverInfo: {
    flex: 1,
  },
  driverName: {
    fontWeight: '600',
    color: '#111827',
    fontSize: 18,
  },
  driverRating: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
    marginTop: 4,
  },
  driverRatingText: {
    fontWeight: '500',
    color: '#111827',
  },
  driverTrips: {
    fontSize: 14,
    color: '#6b7280',
  },
  driverEmail: {
    fontSize: 14,
    color: '#6b7280',
    marginTop: 4,
  },

  // Description
  descriptionSection: {
    paddingVertical: 24,
    borderTopWidth: 1,
    borderTopColor: '#e5e7eb',
  },
  sectionTitle: {
    fontWeight: '600',
    color: '#111827',
    marginBottom: 12,
    fontSize: 18,
  },
  descriptionText: {
    color: '#374151',
  },

  // Preferences
  preferencesSection: {
    paddingVertical: 24,
    borderTopWidth: 1,
    borderTopColor: '#e5e7eb',
  },
  preferencesBadges: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 12,
  },

  // Actions
  actionSection: {
    paddingVertical: 16,
    backgroundColor: '#f9fafb',
    borderTopWidth: 1,
    borderTopColor: '#e5e7eb',
  },
  ownerActions: {
    gap: 12,
  },
  editButton: {
    width: '100%',
  },
  deleteButton: {
    width: '100%',
  },
  bookButton: {
    width: '100%',
  },

  warningBox: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
    padding: 16,
    backgroundColor: '#fef3c7',
    borderWidth: 1,
    borderColor: '#fde68a',
    borderRadius: 8,
  },
  warningText: {
    fontSize: 14,
    color: '#92400e',
    fontWeight: '500',
    flex: 1,
  },

  unavailableBox: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: 8,
    padding: 16,
    backgroundColor: '#f3f4f6',
    borderWidth: 1,
    borderColor: '#e5e7eb',
    borderRadius: 8,
  },
  unavailableText: {
    fontSize: 14,
    color: '#374151',
    fontWeight: '500',
  },
})
