<script setup lang="ts">
import AppButton from '../../components/ui/AppButton.vue'
import AppCard from '../../components/ui/AppCard.vue'
import AppInput from '../../components/ui/AppInput.vue'

const router = useRouter()
const notifications = useNotifications()
const name = ref('')
const loading = ref(false)
const errorMessage = ref('')

const submit = async () => {
  const value = name.value.trim()

  if (!value) {
    errorMessage.value = 'Введите название типа обучения'
    return
  }

  loading.value = true
  errorMessage.value = ''

  try {
    await createEducationTypeRequest(value)
    notifications.success('Тип обучения создан.', 'Типы обучения')
    await router.push('/types')
  } catch (error) {
    errorMessage.value =
      error instanceof Error ? error.message : 'Не удалось создать тип обучения'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <section class="stack">
    <AppCard title="Добавить тип обучения">
      <form class="catalog-form" @submit.prevent="submit">
        <AppInput v-model="name" label="Название типа обучения" placeholder="Например: Индивидуально с преподавателем" />

        <p v-if="errorMessage" class="catalog-form__error">{{ errorMessage }}</p>

        <div class="catalog-form__actions">
          <AppButton type="button" variant="ghost" @click="router.push('/types')">Назад</AppButton>
          <AppButton type="submit" :disabled="loading">
            {{ loading ? 'Сохраняем...' : 'Создать' }}
          </AppButton>
        </div>
      </form>
    </AppCard>
  </section>
</template>

<style scoped>
.catalog-form {
  display: grid;
  gap: 1rem;
}

.catalog-form__error {
  margin: 0;
  padding: 0.85rem 1rem;
  border-radius: 1rem;
  background: rgba(254, 226, 226, 0.9);
  color: #991b1b;
}

.catalog-form__actions {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
}
</style>
