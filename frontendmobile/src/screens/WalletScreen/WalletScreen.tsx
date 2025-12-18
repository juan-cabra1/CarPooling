/**
 * WalletScreen
 * Display driver wallet balance, transactions, and withdrawal options
 */

import React, { useState, useEffect, useCallback } from 'react'
import { View, Text, ScrollView, TouchableOpacity, RefreshControl, ActivityIndicator } from 'react-native'
import { useNavigation, useFocusEffect } from '@react-navigation/native'
import { type StackNavigationProp } from '@react-navigation/stack'
import { Ionicons } from '@expo/vector-icons'
import { Button, Card } from '@/components/ui'
import { useAuth } from '@/context/AuthContext'
import * as paymentsService from '@/services/paymentsService'
import { getErrorMessage, formatMoney, getTransactionTypeLabel } from '@/types'
import type { Wallet, WalletTransaction, WalletStats, SellerStatusResponse } from '@/types'
import { styles } from './WalletScreen.styles'
import type { RootStackParamList } from '@/navigation/types'

type WalletScreenNavigationProp = StackNavigationProp<RootStackParamList>

export default function WalletScreen() {
  const navigation = useNavigation<WalletScreenNavigationProp>()
  const { user } = useAuth()

  const [wallet, setWallet] = useState<Wallet | null>(null)
  const [transactions, setTransactions] = useState<WalletTransaction[]>([])
  const [stats, setStats] = useState<WalletStats | null>(null)
  const [sellerStatus, setSellerStatus] = useState<SellerStatusResponse['data'] | null>(null)
  const [loading, setLoading] = useState(true)
  const [refreshing, setRefreshing] = useState(false)
  const [error, setError] = useState('')

  const loadWalletData = async () => {
    try {
      setError('')
      const [walletData, transactionsData, statsData, sellerData] = await Promise.all([
        paymentsService.getMyWallet(),
        paymentsService.getWalletTransactions({ page: 1, page_size: 10 }),
        paymentsService.getWalletStats(),
        paymentsService.getSellerStatus(),
      ])

      setWallet(walletData)
      setTransactions(transactionsData.transactions)
      setStats(statsData)
      setSellerStatus(sellerData)
    } catch (err) {
      console.error('Error loading wallet data:', err)
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
      setRefreshing(false)
    }
  }

  useFocusEffect(
    useCallback(() => {
      loadWalletData()
    }, [])
  )

  const onRefresh = () => {
    setRefreshing(true)
    loadWalletData()
  }

  const handleWithdraw = () => {
    if (!sellerStatus?.is_linked) {
      // Navigate to seller onboarding
      navigation.navigate('SellerOnboarding' as any)
      return
    }
    navigation.navigate('Withdraw' as any)
  }

  const handleLinkAccount = () => {
    navigation.navigate('SellerOnboarding' as any)
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
      default:
        return { name: 'ellipse', color: '#6b7280' }
    }
  }

  if (loading) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color="#0ea5e9" />
        <Text style={styles.loadingText}>Cargando billetera...</Text>
      </View>
    )
  }

  if (error) {
    return (
      <View style={styles.errorContainer}>
        <Ionicons name="alert-circle" size={64} color="#dc2626" />
        <Text style={styles.errorTitle}>Error al cargar</Text>
        <Text style={styles.errorText}>{error}</Text>
        <Button title="Reintentar" onPress={loadWalletData} style={styles.retryButton} />
      </View>
    )
  }

  return (
    <View style={styles.container}>
      <ScrollView
        refreshControl={<RefreshControl refreshing={refreshing} onRefresh={onRefresh} />}
      >
        {/* Balance Card */}
        <View style={styles.balanceCard}>
          <Text style={styles.balanceLabel}>Saldo disponible</Text>
          <Text style={styles.balanceAmount}>
            {wallet?.available_balance_formatted || formatMoney(0)}
          </Text>

          <View style={styles.balanceDetails}>
            <View style={styles.balanceDetailItem}>
              <Text style={styles.balanceDetailLabel}>Pendiente</Text>
              <Text style={styles.balanceDetailValue}>
                {formatMoney(wallet?.pending_balance || 0)}
              </Text>
            </View>
            <View style={styles.balanceDetailDivider} />
            <View style={styles.balanceDetailItem}>
              <Text style={styles.balanceDetailLabel}>Retenido</Text>
              <Text style={styles.balanceDetailValue}>
                {formatMoney(wallet?.held_balance || 0)}
              </Text>
            </View>
          </View>

          <Button
            title="Retirar fondos"
            onPress={handleWithdraw}
            disabled={(wallet?.available_balance || 0) <= 0}
            style={styles.withdrawButton}
          />
        </View>

        {/* MercadoPago Account Status */}
        {!sellerStatus?.is_linked && (
          <Card style={styles.linkAccountCard}>
            <View style={styles.linkAccountHeader}>
              <Ionicons name="warning" size={24} color="#f59e0b" />
              <Text style={styles.linkAccountTitle}>Vincula tu cuenta</Text>
            </View>
            <Text style={styles.linkAccountText}>
              Para recibir pagos y retirar fondos, necesitas vincular tu cuenta de MercadoPago.
            </Text>
            <Button
              title="Vincular MercadoPago"
              variant="outline"
              onPress={handleLinkAccount}
              style={styles.linkAccountButton}
            />
          </Card>
        )}

        {/* Stats */}
        {stats && (
          <Card style={styles.statsCard}>
            <Text style={styles.cardTitle}>Estadisticas</Text>
            <View style={styles.statsGrid}>
              <View style={styles.statItem}>
                <Text style={styles.statValue}>{stats.total_trips}</Text>
                <Text style={styles.statLabel}>Viajes</Text>
              </View>
              <View style={styles.statItem}>
                <Text style={styles.statValue}>
                  {wallet?.total_earned_formatted || formatMoney(stats.total_earnings)}
                </Text>
                <Text style={styles.statLabel}>Total ganado</Text>
              </View>
              <View style={styles.statItem}>
                <Text style={styles.statValue}>{formatMoney(stats.total_withdrawals)}</Text>
                <Text style={styles.statLabel}>Retirado</Text>
              </View>
              <View style={styles.statItem}>
                <Text style={styles.statValue}>{formatMoney(stats.average_per_trip)}</Text>
                <Text style={styles.statLabel}>Promedio/viaje</Text>
              </View>
            </View>
          </Card>
        )}

        {/* Recent Transactions */}
        <Card style={styles.transactionsCard}>
          <View style={styles.transactionsHeader}>
            <Text style={styles.cardTitle}>Movimientos recientes</Text>
            <TouchableOpacity onPress={() => navigation.navigate('Transactions' as any)}>
              <Text style={styles.seeAllLink}>Ver todos</Text>
            </TouchableOpacity>
          </View>

          {transactions.length === 0 ? (
            <View style={styles.emptyTransactions}>
              <Ionicons name="receipt-outline" size={48} color="#9ca3af" />
              <Text style={styles.emptyText}>No hay movimientos aun</Text>
            </View>
          ) : (
            <View style={styles.transactionsList}>
              {transactions.map((tx) => {
                const icon = getTransactionIcon(tx.type)
                return (
                  <View key={tx.transaction_uuid} style={styles.transactionItem}>
                    <View style={styles.transactionIcon}>
                      <Ionicons name={icon.name as any} size={24} color={icon.color} />
                    </View>
                    <View style={styles.transactionInfo}>
                      <Text style={styles.transactionType}>
                        {getTransactionTypeLabel(tx.type)}
                      </Text>
                      <Text style={styles.transactionDescription}>{tx.description}</Text>
                      <Text style={styles.transactionDate}>
                        {new Date(tx.created_at).toLocaleDateString('es-AR')}
                      </Text>
                    </View>
                    <Text
                      style={[
                        styles.transactionAmount,
                        { color: tx.signed_amount >= 0 ? '#22c55e' : '#ef4444' },
                      ]}
                    >
                      {tx.signed_amount >= 0 ? '+' : ''}
                      {tx.amount_formatted}
                    </Text>
                  </View>
                )
              })}
            </View>
          )}
        </Card>

        {/* Wallet Settings */}
        <Card style={styles.settingsCard}>
          <Text style={styles.cardTitle}>Configuracion</Text>
          <TouchableOpacity style={styles.settingItem}>
            <Ionicons name="repeat" size={20} color="#6b7280" />
            <View style={styles.settingInfo}>
              <Text style={styles.settingLabel}>Retiro automatico</Text>
              <Text style={styles.settingValue}>
                {wallet?.auto_withdraw_enabled ? 'Activado' : 'Desactivado'}
              </Text>
            </View>
            <Ionicons name="chevron-forward" size={20} color="#9ca3af" />
          </TouchableOpacity>

          {sellerStatus?.is_linked && (
            <TouchableOpacity style={styles.settingItem}>
              <Ionicons name="card" size={20} color="#6b7280" />
              <View style={styles.settingInfo}>
                <Text style={styles.settingLabel}>Cuenta MercadoPago</Text>
                <Text style={styles.settingValue}>Vinculada</Text>
              </View>
              <Ionicons name="checkmark-circle" size={20} color="#22c55e" />
            </TouchableOpacity>
          )}
        </Card>
      </ScrollView>
    </View>
  )
}
