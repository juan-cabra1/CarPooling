/**
 * TransactionsScreen
 * Full transaction history for wallet
 */

import React, { useState, useEffect } from 'react'
import { View, Text, ScrollView, RefreshControl, ActivityIndicator, FlatList } from 'react-native'
import { Ionicons } from '@expo/vector-icons'
import { Card } from '@/components/ui'
import * as paymentsService from '@/services/paymentsService'
import { getErrorMessage, formatMoney, getTransactionTypeLabel } from '@/types'
import type { WalletTransaction } from '@/types'
import { styles } from './TransactionsScreen.styles'

export default function TransactionsScreen() {
  const [transactions, setTransactions] = useState<WalletTransaction[]>([])
  const [loading, setLoading] = useState(true)
  const [refreshing, setRefreshing] = useState(false)
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [error, setError] = useState('')

  useEffect(() => {
    loadTransactions()
  }, [page])

  const loadTransactions = async () => {
    try {
      setError('')
      const response = await paymentsService.getWalletTransactions({
        page,
        page_size: 20,
      })

      if (page === 1) {
        setTransactions(response.transactions)
      } else {
        setTransactions((prev) => [...prev, ...response.transactions])
      }

      setTotalPages(response.total_pages)
    } catch (err) {
      console.error('Error loading transactions:', err)
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
      setRefreshing(false)
    }
  }

  const onRefresh = () => {
    setRefreshing(true)
    setPage(1)
    loadTransactions()
  }

  const loadMore = () => {
    if (page < totalPages && !loading) {
      setPage(page + 1)
    }
  }

  const getTransactionIcon = (type: string) => {
    switch (type) {
      case 'credit':
        return { name: 'arrow-down-circle', color: '#22c55e' }
      case 'debit':
      case 'withdrawal':
        return { name: 'arrow-up-circle', color: '#ef4444' }
      case 'hold':
        return { name: 'time', color: '#f59e0b' }
      case 'release':
        return { name: 'checkmark-circle', color: '#22c55e' }
      case 'refund':
        return { name: 'refresh-circle', color: '#3b82f6' }
      case 'fee':
        return { name: 'receipt', color: '#9ca3af' }
      default:
        return { name: 'ellipse', color: '#6b7280' }
    }
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return new Intl.DateTimeFormat('es-AR', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    }).format(date)
  }

  const renderTransaction = ({ item }: { item: WalletTransaction }) => {
    const icon = getTransactionIcon(item.type)

    return (
      <Card style={styles.transactionCard}>
        <View style={styles.transactionHeader}>
          <View style={styles.transactionIcon}>
            <Ionicons name={icon.name as any} size={24} color={icon.color} />
          </View>
          <View style={styles.transactionInfo}>
            <Text style={styles.transactionType}>{getTransactionTypeLabel(item.type)}</Text>
            <Text style={styles.transactionDescription}>{item.description}</Text>
            <Text style={styles.transactionDate}>{formatDate(item.created_at)}</Text>
          </View>
          <Text
            style={[
              styles.transactionAmount,
              { color: item.signed_amount >= 0 ? '#22c55e' : '#ef4444' },
            ]}
          >
            {item.signed_amount >= 0 ? '+' : ''}
            {item.amount_formatted}
          </Text>
        </View>

        <View style={styles.transactionDetails}>
          <View style={styles.detailRow}>
            <Text style={styles.detailLabel}>Saldo anterior</Text>
            <Text style={styles.detailValue}>{formatMoney(item.balance_before)}</Text>
          </View>
          <View style={styles.detailRow}>
            <Text style={styles.detailLabel}>Saldo posterior</Text>
            <Text style={styles.detailValue}>{formatMoney(item.balance_after)}</Text>
          </View>
        </View>

        {item.payment_id && (
          <View style={styles.referenceRow}>
            <Ionicons name="link" size={14} color="#6b7280" />
            <Text style={styles.referenceText}>Pago: {item.payment_id}</Text>
          </View>
        )}

        {item.withdrawal_id && (
          <View style={styles.referenceRow}>
            <Ionicons name="link" size={14} color="#6b7280" />
            <Text style={styles.referenceText}>Retiro: {item.withdrawal_id}</Text>
          </View>
        )}
      </Card>
    )
  }

  if (loading && page === 1) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color="#0ea5e9" />
        <Text style={styles.loadingText}>Cargando transacciones...</Text>
      </View>
    )
  }

  if (error && transactions.length === 0) {
    return (
      <View style={styles.errorContainer}>
        <Ionicons name="alert-circle" size={64} color="#dc2626" />
        <Text style={styles.errorTitle}>Error al cargar</Text>
        <Text style={styles.errorText}>{error}</Text>
      </View>
    )
  }

  return (
    <View style={styles.container}>
      <FlatList
        data={transactions}
        renderItem={renderTransaction}
        keyExtractor={(item) => item.transaction_uuid}
        refreshControl={<RefreshControl refreshing={refreshing} onRefresh={onRefresh} />}
        onEndReached={loadMore}
        onEndReachedThreshold={0.5}
        ListEmptyComponent={
          <View style={styles.emptyContainer}>
            <Ionicons name="receipt-outline" size={64} color="#9ca3af" />
            <Text style={styles.emptyTitle}>Sin transacciones</Text>
            <Text style={styles.emptyText}>Aun no tienes movimientos en tu billetera</Text>
          </View>
        }
        ListFooterComponent={
          loading && page > 1 ? (
            <View style={styles.footerLoader}>
              <ActivityIndicator size="small" color="#0ea5e9" />
            </View>
          ) : null
        }
        contentContainerStyle={styles.listContent}
      />
    </View>
  )
}
