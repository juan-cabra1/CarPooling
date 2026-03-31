/**
 * HomeScreen Styles
 */

import { StyleSheet } from 'react-native'
import { colors } from '@/styles/colors'

export const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.white,
  },
  scrollContent: {
    paddingBottom: 24,
  },

  // Hero Section
  hero: {
    paddingHorizontal: 24,
    paddingVertical: 48,
  },
  heroTitle: {
    fontSize: 36,
    fontWeight: 'bold',
    textAlign: 'center',
    marginBottom: 16,
  },
  heroTitlePrimary: {
    color: colors.primary[600],
  },
  heroTitleSecondary: {
    color: colors.secondary[600],
  },
  heroSubtitle: {
    fontSize: 18,
    color: '#374151',
    textAlign: 'center',
    marginBottom: 32,
  },

  // CTA Buttons
  ctaContainer: {
    gap: 16,
    marginBottom: 40,
  },
  ctaButtonSearch: {
    backgroundColor: colors.primary[500],
    borderRadius: 8,
    paddingVertical: 16,
    paddingHorizontal: 24,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.1,
    shadowRadius: 8,
    elevation: 4,
  },
  ctaButtonCreate: {
    backgroundColor: colors.secondary[500],
    borderRadius: 8,
    paddingVertical: 16,
    paddingHorizontal: 24,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.1,
    shadowRadius: 8,
    elevation: 4,
  },
  ctaButtonActive: {
    opacity: 0.8,
  },
  ctaButtonText: {
    color: colors.white,
    fontSize: 18,
    fontWeight: 'bold',
    marginLeft: 8,
  },

  // Stats
  statsContainer: {
    flexDirection: 'row',
    gap: 12,
  },
  statCard: {
    flex: 1,
    backgroundColor: 'rgba(255, 255, 255, 0.8)',
    borderRadius: 8,
    padding: 16,
    alignItems: 'center',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 2,
  },
  statNumber: {
    fontSize: 24,
    fontWeight: 'bold',
    marginBottom: 4,
  },
  statNumberPrimary: {
    color: colors.primary[500],
  },
  statNumberSecondary: {
    color: colors.secondary[500],
  },
  statNumberAccent: {
    color: colors.accent.DEFAULT,
  },
  statLabel: {
    fontSize: 12,
    color: colors.mutedForeground,
    textAlign: 'center',
  },

  // Featured Trips Section
  section: {
    paddingHorizontal: 16,
    paddingVertical: 24,
  },
  sectionHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 16,
  },
  sectionTitle: {
    fontSize: 24,
    fontWeight: 'bold',
    color: colors.foreground,
  },
  sectionLink: {
    color: colors.primary[500],
    fontWeight: '600',
  },

  // Loading
  loadingContainer: {
    paddingVertical: 32,
    alignItems: 'center',
  },
  loadingText: {
    color: colors.mutedForeground,
    marginTop: 12,
  },

  // Empty State
  emptyContainer: {
    paddingVertical: 48,
    alignItems: 'center',
  },
  emptyText: {
    color: colors.mutedForeground,
    textAlign: 'center',
    marginTop: 12,
  },

  // Features Section
  featuresGrid: {
    paddingHorizontal: 16,
    paddingVertical: 32,
    gap: 24,
  },
  featureCard: {
    backgroundColor: colors.white,
    borderRadius: 8,
    padding: 24,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 2,
    borderWidth: 1,
    borderColor: colors.border,
  },
  featureIconContainer: {
    width: 48,
    height: 48,
    borderRadius: 24,
    alignItems: 'center',
    justifyContent: 'center',
    marginBottom: 16,
  },
  featureIconPrimary: {
    backgroundColor: colors.primary[100],
  },
  featureIconSecondary: {
    backgroundColor: colors.secondary[100],
  },
  featureIconGreen: {
    backgroundColor: colors.green[100],
  },
  featureTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    color: colors.foreground,
    marginBottom: 8,
  },
  featureDescription: {
    fontSize: 14,
    color: colors.mutedForeground,
  },
})
