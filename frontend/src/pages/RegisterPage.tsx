import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { PasswordInput } from '@/components/ui/password-input'
import { Label } from '@/components/ui/label'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { getErrorMessage } from '@/services'
import { authService } from '@/services'
import type { RegisterData, UserSex } from '@/types'
import { UserPlus, Mail, Lock, User, Phone, MapPin, Calendar, AlertCircle, CheckCircle } from 'lucide-react'

export default function RegisterPage() {
  const navigate = useNavigate()
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [success, setSuccess] = useState(false)

  const [formData, setFormData] = useState<RegisterData>({
    email: '',
    password: '',
    name: '',
    lastname: '',
    phone: '',
    street: '',
    number: 0,
    sex: 'otro',
    birthdate: '',
  })

  const [confirmPassword, setConfirmPassword] = useState('')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')

    // Validaciones del lado del cliente
    if (formData.password !== confirmPassword) {
      setError('Las contraseñas no coinciden')
      return
    }

    if (formData.password.length < 6) {
      setError('La contraseña debe tener al menos 6 caracteres')
      return
    }

    // Validar edad mínima (18 años)
    const birthDate = new Date(formData.birthdate)
    const today = new Date()
    const age = today.getFullYear() - birthDate.getFullYear()
    const monthDiff = today.getMonth() - birthDate.getMonth()
    const dayDiff = today.getDate() - birthDate.getDate()

    const actualAge = monthDiff < 0 || (monthDiff === 0 && dayDiff < 0) ? age - 1 : age

    if (actualAge < 18) {
      setError('Debes tener al menos 18 años para registrarte')
      return
    }

    setLoading(true)

    try {
      await authService.register(formData)
      setSuccess(true)

      // Redirigir después de 3 segundos
      setTimeout(() => {
        navigate('/login', {
          state: {
            message: 'Registro exitoso. Por favor, verifica tu correo electrónico antes de iniciar sesión.'
          }
        })
      }, 3000)
    } catch (err) {
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
    }
  }

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement>
  ) => {
    const { name, value, type } = e.target

    setFormData({
      ...formData,
      [name]: type === 'number' ? parseInt(value) || 0 : value,
    })
  }

  const handleSexChange = (value: UserSex) => {
    setFormData({
      ...formData,
      sex: value,
    })
  }

  if (success) {
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
              ¡Registro Exitoso!
            </CardTitle>
            <CardDescription className="text-center">
              Te hemos enviado un correo de verificación a <strong>{formData.email}</strong>
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="p-4 rounded-lg bg-blue-50 border border-blue-200">
              <p className="text-sm text-blue-900">
                Por favor, verifica tu correo electrónico antes de iniciar sesión.
                Si no recibes el correo en unos minutos, revisa tu carpeta de spam.
              </p>
            </div>
            <p className="text-center text-sm text-muted-foreground">
              Serás redirigido a la página de inicio de sesión en unos segundos...
            </p>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="min-h-[calc(100vh-4rem)] flex items-center justify-center bg-gradient-to-br from-primary-50 via-white to-secondary-50 p-4 py-12">
      <Card className="w-full max-w-2xl shadow-xl">
        <CardHeader className="space-y-1">
          <div className="flex items-center justify-center mb-4">
            <div className="w-16 h-16 bg-primary rounded-full flex items-center justify-center">
              <UserPlus className="w-8 h-8 text-white" />
            </div>
          </div>
          <CardTitle className="text-2xl font-bold text-center">
            Crear Cuenta
          </CardTitle>
          <CardDescription className="text-center">
            Completa el formulario para unirte a CarPooling
          </CardDescription>
        </CardHeader>

        <form onSubmit={handleSubmit}>
          <CardContent className="space-y-4">
            {error && (
              <div className="p-3 rounded-lg bg-destructive/10 border border-destructive/20 flex items-start gap-2">
                <AlertCircle className="w-5 h-5 text-destructive flex-shrink-0 mt-0.5" />
                <p className="text-sm text-destructive">{error}</p>
              </div>
            )}

            {/* Información de la cuenta */}
            <div className="space-y-4 p-4 rounded-lg bg-muted/50">
              <h3 className="font-semibold text-sm text-muted-foreground">
                Información de la Cuenta
              </h3>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-2 md:col-span-2">
                  <Label htmlFor="email">
                    <Mail className="w-4 h-4 inline mr-2" />
                    Correo electrónico
                  </Label>
                  <Input
                    id="email"
                    name="email"
                    type="email"
                    placeholder="tu@email.com"
                    value={formData.email}
                    onChange={handleChange}
                    required
                    disabled={loading}
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="password">
                    <Lock className="w-4 h-4 inline mr-2" />
                    Contraseña
                  </Label>
                  <PasswordInput
                    id="password"
                    name="password"
                    placeholder="Mínimo 6 caracteres"
                    value={formData.password}
                    onChange={handleChange}
                    required
                    disabled={loading}
                    minLength={6}
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="confirmPassword">
                    <Lock className="w-4 h-4 inline mr-2" />
                    Confirmar Contraseña
                  </Label>
                  <PasswordInput
                    id="confirmPassword"
                    name="confirmPassword"
                    placeholder="Repite tu contraseña"
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    required
                    disabled={loading}
                    minLength={6}
                  />
                </div>
              </div>
            </div>

            {/* Información personal */}
            <div className="space-y-4 p-4 rounded-lg bg-muted/50">
              <h3 className="font-semibold text-sm text-muted-foreground">
                Información Personal
              </h3>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="name">
                    <User className="w-4 h-4 inline mr-2" />
                    Nombre
                  </Label>
                  <Input
                    id="name"
                    name="name"
                    type="text"
                    placeholder="Juan"
                    value={formData.name}
                    onChange={handleChange}
                    required
                    disabled={loading}
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="lastname">
                    <User className="w-4 h-4 inline mr-2" />
                    Apellido
                  </Label>
                  <Input
                    id="lastname"
                    name="lastname"
                    type="text"
                    placeholder="Pérez"
                    value={formData.lastname}
                    onChange={handleChange}
                    required
                    disabled={loading}
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="sex">
                    <User className="w-4 h-4 inline mr-2" />
                    Sexo
                  </Label>
                  <Select
                    value={formData.sex}
                    onValueChange={handleSexChange}
                    disabled={loading}
                  >
                    <SelectTrigger id="sex">
                      <SelectValue placeholder="Selecciona tu sexo" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="hombre">Hombre</SelectItem>
                      <SelectItem value="mujer">Mujer</SelectItem>
                      <SelectItem value="otro">Otro</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="birthdate">
                    <Calendar className="w-4 h-4 inline mr-2" />
                    Fecha de Nacimiento
                  </Label>
                  <Input
                    id="birthdate"
                    name="birthdate"
                    type="date"
                    value={formData.birthdate}
                    onChange={handleChange}
                    required
                    disabled={loading}
                    max={new Date().toISOString().split('T')[0]}
                  />
                </div>

                <div className="space-y-2 md:col-span-2">
                  <Label htmlFor="phone">
                    <Phone className="w-4 h-4 inline mr-2" />
                    Teléfono
                  </Label>
                  <Input
                    id="phone"
                    name="phone"
                    type="tel"
                    placeholder="+54 9 11 1234-5678"
                    value={formData.phone}
                    onChange={handleChange}
                    required
                    disabled={loading}
                  />
                </div>
              </div>
            </div>

            {/* Dirección */}
            <div className="space-y-4 p-4 rounded-lg bg-muted/50">
              <h3 className="font-semibold text-sm text-muted-foreground">
                Dirección
              </h3>

              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div className="space-y-2 md:col-span-2">
                  <Label htmlFor="street">
                    <MapPin className="w-4 h-4 inline mr-2" />
                    Calle
                  </Label>
                  <Input
                    id="street"
                    name="street"
                    type="text"
                    placeholder="Av. Corrientes"
                    value={formData.street}
                    onChange={handleChange}
                    required
                    disabled={loading}
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="number">
                    <MapPin className="w-4 h-4 inline mr-2" />
                    Número
                  </Label>
                  <Input
                    id="number"
                    name="number"
                    type="number"
                    placeholder="1234"
                    value={formData.number || ''}
                    onChange={handleChange}
                    required
                    disabled={loading}
                    min={1}
                  />
                </div>
              </div>
            </div>
          </CardContent>

          <CardFooter className="flex flex-col space-y-4">
            <Button
              type="submit"
              className="w-full"
              size="lg"
              disabled={loading}
            >
              {loading ? (
                <>
                  <div className="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin mr-2" />
                  Creando cuenta...
                </>
              ) : (
                <>
                  <UserPlus className="w-5 h-5 mr-2" />
                  Crear Cuenta
                </>
              )}
            </Button>

            <div className="text-center text-sm text-muted-foreground">
              ¿Ya tienes cuenta?{' '}
              <Link
                to="/login"
                className="text-primary font-semibold hover:text-primary-600 hover:underline"
              >
                Inicia sesión aquí
              </Link>
            </div>
          </CardFooter>
        </form>
      </Card>
    </div>
  )
}
