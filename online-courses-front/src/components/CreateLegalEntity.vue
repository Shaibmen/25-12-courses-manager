<template>
  <Header title="Создать юридическое лицо" />

  <div style="padding: 100px 20px 20px 20px;">
    <form @submit.prevent="createEntity" class="row g-4">

      <!-- Юридическое лицо -->
      <div class="col-md-4">
        <div class="card p-3 shadow-sm card-block">
          <h5 class="card-title mb-3">Юридическое лицо</h5>
          <div class="d-flex flex-column gap-2">
            <input v-model="legal.name_company" class="form-control" placeholder="Название компании" required>
            <input v-model="legal.inn" @input="numbersOnly('inn', 12)" class="form-control" placeholder="ИНН" required>
            <input v-model="legal.kpp" @input="numbersOnly('kpp', 9)" class="form-control" placeholder="КПП">
            <input v-model="legal.ogrn" @input="numbersOnly('ogrn', 13)" class="form-control" placeholder="ОГРН">
            <input v-model="legal.phone" @input="formatPhone" class="form-control" placeholder="Телефон (+7XXXXXXXXXX)">
            <input v-model="legal.email" @input="validateEmail" :class="['form-control', { 'is-invalid': !emailValid && legal.email }]" placeholder="Email">
            <div v-if="!emailValid && legal.email" class="invalid-feedback d-block">Некорректный email</div>
          </div>
        </div>
      </div>

      <!-- Представитель -->
      <div class="col-md-4">
        <div class="card p-3 shadow-sm card-block">
          <h5 class="card-title mb-3">Представитель</h5>
          <div class="d-flex flex-column gap-2">
            <input v-model="legal.first_name" @input="onlyLetters('first_name')" class="form-control" placeholder="Имя" required>
            <input v-model="legal.second_name" @input="onlyLetters('second_name')" class="form-control" placeholder="Фамилия" required>
            <input v-model="legal.middle_name" @input="onlyLetters('middle_name')" class="form-control" placeholder="Отчество">
          </div>
        </div>
      </div>

      <!-- Адрес регистрации -->
      <div class="col-md-4">
        <div class="card p-3 shadow-sm card-block">
          <h5 class="card-title mb-3">Адрес регистрации</h5>
          <div class="d-flex flex-column gap-2">
            <input v-model="reg.mail_index" @input="formatIndex" class="form-control" placeholder="Индекс (6 цифр)">
            <input v-model="reg.region" @input="onlyLettersAddress('region')" class="form-control" placeholder="Регион">
            <input v-model="reg.city" @input="onlyLettersAddress('city')" class="form-control" placeholder="Город">
            <input v-model="reg.street" class="form-control" placeholder="Улица">
            <input v-model="reg.house" class="form-control" placeholder="Дом">
            <input v-model="reg.building" class="form-control" placeholder="Корпус">
            <input v-model="reg.apartment" class="form-control" placeholder="Квартира">
          </div>
        </div>
      </div>

      <div class="col-12 d-flex justify-content-between">
        <button type="button" class="btn btn-secondary" @click="goBack">Назад</button>

        <button type="submit" class="btn btn-success px-4" :disabled="loading">
          {{ loading ? 'Сохранение...' : 'Создать' }}
        </button>
      </div>

    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { API_URL_CORE } from '../config'
import { toast } from 'vue3-toastify'
import Header from './Header.vue'

const router = useRouter()
const loading = ref(false)

const legal = ref({
  name_company: '',
  inn: '',
  kpp: '',
  ogrn: '',
  phone: '',
  email: '',
  first_name: '',
  second_name: '',
  middle_name: ''
})

const reg = ref({
  mail_index: '',
  region: '',
  city: '',
  street: '',
  house: '',
  building: '',
  apartment: ''
})

const emailValid = ref(true)

const goBack = () => router.push('/legalentities')

const onlyLetters = field => {
  legal.value[field] = legal.value[field].replace(/[^А-Яа-яЁёA-Za-z\s-]/g, '')
}

const onlyLettersAddress = field => {
  reg.value[field] = reg.value[field].replace(/[^А-Яа-яЁёA-Za-z\s-]/g, '')
}

const numbersOnly = (field, max) => {
  legal.value[field] = legal.value[field].replace(/\D/g, '').slice(0, max)
}

const validateEmail = () => {
  const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  emailValid.value = re.test(legal.value.email)
}

const formatPhone = () => {
  let v = legal.value.phone.replace(/\D/g, '')
  if (!v.startsWith('7')) v = '7' + v
  legal.value.phone = '+' + v.slice(0, 11)
}

const formatIndex = () => {
  reg.value.mail_index = reg.value.mail_index.replace(/\D/g, '').slice(0, 6)
}

const createEntity = async () => {
  if (!emailValid.value) return toast.error("Проверьте Email")

  loading.value = true
  const token = localStorage.getItem('access_token')

  const payload = {
    legal_entity: legal.value,
    reg_address: reg.value
  }

  try {
    const res = await fetch(`${API_URL_CORE}/legalentity/`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`
      },
      body: JSON.stringify(payload)
    })

    if (!res.ok) throw new Error("Ошибка создания")

    toast.success("Юридическое лицо создано")
    setTimeout(() => router.push('/legalentities'), 1500)

  } catch (err) {
    toast.error(err.message || "Ошибка создания")
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.invalid {
  border-color: red !important;
}
</style>
