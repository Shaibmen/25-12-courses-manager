<template>
  <div>
    <Header title="Информация о юридическом лице" />

    <div class="container py-4 mt-5">
      <div v-if="loading" class="text-center py-5">Загрузка данных...</div>

      <form v-else class="row g-4">

        <div class="col-md-4">
          <div class="card p-3 shadow-sm">
            <h5 class="card-title mb-3">Компания</h5>
            <p><strong>Название:</strong> {{ legalEntity.name_company }}</p>
            <p><strong>ИНН:</strong> {{ legalEntity.inn }}</p>
            <p><strong>КПП:</strong> {{ legalEntity.kpp }}</p>
            <p><strong>ОГРН:</strong> {{ legalEntity.ogrn }}</p>
            <p><strong>Телефон:</strong> {{ legalEntity.phone }}</p>
            <p><strong>Email:</strong> {{ legalEntity.email }}</p>
          </div>
        </div>

         <div class="col-md-4">
          <div class="card p-3 shadow-sm">
            <h5 class="card-title mb-3">Представитель</h5>
            <p><strong>Имя:</strong> {{ legalEntity.first_name }}</p>
            <p><strong>Фамилия:</strong> {{ legalEntity.second_name }}</p>
            <p><strong>Отчество:</strong> {{ legalEntity.middle_name || '—' }}</p>
          </div>
        </div>

        

        <div class="col-md-4">
          <div class="card p-3 shadow-sm">
            <h5 class="card-title mb-3">Адрес регистрации</h5>
            <p><strong>Почтовый индекс:</strong> {{ regAddress.mail_index }}</p>
            <p><strong>Регион:</strong> {{ regAddress.region }}</p>
            <p><strong>Город:</strong> {{ regAddress.city }}</p>
            <p><strong>Улица:</strong> {{ regAddress.street }}</p>
            <p><strong>Дом:</strong> {{ regAddress.house }}</p>
            <p><strong>Корпус:</strong> {{ regAddress.building }}</p>
            <p><strong>Квартира:</strong> {{ regAddress.apartment }}</p>
          </div>
        </div>

<div class="col-12 d-flex justify-content-between">
    <button type="button" class="btn btn-success mt-2">Добавить слушателя</button>
       <button @click="goBack" class="btn btn-secondary">Назад</button>
      </div>
         <div class="col-12 d-flex justify-content-between">
          
          
        </div>

      </form>
    </div>
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

const loading = ref(true)
const legalEntity = ref({})
const regAddress = ref({})

const loadDetails = async () => {
  try {
    const res = await fetch(`${API_URL_CORE}/legalentity/details/${id}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error(`Ошибка загрузки (${res.status})`)
    const data = await res.json()
    const le = data.data.LegalEntity
    const ra = data.data.RegAddress

    legalEntity.value = {
      name_company: le.NameCompany || '',
      inn: le.Inn || '',
      kpp: le.Kpp || '',
      ogrn: le.Ogrn || '',
      phone: le.Phone || '',
      email: le.Email || '',
      first_name: le.FirstName || '',
      second_name: le.SecondName || '',
      middle_name: le.MiddleName || ''
    }

    regAddress.value = {
      mail_index: (ra.mail_index || '').toString(),
      region: ra.region || '',
      city: ra.city || '',
      street: ra.street || '',
      house: ra.house || '',
      building: ra.building || '',
      apartment: ra.apartment || ''
    }
  } catch (err) {
    toast.error(err.message || 'Ошибка при загрузке данных')
  } finally {
    loading.value = false
  }
}
const goBack = () => router.push('/legalentities')

onMounted(loadDetails)
</script>

<style scoped>
.card {
  border-radius: 12px;
  padding: 16px;
}

.card p {
  display: flex;
  justify-content: flex-start;
  margin: 4px 0;
}

.card p strong {
  width: 150px; 
  text-align: left;
  margin-right: 8px;
}

h5.card-title {
  margin-top: 0;
}
</style>
