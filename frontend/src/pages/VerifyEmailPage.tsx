import { useEffect, useState } from 'react'
import { Link, useSearchParams, useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { authService } from '@/services'
import { getErrorMessage } from '@/services'
import { CheckCircle, XCircle, Loader2, Mail } from 'lucide-react'

export default function VerifyEmailPage() {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const token = searchParams.get('token')

  const [status, setStatus] = useState<'loading' | 'success' | 'error'>('loading')
  const [error, setError] = useState('')

  useEffect(() => {
    const verifyEmail = async () => {
      if (!token) {
        setStatus('error')
        setError('Token de verificación no encontrado en la URL')
        return
      }

      try {
        await authService.verifyEmail(token)
        setStatus('success')

        // Redirigir al login después de 3 segundos
        setTimeout(() => {
          navigate('/login', {
            state: {
              message: 'Email verificado exitosamente. Ahora puedes iniciar sesión.'
            }
          })
        }, 3000)
      } catch (err) {
        setStatus('error')
        setError(getErrorMessage(err))
      }
    }

    verifyEmail()
  }, [token, navigate])

  if (status === 'loading') {
    return (
      <div className="min-h-[calc(100vh-4rem)] flex items-center justify-center bg-gradient-to-br from-primary-50 via-white to-secondary-50 p-4">
        <Card className="w-full max-w-md shadow-xl">
          <CardHeader className="space-y-1">
            <div className="flex items-center justify-center mb-4">
              <div className="w-16 h-16 bg-primary rounded-full flex items-center justify-center">
                <Loader2 className="w-8 h-8 text-white animate-spin" />
              </div>
            </div>
            <CardTitle className="text-2xl font-bold text-center">
              Verificando Email
            </CardTitle>
            <CardDescription className="text-center">
              Por favor espera mientras verificamos tu correo electrónico...
            </CardDescription>
          </CardHeader>
        </Card>
      </div>
    )
  }

  if (status === 'success') {
    return (
      <div className="min-h-[calc(100vh-4rem)] flex items-center justify-center bg-gradient-to-br from-primary-50 via-white to-secondary-50 p-4">
        <Card className="w-full max-w-md shadow-xl">
          <CardHeader className="space-y-1">
            <div className="flex items-center justify-center mb-4">
              <div className="w-16 h-16 bg-green-500 rounded-full flex items-center justify-center">
                <CheckCircle className="w-8 h-8 text-white" />
              </div>
            </div>
            <CardTitle className="text-2xl font-bold text-center">
              ¡Email Verificado!
            </CardTitle>
            <CardDescription className="text-center">
              Tu correo electrónico ha sido verificado exitosamente
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="p-4 rounded-lg bg-green-50 border border-green-200">
              <p className="text-sm text-green-900 text-center">
                Ahora puedes iniciar sesión con tu cuenta
              </p>
            </div>
            <p className="text-center text-sm text-muted-foreground">
              Serás redirigido a la página de inicio de sesión en unos segundos...
            </p>
          </CardContent>
          <CardFooter>
            <Link to="/login" className="w-full">
              <Button className="w-full" size="lg">
                Ir a Iniciar Sesión
              </Button>
            </Link>
          </CardFooter>
        </Card>
      </div>
    )
  }

  // status === 'error'
  return (
    <div className="min-h-[calc(100vh-4rem)] flex items-center justify-center bg-gradient-to-br from-primary-50 via-white to-secondary-50 p-4">
      <Card className="w-full max-w-md shadow-xl">
        <CardHeader className="space-y-1">
          <div className="flex items-center justify-center mb-4">
            <div className="w-16 h-16 bg-red-500 rounded-full flex items-center justify-center">
              <XCircle className="w-8 h-8 text-white" />
            </div>
          </div>
          <CardTitle className="text-2xl font-bold text-center">
            Error de Verificación
          </CardTitle>
          <CardDescription className="text-center">
            No pudimos verificar tu correo electrónico
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="p-4 rounded-lg bg-red-50 border border-red-200">
            <p className="text-sm text-red-900 text-center">
              {error || 'El token de verificación es inválido o ha expirado'}
            </p>
          </div>
          <div className="p-4 rounded-lg bg-blue-50 border border-blue-200">
            <div className="flex items-start gap-2">
              <Mail className="w-5 h-5 text-blue-600 flex-shrink-0 mt-0.5" />
              <p className="text-sm text-blue-900">
                Si el enlace expiró, puedes solicitar un nuevo correo de verificación
              </p>
            </div>
          </div>
        </CardContent>
        <CardFooter className="flex flex-col gap-2">
          <Link to="/resend-verification" className="w-full">
            <Button variant="default" className="w-full" size="lg">
              <Mail className="w-4 h-4 mr-2" />
              Reenviar Email de Verificación
            </Button>
          </Link>
          <Link to="/login" className="w-full">
            <Button variant="outline" className="w-full" size="lg">
              Volver al Inicio de Sesión
            </Button>
          </Link>
        </CardFooter>
      </Card>
    </div>
  )
}
