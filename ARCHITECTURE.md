# Outdoor Equipment Rental Marketplace — Architecture

## Monorepo Structure

```
rent-outdoor/
├── apps/
│   ├── api/                          # Golang Backend (Gin + Clean Architecture)
│   └── web/                          # Nuxt 4 Frontend (Vue 3 + TypeScript)
├── packages/
│   └── shared/                       # Shared constants, types, contracts
├── deployments/
│   ├── docker/
│   │   ├── api.Dockerfile
│   │   ├── web.Dockerfile
│   │   ├── mysql.Dockerfile
│   │   └── nginx.conf
│   ├── docker-compose.yml
│   ├── docker-compose.dev.yml
│   └── docker-compose.prod.yml
├── scripts/
│   ├── setup.sh
│   ├── migrate.sh
│   ├── seed.sh
│   └── generate-swagger.sh
├── .github/
│   └── workflows/
│       ├── ci.yml
│       ├── cd-staging.yml
│       └── cd-production.yml
├── .env.example
├── Makefile
├── README.md
└── ARCHITECTURE.md
```

---

## 1. Backend Structure (`apps/api/`)

```
apps/api/
├── cmd/
│   └── server/
│       └── main.go                   # Entry point
├── internal/
│   ├── domain/                       # Enterprise Business Rules (Entities)
│   │   ├── entity/
│   │   │   ├── user.go
│   │   │   ├── equipment.go
│   │   │   ├── category.go
│   │   │   ├── rental.go
│   │   │   ├── payment.go
│   │   │   ├── review.go
│   │   │   ├── location.go
│   │   │   └── notification.go
│   │   ├── valueobject/
│   │   │   ├── money.go
│   │   │   ├── rental_period.go
│   │   │   ├── coordinate.go
│   │   │   └── rental_status.go
│   │   └── event/
│   │       ├── rental_created.go
│   │       ├── payment_completed.go
│   │       └── equipment_returned.go
│   │
│   ├── usecase/                      # Application Business Rules
│   │   ├── auth/
│   │   │   ├── register.go
│   │   │   ├── login.go
│   │   │   ├── refresh_token.go
│   │   │   └── forgot_password.go
│   │   ├── equipment/
│   │   │   ├── create_equipment.go
│   │   │   ├── update_equipment.go
│   │   │   ├── list_equipment.go
│   │   │   ├── search_equipment.go
│   │   │   └── delete_equipment.go
│   │   ├── rental/
│   │   │   ├── create_rental.go
│   │   │   ├── confirm_rental.go
│   │   │   ├── cancel_rental.go
│   │   │   ├── complete_rental.go
│   │   │   └── list_rentals.go
│   │   ├── payment/
│   │   │   ├── create_payment.go
│   │   │   ├── handle_callback.go
│   │   │   └── check_status.go
│   │   ├── review/
│   │   │   ├── create_review.go
│   │   │   └── list_reviews.go
│   │   └── user/
│   │       ├── get_profile.go
│   │       ├── update_profile.go
│   │       └── upload_avatar.go
│   │
│   ├── repository/                   # Repository Interfaces (Ports)
│   │   ├── user_repository.go
│   │   ├── equipment_repository.go
│   │   ├── rental_repository.go
│   │   ├── payment_repository.go
│   │   ├── review_repository.go
│   │   └── category_repository.go
│   │
│   ├── adapter/                      # Interface Adapters
│   │   ├── handler/                  # HTTP Handlers (Controllers)
│   │   │   ├── auth_handler.go
│   │   │   ├── equipment_handler.go
│   │   │   ├── rental_handler.go
│   │   │   ├── payment_handler.go
│   │   │   ├── review_handler.go
│   │   │   ├── user_handler.go
│   │   │   └── health_handler.go
│   │   ├── middleware/
│   │   │   ├── auth.go
│   │   │   ├── cors.go
│   │   │   ├── ratelimit.go
│   │   │   ├── logger.go
│   │   │   └── recovery.go
│   │   ├── presenter/               # Response formatting
│   │   │   ├── json_presenter.go
│   │   │   └── error_presenter.go
│   │   └── dto/                      # Request/Response DTOs
│   │       ├── request/
│   │       │   ├── auth_request.go
│   │       │   ├── equipment_request.go
│   │       │   ├── rental_request.go
│   │       │   └── payment_request.go
│   │       └── response/
│   │           ├── auth_response.go
│   │           ├── equipment_response.go
│   │           ├── rental_response.go
│   │           ├── pagination_response.go
│   │           └── error_response.go
│   │
│   └── infrastructure/              # Frameworks & Drivers
│       ├── persistence/
│       │   ├── mysql/
│       │   │   ├── connection.go
│       │   │   ├── user_repository_impl.go
│       │   │   ├── equipment_repository_impl.go
│       │   │   ├── rental_repository_impl.go
│       │   │   ├── payment_repository_impl.go
│       │   │   ├── review_repository_impl.go
│       │   │   └── category_repository_impl.go
│       │   └── migration/
│       │       ├── migrator.go
│       │       └── seeds/
│       │           ├── category_seed.go
│       │           └── admin_seed.go
│       ├── storage/
│       │   ├── s3_client.go
│       │   └── upload_service.go
│       ├── payment/
│       │   ├── midtrans_client.go
│       │   └── midtrans_service.go
│       ├── email/
│       │   └── smtp_service.go
│       ├── cache/
│       │   └── redis_client.go
│       └── config/
│           └── config.go
│
├── pkg/                              # Shared internal packages
│   ├── logger/
│   │   └── logger.go
│   ├── validator/
│   │   └── validator.go
│   ├── jwt/
│   │   └── jwt.go
│   ├── hash/
│   │   └── bcrypt.go
│   ├── pagination/
│   │   └── pagination.go
│   └── response/
│       └── response.go
│
├── docs/                             # Swagger documentation
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
│
├── tests/
│   ├── unit/
│   │   ├── usecase/
│   │   │   ├── auth_test.go
│   │   │   ├── equipment_test.go
│   │   │   └── rental_test.go
│   │   └── domain/
│   │       ├── rental_test.go
│   │       └── money_test.go
│   ├── integration/
│   │   ├── repository/
│   │   │   ├── user_repository_test.go
│   │   │   ├── equipment_repository_test.go
│   │   │   └── rental_repository_test.go
│   │   ├── handler/
│   │   │   ├── auth_handler_test.go
│   │   │   └── equipment_handler_test.go
│   │   └── testutil/
│   │       ├── db.go
│   │       └── fixtures.go
│   └── mocks/
│       ├── user_repository_mock.go
│       ├── equipment_repository_mock.go
│       ├── rental_repository_mock.go
│       └── payment_service_mock.go
│
├── go.mod
├── go.sum
└── .air.toml                         # Hot reload config
```

---

## 2. Frontend Structure (`apps/web/`)

```
apps/web/
├── app/
│   ├── components/
│   │   ├── common/
│   │   │   ├── AppHeader.vue
│   │   │   ├── AppFooter.vue
│   │   │   ├── AppNavigation.vue
│   │   │   ├── AppLogo.vue
│   │   │   ├── BaseButton.vue
│   │   │   ├── BaseInput.vue
│   │   │   ├── BaseModal.vue
│   │   │   ├── BaseCard.vue
│   │   │   ├── BasePagination.vue
│   │   │   ├── BaseAlert.vue
│   │   │   ├── BaseLoader.vue
│   │   │   └── BaseEmptyState.vue
│   │   ├── equipment/
│   │   │   ├── EquipmentCard.vue
│   │   │   ├── EquipmentGrid.vue
│   │   │   ├── EquipmentDetail.vue
│   │   │   ├── EquipmentGallery.vue
│   │   │   ├── EquipmentFilter.vue
│   │   │   ├── EquipmentSearch.vue
│   │   │   └── EquipmentMap.vue
│   │   ├── rental/
│   │   │   ├── RentalForm.vue
│   │   │   ├── RentalCard.vue
│   │   │   ├── RentalTimeline.vue
│   │   │   ├── RentalSummary.vue
│   │   │   └── RentalStatus.vue
│   │   ├── payment/
│   │   │   ├── PaymentForm.vue
│   │   │   ├── PaymentSummary.vue
│   │   │   └── PaymentStatus.vue
│   │   ├── review/
│   │   │   ├── ReviewForm.vue
│   │   │   ├── ReviewCard.vue
│   │   │   └── ReviewList.vue
│   │   ├── auth/
│   │   │   ├── LoginForm.vue
│   │   │   ├── RegisterForm.vue
│   │   │   └── ForgotPasswordForm.vue
│   │   └── dashboard/
│   │       ├── DashboardStats.vue
│   │       ├── DashboardRentals.vue
│   │       └── DashboardEquipment.vue
│   │
│   ├── composables/
│   │   ├── useAuth.ts
│   │   ├── useEquipment.ts
│   │   ├── useRental.ts
│   │   ├── usePayment.ts
│   │   ├── useReview.ts
│   │   ├── useNotification.ts
│   │   ├── useGeolocation.ts
│   │   └── useInfiniteScroll.ts
│   │
│   ├── layouts/
│   │   ├── default.vue
│   │   ├── auth.vue
│   │   └── dashboard.vue
│   │
│   ├── pages/
│   │   ├── index.vue                      # Landing / Home
│   │   ├── login.vue
│   │   ├── register.vue
│   │   ├── forgot-password.vue
│   │   ├── equipment/
│   │   │   ├── index.vue                  # Browse / Search
│   │   │   └── [id].vue                   # Equipment Detail
│   │   ├── rental/
│   │   │   ├── index.vue                  # My Rentals
│   │   │   ├── [id].vue                   # Rental Detail
│   │   │   └── checkout.vue               # Checkout Flow
│   │   ├── payment/
│   │   │   ├── [id].vue                   # Payment Page
│   │   │   ├── success.vue
│   │   │   └── failed.vue
│   │   ├── dashboard/
│   │   │   ├── index.vue                  # Owner Dashboard
│   │   │   ├── equipment/
│   │   │   │   ├── index.vue              # Manage Equipment
│   │   │   │   ├── create.vue
│   │   │   │   └── [id]/
│   │   │   │       └── edit.vue
│   │   │   ├── rentals.vue                # Manage Rentals
│   │   │   ├── earnings.vue
│   │   │   └── settings.vue
│   │   └── profile/
│   │       ├── index.vue
│   │       └── edit.vue
│   │
│   ├── middleware/
│   │   ├── auth.ts
│   │   ├── guest.ts
│   │   └── owner.ts
│   │
│   └── plugins/
│       ├── api.ts
│       └── toast.ts
│
├── stores/
│   ├── auth.ts
│   ├── equipment.ts
│   ├── rental.ts
│   ├── payment.ts
│   ├── notification.ts
│   └── ui.ts
│
├── services/
│   ├── api/
│   │   ├── client.ts                     # Axios/ofetch instance
│   │   ├── auth.service.ts
│   │   ├── equipment.service.ts
│   │   ├── rental.service.ts
│   │   ├── payment.service.ts
│   │   ├── review.service.ts
│   │   └── upload.service.ts
│   └── types/
│       ├── auth.types.ts
│       ├── equipment.types.ts
│       ├── rental.types.ts
│       ├── payment.types.ts
│       └── api.types.ts
│
├── utils/
│   ├── formatters.ts
│   ├── validators.ts
│   ├── constants.ts
│   └── helpers.ts
│
├── assets/
│   ├── css/
│   │   └── main.css                      # Tailwind directives + custom styles
│   ├── images/
│   └── icons/
│
├── public/
│   ├── favicon.ico
│   ├── manifest.json                     # PWA manifest
│   └── sw.js
│
├── tests/
│   ├── unit/
│   │   ├── components/
│   │   │   ├── EquipmentCard.spec.ts
│   │   │   └── RentalForm.spec.ts
│   │   ├── composables/
│   │   │   ├── useAuth.spec.ts
│   │   │   └── useRental.spec.ts
│   │   └── stores/
│   │       ├── auth.spec.ts
│   │       └── equipment.spec.ts
│   └── e2e/
│       ├── auth.spec.ts
│       ├── equipment-browse.spec.ts
│       └── rental-flow.spec.ts
│
├── nuxt.config.ts
├── tailwind.config.ts
├── tsconfig.json
├── vitest.config.ts
└── package.json
```

---

## 3. Shared Packages (`packages/shared/`)

```
packages/shared/
├── constants/
│   ├── rental-status.ts
│   ├── payment-status.ts
│   ├── equipment-category.ts
│   └── roles.ts
├── types/
│   ├── api-response.ts
│   ├── pagination.ts
│   └── error-codes.ts
├── validators/
│   └── rental-rules.ts
├── package.json
└── tsconfig.json
```

---

## 4. Docker Structure (`deployments/`)

```
deployments/
├── docker/
│   ├── api.Dockerfile                    # Multi-stage Go build
│   ├── web.Dockerfile                    # Multi-stage Node build + Nginx
│   ├── mysql.Dockerfile                  # Custom MySQL with init scripts
│   └── nginx.conf                        # Reverse proxy config
├── docker-compose.yml                    # Base services definition
├── docker-compose.dev.yml                # Dev overrides (volumes, hot reload)
├── docker-compose.prod.yml               # Prod overrides (resource limits, replicas)
└── init-db/
    └── 01-schema.sql                     # Initial schema for Docker MySQL
```

### Docker Services

| Service     | Image/Build          | Port  | Purpose                      |
|-------------|----------------------|-------|------------------------------|
| `api`       | `docker/api.Dockerfile` | 8080  | Go API server                |
| `web`       | `docker/web.Dockerfile` | 3000  | Nuxt SSR / Static            |
| `mysql`     | `mysql:8.0`          | 3306  | Primary database             |
| `redis`     | `redis:7-alpine`     | 6379  | Cache / Session store        |
| `minio`     | `minio/minio`        | 9000  | S3-compatible object storage |
| `nginx`     | `nginx:alpine`       | 80/443| Reverse proxy / SSL          |
| `swagger-ui`| `swaggerapi/swagger-ui` | 8081 | API docs (dev only)        |

---

## 5. CI/CD Pipeline (`.github/workflows/`)

### `ci.yml` — Triggered on PR

```
Jobs:
├── lint-backend         → golangci-lint
├── lint-frontend        → eslint + prettier
├── test-backend-unit    → go test ./tests/unit/...
├── test-backend-integ   → go test ./tests/integration/... (with MySQL container)
├── test-frontend-unit   → vitest run
├── test-frontend-e2e    → playwright (optional, on label)
├── build-backend        → go build
├── build-frontend       → nuxt build
└── security-scan        → trivy / gosec
```

### `cd-staging.yml` — Triggered on merge to `develop`

```
Jobs:
├── build-and-push-api   → Docker build + push to registry
├── build-and-push-web   → Docker build + push to registry
├── deploy-staging       → Deploy to staging environment
└── smoke-test           → Health check + basic API test
```

### `cd-production.yml` — Triggered on merge to `main`

```
Jobs:
├── build-and-push-api   → Docker build + push (tagged)
├── build-and-push-web   → Docker build + push (tagged)
├── deploy-production    → Rolling deploy
└── notify               → Slack / Discord notification
```

---

## 6. Development Standards

### Git Strategy

| Branch       | Purpose                        |
|--------------|--------------------------------|
| `main`       | Production-ready code          |
| `develop`    | Integration branch             |
| `feature/*`  | New features                   |
| `fix/*`      | Bug fixes                      |
| `hotfix/*`   | Production hotfixes            |
| `release/*`  | Release preparation            |

### Commit Convention

```
<type>(<scope>): <subject>

Types: feat, fix, docs, style, refactor, test, chore, perf
Scopes: api, web, shared, deploy, ci
```

### Backend Standards

| Concern               | Standard                                  |
|-----------------------|-------------------------------------------|
| Architecture          | Clean Architecture (4 layers)             |
| Dependency Direction  | Inward only (domain has zero imports)      |
| Error Handling        | Custom error types with codes             |
| Validation            | DTO-level with `go-playground/validator`  |
| Logging               | Structured (zerolog / zap)                |
| Config                | Environment variables via `viper`         |
| API Versioning        | URL prefix `/api/v1`                      |
| Auth                  | JWT (access + refresh token pair)         |
| Pagination            | Cursor-based for lists, offset for admin  |
| Testing               | Table-driven tests, testify assertions    |
| Mocking               | mockery / gomock for interfaces           |
| Documentation         | swaggo annotations on handlers            |

### Frontend Standards

| Concern               | Standard                                  |
|-----------------------|-------------------------------------------|
| State Management      | Pinia (one store per domain)              |
| API Layer             | Service classes with typed responses      |
| Component Naming      | PascalCase, domain-prefixed               |
| Composables           | `use` prefix, single responsibility      |
| Styling               | TailwindCSS utility-first                 |
| Forms                 | VeeValidate + Zod schemas                 |
| Error Handling        | Global error boundary + toast             |
| PWA                   | Workbox via `@vite-pwa/nuxt`              |
| Testing               | Vitest (unit), Playwright (e2e)           |
| Image Handling        | `<NuxtImage>` with S3 provider           |

### API Contract

```
Base URL: /api/v1

Auth:
  POST   /auth/register
  POST   /auth/login
  POST   /auth/refresh
  POST   /auth/forgot-password

Equipment:
  GET    /equipment                  # List (public, filterable)
  GET    /equipment/:id              # Detail (public)
  POST   /equipment                  # Create (owner)
  PUT    /equipment/:id              # Update (owner)
  DELETE /equipment/:id              # Soft delete (owner)

Rentals:
  GET    /rentals                    # My rentals (renter)
  GET    /rentals/incoming           # Incoming rentals (owner)
  POST   /rentals                    # Create rental request
  PATCH  /rentals/:id/confirm        # Owner confirms
  PATCH  /rentals/:id/cancel         # Cancel
  PATCH  /rentals/:id/complete       # Mark returned

Payments:
  POST   /payments                   # Initiate (Midtrans snap)
  POST   /payments/callback          # Midtrans webhook
  GET    /payments/:id/status        # Check status

Reviews:
  GET    /equipment/:id/reviews      # List reviews
  POST   /equipment/:id/reviews      # Create review (post-rental)

Users:
  GET    /users/me                   # Profile
  PUT    /users/me                   # Update profile
  POST   /users/me/avatar            # Upload avatar

Categories:
  GET    /categories                 # List all
```

### Environment Variables

```env
# API
APP_ENV=development
APP_PORT=8080
APP_SECRET=

# Database
DB_HOST=mysql
DB_PORT=3306
DB_NAME=rent_outdoor
DB_USER=
DB_PASSWORD=

# Redis
REDIS_HOST=redis
REDIS_PORT=6379

# JWT
JWT_ACCESS_SECRET=
JWT_REFRESH_SECRET=
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d

# Storage (S3)
S3_ENDPOINT=
S3_BUCKET=
S3_ACCESS_KEY=
S3_SECRET_KEY=
S3_REGION=

# Midtrans
MIDTRANS_SERVER_KEY=
MIDTRANS_CLIENT_KEY=
MIDTRANS_ENV=sandbox

# Frontend
NUXT_PUBLIC_API_BASE=http://localhost:8080/api/v1
NUXT_PUBLIC_S3_URL=
```

---

## 7. Dependency Flow

```
┌────────────────────────────────────────────────────┐
│                    Presentation                      │
│   (Handlers, Middleware, DTOs, Presenters)          │
└──────────────────────┬─────────────────────────────┘
                       │ depends on
┌──────────────────────▼─────────────────────────────┐
│                    Use Cases                         │
│   (Application logic, orchestration)               │
└──────────────────────┬─────────────────────────────┘
                       │ depends on
┌──────────────────────▼─────────────────────────────┐
│                     Domain                          │
│   (Entities, Value Objects, Repository Interfaces) │
└──────────────────────┬─────────────────────────────┘
                       │ implemented by
┌──────────────────────▼─────────────────────────────┐
│                 Infrastructure                      │
│   (MySQL, S3, Midtrans, Redis, SMTP)               │
└────────────────────────────────────────────────────┘
```

---

## 8. Makefile Targets

```makefile
# Development
make dev              # Start all services (docker-compose dev)
make dev-api          # Start API with hot reload (air)
make dev-web          # Start Nuxt dev server

# Build
make build            # Build all
make build-api        # Build Go binary
make build-web        # Build Nuxt

# Database
make migrate          # Run migrations
make migrate-down     # Rollback last migration
make seed             # Run seeders

# Testing
make test             # Run all tests
make test-api         # Backend unit + integration
make test-web         # Frontend unit
make test-e2e         # End-to-end

# Code Quality
make lint             # Lint all
make fmt              # Format all

# Docs
make swagger          # Generate Swagger docs

# Docker
make docker-build     # Build all images
make docker-up        # Start production stack
make docker-down      # Stop all containers
```

---

## Summary

This architecture enforces separation of concerns at every level:

- **Domain layer** owns business rules with zero external dependencies
- **Use cases** orchestrate domain logic without knowing about HTTP or databases
- **Adapters** translate between external world and use cases
- **Infrastructure** provides concrete implementations swapped via interfaces

The monorepo keeps API and Web in sync while shared packages provide a single source of truth for constants and contracts.
