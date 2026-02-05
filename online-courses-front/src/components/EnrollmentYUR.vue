<template>
  <div>
  <Header title="Запись на курс" />
    </div>

  <div style="padding: 100px 20px 20px 20px; max-width: 1200px; margin: 0 auto;">
    <div>
      <div class="card p-3 shadow-sm card-block mt-4">
        <h5 class="card-title mb-3">Договор</h5>

        <div class="d-flex flex-column gap-3">
          <select
            v-model="selectedContractId"
            class="form-select"
            :disabled="!contracts.length || loadingContracts"
          >
            <option value="">— Выберите договор —</option>
            <option
              v-for="c in filteredContracts"
              :key="c.id_contract"
              :value="c.id_contract"
            >
              {{ c.name }} ({{ c.type === 'bilateral' ? 'двусторонний' : 'трёхсторонний' }})
            </option>
          </select>

          <select v-model="studyLoadOption" class="form-select">
            <option value="">— Недельная учебная нагрузка —</option>
            <option v-for="(label, key) in loadVariantsNotDO" :key="key" :value="key">
              {{ label }}
            </option>
          </select>

          <select v-model="paymentOption" class="form-select">
            <option value="">— Порядок оплаты —</option>
            <option value="full">
              Оплата осуществляется в следующем порядке: 100% предоплата до начала обучения.
            </option>
            <option value="split">
              Оплата осуществляется в следующем порядке: аванс 50 % — предоплата до начала обучения, оставшиеся 50 % — в установленный срок.
            </option>
            <option value="halfsplit">
              Оплата подлежит перечислению на расчётный счёт Исполнителя в срок до 5 (пяти) рабочих дней, считая с момента (даты) подписания Сторонами Акта оказанных услуг.
            </option>
          </select>

          <input
            v-if="paymentOption === 'split'"
            type="date"
            v-model="secondPaymentDate"
            class="form-control"
            placeholder="Срок оплаты оставшихся 50 %"
          />

          <select v-model="ageCategory" class="form-select">
            <option value="">— Возрастная категория —</option>
            <option v-for="(label, key) in ageCategories" :key="key" :value="key">
              {{ label }}
            </option>
          </select>

          <select v-model="optDocumentSelected" class="form-select">
            <option value="">— Итоговый документ и режим выдачи —</option>
            <option v-for="(label, key) in optDocumentOptions" :key="key" :value="key">
              {{ label }}
            </option>
          </select>
        </div>
      </div>

      <!---->
      <div class="card p-3 shadow-sm card-block mt-4">
        <h5 class="card-title mb-3">
          Слушатели <span class="text-danger">*</span>
          <span class="badge bg-secondary ms-2">{{ selectedListenerIds.length }} выбрано</span>
        </h5>

        <div v-if="loadingLegalEntityListeners" class="text-center py-3">
          <div class="spinner-border spinner-border-sm text-primary" role="status">
            <span class="visually-hidden">Загрузка...</span>
          </div>
          <span class="ms-2">Загрузка слушателей...</span>
        </div>

        <div v-else-if="legalEntityListeners.length === 0" class="text-muted">
          Нет доступных слушателей для данного юридического лица
        </div>

        <div v-else class="d-flex flex-column gap-2">
          <div class="form-check">
            <input
              class="form-check-input"
              type="checkbox"
              v-model="selectAllListeners"
              :disabled="saving"
              :id="'select-all-listeners'"
            />
            <label class="form-check-label fw-bold" :for="'select-all-listeners'">
              Выбрать всех
            </label>
          </div>
          
          <div 
            v-for="listenerItem in legalEntityListeners" 
            :key="listenerItem.id_listener" 
            class="form-check"
          >
            <input
              class="form-check-input"
              type="checkbox"
              :value="listenerItem.id_listener"
              v-model="selectedListenerIds"
              :id="'listener-' + listenerItem.id_listener"
              :disabled="saving"
            />
            <label class="form-check-label" :for="'listener-' + listenerItem.id_listener">
              {{ listenerItem.second_name }} {{ listenerItem.first_name }} {{ listenerItem.middle_name || '' }}
              <span v-if="listenerItem.snils" class="text-muted small ms-2">(СНИЛС: {{ listenerItem.snils }})</span>
            </label>
          </div>
        </div>
      </div>

      <div class="card p-3 shadow-sm card-block mt-4">
        <h5 class="card-title mb-3">
          Исполнитель со стороны 25-12 <span class="text-danger">*</span>
        </h5>

        <div class="d-flex flex-column gap-2">
          <div v-for="e in executors" :key="e.id_executor" class="form-check">
            <input
              class="form-check-input"
              type="radio"
              :value="e.id_executor"
              v-model="selectedExecutorId"
              :id="e.id_executor"
              :disabled="saving"
            />
            <label class="form-check-label" :for="e.id_executor">
              {{ e.second_name }} {{ e.first_name }} {{ e.middle_name }}
            </label>
          </div>
          <div v-if="!executors.length" class="text-muted">Исполнители не найдены</div>
        </div>
      </div>

      <div class="card p-3 shadow-sm card-block mt-4">
        <h5 class="card-title mb-3">Программа обучения</h5>
        <div class="d-flex flex-column gap-3">
          <select 
            v-model="selectedProgramId" 
            @change="onProgramChange" 
            class="form-select"
            :disabled="saving"
          >
            <option value="">— Выберите программу обучения —</option>
            <option v-for="p in programs" :key="p.id_program_education" :value="p.id_program_education">
              {{ p.name_prof_education }}
            </option>
          </select>

          <div class="d-flex align-items-center gap-3 flex-wrap">
            <div>
              <label class="form-label">Дата начала</label>
              <input 
                type="date" 
                v-model="startDate" 
                class="form-control" 
                required 
                :disabled="saving"
              />
            </div>
            <div>
              <label class="form-label">Дата окончания</label>
              <input 
                type="date" 
                v-model="endDate" 
                class="form-control" 
                required 
                :disabled="saving"
              />
            </div>
            <div>
              <label class="form-label">Цена</label>
              <select 
                v-model="currentPrice" 
                class="form-select" 
                required
                :disabled="saving || !selectedProgramId"
              >
                <option value="">— Выберите цену —</option>
                <option v-if="price.individual_price !== undefined" :value="price.individual_price">
                  Индивидуальное: {{ price.individual_price }} ₽
                </option>
                <option v-if="price.group_price !== undefined" :value="price.group_price">
                  Групповое: {{ price.group_price }} ₽
                </option>
                <option v-if="price.campus_price !== undefined" :value="price.campus_price">
                  Кампус: {{ price.campus_price }} ₽
                </option>
              </select>
            </div>
          </div>

          <div class="row g-2">
            <div class="col-md-6">
              <label class="form-label">Группа</label>
              <input 
                type="text" 
                v-model="group" 
                class="form-control" 
                placeholder="Например: 0" 
                :disabled="saving"
              />
            </div>
            <div class="col-md-6">
              <label class="form-label">Тип обучения</label>
              <select v-model="typeOfRetraining" class="form-select" disabled>
                <option value="">— Автоподстановка типа обучения —</option>
                <option value="Повышение квалификации">Повышение квалификации</option>
                <option value="Профессиональная переподготовка">Профессиональная переподготовка</option>
                <option value="Дополнительное образование">Дополнительное образование</option>
              </select>
              <div class="form-text">Тип обучения подставляется автоматически по выбранному договору</div>
            </div>
          </div>
        </div>
      </div>

      <div class="mt-4 d-flex gap-2">
        <button 
          class="btn btn-success px-4" 
          @click="createEnrollment" 
          :disabled="!isFormValid || saving"
        >
          <span v-if="saving" class="spinner-border spinner-border-sm me-2" role="status"></span>
          {{ saving ? 'Создание...' : 'Создать запись' }}
        </button>
        <button 
          class="btn btn-secondary" 
          @click="goBack"
          :disabled="saving"
        >
          Назад
        </button>
      </div>
    </div>
  </div>

</template>
<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { toast } from 'vue3-toastify'
import Header from './Header.vue'
import { API_URL_CORE } from '../config'

const route = useRoute()
const router = useRouter()
const legalentityID = route.params.id
const token = localStorage.getItem('access_token')


const legalEntityData = ref(null)
const legalEntityListeners = ref([])
const loadingLegalEntityListeners = ref(true)
const loading = ref(true)


const selectedListenerIds = ref([])
const selectAllListeners = ref(false)

const ageCategory = ref('')
const loadVariant = ref('')
const studyLoadOption = ref('')
const paymentOption = ref('')
const secondPaymentDate = ref('')
const optDocumentSelected = ref('')
const contracts = ref([])
const loadingContracts = ref(false)
const executors = ref([])
const selectedExecutorId = ref('')
const programs = ref([])
const selectedProgramId = ref('')
const price = ref({ individual_price: 0, group_price: 0, campus_price: 0 })
const currentPrice = ref(0)
const startDate = ref('')
const endDate = ref('')
const group = ref('')
const typeOfRetraining = ref('')
const selectedContractId = ref('')
const saving = ref(false)


const loadVariantsNotDO = {
  1: 'с пониженной недельной учебной нагрузкой (3 акад. часа в неделю)',
  2: 'с умеренной недельной учебной нагрузкой (6 акад. часов в неделю)',
  3: 'со стандартной недельной учебной нагрузкой (12 акад. часов в неделю)',
  4: 'с высокой недельной учебной нагрузкой (15 акад. часов в неделю)',
  5: 'с повышенной недельной учебной нагрузкой (30 акад. часов в неделю)',
  6: 'с интенсивной недельной учебной нагрузкой (36 акад. часов в неделю)'
}

const ageCategories = {
  BELOW_EIGHTEEN: 'Меньше восемнадцати',
  FOURTEEN: 'Меньше четырнадцати',
  EIGHTEEN: 'Восемнадцать'
}

const optDocumentOptions = {
  1: 'удостоверение о повышении квалификации вручается по окончании',
  2: 'удостоверение выдаётся одновременно с дипломом СПО/ВО (ч. 16 ст. 76 ФЗ-273). До этого момента удостоверение хранится у Исполнителя.'
}

const hasContractor = computed(() => {
  return legalentityID && legalEntityListeners.value.length > 0
})

const filteredContracts = computed(() => {
  if (hasContractor.value) {
    return contracts.value.filter(c => c.type === 'trilateral')
  } else {
    return contracts.value.filter(c => c.type === 'bilateral')
  }
})



const isPKorPP = computed(() => {
  const c = contracts.value.find(c => c.id_contract === selectedContractId.value)
  if (!c) return false
  const id = String(c.id_contract).toUpperCase()
  return id.startsWith('PK') || id.startsWith('PP')
})


watch(selectAllListeners, (newValue) => {
  if (newValue) {
    selectedListenerIds.value = legalEntityListeners.value.map(l => l.id_listener)
  } else {
    selectedListenerIds.value = []
  }
})

watch(selectedListenerIds, (newIds) => {
  if (newIds.length === legalEntityListeners.value.length) {
    selectAllListeners.value = true
  } else if (newIds.length > 0 && newIds.length < legalEntityListeners.value.length) {
    selectAllListeners.value = false
  } else if (newIds.length === 0) {
    selectAllListeners.value = false
  }
})

const loadLegalEntityDetails = async () => {
  loadingLegalEntityListeners.value = true
  try {
    const res = await fetch(`${API_URL_CORE}/legalentity/details/${legalentityID}`, {
      headers: { 
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })
    
    if (!res.ok) {
      const error = await res.json().catch(() => ({}))
      throw new Error(error.message || `Ошибка ${res.status}: ${res.statusText}`)
    }

    const data = await res.json()
    legalEntityData.value = data.data || null
    
    if (data.data?.legal_entity?.listeners) {
      legalEntityListeners.value = data.data.legal_entity.listeners
    } else {
      legalEntityListeners.value = []
    }
    
  } catch (err) {
    console.error('Ошибка загрузки юридического лица:', err)
    toast.error(`Ошибка загрузки данных: ${err.message}`)
    legalEntityData.value = null
    legalEntityListeners.value = []
  } finally {
    loadingLegalEntityListeners.value = false
  }
}

const loadPrograms = async () => {
  try {
    const res = await fetch(`${API_URL_CORE}/programeducation/?page=1&filter=`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error('Ошибка загрузки программ')
    
    const data = await res.json()
    programs.value = data.data || []
  } catch (err) {
    toast.error('Не удалось загрузить программы обучения')
    console.error('Ошибка загрузки программ:', err)
  }
}

const loadContracts = async () => {
  loadingContracts.value = true
  try {
    contracts.value = [
      { id_contract: 'PK_3_YUR', name: 'ПК с оплатой юридическим лицом', type: 'trilateral' },
      { id_contract: 'PP_3_YUR', name: 'ПП с оплатой юридическим лицом', type: 'trilateral' }
    ]
  } catch (err) {
    toast.error('Не удалось загрузить договоры')
    console.error('Ошибка загрузки договоров:', err)
  } finally {
    loadingContracts.value = false
  }
}

const loadExecutors = async () => {
  try {
    const res = await fetch(`${API_URL_CORE}/executer/`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error('Ошибка загрузки исполнителей')
    
    const data = await res.json()
    executors.value = data.data || []

    if (executors.value.length > 0 && !selectedExecutorId.value) {
      selectedExecutorId.value = executors.value[0].id_executor
    }
  } catch (err) {
    toast.error(err.message || 'Не удалось загрузить исполнителей')
    console.error('Ошибка загрузки исполнителей:', err)
  }
}

const onProgramChange = () => {
  const p = programs.value.find(prog => prog.id_program_education === selectedProgramId.value)
  if (p) {
    price.value = {
      individual_price: p.individual_price || 0,
      group_price: p.group_price || 0,
      campus_price: p.campus_price || 0
    }
    currentPrice.value = p.individual_price || 0
  } else {
    price.value = { individual_price: 0, group_price: 0, campus_price: 0 }
    currentPrice.value = 0
  }
}


const isFormValid = computed(() => {
  return (
    selectedListenerIds.value.length > 0 &&
    selectedProgramId.value &&
    startDate.value &&
    endDate.value &&
    currentPrice.value &&
    currentPrice.value > 0 &&
    selectedExecutorId.value &&
    (!hasContractor.value || selectedContractId.value)
  )
})


//asdfasdf
const createDogovor = async () => {
  if (!isFormValid.value) {
    toast.warn('Заполните все обязательные поля (отмечены *)')
    return
  }

  try {
    saving.value = true

    const selectedListeners = legalEntityListeners.value.filter(
      listener => selectedListenerIds.value.includes(listener.id_listener)
    )

    const listenersData = selectedListeners.map(listener => ({
      first_name: listener.first_name,
      second_name: listener.second_name,
      middle_name: listener.middle_name || '',
      snils: listener.snils || '',
      date_of_birth: listener.date_of_birth || '',
      email: listener.email || '',
    }))

    let optPriceValue = ''
    if (paymentOption.value === 'split') {
      optPriceValue = `Оплата осуществляется в следующем порядке: аванс 50 % — предоплата до начала обучения, оставшиеся 50 % — в установленный срок${secondPaymentDate.value ? ' ' + secondPaymentDate.value : ''}`
    } else if (paymentOption.value === 'full') {
      optPriceValue = 'Оплата осуществляется в следующем порядке: 100% предоплата до начала обучения.'
    } else if (paymentOption.value === 'halfsplit') {
      optPriceValue = 'Оплата подлежит перечислению на расчётный счёт Исполнителя в срок до 5 (пяти) рабочих дней, считая с момента (даты) подписания Сторонами Акта оказанных услуг.'
    }

    let variantValue = null
    let optNagruzValue = null

    if (studyLoadOption.value) {
      variantValue = Number(studyLoadOption.value)
      optNagruzValue = Number(studyLoadOption.value)
    }

    const legalEntity = legalEntityData.value?.legal_entity || {}
    const regAddress = legalEntityData.value?.reg_address || {}
    
    const address = {
      mail_index: String(regAddress.mail_index || ''),
      city: regAddress.city || '',
      street: regAddress.street || '',
      house: regAddress.house || '',
      apartment: regAddress.apartment || ''
    }

     const selectedProgram = programs.value.find(
      prog => prog.id_program_education === selectedProgramId.value
    )
    const programName = selectedProgram ? selectedProgram.name_prof_education : ''
    const programHours = selectedProgram ? selectedProgram.time_education : 0

    let priceType = ''
    let selectedPriceValue = Number(currentPrice.value)
    
    if (selectedProgram) {
      if (selectedPriceValue === selectedProgram.individual_price) {
        priceType = 'individual'
      } else if (selectedPriceValue === selectedProgram.group_price) {
        priceType = 'group'
      } else if (selectedPriceValue === selectedProgram.campus_price) {
        priceType = 'campus'
      }
    }


    const zakazchikFIO = `${legalEntity.second_name || ''} ${legalEntity.first_name || ''} ${legalEntity.middle_name || ''}`.trim()

    const documentsDataRequest = {
      id_executor: selectedExecutorId.value,
      front_data: {
        legal_entity: {
          listeners: listenersData,
          reg_address: address,
          company_name: legalEntity.name_company || '',
          zakazchikfio: zakazchikFIO,
          status: legalEntity.status || '',
          inn: legalEntity.inn || '',
          kpp: legalEntity.kpp || '',
          ogrn: legalEntity.ogrn || '',
          phone: legalEntity.phone || '',
          email: legalEntity.email || ''
        },
        start_date: startDate.value,
        end_date: endDate.value,
        program_name: programName,
        time_education: programHours,
        price_enrollment: parseFloat(selectedPriceValue),
        variant: variantValue,
        dogovor_type: selectedContractId.value || null,
        opt_nagruz: optNagruzValue,
        opt_document: optDocumentSelected.value ? Number(optDocumentSelected.value) : null,
        dogovor_age: ageCategory.value || null,
        opt_price: optPriceValue || null
      }
    }

    const res = await fetch(`${API_URL_CORE}/document/`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`
      },
      body: JSON.stringify(documentsDataRequest)
    })

    if (!res.ok) {
      const errorData = await res.json().catch(() => ({}))
      throw new Error(errorData.message || `Ошибка создания документов: ${res.status}`)
    }

    toast.success('Документы (договор) успешно созданы')

    const enrollmentsPromises = selectedListenerIds.value.map(async (listenerId) => {
      const enrollmentBody = {
        id_listener: listenerId,
        id_program: selectedProgramId.value,
        start_date: startDate.value,
        end_date: endDate.value,
        opt_nagruz: optNagruzValue,
        current_price: Number(currentPrice.value),
        group: group.value || null,
        type_of_retraining: typeOfRetraining.value || null,
        is_active: true
      }

      const enrollmentRes = await fetch(`${API_URL_CORE}/enrollment/`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`
        },
        body: JSON.stringify(enrollmentBody)
      })

      if (!enrollmentRes.ok) {
        const error = await enrollmentRes.json().catch(() => ({}))
        throw new Error(`Ошибка создания записи для слушателя ${listenerId}: ${error.message || 'неизвестная ошибка'}`)
      }

      return enrollmentRes.json()
    })

    const enrollmentResults = await Promise.allSettled(enrollmentsPromises)

    const failedEnrollments = enrollmentResults.filter(result => result.status === 'rejected')
    const successfulEnrollments = enrollmentResults.filter(result => result.status === 'fulfilled')

    if (successfulEnrollments.length > 0) {
      toast.success(`Создано записей на обучение: ${successfulEnrollments.length}`)
    }

    if (failedEnrollments.length > 0) {
      console.error('Не удалось создать некоторые записи:', failedEnrollments)
      toast.warning(`Не удалось создать ${failedEnrollments.length} записей. Проверьте консоль для деталей.`)
    }

    const frontData2 = {
      dogovor_type: selectedContractId.value || null,
      dogovor_age: ageCategory.value || null,
      opt_document: optDocumentSelected.value ? Number(optDocumentSelected.value) : null,
      opt_price: optPriceValue || null,
      opt_nagruz: optNagruzValue  || null
    }

    const docsPromises = selectedListenerIds.value.map(async (listenerId) => {
      const documentPayload = {
        id_listener: listenerId,
        id_program: selectedProgramId.value,
        id_executor: selectedExecutorId.value || null,
        front_data: frontData2
      }

      const docRes = await fetch(`${API_URL_CORE}/document/`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`
        },
        body: JSON.stringify(documentPayload)
      })

      if (!docRes.ok) {
        const error = await docRes.json().catch(() => ({}))
        throw new Error(`Ошибка создания документов для слушателя ${listenerId}: ${error.message || 'неизвестная ошибка'}`)
      }

      return docRes.json()
    })

    const docsResults = await Promise.allSettled(docsPromises)

    const failedDocs = docsResults.filter(result => result.status === 'rejected')
    const successfulDocs = docsResults.filter(result => result.status === 'fulfilled')

    if (successfulDocs.length > 0) {
      toast.success(`Создано дополнительных документов: ${successfulDocs.length}`)
    }

    if (failedDocs.length > 0) {
      console.error('Не удалось создать некоторые документы:', failedDocs)
      toast.warning(`Не удалось создать ${failedDocs.length} документов. Проверьте консоль для деталей.`)
    }


    if (selectedListenerIds.value.length > 0) {
      toast.success(`Обработка завершена. Всего слушателей: ${selectedListenerIds.value.length}`)
      

    }

  } catch (err) {
    console.error('Критическая ошибка при создании:', err)
    toast.error(err.message || 'Произошла критическая ошибка при создании документов')
  } finally {
    saving.value = false
  }
}

const createEnrollment = () => {
  createDogovor()
}

const goBack = () => {
  if (!saving.value) {
    router.push(`/legalentities`)
  }
}

onMounted(async () => {
  try {
    await Promise.all([
      loadLegalEntityDetails(), 
      loadPrograms(),
      loadContracts(),
      loadExecutors()
    ])
  } catch (error) {
    console.error('Ошибка при инициализации:', error)
    toast.error('Ошибка загрузки данных')
  } finally {
    loading.value = false
  }
})


watch(selectedContractId, (newId) => {
  const c = contracts.value.find(ci => ci.id_contract === newId)
  const id = c?.id_contract?.toUpperCase() || ''
  if (id.startsWith('DO')) typeOfRetraining.value = 'Дополнительное образование'
  else if (id.startsWith('PK')) typeOfRetraining.value = 'Повышение квалификации'
  else if (id.startsWith('PP')) typeOfRetraining.value = 'Профессиональная переподготовка'
  else typeOfRetraining.value = ''
})

</script>

<style scoped>
.card-block {
  border-radius: 8px;
  border: 1px solid #dee2e6;
}
.form-label {
  font-weight: 500;
  margin-bottom: 0.25rem;
}
.form-check-label {
  cursor: pointer;
  user-select: none;
}
.form-check-input:disabled + .form-check-label {
  opacity: 0.6;
  cursor: not-allowed;
}
.spinner-border {
  width: 1rem;
  height: 1rem;
}
.btn:disabled {
  cursor: not-allowed;
}
</style>