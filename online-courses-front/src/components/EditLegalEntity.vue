<template>
  <Header title="Редактирование юридического лица" />

  <div style="padding: 100px 20px 20px 20px;">
    <form @submit.prevent="updateLegalEntity" class="row g-4">

      <div class="col-md-4">
        <div class="card p-3 shadow-sm card-block">
          <h5 class="card-title mb-3">Данные компании</h5>
          <div class="d-flex flex-column gap-2">
            <input v-model="legalEntity.name_company" class="form-control" placeholder="Название компании" required />
            <input v-model="legalEntity.inn" @input="numbersOnly('inn', 10)" class="form-control" placeholder="ИНН (10 цифр)" maxlength="10" />
            <input v-model="legalEntity.kpp" @input="numbersOnly('kpp', 9)" class="form-control" placeholder="КПП (9 цифр)" maxlength="9" />
            <input v-model="legalEntity.ogrn" @input="numbersOnly('ogrn', 13)" class="form-control" placeholder="ОГРН (13 цифр)" maxlength="13" />
            <input v-model="legalEntity.phone" @input="formatPhone" class="form-control" placeholder="Телефон (+7XXXXXXXXXX)" maxlength="12" />
            <input v-model="legalEntity.email" @input="validateEmail"
                   :class="['form-control', { 'is-invalid': !emailValid && legalEntity.email }]"
                   placeholder="Email" />
            <div v-if="!emailValid && legalEntity.email" class="invalid-feedback d-block">
              Некорректный email
            </div>
          </div>
        </div>
      </div>

      <div class="col-md-4">
        <div class="card p-3 shadow-sm card-block">
          <h5 class="card-title mb-3">Представитель</h5>
          <div class="d-flex flex-column gap-2">
            <input v-model="legalEntity.first_name" @input="onlyLetters('first_name')" class="form-control" placeholder="Имя" required />
            <input v-model="legalEntity.second_name" @input="onlyLetters('second_name')" class="form-control" placeholder="Фамилия" required />
            <input v-model="legalEntity.middle_name" @input="onlyLetters('middle_name')" class="form-control" placeholder="Отчество" />
            <input v-model="legalEntity.status" @input="onlyLetters('status')" class="form-control" placeholder="Должность">
          </div>
        </div>
      </div>

      <div class="col-md-4">
        <div class="card p-3 shadow-sm card-block">
          <h5 class="card-title mb-3">Адрес регистрации</h5>
          <div class="d-flex flex-column gap-2">
            <input v-model="regAddress.mail_index" @input="formatIndex" class="form-control" placeholder="Индекс (6 цифр)" maxlength="6" />
            <input v-model="regAddress.region" @input="onlyLettersAddress('region')" class="form-control" placeholder="Регион" />
            <input v-model="regAddress.city" @input="onlyLettersAddress('city')" class="form-control" placeholder="Город" />
            <input v-model="regAddress.street" class="form-control" placeholder="Улица" />
            <input v-model="regAddress.house" class="form-control" placeholder="Дом" />
            <input v-model="regAddress.building" class="form-control" placeholder="Корпус" />
            <input v-model="regAddress.apartment" class="form-control" placeholder="Квартира" />
          </div>
        </div>
      </div>

      <!-- Кнопки -->
      <div class="col-12 d-flex justify-content-between">
        <button type="button" class="btn btn-secondary" @click="goBack" :disabled="loading">Отмена</button>
        <button type="submit" class="btn btn-success px-4" :disabled="loading">
          {{ loading ? 'Сохранение...' : 'Сохранить изменения' }}
        </button>
      </div>

    </form>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { API_URL_CORE } from '../config'
import { toast } from 'vue3-toastify'
import Header from './Header.vue'

const router = useRouter()
const route = useRoute()
const id = route.params.id
const token = localStorage.getItem('access_token')

const loading = ref(false)
const emailValid = ref(true)

const legalEntity = ref({})
const regAddress = ref({})

const onlyLetters = field => {
  legalEntity.value[field] = (legalEntity.value[field] || '').replace(/[^А-Яа-яЁёA-Za-z\s-]/g, '')
}

const onlyLettersAddress = field => {
  regAddress.value[field] = (regAddress.value[field] || '').replace(/[^А-Яа-яЁёA-Za-z\s-]/g, '')
}

const numbersOnly = (field, max) => {
  legalEntity.value[field] = (legalEntity.value[field] || '').replace(/\D/g, '').slice(0, max)
}

const validateEmail = () => {
  const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  emailValid.value = re.test((legalEntity.value.email || '').toString())
}

const formatPhone = () => {
  let v = (legalEntity.value.phone || '').replace(/\D/g, '')
  if (!v.startsWith('7')) v = '7' + v
  legalEntity.value.phone = '+' + v.slice(0, 11)
}
const formatIndex = () => {
  regAddress.value.mail_index = (regAddress.value.mail_index || '').replace(/\D/g, '').slice(0, 6)
}

const loadLegalEntity = async () => {
  try {
    loading.value = true
    const res = await fetch(`${API_URL_CORE}/legalentity/details/${id}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error(`Ошибка загрузки (${res.status})`)
    const data = await res.json()


    const le = data.data.legal_entity
    const addr = data.data.reg_address

    legalEntity.value = {
      name_company: le.name_company || '',
      inn: le.inn || '',
      kpp: le.kpp || '',
      ogrn: le.ogrn || '',
      phone: le.phone || '',
      email: le.email || '',
      first_name: le.first_name || '',
      second_name: le.second_name || '',
      middle_name: le.middle_name || '',
      status: le.status || ''
    }

    regAddress.value = {
      mail_index: (addr.mail_index || '').toString(),
      region: addr.region || '',
      city: addr.city || '',
      street: addr.street || '',
      house: addr.house || '',
      building: addr.building || '',
      apartment: addr.apartment || ''
    }

  } catch (err) {
    toast.error(err.message || 'Ошибка при загрузке данных')
  } finally {
    loading.value = false
  }
}


const updateLegalEntity = async () => {
  if (!emailValid.value) {
    toast.error('Проверьте правильность Email')
    return
  }

  loading.value = true
  try {
    const payload = {
      legal_entity: { ...legalEntity.value },
      reg_address: { ...regAddress.value }
    }

    const res = await fetch(`${API_URL_CORE}/legalentity/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
      body: JSON.stringify(payload)
    })

    if (!res.ok) {
      let text = `Ошибка сохранения (${res.status})`
      try { const body = await res.text(); if (body) text = body } catch {}
      throw new Error(text)
    }

    toast.success('Данные успешно обновлены')
    setTimeout(() => router.push('/legalentities'), 2000)
  } catch (err) {
    toast.error(err.message || 'Ошибка при сохранении')
  } finally {
    loading.value = false
  }
}

const goBack = () => router.push('/legalentities')

onMounted(loadLegalEntity)
</script>

<style scoped>
.card-block { border-radius: 8px; }
.is-invalid { border-color: red !important; }
.invalid-feedback { font-size: 13px; }
</style>
