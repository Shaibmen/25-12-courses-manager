<script setup lang="ts">
import type {
  EnrollmentProgramDetails,
  ListenerDetailsResponse
} from '../../types/listener'
import {
  buildListenerScanDiplomBaseName,
  deleteListenerScanDiplomFile,
  downloadListenerFile,
  downloadListenerScanDiplomFile,
  getListenerDetails,
  getListenerEnrollments,
  getListenerFiles,
  getListenerScanDiplomFiles,
  uploadListenerScanDiplomFile
} from '../../services/api/listeners'
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
const scanDiplomFiles = ref<string[]>([])
const loading = ref(true)
const filesLoading = ref(false)
const scanDiplomLoading = ref(false)
const downloadLoading = ref('')
const scanDownloadLoading = ref('')
const scanDeleteLoading = ref('')
const scanUploadLoading = ref(false)
const errorMessage = ref('')
const scanUploadInput = ref<HTMLInputElement | null>(null)

const getScanDiplomBaseName = () => {
  if (!details.value) {
    return ''
  }

  return buildListenerScanDiplomBaseName(
    details.value.listener.second_name,
    details.value.listener.first_name,
    details.value.listener.middle_name
  )
}

const getNextScanDiplomName = () => {
  const baseName = getScanDiplomBaseName()

  if (!baseName) {
    return ''
  }

  const escapedBase = baseName.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  const matcher = new RegExp(`^${escapedBase}_скан_(\\d+)(\\.pdf)?$`, 'i')
  const numbers = scanDiplomFiles.value
    .map(fileName => {
      const match = fileName.match(matcher)
      return match ? Number(match[1]) : 0
    })
    .filter(value => value > 0)

  const nextNumber = numbers.length ? Math.max(...numbers) + 1 : 1

  return `${baseName}_скан_${nextNumber}`
}

const loadScanDiplomFiles = async () => {
  const baseName = getScanDiplomBaseName()

  scanDiplomLoading.value = Boolean(baseName)

  try {
    scanDiplomFiles.value = baseName
      ? await getListenerScanDiplomFiles(baseName)
      : []
  } finally {
    scanDiplomLoading.value = false
  }
}

const load = async () => {
  loading.value = true
  errorMessage.value = ''
  files.value = []
  scanDiplomFiles.value = []
  enrollments.value = []

  try {
    const listenerDetails = await getListenerDetails(id.value)
    details.value = listenerDetails

    const scanBaseName = buildListenerScanDiplomBaseName(
      listenerDetails.listener.second_name,
      listenerDetails.listener.first_name,
      listenerDetails.listener.middle_name
    )

    filesLoading.value = Boolean(listenerDetails.listener.snils)
    scanDiplomLoading.value = Boolean(scanBaseName)

    const [enrollmentsResult, filesResult, scanDiplomResult] = await Promise.allSettled([
      getListenerEnrollments(id.value),
      listenerDetails.listener.snils
        ? getListenerFiles(listenerDetails.listener.snils)
        : Promise.resolve([]),
      scanBaseName
        ? getListenerScanDiplomFiles(scanBaseName)
        : Promise.resolve([])
    ])

    enrollments.value = enrollmentsResult.status === 'fulfilled' ? enrollmentsResult.value : []
    files.value = filesResult.status === 'fulfilled' ? filesResult.value : []
    scanDiplomFiles.value = scanDiplomResult.status === 'fulfilled' ? scanDiplomResult.value : []
  } catch (error) {
    errorMessage.value =
      error instanceof Error ? error.message : 'Не удалось загрузить данные слушателя'
  } finally {
    filesLoading.value = false
    scanDiplomLoading.value = false
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

const downloadScanDiplom = async (fileName: string) => {
  scanDownloadLoading.value = fileName

  try {
    const blob = await downloadListenerScanDiplomFile(fileName)
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
      error instanceof Error ? error.message : 'Не удалось скачать скан диплома',
      'Скан диплома'
    )
  } finally {
    scanDownloadLoading.value = ''
  }
}

const deleteScanDiplom = async (fileName: string) => {
  scanDeleteLoading.value = fileName

  try {
    await deleteListenerScanDiplomFile(fileName)
    await loadScanDiplomFiles()
    notifications.success(`Скан "${fileName}" удалён`, 'Скан диплома')
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось удалить скан диплома',
      'Скан диплома'
    )
  } finally {
    scanDeleteLoading.value = ''
  }
}

const triggerScanDiplomUpload = () => {
  scanUploadInput.value?.click()
}

const handleScanDiplomUpload = async (event: Event) => {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]

  if (!file) {
    return
  }

  const isPdf = file.type === 'application/pdf' || file.name.toLowerCase().endsWith('.pdf')

  if (!isPdf) {
    notifications.error('Разрешена загрузка только PDF-файлов', 'Скан диплома')
    input.value = ''
    return
  }

  const scanName = getNextScanDiplomName()

  if (!scanName) {
    notifications.error('Не удалось сформировать имя скана диплома', 'Скан диплома')
    input.value = ''
    return
  }

  scanUploadLoading.value = true

  try {
    await uploadListenerScanDiplomFile(scanName, file)
    await loadScanDiplomFiles()
    notifications.success(`Скан "${scanName}" загружен`, 'Скан диплома')
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось загрузить скан диплома',
      'Скан диплома'
    )
  } finally {
    scanUploadLoading.value = false
    input.value = ''
  }
}

onMounted(() => {
  void load()
})
</script>

<template>
  <section class="stack">
    <input
      ref="scanUploadInput"
      type="file"
      accept="application/pdf,.pdf"
      class="hidden-upload"
      @change="handleScanDiplomUpload"
    >

    <AppCard title="Профиль слушателя">
      <p>Здесь собраны личные данные, документы, образование, место работы, курсы, файлы и сканы диплома.</p>
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
        :scan-diplom-files="scanDiplomFiles"
        :files-loading="filesLoading"
        :scan-diplom-loading="scanDiplomLoading"
        :download-loading="downloadLoading"
        :scan-download-loading="scanDownloadLoading"
        :scan-delete-loading="scanDeleteLoading"
        :scan-upload-loading="scanUploadLoading"
        @back="router.push('/listeners')"
        @edit="router.push(`/listeners/edit/${id}`)"
        @enroll="router.push(`/enrollment/create/${id}`)"
        @enrollments="router.push(`/enrollment/details/${id}`)"
        @download="download"
        @download-scan-diplom="downloadScanDiplom"
        @delete-scan-diplom="deleteScanDiplom"
        @upload-scan-diplom="triggerScanDiplomUpload"
      />
    </template>
  </section>
</template>

<style scoped>
.hidden-upload {
  display: none;
}
</style>
