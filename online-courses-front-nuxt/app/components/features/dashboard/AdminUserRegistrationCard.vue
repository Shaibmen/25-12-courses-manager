<script setup lang="ts">
import type { RegisterUserPayload } from '../../../types/auth'
import AppButton from '../../ui/AppButton.vue'
import AppCard from '../../ui/AppCard.vue'
import AppInput from '../../ui/AppInput.vue'
import AppSelect from '../../ui/AppSelect.vue'

const props = defineProps<{
  roles: Array<{
    id: string | number
    role: string
  }>
}>()

const notifications = useNotifications()
const username = ref('')
const password = ref('')
const selectedRole = ref('')
const pending = ref(false)

const roleOptions = computed(() =>
  props.roles.filter((role) => Boolean(role.id) && Boolean(role.role))
)

const submit = async () => {
  const payload: RegisterUserPayload = {
    username: username.value.trim(),
    password: password.value,
    role: selectedRole.value
  }

  if (!payload.username || !payload.password || !payload.role) {
    notifications.error('Заполните логин, пароль и роль', 'Регистрация пользователя')
    return
  }

  pending.value = true

  try {
    const response = await registerUserRequest(payload)

    notifications.success(
      response.message || 'Пользователь успешно зарегистрирован',
      'Регистрация пользователя'
    )
    username.value = ''
    password.value = ''
    selectedRole.value = ''
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось зарегистрировать пользователя',
      'Регистрация пользователя'
    )
  } finally {
    pending.value = false
  }
}
</script>

<template>
  <AppCard title="Регистрация нового пользователя">
    <form class="registration-form" @submit.prevent="submit">
      <div class="registration-form__grid">
        <AppInput
          v-model="username"
          label="Логин"
          placeholder="Введите имя пользователя"
          autocomplete="username"
        />

        <AppInput
          v-model="password"
          label="Пароль"
          type="password"
          placeholder="Введите пароль"
          autocomplete="new-password"
        />

        <AppSelect v-model="selectedRole" label="Роль" placeholder="Выберите роль">
          <option v-for="role in roleOptions" :key="role.id" :value="String(role.id)">
            {{ role.role }}
          </option>
        </AppSelect>
      </div>

      <div class="registration-form__actions">
        <AppButton type="submit" :disabled="pending">
          {{ pending ? 'Сохраняем...' : 'Зарегистрировать' }}
        </AppButton>
      </div>
    </form>
  </AppCard>
</template>

<style scoped>
.registration-form {
  display: grid;
  gap: 1rem;
}

.registration-form__grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 1rem;
}

.registration-form__actions {
  display: flex;
  justify-content: flex-start;
}

@media (max-width: 900px) {
  .registration-form__grid {
    grid-template-columns: 1fr;
  }
}
</style>
