<script setup lang="ts">
import type {
  EnrollmentProgramDetails,
  ListenerDetailsResponse
} from '../../types/listener'
import ListenerDetailsView from '../../components/features/listeners/ListenerDetailsView.vue'
import AppButton from '../../components/ui/AppButton.vue'
import AppCard from '../../components/ui/AppCard.vue'

const route = useRoute()
const router = useRouter()
const notifications = useNotifications()
const id = computed(() => String(route.params.id || ''))

const details = ref<ListenerDetailsResponse | null>(null)
const enrollments = ref<EnrollmentProgramDetails[]>([])
const files = ref<string[]>([])
const loading = ref(true)
const filesLoading = ref(false)
const downloadLoading = ref('')
const errorMessage = ref('')

const load = async () => {
  loading.value = true
  errorMessage.value = ''
  files.value = []
  enrollments.value = []

  try {
    const listenerDetails = await getListenerDetails(id.value)
    details.value = listenerDetails

    filesLoading.value = Boolean(listenerDetails.listener.snils)

    const [enrollmentsResult, filesResult] = await Promise.allSettled([
      getListenerEnrollments(id.value),
      listenerDetails.listener.snils
        ? getListenerFiles(listenerDetails.listener.snils)
        : Promise.resolve([])
    ])

    enrollments.value = enrollmentsResult.status === 'fulfilled' ? enrollmentsResult.value : []
    files.value = filesResult.status === 'fulfilled' ? filesResult.value : []
  } catch (error) {
    errorMessage.value =
      error instanceof Error ? error.message : 'Не удалось загрузить данные слушателя'
  } finally {
    filesLoading.value = false
    loading.value = false
  }
}

const download = async (fileName: string) => {
  downloadLoading.value = fileName

  try {
    const blob = await downloadListenerFile(fileName)
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = fileName
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось скачать файл',
      'Файлы слушателя'
    )
  } finally {
    downloadLoading.value = ''
  }
}

onMounted(() => {
  void load()
})
</script>

<template>
  <section class="stack">
    <AppCard title="Профиль слушателя">
      <p>Здесь собраны личные данные, документы, образование, место работы, курсы и файлы.</p>
    </AppCard>

    <AppCard v-if="loading" title="Загрузка">
      <p>Подтягиваю карточку слушателя.</p>
    </AppCard>

    <AppCard v-else-if="errorMessage && !details" title="Ошибка">
      <p>{{ errorMessage }}</p>
      <AppButton variant="secondary" @click="load">Повторить запрос</AppButton>
    </AppCard>

    <template v-else-if="details">
      <AppCard v-if="errorMessage" title="Предупреждение">
        <p>{{ errorMessage }}</p>
      </AppCard>

      <ListenerDetailsView
        :details="details"
        :enrollments="enrollments"
        :files="files"
        :files-loading="filesLoading"
        :download-loading="downloadLoading"
        @back="router.push('/listeners')"
        @edit="router.push(`/listeners/edit/${id}`)"
        @enroll="router.push(`/enrollment/create/${id}`)"
        @enrollments="router.push(`/enrollment/details/${id}`)"
        @download="download"
      />
    </template>
  </section>
</template>
