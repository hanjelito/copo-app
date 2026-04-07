# CoPo App

App de carpooling similar a BlaBlaCar.

## Stack

- **API Gateway**: Go + chi + JWT + rate limiting
- **Microservicios**: Go + chi
- **Base de datos**: PostgreSQL 16 (source of truth) + Redis 7 (cache + sessions + geo)
- **Mensajería**: Kafka 4.2
- **Workers**: Go

## Servicios e Infraestructura

| Servicio | Puerto | Estado |
|----------|--------|--------|
| Auth Service | 8081 | ✅ |
| PostgreSQL | 5432 | ✅ |
| Redis | 6379 | ✅ |
| Kafka | 9092 | ✅ |

## Auth Service (`services/auth`)

### Endpoints

| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/auth/register` | Registro de usuario |
| POST | `/auth/login` | Login → devuelve JWT |
| POST | `/auth/refresh` | Refresh token |

### Variables de entorno (`services/auth/.env`)

```env
DATABASE_URL=postgres://admin:admin@localhost:5432/tienda
JWT_SECRET=un_secreto_muy_largo_y_seguro
PORT=8081
```

### Schema PostgreSQL

```sql
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email      TEXT UNIQUE NOT NULL,
  password   TEXT NOT NULL,
  name       TEXT NOT NULL,
  role       TEXT NOT NULL DEFAULT 'passenger',
  created_at TIMESTAMP DEFAULT NOW()
);
```

### Estructura

```
services/auth/
├── cmd/
│   └── main.go              ← entry point, arranca el server en :8081
├── internal/
│   ├── model/
│   │   └── user.go          ← structs: User, RegisterRequest, LoginRequest, AuthResponse
│   ├── repository/
│   │   └── user.go          ← queries SQL: Create, FindByEmail
│   ├── service/
│   │   └── auth.go          ← lógica: bcrypt, generar/validar JWT
│   └── handler/
│       └── auth.go          ← HTTP handlers: Register, Login, Refresh
└── go.mod
```

### Dependencias

```
github.com/go-chi/chi/v5       ← router HTTP
github.com/golang-jwt/jwt/v5   ← JWT tokens
golang.org/x/crypto            ← bcrypt
github.com/jackc/pgx/v5        ← driver PostgreSQL
github.com/joho/godotenv       ← variables de entorno
```

### Arrancar

```bash
cd services/auth && go run cmd/main.go
```

### Ejemplos curl

```bash
# Registro
curl -X POST http://localhost:8081/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@copo.com","password":"123456","name":"Angel","role":"driver"}'

# Login
curl -X POST http://localhost:8081/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@copo.com","password":"123456"}'

# Refresh token
curl -X POST http://localhost:8081/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"<tu_refresh_token>"}'
```

## Arquitectura

![Arquitectura CoPo](images/Copo.jpg)

```
Clientes (Mobile App / Web / CLI)
          ↓
    API Gateway
    (JWT Auth + rate limiting)
          ↓
  ┌───────┼───────┐
Auth    User    Rides    Bookings
         └───────┴──────────┘
                  ↓
          PostgreSQL + Redis
                  ↓
               Kafka
                  ↓
       Notification Worker (email/push)
```

## Flujo: Reservar un Viaje (POST /bookings)

```
Cliente → API Gateway (JWT + rate limit)
            ↓
         Try Redis (caché de disponibilidad)
           ├── Hit  → INSERT booking
           └── Miss → PostgreSQL (check ride & seats)
                        ├── No existe → error
                        └── Existe   → PostgreSQL INSERT booking (transacción)
                                          ↓
                                Kafka: publish booking.created
                                          ↓
                                API Gateway → 201 Created
                                          ↓
                                Notification Worker → email/push

         Si PG falla → Booking Retry Worker reintenta y notifica al usuario
```

---

## Roadmap de Implementación

### Fase 1 — Base
- [x] Infraestructura Docker (PostgreSQL, Redis, Kafka)
- [x] Auth Service (`/auth/*`) — registro, login, refresh token ✅
  - [x] `model/user.go` — structs
  - [x] `repository/user.go` — queries SQL
  - [x] `service/auth.go` — bcrypt + JWT
  - [x] `handler/auth.go` — HTTP handlers
  - [x] `cmd/main.go` — entry point
  - [x] `middleware/jwt.go` — validar JWT en rutas protegidas ✅
- [ ] User Service (`/users`) — perfil básico
  - [ ] GET `/users/me` — obtener perfil
  - [ ] PUT `/users/me` — actualizar perfil
- [ ] API Gateway — enrutar hacia servicios internos vía gRPC
  - [ ] JWT Auth + rate limiting
  - [ ] Rutas: `/auth/*`, `/users/*`, `/rides/*`, `/bookings/*`

### Fase 2 — Core del negocio
- [ ] Rides Service (`/rides/*`) — viajes del conductor
  - [ ] POST `/rides` — crear viaje
  - [ ] GET `/rides` — listar viajes disponibles
  - [ ] GET `/rides/:id` — detalle de viaje
- [ ] Bookings Service (`/bookings/*`) — reservas del pasajero
  - [ ] POST `/bookings` — reservar viaje (flujo Redis + PG)
  - [ ] GET `/bookings/me` — mis reservas
  - [ ] DELETE `/bookings/:id` — cancelar reserva

### Fase 3 — Async
- [ ] Kafka topics: `booking.created`, `booking.pending`
- [ ] Notification Worker — email + push al conductor y pasajero

### Fase 4 — Extras
- [ ] Booking Retry Worker — reintenta INSERT si PG falla
- [ ] Geo / tracking en tiempo real (Redis Geo + WebSockets)
- [ ] Pagos
- [ ] Ratings post-viaje
- [ ] gRPC entre API Gateway y microservicios
