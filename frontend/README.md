# 🎨 Frontend - CarPooling Platform

Aplicación web SPA (Single Page Application) desarrollada con React y TypeScript para la plataforma CarPooling. Proporciona una interfaz de usuario moderna e intuitiva para conductores, pasajeros y administradores.

## 📋 Descripción

El frontend de CarPooling es una aplicación React moderna que permite:
- **Usuarios**: Buscar viajes, realizar reservas, gestionar perfil
- **Conductores**: Publicar y gestionar viajes, ver reservas
- **Administradores**: Panel de control para gestionar usuarios, viajes y reservas
- **Autenticación**: Sistema de login/registro con verificación de email
- **Responsive**: Diseño adaptable a dispositivos móviles y desktop

### Características Principales

- ✅ SPA con React Router para navegación
- ✅ Autenticación con JWT y Context API
- ✅ Rutas protegidas y públicas
- ✅ Panel de administración completo
- ✅ Sistema de búsqueda de viajes con filtros
- ✅ Gestión de reservas en tiempo real
- ✅ Perfil de usuario editable
- ✅ Diseño responsive con Tailwind CSS
- ✅ Animaciones fluidas con Framer Motion
- ✅ TypeScript para type safety

---

## 🚀 Tecnologías

| Tecnología | Versión | Propósito |
|------------|---------|-----------|
| **React** | 19.2.0 | Framework UI |
| **TypeScript** | 5.9.3 | Lenguaje tipado |
| **Vite** | 7.2.2 (Rolldown) | Build tool ultra-rápido |
| **React Router** | 7.9.6 | Routing y navegación |
| **Tailwind CSS** | 4.1.17 | Framework CSS utility-first |
| **Axios** | 1.13.2 | Cliente HTTP |
| **Framer Motion** | 12.23.24 | Librería de animaciones |
| **Tabler Icons** | 3.35.0 | Biblioteca de iconos |
| **Lucide React** | 0.553.0 | Iconos adicionales |
| **Radix UI** | Latest | Componentes accesibles |

### Herramientas de Desarrollo

- **ESLint**: Linting de código
- **TypeScript ESLint**: Reglas específicas de TypeScript
- **Vite Plugin React**: Hot Module Replacement (HMR)
- **Autoprefixer**: Prefijos CSS automáticos

---

## 📁 Estructura del Proyecto

```
frontend/
├── public/                    # Archivos estáticos
│   └── vite.svg              # Favicon
├── src/
│   ├── components/           # Componentes reutilizables
│   │   ├── layout/          # Layout components (Navbar, Footer)
│   │   ├── admin/           # Componentes del admin panel
│   │   ├── routes/          # Componentes de rutas (AdminRoute)
│   │   └── ui/              # UI components (Button, Input, etc.)
│   │
│   ├── pages/               # Páginas de la aplicación
│   │   ├── HomePage.tsx
│   │   ├── LoginPage.tsx
│   │   ├── RegisterPage.tsx
│   │   ├── SearchPage.tsx
│   │   ├── CreateTripPage.tsx
│   │   ├── MyTripsPage.tsx
│   │   ├── MyBookingsPage.tsx
│   │   ├── ProfilePage.tsx
│   │   ├── TripDetailPage.tsx
│   │   ├── EditTripPage.tsx
│   │   └── admin/           # Páginas del admin
│   │       ├── AdminDashboardPage.tsx
│   │       ├── AdminUsersPage.tsx
│   │       ├── AdminTripsPage.tsx
│   │       └── AdminBookingsPage.tsx
│   │
│   ├── services/            # API clients
│   │   ├── api.ts          # Cliente Axios configurado
│   │   ├── authService.ts  # Servicios de autenticación
│   │   ├── tripService.ts  # Servicios de viajes
│   │   ├── bookingService.ts # Servicios de reservas
│   │   └── searchService.ts # Servicios de búsqueda
│   │
│   ├── context/             # Context providers
│   │   └── AuthContext.tsx # Context de autenticación
│   │
│   ├── types/              # TypeScript types
│   │   └── index.ts        # Definiciones de tipos
│   │
│   ├── lib/                # Utilidades
│   │   └── utils.ts        # Funciones helper
│   │
│   ├── App.tsx             # Componente principal
│   ├── main.tsx            # Entry point
│   └── index.css           # Estilos globales
│
├── nginx.conf              # Configuración de Nginx para producción
├── Dockerfile              # Build de producción
├── vite.config.ts          # Configuración de Vite
├── tsconfig.json           # Configuración de TypeScript
├── tailwind.config.js      # Configuración de Tailwind
├── package.json            # Dependencias y scripts
└── README.md               # Este archivo
```

---

## ⚙️ Configuración

### Variables de Entorno

El frontend usa un proxy de Vite en desarrollo. En producción, Nginx maneja el routing de APIs.

**Desarrollo (`vite.config.ts`)**:
```typescript
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:8001', // users-api
      changeOrigin: true,
    }
  }
}
```

**Producción**: Nginx maneja el proxy (ver [nginx.conf](nginx.conf))

---

## 📦 Instalación

### Prerrequisitos

- Node.js 18+ y npm
- Git

### Pasos de Instalación

1. **Clonar el repositorio**
   ```bash
   git clone https://github.com/your-org/CarPooling.git
   cd CarPooling/frontend
   ```

2. **Instalar dependencias**
   ```bash
   npm install
   ```

3. **Ejecutar en modo desarrollo**
   ```bash
   npm run dev
   ```

4. **Acceder a la aplicación**
   ```
   http://localhost:3000
   ```

---

## 🛠️ Scripts Disponibles

```bash
# Desarrollo con Hot Reload
npm run dev

# Build de producción
npm run build

# Preview del build
npm run preview

# Linting
npm run lint
```

---

## 🎨 Características del UI

### Diseño Responsive

- ✅ Mobile-first design
- ✅ Breakpoints: sm, md, lg, xl, 2xl
- ✅ Componentes adaptables con Tailwind CSS

### Temas y Estilos

- **Colores**: Palette personalizada con variables CSS
- **Tipografía**: Inter font (via @fontsource)
- **Componentes**: Radix UI para accesibilidad
- **Animaciones**: Framer Motion para transiciones suaves

### Componentes Principales

#### Layout
- **Navbar**: Navegación principal con scroll effects
- **AdminLayout**: Layout específico para panel admin
- **Layout**: Layout general con Navbar y Outlet

#### Rutas Protegidas
```typescript
<ProtectedRoute>
  <MyTripsPage />
</ProtectedRoute>
```

#### Rutas de Admin
```typescript
<AdminRoute>
  <AdminDashboardPage />
</AdminRoute>
```

---

## 🔐 Autenticación

### Sistema de Autenticación

El frontend usa **JWT tokens** almacenados en `localStorage`:

```typescript
// AuthContext.tsx
const AuthContext = createContext({
  user: null,
  isAuthenticated: false,
  login: (token: string) => {},
  logout: () => {},
})
```

### Flujo de Autenticación

1. Usuario ingresa credenciales en LoginPage
2. authService hace POST a `/api/login`
3. Backend retorna JWT token
4. Token se guarda en localStorage
5. AuthContext actualiza estado global
6. Rutas protegidas verifican isAuthenticated

### Verificación de Email

- Usuario recibe email con token de verificación
- Clic en link redirige a `/verify-email?token=xxx`
- VerifyEmailPage valida el token con users-api

---

## 🚏 Routing

### Rutas Públicas

- `/` - HomePage (landing page)
- `/login` - LoginPage
- `/register` - RegisterPage
- `/search` - SearchPage (búsqueda de viajes)
- `/trips/:id` - TripDetailPage
- `/verify-email` - VerifyEmailPage
- `/forgot-password` - ForgotPasswordPage
- `/reset-password` - ResetPasswordPage

### Rutas Protegidas (requieren login)

- `/create-trip` - CreateTripPage
- `/my-trips` - MyTripsPage
- `/my-bookings` - MyBookingsPage
- `/profile` - ProfilePage
- `/trips/:id/edit` - EditTripPage

### Rutas de Admin

- `/admin` - AdminDashboardPage
- `/admin/users` - AdminUsersPage
- `/admin/trips` - AdminTripsPage
- `/admin/bookings` - AdminBookingsPage

---

## 🌐 API Integration

### API Services

Todos los servicios usan Axios configurado en `services/api.ts`:

```typescript
// api.ts
const apiClient = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
  },
})

// Interceptor para agregar JWT token
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})
```

### Servicios Disponibles

#### authService
```typescript
login(email, password): Promise<AuthResponse>
register(userData): Promise<User>
verifyEmail(token): Promise<void>
forgotPassword(email): Promise<void>
resetPassword(token, newPassword): Promise<void>
```

#### tripService
```typescript
createTrip(tripData): Promise<Trip>
getTrip(id): Promise<Trip>
listTrips(filters): Promise<Trip[]>
updateTrip(id, data): Promise<Trip>
deleteTrip(id): Promise<void>
```

#### bookingService
```typescript
createBooking(tripId, seats): Promise<Booking>
getMyBookings(): Promise<Booking[]>
confirmBooking(id): Promise<Booking>
cancelBooking(id): Promise<void>
```

#### searchService
```typescript
searchTrips(query): Promise<SearchResponse>
autocomplete(query): Promise<string[]>
getTripDetails(id): Promise<SearchTrip>
```

---

## 🏗️ Build de Producción

### Build con Vite

```bash
npm run build
```

Esto genera archivos optimizados en `dist/`:
- HTML minificado
- CSS minificado y concatenado
- JS chunks con code splitting
- Assets con hashes para cache busting

### Dockerfile

El frontend incluye un Dockerfile multi-stage:

```dockerfile
# Stage 1: Build
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# Stage 2: Serve with Nginx
FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

### Nginx Configuration

El archivo [nginx.conf](nginx.conf) configura:
- Servidor del frontend en puerto 80
- Reverse proxy para APIs (/api/*)
- Gzip compression
- Security headers
- Cache headers para assets estáticos
- Fallback a index.html para React Router

**Rutas del proxy**:
```nginx
location /api/users {
    proxy_pass http://users-api:8001;
}

location /api/trips {
    proxy_pass http://trips-api:8002;
}

location /api/bookings {
    proxy_pass http://bookings-api:8003;
}

location /api/search {
    proxy_pass http://search-api:8004;
}
```

---

## 🐳 Docker

### Build y Run

```bash
# Build de la imagen
docker build -t frontend:latest .

# Run del contenedor
docker run -p 80:80 frontend:latest
```

### Docker Compose (Recomendado)

```bash
# Desde la raíz del proyecto
docker-compose up -d frontend

# Ver logs
docker-compose logs -f frontend

# Reconstruir tras cambios
docker-compose build frontend
docker-compose up -d frontend
```

---

## 🎯 Páginas Principales

### HomePage
- Landing page con llamada a la acción
- Preview de funcionalidades
- Links a register/login

### SearchPage
- Buscador de viajes con filtros
- Filtros: origen, destino, fecha, precio, asientos
- Resultados paginados
- Card de cada viaje con detalles

### CreateTripPage
- Formulario para publicar un viaje
- Validación de campos
- Integración con Google Places (opcional)

### MyTripsPage
- Lista de viajes publicados por el conductor
- Opciones: editar, eliminar
- Ver reservas de cada viaje

### MyBookingsPage
- Lista de reservas del usuario
- Estado: pending, confirmed, cancelled
- Opciones: confirmar, cancelar

### AdminDashboardPage
- Estadísticas generales
- Gráficos de uso
- Acciones rápidas

---

## 🧪 Testing

```bash
# Ejecutar tests (cuando estén configurados)
npm test

# Tests con coverage
npm run test:coverage
```

---

## 🔧 Desarrollo

### Agregar una Nueva Página

1. Crear archivo en `src/pages/`
2. Importar en `App.tsx`
3. Agregar ruta en el componente `Routes`
4. Si es protegida, envolver en `<ProtectedRoute>` o `<AdminRoute>`

### Agregar un Nuevo Servicio API

1. Crear archivo en `src/services/`
2. Importar `apiClient` de `api.ts`
3. Exportar funciones que usen `apiClient.get/post/put/delete`

### Agregar Tipos TypeScript

1. Agregar types en `src/types/index.ts`
2. Usar en componentes y servicios

---

## 📱 Responsive Breakpoints

```css
/* Tailwind CSS Breakpoints */
sm: 640px   /* Teléfonos grandes */
md: 768px   /* Tablets */
lg: 1024px  /* Laptops */
xl: 1280px  /* Desktops */
2xl: 1536px /* Pantallas grandes */
```

Ejemplo de uso:
```tsx
<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
  {/* 1 col en móvil, 2 en tablet, 3 en desktop */}
</div>
```

---

## 🚀 Performance Optimizations

- ✅ **Code Splitting**: React.lazy y Suspense
- ✅ **Tree Shaking**: Vite elimina código no usado
- ✅ **Asset Optimization**: Imágenes y fonts optimizados
- ✅ **Gzip Compression**: Configurado en Nginx
- ✅ **Cache Headers**: Assets con cache de 1 año
- ✅ **Lazy Loading**: Componentes cargados on-demand

---

## 🤝 Contribución

1. Crear una rama feature (`git checkout -b feature/nueva-pagina`)
2. Seguir las convenciones de código (ESLint)
3. Usar TypeScript para todos los nuevos archivos
4. Escribir componentes reutilizables
5. Hacer commit con conventional commits
6. Crear Pull Request

---

## 📄 Licencia

Este proyecto es parte del sistema CarPooling desarrollado para fines educativos.

---

## 👥 Equipo

Desarrollado por el equipo de CarPooling - Arquitectura de Software II

---

**Versión**: 1.0.0
**Puerto de desarrollo**: 3000
**Puerto de producción (Nginx)**: 80
**Estado**: ✅ En producción en AWS EC2
