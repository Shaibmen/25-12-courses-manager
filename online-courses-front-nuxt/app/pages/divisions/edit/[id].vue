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
    const items = await getDivisions('')
    const current = items.find((item) => item.id_divisionsEducation === id.value)

    if (!current) {
      throw new Error('Подразделение не найдено')
    }

    name.value = current.divisions
  } catch (error) {
    errorMessage.value =
      error instanceof Error ? error.message : 'Не удалось загрузить подразделение'
  } finally {
    loading.value = false
  }
}

const submit = async () => {
  const value = name.value.trim()

  if (!value) {
    errorMessage.value = 'Введите название подразделения'
    return
  }

  saving.value = true
  errorMessage.value = ''

  try {
    await updateDivisionRequest(id.value, value)
    notifications.success('Изменения по подразделению сохранены.', 'Подразделения')
    await router.push('/divisions')
  } catch (error) {
    errorMessage.value =
      error instanceof Error ? error.message : 'Не удалось обновить подразделение'
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
      <p>Подтягиваю данные подразделения.</p>
    </AppCard>

    <AppCard v-else title="Редактировать подразделение">
      <form class="catalog-form" @submit.prevent="submit">
        <AppInput v-model="name" label="Название подразделения" placeholder="Введите название подразделения" />

        <p v-if="errorMessage" class="catalog-form__error">{{ errorMessage }}</p>

        <div class="catalog-form__actions">
          <AppButton type="button" variant="ghost" @click="router.push('/divisions')">Назад</AppButton>
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
