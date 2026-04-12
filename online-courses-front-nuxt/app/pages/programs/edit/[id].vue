<script setup lang="ts">
import type { DivisionItem, EducationTypeItem } from '../../../types/catalogs'
import type { ProgramFormState, ProgramUpdatePayload } from '../../../types/program'
import { createEmptyProgramFormState } from '../../../types/program'
import ProgramForm from '../../../components/features/programs/ProgramForm.vue'
import AppCard from '../../../components/ui/AppCard.vue'

const route = useRoute()
const router = useRouter()
const notifications = useNotifications()
const id = computed(() => String(route.params.id || ''))
const educationTypes = ref<EducationTypeItem[]>([])
const divisions = ref<DivisionItem[]>([])
const initialState = ref<ProgramFormState | null>(null)
const loading = ref(true)
const saving = ref(false)
const loadError = ref('')

const load = async () => {
  loading.value = true
  loadError.value = ''

  try {
    const [typeData, divisionData, program] = await Promise.all([
      getEducationTypes(''),
      getDivisions(''),
      getProgramDetails(id.value)
    ])

    educationTypes.value = typeData
    divisions.value = divisionData
    initialState.value = {
      ...createEmptyProgramFormState(),
      name_prof_education: program.name_prof_education || '',
      time_education: String(program.time_education ?? ''),
      individual_price: String(program.individual_price ?? ''),
      group_price: String(program.group_price ?? ''),
      campus_price: String(program.campus_price ?? ''),
      id_education_type: program.id_education_type || '',
      id_divisions_education: program.id_divisions_education || ''
    }
  } catch (error) {
    loadError.value =
      error instanceof Error ? error.message : 'Не удалось загрузить программу обучения'
  } finally {
    loading.value = false
  }
}

const submit = async (payload: ProgramUpdatePayload) => {
  saving.value = true

  try {
    await updateProgramRequest(id.value, payload)
    notifications.success('Программа обучения обновлена.', 'Программы обучения')
    await router.push('/programs')
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось обновить программу обучения',
      'Программы обучения'
    )
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
      <p>Подтягиваю программу обучения и связанные справочники.</p>
    </AppCard>

    <AppCard v-else-if="loadError" title="Ошибка">
      <p>{{ loadError }}</p>
    </AppCard>

    <ProgramForm
      v-else
      title="Редактировать программу обучения"
      submit-label="Сохранить изменения"
      :loading="saving"
      :initial-state="initialState"
      :education-types="educationTypes"
      :divisions="divisions"
      mode="edit"
      @submit="submit"
      @cancel="router.push('/programs')"
    />
  </section>
</template>
