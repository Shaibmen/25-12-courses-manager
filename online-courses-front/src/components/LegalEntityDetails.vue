<template>
  <div>
    <Header title="Информация о юридическом лице" />

    <div class="container py-4 mt-5">
      <div v-if="loading" class="text-center py-5">Загрузка данных...</div>

      <form v-else class="row g-4">

        <div class="col-md-4">
          <div class="card p-3 shadow-sm">
            <h5 class="card-title mb-3">Компания</h5>
            <p><strong>Название:</strong> {{ legal_entity.name_company}}</p>
            <p><strong>ИНН:</strong> {{ legal_entity.inn }}</p>
            <p><strong>КПП:</strong> {{ legal_entity.kpp }}</p>
            <p><strong>ОГРН:</strong> {{ legal_entity.ogrn }}</p>
            <p><strong>Телефон:</strong> {{ legal_entity.phone }}</p>
            <p><strong>Email:</strong> {{ legal_entity.email }}</p>
          </div>
        </div>

         <div class="col-md-4">
          <div class="card p-3 shadow-sm">
            <h5 class="card-title mb-3">Представитель</h5>
            <p><strong>Имя:</strong> {{ legal_entity.first_name }}</p>
            <p><strong>Фамилия:</strong> {{ legal_entity.second_name }}</p>
            <p><strong>Отчество:</strong> {{ legal_entity.middle_name || '—' }}</p>
             <p><strong>Должность:</strong> {{ legal_entity.status }}</p>
          </div>
        </div>

        

        <div class="col-md-4">
          <div class="card p-3 shadow-sm">
            <h5 class="card-title mb-3">Адрес регистрации</h5>
            <p><strong>Почтовый индекс:</strong> {{ reg_address.mail_index }}</p>
            <p><strong>Регион:</strong> {{ reg_address.region }}</p>
            <p><strong>Город:</strong> {{ reg_address.city }}</p>
            <p><strong>Улица:</strong> {{ reg_address.street }}</p>
            <p><strong>Дом:</strong> {{ reg_address.house }}</p>
            <p><strong>Корпус:</strong> {{ reg_address.building }}</p>
            <p><strong>Квартира:</strong> {{ reg_address.apartment }}</p>
          </div>
        </div>

        <div class="col-12">
  <div class="card p-3 shadow-sm">
    <h5 class="card-title mb-3">Слушатели</h5>

    <div v-if="listeners.length === 0" class="text-muted">
      У юридического лица нет слушателей
    </div>

    <div v-else class="table-responsive">
      <table class="table table-sm align-middle mb-0">
        <thead>
          <tr>
            <th>ФИО</th>
            <th>СНИЛС</th>
            <th class="text-end">Действия</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="l in listeners" :key="l.id">
            <td>
              {{ l.second_name }} {{ l.first_name }} {{ l.middle_name }}
            </td>
            <td>{{ l.snils }}</td>
            <td class="text-end">
               <button 
        type="button" 
        class="btn btn-success mt-2"
        @click="pushListener(l.id_listener)"
      >
        Перейти
      </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</div>


<div class="col-12 d-flex justify-content-between">
          <button 
        type="button" 
        class="btn btn-success mt-2"
        @click="goCreateListener(legal_entity.id_legalentity)"
      >
        Добавить слушателя
      </button>

      <button 
        type="button" 
        class="btn btn-success mt-2"
        @click="goEnrollmentYUR(legal_entity.id_legalentity)"
      >
        Записать слушателей на курс
      </button>

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
const legal_entity = ref({})
const reg_address = ref({})
const listeners = ref([])


const goCreateListener = () => {
  router.push({
    path: '/listeners/create',
    query: { id_legalentity: legal_entity.value.id_legalentity }
  })
}



const loadDetails = async () => {
  try {
    const res = await fetch(`${API_URL_CORE}/legalentity/details/${id}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error(`Ошибка загрузки (${res.status})`)

    const data = await res.json()

    const le = data.data.legal_entity
    const ra = data.data.reg_address


     listeners.value = (le.listeners || []).map(l => ({
      id_listener: l.id_listener,
      first_name: l.first_name,
      second_name: l.second_name,
      middle_name: l.middle_name,
      snils: l.snils
    }))

    legal_entity.value = {
      listeners: listeners.value,
      id_legalentity: le.id_legalentity || '',
      name_company: le.name_company || '',
      inn: le.inn || '',
      kpp: le.kpp || '',
      ogrn: le.ogrn || '',
      phone: le.phone || '',
      email: le.email || '',
      first_name: le.first_name || '',
      second_name: le.second_name || '',
      middle_name: le.middle_name || '',
      id_regaddress: le.id_regaddress || '',
      status: le.status || ''
     
    }
    

    reg_address.value = {
      mail_index: ra.mail_index || '',
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
const pushListener = (id) => router.push(`/listeners/${id}`)

const goEnrollmentYUR = (id) => router.push(`/enrollment/yur/${id}`)
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
