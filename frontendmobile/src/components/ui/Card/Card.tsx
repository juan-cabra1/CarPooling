/**
 * Card Component
 * Reusable card container with header, content, and footer sections
 */

import React from 'react'
import { View, Text, type ViewProps, type TextStyle, type ViewStyle } from 'react-native'
import { styles } from './Card.styles'

export interface CardProps extends ViewProps {
  /** Use shadow instead of border */
  shadow?: boolean
}

export interface CardHeaderProps extends ViewProps {}
export interface CardTitleProps {
  children: React.ReactNode
  style?: TextStyle
}
export interface CardDescriptionProps {
  children: React.ReactNode
  style?: TextStyle
}
export interface CardContentProps extends ViewProps {}
export interface CardFooterProps extends ViewProps {}

export function Card({ shadow = false, style, children, ...props }: CardProps) {
  const cardStyle = shadow ? styles.cardShadow : styles.card

  return (
    <View style={[cardStyle, style]} {...props}>
      {children}
    </View>
  )
}

export function CardHeader({ style, children, ...props }: CardHeaderProps) {
  return (
    <View style={[styles.header, style]} {...props}>
      {children}
    </View>
  )
}

export function CardTitle({ children, style }: CardTitleProps) {
  return <Text style={[styles.title, style]}>{children}</Text>
}

export function CardDescription({ children, style }: CardDescriptionProps) {
  return <Text style={[styles.description, style]}>{children}</Text>
}

export function CardContent({ style, children, ...props }: CardContentProps) {
  return (
    <View style={[styles.content, style]} {...props}>
      {children}
    </View>
  )
}

export function CardFooter({ style, children, ...props }: CardFooterProps) {
  return (
    <View style={[styles.footer, style]} {...props}>
      {children}
    </View>
  )
}

export default Card
