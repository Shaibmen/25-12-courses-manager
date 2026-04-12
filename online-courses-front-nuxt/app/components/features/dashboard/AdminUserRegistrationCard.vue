<script setup lang="ts">
import type { RegisterUserPayload } from '../../../types/auth'
import AppButton from '../../ui/AppButton.vue'
import AppCard from '../../ui/AppCard.vue'
import AppInput from '../../ui/AppInput.vue'

const props = defineProps<{
  roles: Array<{
    id: string | number
    role: string
  }>
}>()

const username = ref('')
const password = ref('')
const selectedRole = ref('')
const pending = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const roleOptions = computed(() =>
  props.roles.filter((role) => Boolean(role.id) && Boolean(role.role))
)

const submit = async () => {
  errorMessage.value = ''
  successMessage.value = ''

  const payload: RegisterUserPayload = {
    username: username.value.trim(),
    password: password.value,
    role: selectedRole.value
  }

  if (!payload.username || !payload.password || !payload.role) {
    errorMessage.value = 'Заполните логин, пароль и роль'
    return
  }

  pending.value = true

  try {
    const response = await registerUserRequest(payload)

    successMessage.value = response.message || 'Пользователь успешно зарегистрирован'
    username.value = ''
    password.value = ''
    selectedRole.value = ''
  } catch (error) {
    errorMessage.value =
      error instanceof Error ? error.message : 'Не удалось зарегистрировать пользователя'
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

        <label class="app-select">
          <span class="app-select__label">Роль</span>

          <div class="app-select__field" :class="{ 'app-select__field--placeholder': !selectedRole }">
            <select v-model="selectedRole" class="app-select__control">
              <option value="">Выберите роль</option>
              <option v-for="role in roleOptions" :key="role.id" :value="String(role.id)">
                {{ role.role }}
              </option>
            </select>

            <span class="app-select__icon" aria-hidden="true">
              <svg viewBox="0 0 20 20" fill="none">
                <path
                  d="M5 7.5L10 12.5L15 7.5"
                  stroke="currentColor"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="1.8"
                />
              </svg>
            </span>
          </div>
        </label>
      </div>

      <p v-if="errorMessage" class="registration-form__message registration-form__message--error">
        {{ errorMessage }}
      </p>

      <p v-if="successMessage" class="registration-form__message registration-form__message--success">
        {{ successMessage }}
      </p>

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

.registration-form__message {
  margin: 0;
  padding: 0.85rem 1rem;
  border-radius: 1rem;
  font-size: 0.92rem;
}

.registration-form__message--error {
  background: rgba(254, 226, 226, 0.9);
  color: #991b1b;
}

.registration-form__message--success {
  background: rgba(220, 252, 231, 0.92);
  color: #166534;
}

.app-select {
  display: grid;
  gap: 0.45rem;
}

.app-select__label {
  font-size: 0.92rem;
  font-weight: 600;
  color: #0f172a;
}

.app-select__field {
  position: relative;
  border-radius: 1rem;
  border: 1px solid rgba(15, 23, 42, 0.12);
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.97) 0%, rgba(248, 250, 252, 0.94) 100%);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.7);
  transition:
    border-color 180ms ease,
    box-shadow 180ms ease,
    transform 180ms ease;
}

.app-select__control {
  appearance: none;
  -webkit-appearance: none;
  position: relative;
  z-index: 1;
  width: 100%;
  min-height: 3rem;
  padding: 0.85rem 3rem 0.85rem 1rem;
  border: none;
  border-radius: 1rem;
  background: transparent;
  color: #0f172a;
  font: inherit;
  cursor: pointer;
}

.app-select__field:focus-within {
  outline: none;
  border-color: rgba(59, 130, 246, 0.5);
  box-shadow: 0 0 0 5px rgba(59, 130, 246, 0.12);
  transform: translateY(-1px);
}

.app-select__field--placeholder .app-select__control {
  color: #64748b;
}

.app-select__icon {
  position: absolute;
  top: 50%;
  right: 1rem;
  z-index: 0;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 1.1rem;
  height: 1.1rem;
  color: #475569;
  pointer-events: none;
  transform: translateY(-50%);
}

.app-select__icon svg {
  width: 100%;
  height: 100%;
}

@media (max-width: 900px) {
  .registration-form__grid {
    grid-template-columns: 1fr;
  }
}
</style>
