# ============================================
# Stage 1: Build the React application
# ============================================
FROM node:20-alpine AS builder

# Set working directory
WORKDIR /app

# Install dependencies first (for better caching)
# Copy only package files to leverage Docker layer caching
COPY package*.json ./

# Install all dependencies (including devDependencies needed for build)
RUN npm ci

# Copy the rest of the application code
COPY . .

# Build arguments for Vite environment variables
ARG VITE_USERS_API_URL=http://localhost:8001
ARG VITE_TRIPS_API_URL=http://localhost:8002
ARG VITE_BOOKINGS_API_URL=http://localhost:8003
ARG VITE_SEARCH_API_URL=http://localhost:8004

# Set environment variables for the build
ENV VITE_USERS_API_URL=${VITE_USERS_API_URL}
ENV VITE_TRIPS_API_URL=${VITE_TRIPS_API_URL}
ENV VITE_BOOKINGS_API_URL=${VITE_BOOKINGS_API_URL}
ENV VITE_SEARCH_API_URL=${VITE_SEARCH_API_URL}

# Build the application
# This creates an optimized production build in the dist/ directory
RUN npm run build

# ============================================
# Stage 2: Serve the application with Nginx
# ============================================
FROM nginx:1.25-alpine

# Install curl for health checks
RUN apk add --no-cache curl

# Copy custom nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy built assets from builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Create a non-root user for better security
RUN adduser -D -u 1000 frontend && \
    chown -R frontend:frontend /usr/share/nginx/html && \
    chown -R frontend:frontend /var/cache/nginx && \
    chown -R frontend:frontend /var/log/nginx && \
    chown -R frontend:frontend /etc/nginx/conf.d && \
    touch /var/run/nginx.pid && \
    chown -R frontend:frontend /var/run/nginx.pid

# Switch to non-root user
USER frontend

# Expose port 80
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:80/ || exit 1

# Start nginx in the foreground
CMD ["nginx", "-g", "daemon off;"]
