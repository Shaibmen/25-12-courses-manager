<script setup lang="ts">
import type {
  EnrollmentProgramDetails,
  ListenerDetailsResponse
} from '../../../types/listener'
import { formatApiDate } from '../../../utils/date'
import AppButton from '../../ui/AppButton.vue'
import AppCard from '../../ui/AppCard.vue'

const props = defineProps<{
  details: ListenerDetailsResponse
  enrollments: EnrollmentProgramDetails[]
  files: string[]
  filesLoading?: boolean
  downloadLoading?: string
}>()

defineEmits<{
  back: []
  edit: []
  enroll: []
  download: [fileName: string]
}>()

const renderValue = (value?: string | number | null) =>
  value === undefined || value === null || value === '' ? '—' : String(value)
</script>

<template>
  <section class="stack">
    <div class="details-grid details-grid--triple">
      <AppCard title="Личные данные">
        <div class="fact-list">
          <p><strong>Фамилия:</strong> {{ renderValue(details.listener.second_name) }}</p>
          <p><strong>Имя:</strong> {{ renderValue(details.listener.first_name) }}</p>
          <p><strong>Отчество:</strong> {{ renderValue(details.listener.middle_name) }}</p>
          <p><strong>Дата рождения:</strong> {{ formatApiDate(details.listener.date_of_birth) }}</p>
          <p><strong>СНИЛС:</strong> {{ renderValue(details.listener.snils) }}</p>
          <p><strong>Телефон:</strong> {{ renderValue(details.listener.contact_phone) }}</p>
          <p><strong>Email:</strong> {{ renderValue(details.listener.email) }}</p>
        </div>
      </AppCard>

      <AppCard title="Паспорт">
        <div class="fact-list">
          <p><strong>Место рождения:</strong> {{ renderValue(details.passport?.place_birth) }}</p>
          <p><strong>Гражданство:</strong> {{ renderValue(details.passport?.citizenship) }}</p>
          <p><strong>Пол:</strong> {{ renderValue(details.passport?.gender) }}</p>
          <p><strong>Серия и номер:</strong> {{ `${details.passport?.seria || '—'} ${details.passport?.number || ''}`.trim() }}</p>
          <p><strong>Кем выдан:</strong> {{ renderValue(details.passport?.passport_given) }}</p>
          <p><strong>Дата выдачи:</strong> {{ formatApiDate(details.passport?.date_given) }}</p>
          <p><strong>Код подразделения:</strong> {{ renderValue(details.passport?.code) }}</p>
        </div>
      </AppCard>

      <AppCard title="Адрес регистрации">
        <div class="fact-list">
          <p><strong>Индекс:</strong> {{ renderValue(details.regaddress.mail_index) }}</p>
          <p><strong>Регион:</strong> {{ renderValue(details.regaddress.region) }}</p>
          <p><strong>Город:</strong> {{ renderValue(details.regaddress.city) }}</p>
          <p><strong>Улица:</strong> {{ renderValue(details.regaddress.street) }}</p>
          <p><strong>Дом:</strong> {{ renderValue(details.regaddress.house) }}</p>
          <p><strong>Корпус:</strong> {{ renderValue(details.regaddress.building) }}</p>
          <p><strong>Квартира:</strong> {{ renderValue(details.regaddress.apartment) }}</p>
        </div>
      </AppCard>
    </div>

    <div class="details-grid details-grid--double">
      <AppCard title="Образование">
        <div v-if="details.listener.looting_education" class="fact-list">
          <p><strong>Статус:</strong> Получает образование сейчас</p>
        </div>

        <div v-else class="fact-list">
          <p><strong>Учреждение:</strong> {{ renderValue(details.education_listener?.educational_institution) }}</p>
          <p><strong>Город:</strong> {{ renderValue(details.education_listener?.city) }}</p>
          <p><strong>Регион:</strong> {{ renderValue(details.education_listener?.region) }}</p>
          <p><strong>Специальность:</strong> {{ renderValue(details.education_listener?.speciality) }}</p>
          <p><strong>Уровень:</strong> {{ renderValue(details.education_listener?.level_education) }}</p>
          <p><strong>Диплом:</strong> {{ `${details.education_listener?.diplom_seria || '—'} ${details.education_listener?.diplom_number || ''}`.trim() }}</p>
          <p><strong>Дата выдачи:</strong> {{ formatApiDate(details.education_listener?.date_given) }}</p>
        </div>
      </AppCard>

      <AppCard title="Место работы">
        <div class="fact-list">
          <p><strong>Компания:</strong> {{ renderValue(details.placework?.name_company) }}</p>
          <p><strong>Должность:</strong> {{ renderValue(details.placework?.job_title) }}</p>
          <p><strong>Общий стаж:</strong> {{ renderValue(details.placework?.all_experience) }}</p>
          <p>
            <strong>Стаж по должности:</strong>
            {{ renderValue(details.placework?.job_title_experience ?? details.placework?.job_title_expirience) }}
          </p>
        </div>
      </AppCard>
    </div>

    <AppCard title="Курсы слушателя">
      <p v-if="!enrollments.length" class="empty-state">
        Курсов не найдено.
      </p>

      <div v-else class="course-list">
        <article v-for="course in enrollments" :key="`${course.id_program_education}-${course.start_date}`" class="course-card">
          <strong>{{ course.name_prof_education }}</strong>
          <span>{{ formatApiDate(course.start_date) }} - {{ formatApiDate(course.end_date) }}</span>
          <span>{{ course.current_price }} ₽</span>
        </article>
      </div>
    </AppCard>

    <AppCard title="Файлы слушателя">
      <p v-if="filesLoading" class="empty-state">Загрузка файлов...</p>
      <p v-else-if="!files.length" class="empty-state">Файлов не найдено.</p>

      <div v-else class="file-list">
        <div v-for="file in files" :key="file" class="file-row">
          <span class="file-row__name">{{ file }}</span>
          <AppButton
            variant="ghost"
            :disabled="downloadLoading === file"
            @click="$emit('download', file)"
          >
            {{ downloadLoading === file ? 'Скачиваем...' : 'Скачать' }}
          </AppButton>
        </div>
      </div>
    </AppCard>

    <div class="details-actions">
      <AppButton variant="ghost" @click="$emit('back')">Назад</AppButton>
      <div class="details-actions__group">
        <AppButton variant="secondary" @click="$emit('enroll')">Записать на курс</AppButton>
        <AppButton @click="$emit('edit')">Редактировать</AppButton>
      </div>
    </div>
  </section>
</template>

<style scoped>
.details-grid {
  display: grid;
  gap: 1rem;
}

.details-grid--triple {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.details-grid--double {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.fact-list {
  display: grid;
  gap: 0.55rem;
}

.fact-list p {
  margin: 0;
  color: #334155;
}

.fact-list strong {
  color: #0f172a;
}

.course-list,
.file-list {
  display: grid;
  gap: 0.8rem;
}

.course-card,
.file-row {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
  align-items: center;
  padding: 0.95rem 1rem;
  border-radius: 1rem;
  background: rgba(255, 255, 255, 0.7);
  border: 1px solid rgba(15, 23, 42, 0.08);
}

.course-card {
  flex-wrap: wrap;
}

.file-row__name {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.empty-state {
  margin: 0;
  color: #64748b;
}

.details-actions {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
}

.details-actions__group {
  display: flex;
  gap: 1rem;
}

@media (max-width: 1100px) {
  .details-grid--triple,
  .details-grid--double {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 720px) {
  .file-row,
  .details-actions,
  .details-actions__group {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
