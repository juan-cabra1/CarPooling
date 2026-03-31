# CarPooling Mobile App

Aplicación móvil nativa de CarPooling construida con React Native y Expo, diseñada para igualar completamente el diseño y funcionalidad de la versión web.

## Características

- **Autenticación completa**: Login, registro, recuperación de contraseña y verificación de email
- **Búsqueda de viajes**: Sistema de búsqueda con autocompletado de ubicaciones usando Photon API (OpenStreetMap)
- **Gestión de viajes**: Crear, editar y publicar viajes como conductor
- **Reservas**: Sistema completo de reservas para pasajeros
- **Perfil de usuario**: Gestión de información personal y configuración
- **Diseño responsive**: Interfaz adaptada para dispositivos móviles
- **Sistema de colores consistente**: Coincide exactamente con la versión web

## Requisitos Previos

- Node.js >= 18.x
- npm o yarn
- Expo CLI
- Para desarrollo:
  - Expo Go app en tu dispositivo móvil (iOS/Android)
  - O un emulador de Android/iOS

## Instalación

1. **Navegar al directorio del proyecto:**
   ```bash
   cd frontendmobile
   ```

2. **Instalar dependencias:**
   ```bash
   npm install
   ```

3. **Configurar variables de entorno:**

   Las APIs se conectan a:
   - Users API: `http://localhost:8080/api`
   - Trips API: `http://localhost:8081/api`
   - Bookings API: `http://localhost:8082/api`
   - Search API: `http://localhost:8083/api`

   Asegúrate de que los servicios backend estén ejecutándose.

## Ejecución

### Desarrollo

```bash
# Iniciar el servidor de desarrollo de Expo
npm start

# O con Expo CLI
npx expo start
```

Esto abrirá Expo DevTools en tu navegador. Desde aquí puedes:
- Escanear el código QR con la app Expo Go (iOS/Android)
- Presionar `a` para abrir en emulador de Android
- Presionar `i` para abrir en simulador de iOS
- Presionar `w` para abrir en navegador web (limitado)

### Plataformas Específicas

```bash
# Android
npm run android
# o
npx expo start --android

# iOS (solo en macOS)
npm run ios
# o
npx expo start --ios

# Web
npm run web
# o
npx expo start --web
```

## Estructura del Proyecto

```
frontendmobile/
├── App.tsx                      # Punto de entrada principal
├── index.ts                     # Punto de entrada de Expo
├── src/
│   ├── components/              # Componentes reutilizables
│   │   ├── ui/                  # Componentes de interfaz base
│   │   │   ├── Button/          # Botón con variantes
│   │   │   ├── Input/           # Input de texto con estados
│   │   │   ├── Card/            # Contenedor de tarjeta
│   │   │   ├── Badge/           # Indicador de estado
│   │   │   ├── TextArea/        # Área de texto multilinea
│   │   │   └── Switch/          # Interruptor toggle
│   │   ├── search/              # Componentes de búsqueda
│   │   │   └── LocationInput/   # Input con autocompletado
│   │   ├── trips/               # Componentes de viajes
│   │   │   └── TripCard/        # Tarjeta de viaje
│   │   └── modals/              # Componentes modales
│   │       └── BookingModal/    # Modal de reserva
│   ├── screens/                 # Pantallas de la aplicación
│   │   ├── LoginScreen/         # Pantalla de inicio de sesión
│   │   ├── RegisterScreen/      # Pantalla de registro
│   │   ├── ForgotPasswordScreen/ # Recuperar contraseña
│   │   ├── HomeScreen/          # Pantalla principal
│   │   ├── SearchScreen/        # Búsqueda de viajes
│   │   ├── TripDetailScreen/    # Detalles del viaje
│   │   ├── CreateTripScreen/    # Crear nuevo viaje
│   │   ├── EditTripScreen/      # Editar viaje existente
│   │   ├── MyTripsScreen/       # Mis viajes como conductor
│   │   ├── MyBookingsScreen/    # Mis reservas como pasajero
│   │   └── ProfileScreen/       # Perfil de usuario
│   ├── navigation/              # Configuración de navegación
│   │   ├── AppNavigator.tsx     # Navegador principal
│   │   └── types.ts             # Tipos de navegación
│   ├── context/                 # Contextos de React
│   │   └── AuthContext.tsx      # Contexto de autenticación
│   ├── services/                # Servicios y API
│   │   └── api.ts               # Cliente de API con Axios
│   ├── styles/                  # Estilos globales
│   │   └── colors.ts            # Paleta de colores
│   └── types/                   # Definiciones de TypeScript
│       └── index.ts             # Tipos de datos
├── assets/                      # Recursos estáticos
├── package.json                 # Dependencias del proyecto
├── tsconfig.json                # Configuración de TypeScript
├── babel.config.js              # Configuración de Babel
└── README.md                    # Este archivo
```

## Arquitectura de Estilos

Este proyecto utiliza **StyleSheet nativo de React Native** para todos los estilos:

- **Separación de estilos y lógica**: Cada componente tiene su archivo `.styles.ts` separado
- **Sistema de colores centralizado**: Todos los colores definidos en `src/styles/colors.ts`
- **Coincidencia exacta con web**: Los colores y diseño replican exactamente la versión web
- **Type-safe**: StyleSheet.create() proporciona validación y autocompletado

### Ejemplo de Uso de Estilos

```typescript
// Button.styles.ts
import { StyleSheet } from 'react-native'
import { colors } from '@/styles/colors'

export const styles = StyleSheet.create({
  primary: {
    backgroundColor: colors.primary[500],
    paddingHorizontal: 24,
    paddingVertical: 12,
    borderRadius: 8,
  },
})

// Button.tsx
import { styles } from './Button.styles'

<TouchableOpacity style={styles.primary}>
  <Text>Click me</Text>
</TouchableOpacity>
```

## Paleta de Colores

La aplicación utiliza una paleta de colores consistente definida en `src/styles/colors.ts`:

- **Primary (Azul)**: `#0ea5e9` - Acciones principales y elementos destacados
- **Secondary (Teal)**: `#14b8a6` - Acciones secundarias
- **Destructive (Rojo)**: `#dc2626` - Acciones destructivas y errores
- **Success (Verde)**: `#22c55e` - Estados exitosos
- **Warning (Amarillo)**: `#eab308` - Advertencias
- **Muted**: `#f1f5f9` - Fondos y elementos deshabilitados
- **Foreground**: `#1e293b` - Texto principal
- **Border**: `#cbd5e1` - Bordes y divisores

## Navegación

La aplicación usa React Navigation v7 con:

- **Stack Navigator**: Para navegación principal con transiciones de pantalla
- **Tab Navigator**: Para la navegación inferior (Home, Search, My Trips, Profile)
- **Auth Flow**: Navegación condicional basada en el estado de autenticación

## Gestión de Estado

- **AuthContext**: Maneja el estado de autenticación global
- **Local State**: useState para estado de componente local
- **AsyncStorage**: Persistencia de tokens y datos de usuario

## APIs y Servicios

El servicio de API (`src/services/api.ts`) proporciona:

- **Interceptores de autenticación**: Añade tokens JWT automáticamente
- **Manejo de errores**: Gestión centralizada de errores de red
- **Redirección automática**: Redirige al login en caso de tokens expirados
- **Múltiples servicios**: Conecta con Users, Trips, Bookings y Search APIs

### Endpoints Disponibles

#### Autenticación
- `POST /auth/register` - Registro de usuario
- `POST /auth/login` - Inicio de sesión
- `POST /auth/logout` - Cierre de sesión
- `POST /auth/forgot-password` - Solicitar recuperación de contraseña
- `POST /auth/reset-password` - Restablecer contraseña
- `POST /auth/verify-email` - Verificar email

#### Viajes
- `GET /trips` - Listar viajes
- `GET /trips/:id` - Obtener detalles de viaje
- `POST /trips` - Crear viaje
- `PUT /trips/:id` - Actualizar viaje
- `DELETE /trips/:id` - Eliminar viaje
- `GET /trips/user` - Viajes del usuario actual

#### Reservas
- `GET /bookings` - Listar reservas del usuario
- `POST /bookings` - Crear reserva
- `PUT /bookings/:id` - Actualizar reserva
- `DELETE /bookings/:id` - Cancelar reserva

#### Búsqueda
- `GET /search` - Buscar viajes con filtros
- Photon API para autocompletado de ubicaciones

## Componentes Principales

### Componentes UI Base

#### Button
Botón reutilizable con múltiples variantes:
- `primary` - Botón principal azul
- `secondary` - Botón secundario teal
- `outline` - Botón con borde
- `destructive` - Botón rojo para acciones destructivas
- `ghost` - Botón transparente
- `link` - Botón tipo enlace

#### Input
Input de texto con soporte para:
- Labels
- Mensajes de error
- Texto de ayuda
- Estados de focus
- Validación visual

#### Card
Contenedor de tarjeta con secciones:
- CardHeader
- CardTitle
- CardDescription
- CardContent
- CardFooter

### Componentes Especializados

#### LocationInput
Input con autocompletado usando Photon API (OpenStreetMap) para búsqueda de ubicaciones.

#### TripCard
Tarjeta que muestra información del viaje:
- Fecha y hora
- Origen y destino
- Conductor
- Precio
- Asientos disponibles
- Estado

#### BookingModal
Modal para crear reservas con:
- Selección de asientos
- Confirmación de detalles
- Validación

## TypeScript

El proyecto está completamente tipado con TypeScript:

- **Strict mode** habilitado
- **Tipos de navegación** definidos en `src/navigation/types.ts`
- **Tipos de datos** definidos en `src/types/index.ts`
- **Props de componentes** completamente tipados
- **Respuestas de API** tipadas

## Testing

Para ejecutar la aplicación en modo de prueba:

1. Asegúrate de que los servicios backend estén corriendo
2. Ejecuta `npm start`
3. Prueba las siguientes funcionalidades:
   - Registro y login de usuario
   - Búsqueda de viajes
   - Creación de viajes
   - Realización de reservas
   - Edición de perfil

## Solución de Problemas

### El bundler no inicia
```bash
# Limpiar caché de Expo
npx expo start -c
```

### Errores de módulos no encontrados
```bash
# Reinstalar dependencias
rm -rf node_modules
npm install
```

### Problemas de conexión con APIs
- Verifica que los servicios backend estén ejecutándose
- En dispositivos físicos, asegúrate de usar la IP local en lugar de localhost
- Verifica que el firewall permita las conexiones

### Puerto ocupado
```bash
# Usar un puerto diferente
npx expo start --port 8082
```

## Compatibilidad

- **Expo SDK**: ~54.0.25
- **React Native**: 0.81.5
- **React**: 19.1.0
- **TypeScript**: ~5.9.2

### Plataformas Soportadas
- iOS 13.4+
- Android 5.0+ (API 21+)
- Web (limitado)

## Próximas Mejoras

- [ ] Implementar notificaciones push
- [ ] Añadir modo oscuro
- [ ] Integrar mapas nativos
- [ ] Añadir chat entre conductor y pasajeros
- [ ] Implementar sistema de calificaciones
- [ ] Añadir soporte offline
- [ ] Optimizar imágenes y rendimiento
- [ ] Añadir tests unitarios y e2e

## Contribución

1. Clona el repositorio
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## Licencia

Este proyecto es parte del sistema CarPooling.

## Contacto

Para preguntas o soporte, contacta al equipo de desarrollo.

---

**Nota**: Esta aplicación móvil está diseñada para replicar exactamente la funcionalidad y diseño de la versión web de CarPooling, proporcionando una experiencia nativa y optimizada para dispositivos móviles.
