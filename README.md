# CoPo App

App de carpooling.

## Stack

- **API Gateway**: Go + chi + JWT + rate limiting
- **Microservicios**: Go + chi
- **Base de datos**: PostgreSQL 16 (source of truth) + Redis 7 (cache + sessions + geo)
- **Mensajería**: Kafka 4.2
- **Workers**: Go

## Servicios e Infraestructura

| Servicio | Puerto | Estado |
|----------|--------|--------|
| API Gateway | 8080 | ✅ |
| Auth + User Service | 8081 | ✅ |
| Rides Service | 8083 | ✅ |
| Bookings Service | 8084 | ✅ |
| PostgreSQL | 5432 | ✅ |
| Redis | 6379 | ✅ |
| Kafka | 9092 | ✅ |

---

## Requisitos

- [Go 1.25+](https://golang.org/dl/)
- [Docker](https://www.docker.com/)
- [goose](https://github.com/pressly/goose) — migraciones

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

---

## Setup

### 1. Levantar infraestructura

```bash
docker compose up -d
```

### 2. Crear schema en PostgreSQL

```bash
docker exec -it gotest_postgres psql -U admin -d tienda -c "CREATE SCHEMA public;"
```

### 3. Correr migraciones

```bash
goose -dir migrations postgres "host=localhost port=5432 user=admin password=admin dbname=tienda sslmode=disable" up
```

### 4. Levantar servicios (cada uno en su terminal)

```bash
cd services/auth && go run cmd/main.go
cd services/rides && go run cmd/main.go
cd services/bookings && go run cmd/main.go
cd services/gateway && go run cmd/main.go
cd workers/notification && go run cmd/main.go
```

---

## Variables de entorno

### `services/auth/.env`
```env
DATABASE_URL=postgres://admin:admin@localhost:5432/tienda
JWT_SECRET=un_secreto_muy_largo_y_seguro
PORT=8081
```

### `services/rides/.env`
```env
DATABASE_URL=postgres://admin:admin@localhost:5432/tienda
PORT=8083
```

### `services/bookings/.env`
```env
DATABASE_URL=postgres://admin:admin@localhost:5432/tienda
REDIS_URL=localhost:6379
KAFKA_URL=localhost:9092
PORT=8084
```

### `workers/notification/.env`
```env
KAFKA_URL=localhost:9092
```

### `services/gateway/.env`
```env
PORT=8080
JWT_SECRET=un_secreto_muy_largo_y_seguro
AUTH_SERVICE_URL=http://localhost:8081
USER_SERVICE_URL=http://localhost:8081
RIDES_SERVICE_URL=http://localhost:8083
BOOKINGS_SERVICE_URL=http://localhost:8084
```

---

## Migraciones

```bash
# Aplicar
goose -dir migrations postgres "host=localhost port=5432 user=admin password=admin dbname=tienda sslmode=disable" up

# Revertir última
goose -dir migrations postgres "host=localhost port=5432 user=admin password=admin dbname=tienda sslmode=disable" down
```

---

## API — Ejemplos curl (vía gateway en :8080)

```bash
# ─── AUTH ───────────────────────────────────────────

# Registro
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@copo.com","password":"123456","name":"Angel","role":"driver"}'

# Login
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@copo.com","password":"123456"}'

# Guardar token en variable
TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@copo.com","password":"123456"}' | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)

# Refresh token
curl -X POST http://localhost:8080/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"<tu_refresh_token>"}'

# ─── USERS ──────────────────────────────────────────

# Obtener perfil
curl http://localhost:8080/users/me \
  -H "Authorization: Bearer $TOKEN"

# Actualizar nombre
curl -X PUT http://localhost:8080/users/me \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Angel Updated"}'

# ─── RIDES ──────────────────────────────────────────

# Crear viaje
curl -X POST http://localhost:8080/rides \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"origin":"Monterrey","destination":"CDMX","departure":"2026-04-10T08:00:00Z","seats":3}'

# Listar viajes
curl http://localhost:8080/rides \
  -H "Authorization: Bearer $TOKEN"

# Detalle de viaje
curl http://localhost:8080/rides/<id> \
  -H "Authorization: Bearer $TOKEN"

# ─── BOOKINGS ───────────────────────────────────────

# Crear reserva
curl -X POST http://localhost:8080/bookings \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"ride_id":"<ride_id>"}'

# Mis reservas
curl http://localhost:8080/bookings/me \
  -H "Authorization: Bearer $TOKEN"

# Cancelar reserva
curl -X DELETE http://localhost:8080/bookings/<id> \
  -H "Authorization: Bearer $TOKEN"
```

---

## Arquitectura

![Arquitectura CoPo](images/Arquitectura\ general.jpg)

```
Clientes (Mobile App / Web / CLI)
          ↓
    API Gateway :8080
    (JWT Auth + rate limiting)
          ↓
  ┌───────┼─────────┐
Auth   Users   Rides   Bookings
:8081  :8081  :8083    :8084
                         ↓
               PostgreSQL + Redis + Kafka
                         ↓
              Notification Worker (email/push)
```

## Flujo: Reservar un Viaje (POST /bookings)

![Flujo reserva](images/flujo-reserva-trip.jpg)

```
Cliente → API Gateway (JWT + rate limit)
            ↓
         Try Redis (caché de disponibilidad)
           ├── Hit  → INSERT booking en PostgreSQL
           └── Miss → PostgreSQL (check ride & seats)
                        └── Existe → INSERT booking
                                       ↓
                              Kafka: publish booking.created
                                       ↓
                              API Gateway → 201 Created
                                       ↓
                              Notification Worker → email/push

         Si PG falla → Kafka: booking.pending → Retry Worker
```

---

## Roadmap

### Fase 1 — Base ✅
- [x] Infraestructura Docker (PostgreSQL, Redis, Kafka)
- [x] Auth Service — registro, login, refresh token
- [x] User Service — GET/PUT `/users/me`
- [x] API Gateway — proxy HTTP + JWT + rate limiting

### Fase 2 — Core del negocio ✅
- [x] Rides Service — crear, listar, detalle de viajes
- [x] Bookings Service — reservar, mis reservas, cancelar (Redis + PG + Kafka)

### Fase 3 — Async
- [x] Notification Worker — consume `booking.created` de Kafka ✅
- [x] Control de capacidad — Redis controla asientos disponibles por ride ✅

### Fase 4 — Extras
- [x] Migraciones con `goose` ✅
- [ ] Booking Retry Worker — reintenta INSERT si PG falla
- [ ] Geo / tracking en tiempo real (Redis Geo + WebSockets)
- [ ] Pagos
- [ ] Ratings post-viaje
- [ ] gRPC entre API Gateway y microservicios
