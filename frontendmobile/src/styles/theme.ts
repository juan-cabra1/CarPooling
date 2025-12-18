/**
 * CarPooling Mobile - Theme System
 * Reusable style classes for consistent design across the app
 * Matches exactly the web frontend design
 */

export const theme = {
  // Container styles
  containers: {
    page: 'flex-1 bg-background',
    pageWithPadding: 'flex-1 bg-background p-4',
    card: 'bg-card p-4 rounded-lg border border-border',
    cardShadow: 'bg-card p-4 rounded-lg shadow-md',
  },

  // Text styles
  text: {
    // Headings
    h1: 'text-3xl font-bold text-foreground',
    h2: 'text-2xl font-bold text-foreground',
    h3: 'text-xl font-bold text-foreground',
    h4: 'text-lg font-semibold text-foreground',

    // Body
    body: 'text-base text-foreground',
    bodySmall: 'text-sm text-foreground',
    bodyLarge: 'text-lg text-foreground',

    // Muted
    muted: 'text-sm text-muted-foreground',
    mutedSmall: 'text-xs text-muted-foreground',

    // Special
    label: 'text-sm font-medium text-foreground mb-1',
    link: 'text-primary-500 font-semibold',
    error: 'text-destructive text-sm',
    success: 'text-success text-sm',
    white: 'text-white',
  },

  // Button styles
  buttons: {
    // Primary button
    primary: 'bg-primary-500 px-6 py-3 rounded-lg active:opacity-80',
    primaryText: 'text-white text-center font-bold text-base',

    // Secondary button
    secondary: 'bg-secondary-500 px-6 py-3 rounded-lg active:opacity-80',
    secondaryText: 'text-white text-center font-bold text-base',

    // Outline button
    outline: 'border-2 border-primary-500 px-6 py-3 rounded-lg active:bg-primary-50',
    outlineText: 'text-primary-500 text-center font-bold text-base',

    // Destructive button
    destructive: 'bg-destructive px-6 py-3 rounded-lg active:opacity-80',
    destructiveText: 'text-white text-center font-bold text-base',

    // Ghost button
    ghost: 'px-6 py-3 rounded-lg active:bg-muted',
    ghostText: 'text-foreground text-center font-semibold text-base',

    // Disabled state
    disabled: 'bg-muted px-6 py-3 rounded-lg',
    disabledText: 'text-muted-foreground text-center font-bold text-base',

    // Small button
    small: 'px-4 py-2 rounded-lg',
    smallText: 'text-sm font-semibold',
  },

  // Input styles
  inputs: {
    base: 'bg-white border border-input rounded-lg px-4 py-3 text-base text-foreground',
    error: 'bg-white border-2 border-destructive rounded-lg px-4 py-3 text-base text-foreground',
    focused: 'bg-white border-2 border-primary-500 rounded-lg px-4 py-3 text-base text-foreground',
    disabled: 'bg-muted border border-input rounded-lg px-4 py-3 text-base text-muted-foreground',
    textarea: 'bg-white border border-input rounded-lg px-4 py-3 text-base text-foreground min-h-[100px]',
  },

  // Form styles
  forms: {
    group: 'mb-4',
    label: 'text-sm font-medium text-foreground mb-2',
    error: 'text-destructive text-sm mt-1',
    helper: 'text-muted-foreground text-sm mt-1',
  },

  // Badge styles
  badges: {
    default: 'bg-primary-100 px-3 py-1 rounded-full',
    defaultText: 'text-primary-700 text-xs font-semibold',

    success: 'bg-green-100 px-3 py-1 rounded-full',
    successText: 'text-green-700 text-xs font-semibold',

    warning: 'bg-yellow-100 px-3 py-1 rounded-full',
    warningText: 'text-yellow-700 text-xs font-semibold',

    destructive: 'bg-red-100 px-3 py-1 rounded-full',
    destructiveText: 'text-red-700 text-xs font-semibold',

    secondary: 'bg-secondary-100 px-3 py-1 rounded-full',
    secondaryText: 'text-secondary-700 text-xs font-semibold',
  },

  // Card styles
  cards: {
    base: 'bg-card rounded-lg border border-border p-4',
    shadow: 'bg-card rounded-lg shadow-md p-4',
    header: 'mb-3',
    title: 'text-lg font-bold text-card-foreground',
    description: 'text-sm text-muted-foreground mt-1',
    content: 'mt-2',
    footer: 'mt-4 pt-4 border-t border-border',
  },

  // Trip Card specific styles
  tripCard: {
    container: 'bg-card rounded-lg border border-border p-4 mb-3',
    header: 'flex-row justify-between items-start mb-3',
    route: 'flex-1 mr-2',
    price: 'text-2xl font-bold text-primary-500',
    priceLabel: 'text-xs text-muted-foreground',
    origin: 'text-base font-semibold text-foreground',
    destination: 'text-base font-semibold text-foreground',
    arrow: 'text-muted-foreground text-sm my-1',
    datetime: 'text-sm text-muted-foreground mt-2',
    driver: 'flex-row items-center mt-3 pt-3 border-t border-border',
    driverAvatar: 'w-10 h-10 rounded-full bg-primary-100 items-center justify-center mr-3',
    driverInfo: 'flex-1',
    driverName: 'text-sm font-semibold text-foreground',
    driverRating: 'text-xs text-muted-foreground',
    preferences: 'flex-row flex-wrap gap-2 mt-3',
    seats: 'text-sm font-semibold text-foreground',
    seatsLabel: 'text-sm text-muted-foreground',
  },

  // Modal styles
  modals: {
    overlay: 'flex-1 bg-black/50 justify-center items-center',
    container: 'bg-card rounded-lg p-6 m-4 max-w-md w-full',
    header: 'mb-4',
    title: 'text-xl font-bold text-foreground',
    content: 'mb-6',
    footer: 'flex-row justify-end gap-3',
  },

  // Header/Navbar styles
  headers: {
    container: 'bg-primary-500 px-4 py-3 flex-row items-center',
    title: 'text-white text-lg font-bold flex-1 text-center',
    backButton: 'mr-3',
    rightButton: 'ml-3',
  },

  // Tab Bar styles
  tabs: {
    container: 'bg-card border-t border-border',
    label: 'text-xs mt-1',
  },

  // List styles
  lists: {
    container: 'flex-1',
    separator: 'h-3',
    empty: 'flex-1 items-center justify-center p-8',
    emptyText: 'text-muted-foreground text-center text-base',
    emptyIcon: 'mb-3',
  },

  // Loading styles
  loading: {
    container: 'flex-1 items-center justify-center',
    text: 'text-muted-foreground mt-3',
  },

  // Spacing helpers
  spacing: {
    xs: 'gap-1',
    sm: 'gap-2',
    md: 'gap-3',
    lg: 'gap-4',
    xl: 'gap-6',
  },

  // Gradient backgrounds (for LinearGradient component)
  gradients: {
    primary: ['#0ea5e9', '#0284c7'] as const,
    secondary: ['#14b8a6', '#0d9488'] as const,
    accent: ['#f97316', '#ea580c'] as const,
  },
}

// Export individual sections for easier imports
export const { containers, text, buttons, inputs, forms, badges, cards, tripCard, modals, headers, tabs, lists, loading, spacing, gradients } = theme

export default theme
