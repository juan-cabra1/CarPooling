/**
 * ProfileScreen Styles
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
  },
  loadingText: {
    fontSize: 16,
    color: '#6b7280',
  },

  header: {
    backgroundColor: colors.white,
    alignItems: 'center',
    paddingVertical: 32,
    paddingHorizontal: 16,
    borderBottomWidth: 1,
    borderBottomColor: colors.border,
  },
  avatarContainer: {
    marginBottom: 16,
  },
  avatar: {
    width: 96,
    height: 96,
    backgroundColor: colors.primary[500],
    borderRadius: 48,
    alignItems: 'center',
    justifyContent: 'center',
  },
  avatarText: {
    color: colors.white,
    fontWeight: 'bold',
    fontSize: 36,
  },
  userName: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#111827',
    marginBottom: 4,
  },
  userEmail: {
    fontSize: 16,
    color: '#6b7280',
  },

  card: {
    marginHorizontal: 16,
    marginBottom: 0,
    marginTop: 16,
  },
  cardHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
    marginBottom: 16,
  },
  cardTitle: {
    fontSize: 18,
    fontWeight: '600',
    color: '#111827',
  },

  infoSection: {
    gap: 12,
  },
  infoRow: {
    gap: 4,
  },
  infoLabel: {
    fontSize: 14,
    fontWeight: '500',
    color: '#6b7280',
  },
  infoValue: {
    fontSize: 16,
    color: '#111827',
  },

  verificationBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
  },
  verificationTextSuccess: {
    fontSize: 16,
    color: colors.success,
    fontWeight: '500',
  },
  verificationTextError: {
    fontSize: 16,
    color: colors.destructive,
    fontWeight: '500',
  },

  statsGrid: {
    flexDirection: 'row',
    gap: 16,
  },
  statItem: {
    flex: 1,
    alignItems: 'center',
    padding: 16,
    backgroundColor: colors.muted,
    borderRadius: 8,
  },
  statValue: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#111827',
  },
  statLabel: {
    fontSize: 14,
    color: '#6b7280',
    marginTop: 4,
  },
  ratingContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
  },

  actionsSection: {
    gap: 4,
  },
  actionButton: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 12,
    paddingVertical: 16,
    borderBottomWidth: 1,
    borderBottomColor: '#f3f4f6',
  },
  actionButtonText: {
    flex: 1,
    fontSize: 16,
    color: '#111827',
  },

  logoutContainer: {
    paddingHorizontal: 16,
    paddingVertical: 24,
  },
  logoutButton: {
    width: '100%',
  },

  versionContainer: {
    alignItems: 'center',
    paddingVertical: 24,
    paddingBottom: 32,
  },
  versionText: {
    fontSize: 14,
    color: '#9ca3af',
  },
})
