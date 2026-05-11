<script setup lang="ts">
import type { LegalEntityDetailsResponse } from '../../types/legalentity'
import LegalEntityDetailsView from '../../components/features/legalentities/LegalEntityDetailsView.vue'
import AppCard from '../../components/ui/AppCard.vue'

const route = useRoute()
const router = useRouter()
const notifications = useNotifications()
const id = computed(() => String(route.params.id || ''))

const details = ref<LegalEntityDetailsResponse | null>(null)
const files = ref<string[]>([])
const loading = ref(true)
const loadError = ref('')
const filesLoading = ref(false)
const downloadLoading = ref('')

const load = async () => {
  loading.value = true
  loadError.value = ''
  filesLoading.value = true

  try {
    const legalEntityDetails = await getLegalEntityDetails(id.value)
    details.value = legalEntityDetails
    files.value = legalEntityDetails.legal_entity.name_company
      ? await getLegalEntityFiles(legalEntityDetails.legal_entity.name_company)
      : []
  } catch (error) {
    loadError.value =
      error instanceof Error ? error.message : 'Не удалось загрузить юридическое лицо'
  } finally {
    loading.value = false
    filesLoading.value = false
  }
}

const downloadFile = async (fileName: string) => {
  downloadLoading.value = fileName

  try {
    const blob = await downloadLegalEntityFile(fileName)
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
      'Юридические лица'
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
    <AppCard v-if="loading" title="Загрузка">
      <p>Подтягиваю юридическое лицо.</p>
    </AppCard>

    <AppCard v-else-if="loadError" title="Ошибка">
      <p>{{ loadError }}</p>
    </AppCard>

    <LegalEntityDetailsView
      v-else-if="details"
      :details="details"
      :files="files"
      :files-loading="filesLoading"
      :download-loading="downloadLoading"
      @back="router.push('/legalentities')"
      @edit="router.push(`/legalentities/edit/${id}`)"
      @add-listener="router.push({ path: '/listeners/create', query: { id_legalentity: details.legal_entity.id_legalentity } })"
      @enroll="router.push(`/enrollment/yur/${id}`)"
      @open-listener="router.push(`/listeners/${$event}`)"
      @download="downloadFile"
    />
  </section>
</template>
