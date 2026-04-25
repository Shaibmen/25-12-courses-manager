<script setup lang="ts">
import type { DivisionItem, EducationTypeItem } from '../../types/catalogs'
import type { ProgramPayload } from '../../types/program'
import ProgramForm from '../../components/features/programs/ProgramForm.vue'
import AppCard from '../../components/ui/AppCard.vue'

const router = useRouter()
const notifications = useNotifications()
const educationTypes = ref<EducationTypeItem[]>([])
const divisions = ref<DivisionItem[]>([])
const loading = ref(true)
const saving = ref(false)
const loadError = ref('')

const loadCatalogs = async () => {
  loading.value = true
  loadError.value = ''

  try {
    const [typeData, divisionData] = await Promise.all([getEducationTypes(''), getDivisions('')])
    educationTypes.value = typeData
    divisions.value = divisionData
  } catch (error) {
    loadError.value =
      error instanceof Error ? error.message : 'Не удалось загрузить справочники программы'
  } finally {
    loading.value = false
  }
}

const submit = async (payload: ProgramPayload) => {
  saving.value = true

  try {
    await createProgramRequest(payload)
    notifications.success('Программа обучения создана.', 'Программы обучения')
    await router.push('/programs')
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось создать программу обучения',
      'Программы обучения'
    )
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  void loadCatalogs()
})
</script>

<template>
  <section class="stack content-shell">
    <AppCard v-if="loading" title="Загрузка">
      <p>Подтягиваю типы обучения и подразделения.</p>
    </AppCard>

    <AppCard v-else-if="loadError" title="Ошибка">
      <p>{{ loadError }}</p>
    </AppCard>

    <ProgramForm
      v-else
      title="Добавить программу обучения"
      submit-label="Создать"
      :loading="saving"
      :education-types="educationTypes"
      :divisions="divisions"
      @submit="submit"
      @cancel="router.push('/programs')"
    />
  </section>
</template>
