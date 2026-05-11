<script setup lang="ts">
import AppButton from '../components/ui/AppButton.vue'
import AppNotifications from '../components/ui/AppNotifications.vue'

const auth = useAuth()

const isAuthenticated = computed(() => auth.isAuthenticated.value)
const userName = computed(() => auth.state.value.username || 'Гость')
const roleLabel = computed(() => auth.roleLabel.value || 'Не авторизован')
</script>

<template>
  <div class="shell">
    <AppNotifications />

    <header class="shell__header">
      <div class="shell__brand">
        <h1 class="shell__title">Система управления курсами 25-12</h1>
      </div>

      <div class="shell__actions">
        <div class="shell__user">
          <span class="shell__user-label">Пользователь</span>
          <strong>{{ userName }}</strong>
          <span>{{ roleLabel }}</span>
        </div>

        <nav class="shell__nav" aria-label="Основная навигация">
          <AppButton v-if="!isAuthenticated" to="/login" variant="ghost">
            Вход
          </AppButton>
          <AppButton to="/dashboard" variant="ghost">
            Панель
          </AppButton>
          <AppButton v-if="isAuthenticated" to="/graphics" variant="ghost">
            Аналитика
          </AppButton>
          <AppButton v-if="isAuthenticated" variant="secondary" @click="auth.logout">
            Выйти
          </AppButton>
        </nav>
      </div>
    </header>

    <main class="shell__main">
      <slot />
    </main>
  </div>
</template>

<style scoped>
.shell {
  min-height: 100vh;
  background:
    radial-gradient(circle at top left, rgba(56, 189, 248, 0.16), transparent 24rem),
    radial-gradient(circle at top right, rgba(249, 115, 22, 0.14), transparent 20rem),
    linear-gradient(180deg, #f8fafc 0%, #eef2ff 100%);
}

.shell__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 2rem;
  padding: 1.5rem 2rem;
  border-bottom: 1px solid rgba(15, 23, 42, 0.08);
  backdrop-filter: blur(14px);
}

.shell__brand {
  max-width: 44rem;
}

.shell__eyebrow {
  margin: 0;
  font-size: 0.75rem;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: #475569;
}

.shell__title {
  margin: 0.25rem 0 0;
  font-size: clamp(1.5rem, 3vw, 2.2rem);
  color: #0f172a;
}

.shell__subtitle {
  margin: 0.5rem 0 0;
  color: #475569;
}

.shell__actions {
  display: flex;
  align-items: center;
  gap: 1rem;
  flex-wrap: wrap;
  justify-content: flex-end;
}

.shell__user {
  min-width: 12rem;
  padding: 0.75rem 0.95rem;
  border-radius: 1rem;
  background: rgba(255, 255, 255, 0.72);
  border: 1px solid rgba(15, 23, 42, 0.08);
  display: grid;
  gap: 0.15rem;
  color: #334155;
}

.shell__user-label {
  font-size: 0.72rem;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: #64748b;
}

.shell__nav {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
}

.shell__main {
  width: min(1440px, calc(100% - 2rem));
  margin: 0 auto;
  padding: 2rem 0 3rem;
}

@media (max-width: 720px) {
  .shell__header {
    align-items: flex-start;
    flex-direction: column;
    padding: 1.25rem 1rem;
  }

  .shell__actions {
    width: 100%;
    justify-content: stretch;
  }

  .shell__user {
    width: 100%;
  }

  .shell__main {
    width: min(100%, calc(100% - 1rem));
    padding-top: 1.25rem;
  }
}
</style>
