<script setup lang="ts">
import AppButton from '../../ui/AppButton.vue'
import AppCard from '../../ui/AppCard.vue'
import AppInput from '../../ui/AppInput.vue'

const auth = useAuth()

const username = ref('')
const password = ref('')
const pending = ref(false)
const errorMessage = ref('')

const submit = async () => {
  errorMessage.value = ''

  if (!username.value.trim() || !password.value.trim()) {
    errorMessage.value = 'Введите логин и пароль'
    return
  }

  pending.value = true

  try {
    await auth.login({
      username: username.value.trim(),
      password: password.value
    })

    await navigateTo('/dashboard')
  } catch (error) {
    errorMessage.value =
      error instanceof Error ? error.message : 'Не удалось выполнить вход'
  } finally {
    pending.value = false
  }
}
</script>

<template>
  <AppCard title="Вход в систему">
    <form class="login-form" @submit.prevent="submit">
      <AppInput
        v-model="username"
        label="Логин"
        placeholder="Введите логин"
        autocomplete="username"
      />

      <AppInput
        v-model="password"
        label="Пароль"
        type="password"
        placeholder="Введите пароль"
        autocomplete="current-password"
      />

      <p v-if="errorMessage" class="login-form__error">
        {{ errorMessage }}
      </p>

      <AppButton type="submit" :disabled="pending" block>
        {{ pending ? 'Входим...' : 'Войти' }}
      </AppButton>

      <p class="login-form__hint">
        Страница уже использует новый auth-слой: API вызов, хранение сессии и редиректы
        вынесены из компонента страницы.
      </p>
    </form>
  </AppCard>
</template>

<style scoped>
.login-form {
  display: grid;
  gap: 1rem;
}

.login-form__error {
  margin: 0;
  padding: 0.85rem 1rem;
  border-radius: 1rem;
  background: rgba(254, 226, 226, 0.9);
  color: #991b1b;
}

.login-form__hint {
  margin: 0;
  color: #64748b;
  font-size: 0.88rem;
}
</style>
