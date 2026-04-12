<script setup lang="ts">
import AppButton from '../../../components/ui/AppButton.vue'
import AppCard from '../../../components/ui/AppCard.vue'
import AppInput from '../../../components/ui/AppInput.vue'

const route = useRoute()
const router = useRouter()
const notifications = useNotifications()
const id = computed(() => String(route.params.id || ''))
const name = ref('')
const loading = ref(true)
const saving = ref(false)
const errorMessage = ref('')

const load = async () => {
  loading.value = true
  errorMessage.value = ''

  try {
    const items = await getEducationTypes('')
    const current = items.find((item) => item.id_educationType === id.value)

    if (!current) {
      throw new Error('Тип обучения не найден')
    }

    name.value = current.typeName
  } catch (error) {
    errorMessage.value =
      error instanceof Error ? error.message : 'Не удалось загрузить тип обучения'
  } finally {
    loading.value = false
  }
}

const submit = async () => {
  const value = name.value.trim()

  if (!value) {
    errorMessage.value = 'Введите название типа обучения'
    return
  }

  saving.value = true
  errorMessage.value = ''

  try {
    await updateEducationTypeRequest(id.value, value)
    notifications.success('Изменения по типу обучения сохранены.', 'Типы обучения')
    await router.push('/types')
  } catch (error) {
    errorMessage.value =
      error instanceof Error ? error.message : 'Не удалось обновить тип обучения'
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  void load()
})
</script>

<template>
  <section class="stack">
    <AppCard v-if="loading" title="Загрузка">
      <p>Подтягиваю данные типа обучения.</p>
    </AppCard>

    <AppCard v-else title="Редактировать тип обучения">
      <form class="catalog-form" @submit.prevent="submit">
        <AppInput v-model="name" label="Название типа обучения" placeholder="Введите новое название" />

        <p v-if="errorMessage" class="catalog-form__error">{{ errorMessage }}</p>

        <div class="catalog-form__actions">
          <AppButton type="button" variant="ghost" @click="router.push('/types')">Назад</AppButton>
          <AppButton type="submit" :disabled="saving">
            {{ saving ? 'Сохраняем...' : 'Сохранить изменения' }}
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
