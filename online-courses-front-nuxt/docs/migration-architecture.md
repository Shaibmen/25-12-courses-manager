# Migration Architecture

## Why a parallel Nuxt app

The current frontend has two separate problems:

1. the visual layer is tightly coupled to Bootstrap classes and many inline styles
2. page components contain API calls, auth state, formatting, validation, and navigation all at once

A parallel Nuxt workspace lets us improve structure without freezing the existing app.

## Architecture rules

### 1. Pages stay thin

Pages should compose feature blocks and trigger loading, but they should not contain raw API wiring for every action.

### 2. API access is centralized

All network access should move toward `app/services/api`.

Examples:

- `auth/login`
- `dashboard/getWorkerStats`
- `listeners/list`
- `listeners/remove`

### 3. Auth logic is not page-owned

The old app spreads auth concerns between `Login.vue`, `Header.vue`, and the router guard.

The new app should move auth into:

- `composables/useAuthState.ts`
- `services/api/client.ts`
- `middleware/auth.global.ts`
- later: `stores/auth.ts` if richer shared flows are needed

### 4. UI is not tied to business entities

Reusable pieces live in `components/ui`.

Examples:

- `AppButton`
- `AppInput`
- `AppSelect`
- `AppTable`
- `AppModal`
- `AppStatCard`

Entity-specific blocks live in `components/features`.

Examples:

- `features/auth/LoginForm`
- `features/dashboard/WorkerStats`
- `features/listeners/ListenersTable`

### 5. Migration order follows risk

Move low-risk screens first.

Recommended sequence:

1. login, auth shell, middleware
2. admin/worker/accountant dashboards
3. simple dictionaries: divisions, education types, levels, executers
4. medium CRUD screens: listeners, programs, legal entities
5. heavy flows: enrollment create/edit/yur and document generation

## Mapping from the old app

### Old source hotspots

- `src/components/Login.vue`
- `src/components/Header.vue`
- `src/router/index.js`
- `src/components/UserDashboard.vue`
- `src/components/AdminDashboard.vue`
- `src/components/EnrollmentCreate.vue`
- `src/components/EnrollmentYUR.vue`

### New home

- auth concerns -> `app/composables`, `app/middleware`, `app/services/api`
- shell/layout concerns -> `app/layouts`
- route screens -> `app/pages`
- shared visuals -> `app/components/ui`
- entity workflows -> `app/components/features`

## Exit criteria for stage 1

- Nuxt project exists next to the Vite app
- folder structure is fixed
- runtime config uses environment variables instead of hardcoded localhost-only constants
- auth and API direction are captured in code placeholders
- migration slice order is documented
