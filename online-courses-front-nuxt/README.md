# online-courses-front-nuxt

Nuxt workspace for the gradual migration of `online-courses-front`.

## Goals

- keep the current Vite app running while the new frontend grows nearby
- move shared logic out of page components before large screen rewrites
- replace Bootstrap-first styling with a small design system and reusable UI layer

## Target structure

- `app/pages`: route screens grouped by domain
- `app/layouts`: app shells for auth and role-based areas
- `app/components/ui`: small reusable building blocks
- `app/components/features`: domain components assembled from UI pieces
- `app/composables`: shared stateful logic
- `app/services/api`: API clients and endpoint adapters
- `app/middleware`: auth and role guards
- `app/stores`: Pinia stores where shared app state truly belongs
- `docs`: migration notes and slice plan

## First migration slices

1. auth shell, login, route guards
2. dashboards by role
3. simple CRUD directories
4. lists/details with filters and pagination
5. heavy enrollment flows
