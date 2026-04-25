<script setup lang="ts">
import type { DivisionItem, EducationTypeItem } from '../../../types/catalogs'
import type { ContractorPayload } from '../../../types/enrollment'
import type { GroupItem } from '../../../types/group'
import {
  ageCategories,
  loadVariantsDO,
  loadVariantsNotDO,
  optDocumentOptions,
  singleEnrollmentContracts
} from '../../../types/enrollment'
import ContractorForm from '../../../components/features/enrollments/ContractorForm.vue'
import GroupSearchField from '../../../components/features/enrollments/GroupSearchField.vue'
import AppButton from '../../../components/ui/AppButton.vue'
import AppCard from '../../../components/ui/AppCard.vue'
import AppCheckbox from '../../../components/ui/AppCheckbox.vue'
import AppConfirmDialog from '../../../components/ui/AppConfirmDialog.vue'
import AppInput from '../../../components/ui/AppInput.vue'
import AppSelect from '../../../components/ui/AppSelect.vue'

const route = useRoute()
const router = useRouter()
const notifications = useNotifications()
const listenerId = computed(() => String(route.params.listenerId || ''))

const listenerContext = ref<Awaited<ReturnType<typeof getListenerEnrollmentContext>> | null>(null)
const programs = ref<Awaited<ReturnType<typeof getPrograms>>>([])
const executers = ref<Awaited<ReturnType<typeof getExecuters>>>([])
const divisions = ref<DivisionItem[]>([])
const educationTypes = ref<EducationTypeItem[]>([])
const groupOptions = ref<GroupItem[]>([])
const loading = ref(true)
const saving = ref(false)
const contractorSaving = ref(false)
const groupsLoading = ref(false)
const showContractorForm = ref(false)
const confirmContractorDelete = ref(false)
const contractorEnabled = ref(false)

const selectedDivisionId = ref('')
const selectedEducationTypeId = ref('')
const selectedContractId = ref('')
const selectedExecutorId = ref('')
const selectedProgramId = ref('')
const selectedGroupId = ref('')
const groupSearch = ref('')
const startDate = ref('')
const endDate = ref('')
const typeOfRetraining = ref('')
const paymentOption = ref('')
const secondPaymentDate = ref('')
const ageCategory = ref('')
const loadVariant = ref('')
const studyLoadOption = ref('')
const optDocumentSelected = ref('')

let groupSearchTimer: ReturnType<typeof setTimeout> | null = null

const filteredContracts = computed(() =>
  singleEnrollmentContracts.filter((item) => item.type === (contractorEnabled.value ? 'trilateral' : 'bilateral'))
)

const filteredPrograms = computed(() =>
  programs.value.filter((item) =>
    item.id_divisions_education === selectedDivisionId.value &&
    item.id_education_type === selectedEducationTypeId.value
  )
)

const selectedProgram = computed(() =>
  filteredPrograms.value.find((item) => item.id_program_education === selectedProgramId.value) || null
)

const selectedGroup = computed(() =>
  groupOptions.value.find((item) => item.group === selectedGroupId.value) || null
)

const isDO = computed(() => selectedContractId.value.toUpperCase().startsWith('DO'))
const isPKorPP = computed(() => {
  const id = selectedContractId.value.toUpperCase()
  return id.startsWith('PK') || id.startsWith('PP')
})

const isFormValid = computed(() =>
  Boolean(
    selectedDivisionId.value &&
    selectedEducationTypeId.value &&
    selectedProgramId.value &&
    selectedGroupId.value &&
    startDate.value &&
    endDate.value &&
    selectedExecutorId.value &&
    typeOfRetraining.value &&
    (!contractorEnabled.value || selectedContractId.value)
  )
)

const syncSelectedGroup = (groups: GroupItem[]) => {
  if (!selectedGroupId.value || groups.some((item) => item.group === selectedGroupId.value)) {
    return groups
  }

  if (!groupSearch.value.trim()) {
    selectedGroupId.value = ''
    return groups
  }

  return [
    {
      group: selectedGroupId.value,
      name_group: groupSearch.value,
      rapspisanie: []
    },
    ...groups
  ]
}

const loadGroups = async (filter: string) => {
  groupsLoading.value = true

  try {
    const groups = await getGroups(1, filter.trim())
    groupOptions.value = syncSelectedGroup(groups)
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось загрузить список групп',
      'Запись на курс'
    )
  } finally {
    groupsLoading.value = false
  }
}

const getPaymentText = () => {
  if (paymentOption.value === 'split') {
    return `Оплата в два этапа, вторая часть${secondPaymentDate.value ? ` до ${secondPaymentDate.value}` : ''}`
  }

  if (paymentOption.value === 'full') {
    return 'Полная предоплата до начала обучения'
  }

  if (paymentOption.value === 'halfsplit') {
    return 'Оплата после подписания акта'
  }

  return null
}

const getNagruzValue = () => {
  if (isDO.value) {
    return loadVariant.value ? Number(loadVariant.value) : null
  }

  if (isPKorPP.value) {
    return studyLoadOption.value ? Number(studyLoadOption.value) : null
  }

  return null
}

const load = async () => {
  loading.value = true

  try {
    const [listenerData, programData, executerData, divisionData, educationTypeData] = await Promise.all([
      getListenerEnrollmentContext(listenerId.value),
      getPrograms(1, ''),
      getExecuters(1, ''),
      getDivisions(''),
      getEducationTypes('')
    ])

    listenerContext.value = listenerData
    programs.value = programData
    executers.value = executerData
    divisions.value = divisionData
    educationTypes.value = educationTypeData
    selectedExecutorId.value = executerData[0]?.id_executor || ''
    await loadGroups('')
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось открыть форму записи на курс',
      'Запись на курс'
    )
  } finally {
    loading.value = false
  }
}

const saveContractor = async (payload: ContractorPayload) => {
  contractorSaving.value = true

  try {
    await upsertContractorRequest(listenerId.value, payload)
    notifications.success('Заказчик сохранён.', 'Запись на курс')
    showContractorForm.value = false
    contractorEnabled.value = true
    await load()
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось сохранить заказчика',
      'Запись на курс'
    )
  } finally {
    contractorSaving.value = false
  }
}

const removeContractor = async () => {
  const contractorId = listenerContext.value?.contractor?.contractor.id_contractor

  if (!contractorId) {
    return
  }

  contractorSaving.value = true

  try {
    await deleteContractorRequest(contractorId)
    confirmContractorDelete.value = false
    contractorEnabled.value = false
    showContractorForm.value = false
    notifications.success('Заказчик удалён.', 'Запись на курс')
    await load()
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось удалить заказчика',
      'Запись на курс'
    )
  } finally {
    contractorSaving.value = false
  }
}

const selectGroup = (group: GroupItem) => {
  selectedGroupId.value = group.group
  groupSearch.value = group.name_group
}

const createEnrollment = async () => {
  if (!isFormValid.value || !listenerContext.value || !selectedProgram.value) {
    notifications.error('Заполните обязательные поля записи.', 'Запись на курс')
    return
  }

  saving.value = true

  try {
    const optNagruz = getNagruzValue()

    await createEnrollmentRequest({
      id_listener: listenerId.value,
      id_program: selectedProgramId.value,
      start_date: startDate.value,
      end_date: endDate.value,
      id_group: selectedGroupId.value,
      type_of_retraining: typeOfRetraining.value,
      is_active: true
    })

    await createEnrollmentDocumentRequest({
      id_listener: listenerId.value,
      id_program: selectedProgramId.value,
      id_executor: selectedExecutorId.value || null,
      front_data: {
        dogovor_type: selectedContractId.value || null,
        dogovor_age: ageCategory.value || null,
        opt_document: optDocumentSelected.value ? Number(optDocumentSelected.value) : null,
        opt_price: getPaymentText(),
        opt_nagruz: optNagruz
      }
    })

    notifications.success('Запись и документы успешно созданы.', 'Запись на курс')
    await router.push(`/enrollment/details/${listenerId.value}`)
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось создать запись на курс',
      'Запись на курс'
    )
  } finally {
    saving.value = false
  }
}

watch(selectedContractId, (value) => {
  const id = value.toUpperCase()

  if (id.startsWith('DO')) {
    typeOfRetraining.value = 'Дополнительное образование'
  } else if (id.startsWith('PK')) {
    typeOfRetraining.value = 'Повышение квалификации'
  } else if (id.startsWith('PP')) {
    typeOfRetraining.value = 'Профессиональная переподготовка'
  } else {
    typeOfRetraining.value = ''
  }
})

watch(filteredContracts, (items) => {
  if (selectedContractId.value && !items.some((item) => item.id_contract === selectedContractId.value)) {
    selectedContractId.value = ''
  }
}, { immediate: true })

watch(filteredPrograms, (items) => {
  if (selectedProgramId.value && !items.some((item) => item.id_program_education === selectedProgramId.value)) {
    selectedProgramId.value = ''
  }
}, { immediate: true })

watch(selectedProgramId, () => {
  selectedGroupId.value = ''
  groupSearch.value = ''
})

watch(groupSearch, (value) => {
  if (groupSearchTimer) {
    clearTimeout(groupSearchTimer)
  }

  if (selectedGroupId.value && value.trim() !== (selectedGroup.value?.name_group || '')) {
    selectedGroupId.value = ''
  }

  groupSearchTimer = setTimeout(() => {
    void loadGroups(value)
  }, 300)
})

watch(contractorEnabled, (enabled) => {
  showContractorForm.value = enabled ? showContractorForm.value : false
})

onBeforeUnmount(() => {
  if (groupSearchTimer) {
    clearTimeout(groupSearchTimer)
  }
})

onMounted(() => {
  void load()
})
</script>

<template>
  <section class="stack content-shell">
    <AppCard v-if="loading" title="Загрузка">
      <p>Подтягиваю данные для записи на курс.</p>
    </AppCard>

    <template v-else-if="listenerContext">
      <AppCard title="Слушатель">
        <p class="summary-line">
          {{ listenerContext.listener.second_name }} {{ listenerContext.listener.first_name }} {{ listenerContext.listener.middle_name || '' }}
        </p>
      </AppCard>

      <AppCard title="Заказчик">
        <div class="summary-stack">
          <AppCheckbox v-model="contractorEnabled" label="Есть заказчик" />

          <template v-if="contractorEnabled">
            <div v-if="listenerContext.contractor" class="summary-stack">
              <p class="summary-line">
                {{ listenerContext.contractor.contractor.second_name }} {{ listenerContext.contractor.contractor.first_name }} {{ listenerContext.contractor.contractor.middle_name }}
              </p>
              <p class="summary-subline">
                {{ listenerContext.contractor.contractor.contact_phone }} · {{ listenerContext.contractor.contractor.email }}
              </p>
              <div class="summary-actions">
                <AppButton variant="secondary" @click="showContractorForm = !showContractorForm">
                  {{ showContractorForm ? 'Скрыть форму' : 'Изменить заказчика' }}
                </AppButton>
                <AppButton variant="ghost" @click="confirmContractorDelete = true">Удалить заказчика</AppButton>
              </div>
            </div>

            <div v-else class="summary-actions">
              <p class="summary-subline">Заказчик пока не добавлен.</p>
              <AppButton @click="showContractorForm = !showContractorForm">
                {{ showContractorForm ? 'Скрыть форму' : 'Добавить заказчика' }}
              </AppButton>
            </div>
          </template>

          <p v-else class="summary-subline">
            Без галочки данные заказчика в запись не отправляем.
          </p>
        </div>
      </AppCard>

      <ContractorForm
        v-if="contractorEnabled && showContractorForm"
        :initial-state="listenerContext.contractor || undefined"
        :loading="contractorSaving"
        @submit="saveContractor"
        @cancel="showContractorForm = false"
      />

      <AppCard title="Параметры договора">
        <div class="enrollment-grid">
          <AppSelect v-model="selectedContractId" label="Договор" placeholder="Выберите договор">
            <option v-for="item in filteredContracts" :key="item.id_contract" :value="item.id_contract">
              {{ item.name }}
            </option>
          </AppSelect>

          <AppSelect v-model="paymentOption" label="Порядок оплаты" placeholder="Выберите порядок оплаты">
            <option value="full">100% предоплата</option>
            <option value="split">50/50</option>
            <option value="halfsplit">Оплата после акта</option>
          </AppSelect>

          <AppInput v-if="paymentOption === 'split'" v-model="secondPaymentDate" label="Срок второй оплаты" type="date" />

          <AppSelect v-model="ageCategory" label="Возрастная категория" placeholder="Выберите категорию">
            <option v-for="(label, key) in ageCategories" :key="key" :value="key">{{ label }}</option>
          </AppSelect>

          <AppSelect v-model="optDocumentSelected" label="Итоговый документ" placeholder="Выберите режим выдачи">
            <option v-for="(label, key) in optDocumentOptions" :key="key" :value="String(key)">{{ label }}</option>
          </AppSelect>

          <AppSelect v-model="loadVariant" label="Нагрузка ДО" placeholder="Выберите вариант" :disabled="!isDO">
            <option v-for="(label, key) in loadVariantsDO" :key="key" :value="String(key)">{{ label }}</option>
          </AppSelect>

          <AppSelect v-model="studyLoadOption" label="Нагрузка ПК/ПП" placeholder="Выберите вариант" :disabled="!isPKorPP">
            <option v-for="(label, key) in loadVariantsNotDO" :key="key" :value="String(key)">{{ label }}</option>
          </AppSelect>
        </div>
      </AppCard>

      <AppCard title="Исполнитель">
        <div class="radio-grid">
          <label v-for="executer in executers" :key="executer.id_executor" class="radio-card">
            <input v-model="selectedExecutorId" type="radio" :value="executer.id_executor">
            <span>{{ executer.second_name }} {{ executer.first_name }} {{ executer.middle_name || '' }}</span>
          </label>
        </div>
      </AppCard>

      <AppCard title="Программа обучения">
        <div class="enrollment-grid">
          <AppSelect v-model="selectedDivisionId" label="Подразделение" placeholder="Сначала выберите подразделение">
            <option v-for="division in divisions" :key="division.id_divisionsEducation" :value="division.id_divisionsEducation">
              {{ division.divisions }}
            </option>
          </AppSelect>

          <AppSelect
            v-model="selectedEducationTypeId"
            label="Тип обучения"
            placeholder="Выберите тип обучения"
            :disabled="!selectedDivisionId"
          >
            <option v-for="type in educationTypes" :key="type.id_educationType" :value="type.id_educationType">
              {{ type.typeName }}
            </option>
          </AppSelect>

          <AppSelect
            v-model="selectedProgramId"
            label="Программа"
            placeholder="Выберите программу"
            :disabled="!selectedDivisionId || !selectedEducationTypeId || !filteredPrograms.length"
            :help="selectedDivisionId && selectedEducationTypeId && !filteredPrograms.length ? 'По выбранным фильтрам программ не найдено.' : ''"
          >
            <option v-for="program in filteredPrograms" :key="program.id_program_education" :value="program.id_program_education">
              {{ program.name_prof_education }}
            </option>
          </AppSelect>

          <GroupSearchField
            v-model="groupSearch"
            label="Группа"
            placeholder="Введите часть названия группы"
            :options="groupOptions"
            :loading="groupsLoading"
            :disabled="!selectedProgramId"
            :help="selectedGroupId ? `Выбрана группа: ${selectedGroup?.group || ''}` : 'Поиск показывает похожие группы по введённому тексту.'"
            @select="selectGroup"
          />

          <AppInput v-model="startDate" label="Дата начала" type="date" />
          <AppInput v-model="endDate" label="Дата окончания" type="date" />
          <AppInput :model-value="selectedProgram ? `${new Intl.NumberFormat('ru-RU').format(selectedProgram.price)} ₽` : ''" label="Цена программы" disabled />
          <AppInput v-model="typeOfRetraining" label="Тип обучения" disabled />
        </div>
      </AppCard>

      <div class="page-actions">
        <AppButton variant="ghost" @click="router.push(`/listeners/${listenerId}`)">Назад</AppButton>
        <AppButton :disabled="saving" @click="createEnrollment">
          {{ saving ? 'Создаём...' : 'Создать запись' }}
        </AppButton>
      </div>

      <AppConfirmDialog
        :open="confirmContractorDelete"
        title="Удаление заказчика"
        message="Удалить заказчика у этого слушателя?"
        :loading="contractorSaving"
        confirm-label="Удалить"
        @cancel="confirmContractorDelete = false"
        @confirm="removeContractor"
      />
    </template>
  </section>
</template>

<style scoped>
.summary-line {
  margin: 0;
  font-weight: 700;
  color: #0f172a;
}

.summary-subline {
  margin: 0;
  color: #64748b;
}

.summary-stack,
.summary-actions,
.enrollment-grid,
.radio-grid {
  display: grid;
  gap: 1rem;
}

.summary-actions,
.page-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
}

.enrollment-grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.radio-grid {
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
}

.radio-card {
  display: flex;
  gap: 0.75rem;
  align-items: flex-start;
  padding: 1rem;
  border-radius: 1rem;
  border: 1px solid rgba(15, 23, 42, 0.08);
  background: rgba(255, 255, 255, 0.7);
}

.page-actions {
  justify-content: space-between;
}

@media (max-width: 900px) {
  .enrollment-grid {
    grid-template-columns: 1fr;
  }

  .page-actions {
    flex-direction: column;
  }
}
</style>
