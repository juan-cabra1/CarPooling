/**
 * WithdrawScreen
 * Request withdrawal from wallet to MercadoPago account
 */

import React, { useState, useEffect } from 'react'
import { View, Text, ScrollView, TextInput, Alert, ActivityIndicator } from 'react-native'
import { useNavigation } from '@react-navigation/native'
import { type StackNavigationProp } from '@react-navigation/stack'
import { Ionicons } from '@expo/vector-icons'
import { Button, Card } from '@/components/ui'
import * as paymentsService from '@/services/paymentsService'
import { getErrorMessage, formatMoney } from '@/types'
import type { Wallet } from '@/types'
import { styles } from './WithdrawScreen.styles'
import type { RootStackParamList } from '@/navigation/types'

type WithdrawScreenNavigationProp = StackNavigationProp<RootStackParamList>

export default function WithdrawScreen() {
  const navigation = useNavigation<WithdrawScreenNavigationProp>()

  const [wallet, setWallet] = useState<Wallet | null>(null)
  const [amount, setAmount] = useState('')
  const [loading, setLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    loadWallet()
  }, [])

  const loadWallet = async () => {
    try {
      const walletData = await paymentsService.getMyWallet()
      setWallet(walletData)
    } catch (err) {
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
    }
  }

  const handleAmountChange = (text: string) => {
    // Only allow numbers and one decimal point
    const cleaned = text.replace(/[^0-9.]/g, '')
    const parts = cleaned.split('.')
    if (parts.length > 2) return
    if (parts[1] && parts[1].length > 2) return
    setAmount(cleaned)
  }

  const handleWithdrawAll = () => {
    if (wallet) {
      const availablePesos = wallet.available_balance / 100
      setAmount(availablePesos.toFixed(2))
    }
  }

  const handleSubmit = async () => {
    if (!amount || parseFloat(amount) <= 0) {
      setError('Ingresa un monto valido')
      return
    }

    const amountCentavos = Math.round(parseFloat(amount) * 100)

    if (wallet && amountCentavos > wallet.available_balance) {
      setError('El monto excede tu saldo disponible')
      return
    }

    if (amountCentavos < 1000) {
      setError('El monto minimo de retiro es $10')
      return
    }

    Alert.alert(
      'Confirmar retiro',
      `Vas a retirar ${formatMoney(amountCentavos)} a tu cuenta de MercadoPago. El dinero estara disponible en 1-2 dias habiles.`,
      [
        { text: 'Cancelar', style: 'cancel' },
        {
          text: 'Confirmar',
          onPress: async () => {
            try {
              setSubmitting(true)
              setError('')

              await paymentsService.requestWithdrawal({ amount: amountCentavos })

              Alert.alert(
                'Retiro solicitado',
                'Tu solicitud de retiro fue enviada. Recibiras el dinero en tu cuenta de MercadoPago en 1-2 dias habiles.',
                [{ text: 'OK', onPress: () => navigation.goBack() }]
              )
            } catch (err) {
              setError(getErrorMessage(err))
            } finally {
              setSubmitting(false)
            }
          },
        },
      ]
    )
  }

  if (loading) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color="#0ea5e9" />
        <Text style={styles.loadingText}>Cargando...</Text>
      </View>
    )
  }

  const availablePesos = (wallet?.available_balance || 0) / 100
  const amountPesos = parseFloat(amount) || 0

  return (
    <View style={styles.container}>
      <ScrollView>
        {/* Balance Info */}
        <Card style={styles.balanceCard}>
          <Text style={styles.balanceLabel}>Saldo disponible para retirar</Text>
          <Text style={styles.balanceAmount}>
            {wallet?.available_balance_formatted || formatMoney(0)}
          </Text>
        </Card>

        {/* Amount Input */}
        <Card style={styles.amountCard}>
          <Text style={styles.inputLabel}>Monto a retirar</Text>
          <View style={styles.amountInputContainer}>
            <Text style={styles.currencySymbol}>$</Text>
            <TextInput
              style={styles.amountInput}
              value={amount}
              onChangeText={handleAmountChange}
              placeholder="0.00"
              keyboardType="decimal-pad"
              editable={!submitting}
            />
          </View>

          <Button
            title="Retirar todo"
            variant="outline"
            onPress={handleWithdrawAll}
            disabled={submitting || availablePesos <= 0}
            style={styles.withdrawAllButton}
          />
        </Card>

        {/* Summary */}
        <Card style={styles.summaryCard}>
          <Text style={styles.summaryTitle}>Resumen</Text>

          <View style={styles.summaryRow}>
            <Text style={styles.summaryLabel}>Monto a retirar</Text>
            <Text style={styles.summaryValue}>{formatMoney(amountPesos * 100)}</Text>
          </View>

          <View style={styles.summaryRow}>
            <Text style={styles.summaryLabel}>Comision</Text>
            <Text style={styles.summaryValue}>$0.00</Text>
          </View>

          <View style={styles.summaryDivider} />

          <View style={styles.summaryRow}>
            <Text style={styles.summaryLabelBold}>Total a recibir</Text>
            <Text style={styles.summaryValueBold}>{formatMoney(amountPesos * 100)}</Text>
          </View>
        </Card>

        {/* Info Card */}
        <Card style={styles.infoCard}>
          <View style={styles.infoItem}>
            <Ionicons name="time-outline" size={20} color="#6b7280" />
            <Text style={styles.infoText}>
              El dinero estara disponible en tu cuenta de MercadoPago en 1-2 dias habiles.
            </Text>
          </View>
          <View style={styles.infoItem}>
            <Ionicons name="information-circle-outline" size={20} color="#6b7280" />
            <Text style={styles.infoText}>El monto minimo de retiro es $10.</Text>
          </View>
        </Card>

        {/* Error Message */}
        {error && (
          <View style={styles.errorContainer}>
            <Ionicons name="alert-circle" size={20} color="#dc2626" />
            <Text style={styles.errorText}>{error}</Text>
          </View>
        )}

        {/* Submit Button */}
        <View style={styles.submitContainer}>
          <Button
            title={submitting ? 'Procesando...' : 'Solicitar retiro'}
            onPress={handleSubmit}
            loading={submitting}
            disabled={submitting || amountPesos <= 0 || amountPesos > availablePesos}
            style={styles.submitButton}
          />
        </View>
      </ScrollView>
    </View>
  )
}
