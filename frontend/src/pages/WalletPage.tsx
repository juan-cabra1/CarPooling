import { useState, useEffect } from 'react'
import { Wallet, ArrowDownCircle, TrendingUp, Clock, AlertCircle, ExternalLink, Settings } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { paymentsService, getErrorMessage } from '@/services'
import type {
  WalletResponse,
  TransactionResponse,
  WalletStatsResponse,
  SellerStatusResponse,
} from '@/types/payment'

export default function WalletPage() {
  const [wallet, setWallet] = useState<WalletResponse | null>(null)
  const [transactions, setTransactions] = useState<TransactionResponse[]>([])
  const [stats, setStats] = useState<WalletStatsResponse | null>(null)
  const [seller, setSeller] = useState<SellerStatusResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  // Withdrawal modal state
  const [showWithdraw, setShowWithdraw] = useState(false)
  const [withdrawAmount, setWithdrawAmount] = useState('')
  const [withdrawing, setWithdrawing] = useState(false)
  const [withdrawError, setWithdrawError] = useState('')
  const [withdrawSuccess, setWithdrawSuccess] = useState(false)

  // Link MP modal state
  const [showLinkMP, setShowLinkMP] = useState(false)
  const [mpToken, setMpToken] = useState('')
  const [linking, setLinking] = useState(false)
  const [linkError, setLinkError] = useState('')

  useEffect(() => {
    loadAll()
  }, [])

  const loadAll = async () => {
    try {
      setLoading(true)
      setError('')
      const [w, tx, s, sel] = await Promise.all([
        paymentsService.getWallet(),
        paymentsService.getTransactions(1, 10),
        paymentsService.getWalletStats(),
        paymentsService.getSellerStatus(),
      ])
      setWallet(w)
      setTransactions(tx.transactions)
      setStats(s)
      setSeller(sel)
    } catch (err) {
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
    }
  }

  const handleWithdraw = async (e: React.FormEvent) => {
    e.preventDefault()
    const amountPesos = parseFloat(withdrawAmount)
    if (isNaN(amountPesos) || amountPesos <= 0) {
      setWithdrawError('Ingresá un monto válido')
      return
    }
    if (!wallet || amountPesos * 100 > wallet.available_balance) {
      setWithdrawError('Saldo insuficiente')
      return
    }
    try {
      setWithdrawing(true)
      setWithdrawError('')
      await paymentsService.requestWithdrawal(Math.round(amountPesos * 100))
      setWithdrawSuccess(true)
      setWithdrawAmount('')
      await loadAll()
      setTimeout(() => {
        setShowWithdraw(false)
        setWithdrawSuccess(false)
      }, 2000)
    } catch (err) {
      setWithdrawError(getErrorMessage(err))
    } finally {
      setWithdrawing(false)
    }
  }

  const handleLinkManual = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      setLinking(true)
      setLinkError('')
      await paymentsService.linkManualToken(mpToken)
      setMpToken('')
      setShowLinkMP(false)
      await loadAll()
    } catch (err) {
      setLinkError(getErrorMessage(err))
    } finally {
      setLinking(false)
    }
  }

  const handleLinkOAuth = async () => {
    try {
      const url = await paymentsService.getOAuthURL()
      window.location.href = url
    } catch (err) {
      setLinkError(getErrorMessage(err))
    }
  }

  const formatARS = (centavos: number) =>
    new Intl.NumberFormat('es-AR', {
      style: 'currency',
      currency: 'ARS',
      minimumFractionDigits: 0,
    }).format(centavos / 100)

  const txTypeLabel: Record<string, string> = {
    credit: 'Crédito',
    debit: 'Débito',
    withdrawal: 'Retiro',
    refund: 'Reembolso',
  }

  const txTypeVariant: Record<string, 'default' | 'secondary' | 'destructive' | 'outline'> = {
    credit: 'default',
    debit: 'destructive',
    withdrawal: 'secondary',
    refund: 'outline',
  }

  if (loading) {
    return (
      <div className="min-h-[calc(100vh-4rem)] flex items-center justify-center">
        <div className="text-center">
          <div className="w-16 h-16 border-4 border-primary border-t-transparent rounded-full animate-spin mx-auto mb-4" />
          <p className="text-gray-600">Cargando billetera...</p>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="min-h-[calc(100vh-4rem)] flex items-center justify-center p-4">
        <Card className="max-w-md w-full">
          <CardContent className="pt-6 text-center">
            <AlertCircle className="w-12 h-12 text-destructive mx-auto mb-4" />
            <p className="text-gray-700">{error}</p>
            <Button onClick={loadAll} className="mt-4">Reintentar</Button>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="min-h-[calc(100vh-4rem)] bg-gradient-to-br from-primary-50 via-white to-secondary-50 py-8">
      <div className="container mx-auto px-4 max-w-4xl">
        <h1 className="text-3xl font-bold text-gray-900 mb-6 flex items-center gap-2">
          <Wallet className="w-8 h-8 text-primary" />
          Mi Billetera
        </h1>

        {/* ── MercadoPago account banner ── */}
        {seller && !seller.linked && (
          <div className="mb-6 p-4 bg-yellow-50 border border-yellow-200 rounded-lg flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3">
            <div>
              <p className="font-semibold text-yellow-800">Vinculá tu cuenta de MercadoPago</p>
              <p className="text-sm text-yellow-700">Para recibir pagos de los pasajeros necesitás vincular tu cuenta.</p>
            </div>
            <div className="flex gap-2 shrink-0">
              <Button size="sm" onClick={handleLinkOAuth} className="gap-1">
                <ExternalLink className="w-3 h-3" />
                Vincular con OAuth
              </Button>
              <Button size="sm" variant="outline" onClick={() => setShowLinkMP(true)}>
                Token manual
              </Button>
            </div>
          </div>
        )}
        {seller?.linked && (
          <div className="mb-6 p-3 bg-green-50 border border-green-200 rounded-lg flex items-center justify-between">
            <p className="text-sm text-green-800">
              MercadoPago vinculado: <span className="font-medium">{seller.mp_email}</span>
            </p>
            <Button
              size="sm"
              variant="ghost"
              className="text-red-600 hover:text-red-700 hover:bg-red-50"
              onClick={async () => {
                await paymentsService.disconnectSeller()
                await loadAll()
              }}
            >
              Desvincular
            </Button>
          </div>
        )}

        {/* ── Balance cards ── */}
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-6">
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>Disponible</CardDescription>
              <CardTitle className="text-2xl text-green-600">
                {wallet ? formatARS(wallet.available_balance) : '-'}
              </CardTitle>
            </CardHeader>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>En proceso</CardDescription>
              <CardTitle className="text-2xl text-yellow-600">
                {wallet ? formatARS(wallet.pending_balance) : '-'}
              </CardTitle>
            </CardHeader>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>Total ganado</CardDescription>
              <CardTitle className="text-2xl text-primary">
                {wallet ? wallet.total_earned_formatted || formatARS(wallet.total_earned) : '-'}
              </CardTitle>
            </CardHeader>
          </Card>
        </div>

        {/* ── Stats row ── */}
        {stats && (
          <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 mb-6">
            <Card>
              <CardContent className="pt-4 text-center">
                <p className="text-2xl font-bold text-gray-900">{stats.current_month_trips}</p>
                <p className="text-xs text-gray-500">viajes este mes</p>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="pt-4 text-center">
                <p className="text-2xl font-bold text-gray-900">{formatARS(stats.current_month_earnings)}</p>
                <p className="text-xs text-gray-500">ganado este mes</p>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="pt-4 text-center">
                <p className={`text-2xl font-bold ${stats.earnings_growth_percent >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                  {stats.earnings_growth_percent >= 0 ? '+' : ''}{stats.earnings_growth_percent.toFixed(1)}%
                </p>
                <p className="text-xs text-gray-500">vs mes anterior</p>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="pt-4 text-center">
                <p className="text-2xl font-bold text-gray-900">{formatARS(stats.average_per_trip)}</p>
                <p className="text-xs text-gray-500">promedio por viaje</p>
              </CardContent>
            </Card>
          </div>
        )}

        {/* ── Actions ── */}
        <div className="flex gap-3 mb-6">
          <Button
            onClick={() => setShowWithdraw(true)}
            disabled={!wallet || wallet.available_balance === 0}
            className="gap-2"
          >
            <ArrowDownCircle className="w-4 h-4" />
            Retirar fondos
          </Button>
          <Button variant="outline" className="gap-2" onClick={() => {}}>
            <Settings className="w-4 h-4" />
            Configuración
          </Button>
        </div>

        {/* ── Transactions ── */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <TrendingUp className="w-5 h-5" />
              Últimas transacciones
            </CardTitle>
          </CardHeader>
          <CardContent>
            {transactions.length === 0 ? (
              <div className="text-center py-8 text-gray-500">
                <Clock className="w-12 h-12 mx-auto mb-3 opacity-50" />
                <p>No hay transacciones aún.</p>
              </div>
            ) : (
              <div className="space-y-3">
                {transactions.map((tx) => (
                  <div key={tx.transaction_uuid} className="flex items-center justify-between py-3 border-b last:border-0">
                    <div className="flex items-center gap-3">
                      <Badge variant={txTypeVariant[tx.type] || 'outline'} className="text-xs">
                        {txTypeLabel[tx.type] || tx.type}
                      </Badge>
                      <div>
                        <p className="text-sm font-medium text-gray-900">{tx.description}</p>
                        <p className="text-xs text-gray-500">
                          {new Date(tx.created_at).toLocaleDateString('es-AR', {
                            day: 'numeric',
                            month: 'short',
                            year: 'numeric',
                            hour: '2-digit',
                            minute: '2-digit',
                          })}
                        </p>
                      </div>
                    </div>
                    <span className={`font-semibold ${tx.signed_amount >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                      {tx.signed_amount >= 0 ? '+' : ''}{tx.amount_formatted || formatARS(Math.abs(tx.amount))}
                    </span>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>

        {/* ── Withdrawal modal ── */}
        {showWithdraw && (
          <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black bg-opacity-50">
            <div className="bg-white rounded-lg shadow-xl max-w-sm w-full p-6">
              <h3 className="text-xl font-bold text-gray-900 mb-4">Retirar fondos</h3>
              {withdrawSuccess ? (
                <div className="text-center py-4">
                  <p className="text-green-600 font-medium">Retiro solicitado exitosamente</p>
                </div>
              ) : (
                <form onSubmit={handleWithdraw} className="space-y-4">
                  <div>
                    <Label htmlFor="withdraw-amount">Monto a retirar (ARS)</Label>
                    <Input
                      id="withdraw-amount"
                      type="number"
                      min="1"
                      step="1"
                      value={withdrawAmount}
                      onChange={(e) => setWithdrawAmount(e.target.value)}
                      placeholder="Ej: 5000"
                      disabled={withdrawing}
                    />
                    {wallet && (
                      <p className="text-xs text-gray-500 mt-1">
                        Disponible: {formatARS(wallet.available_balance)}
                      </p>
                    )}
                  </div>
                  {withdrawError && (
                    <p className="text-sm text-red-600">{withdrawError}</p>
                  )}
                  <div className="flex gap-3">
                    <Button
                      type="button"
                      variant="outline"
                      onClick={() => { setShowWithdraw(false); setWithdrawAmount(''); setWithdrawError('') }}
                      disabled={withdrawing}
                      className="flex-1"
                    >
                      Cancelar
                    </Button>
                    <Button type="submit" disabled={withdrawing} className="flex-1">
                      {withdrawing ? 'Procesando...' : 'Retirar'}
                    </Button>
                  </div>
                </form>
              )}
            </div>
          </div>
        )}

        {/* ── Link MP manual token modal ── */}
        {showLinkMP && (
          <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black bg-opacity-50">
            <div className="bg-white rounded-lg shadow-xl max-w-sm w-full p-6">
              <h3 className="text-xl font-bold text-gray-900 mb-2">Vincular con access token</h3>
              <p className="text-sm text-gray-600 mb-4">
                Pegá tu Access Token de producción de MercadoPago.
              </p>
              <form onSubmit={handleLinkManual} className="space-y-4">
                <div>
                  <Label htmlFor="mp-token">Access Token</Label>
                  <Input
                    id="mp-token"
                    type="password"
                    value={mpToken}
                    onChange={(e) => setMpToken(e.target.value)}
                    placeholder="APP_USR-..."
                    disabled={linking}
                  />
                </div>
                {linkError && <p className="text-sm text-red-600">{linkError}</p>}
                <div className="flex gap-3">
                  <Button
                    type="button"
                    variant="outline"
                    onClick={() => { setShowLinkMP(false); setMpToken(''); setLinkError('') }}
                    disabled={linking}
                    className="flex-1"
                  >
                    Cancelar
                  </Button>
                  <Button type="submit" disabled={linking || !mpToken} className="flex-1">
                    {linking ? 'Vinculando...' : 'Vincular'}
                  </Button>
                </div>
              </form>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
