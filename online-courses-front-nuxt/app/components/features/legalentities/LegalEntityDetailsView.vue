<script setup lang="ts">
import type { LegalEntityDetailsResponse } from '../../../types/legalentity'
import AppButton from '../../ui/AppButton.vue'
import AppCard from '../../ui/AppCard.vue'

defineProps<{
  details: LegalEntityDetailsResponse
  files: string[]
  filesLoading?: boolean
  downloadLoading?: string
}>()

defineEmits<{
  back: []
  edit: []
  addListener: []
  enroll: []
  openListener: [id: string]
  download: [fileName: string]
}>()

const renderValue = (value?: string | number | null) =>
  value === undefined || value === null || value === '' ? '—' : String(value)
</script>

<template>
  <section class="stack">
    <div class="details-grid">
      <AppCard title="Компания">
        <div class="fact-list">
          <p><strong>Название:</strong> {{ renderValue(details.legal_entity.name_company) }}</p>
          <p><strong>ИНН:</strong> {{ renderValue(details.legal_entity.inn) }}</p>
          <p><strong>КПП:</strong> {{ renderValue(details.legal_entity.kpp) }}</p>
          <p><strong>ОГРН:</strong> {{ renderValue(details.legal_entity.ogrn) }}</p>
          <p><strong>Телефон:</strong> {{ renderValue(details.legal_entity.phone) }}</p>
          <p><strong>Email:</strong> {{ renderValue(details.legal_entity.email) }}</p>
        </div>
      </AppCard>

      <AppCard title="Представитель">
        <div class="fact-list">
          <p><strong>Фамилия:</strong> {{ renderValue(details.legal_entity.second_name) }}</p>
          <p><strong>Имя:</strong> {{ renderValue(details.legal_entity.first_name) }}</p>
          <p><strong>Отчество:</strong> {{ renderValue(details.legal_entity.middle_name) }}</p>
          <p><strong>Должность:</strong> {{ renderValue(details.legal_entity.status) }}</p>
        </div>
      </AppCard>

      <AppCard title="Адрес регистрации">
        <div class="fact-list">
          <p><strong>Индекс:</strong> {{ renderValue(details.reg_address.mail_index) }}</p>
          <p><strong>Регион:</strong> {{ renderValue(details.reg_address.region) }}</p>
          <p><strong>Город:</strong> {{ renderValue(details.reg_address.city) }}</p>
          <p><strong>Улица:</strong> {{ renderValue(details.reg_address.street) }}</p>
          <p><strong>Дом:</strong> {{ renderValue(details.reg_address.house) }}</p>
          <p><strong>Корпус:</strong> {{ renderValue(details.reg_address.building) }}</p>
          <p><strong>Квартира:</strong> {{ renderValue(details.reg_address.apartment) }}</p>
        </div>
      </AppCard>
    </div>

    <AppCard title="Слушатели организации">
      <p v-if="!details.legal_entity.listeners?.length" class="empty-state">
        У юридического лица пока нет слушателей.
      </p>

      <div v-else class="listener-list">
        <article
          v-for="listener in details.legal_entity.listeners"
          :key="listener.id_listener"
          class="listener-card"
        >
          <div class="listener-card__content">
            <strong>{{ listener.second_name }} {{ listener.first_name }} {{ listener.middle_name || '' }}</strong>
            <span>{{ listener.snils }}</span>
          </div>

          <AppButton variant="ghost" @click="$emit('openListener', listener.id_listener)">
            Перейти
          </AppButton>
        </article>
      </div>
    </AppCard>

    <AppCard title="Файлы юридического лица">
      <p v-if="filesLoading" class="empty-state">Загрузка файлов...</p>
      <p v-else-if="!files.length" class="empty-state">Файлы не найдены.</p>

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
        <AppButton variant="secondary" @click="$emit('addListener')">Добавить слушателя</AppButton>
        <AppButton variant="secondary" @click="$emit('enroll')">Записать на курс</AppButton>
        <AppButton @click="$emit('edit')">Редактировать</AppButton>
      </div>
    </div>
  </section>
</template>

<style scoped>
.details-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 1rem;
}

.fact-list,
.listener-list,
.file-list {
  display: grid;
  gap: 0.8rem;
}

.fact-list p {
  margin: 0;
  color: #334155;
}

.fact-list strong {
  color: #0f172a;
}

.listener-card,
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

.listener-card__content {
  display: grid;
  gap: 0.25rem;
  min-width: 0;
}

.listener-card__content strong,
.file-row__name {
  overflow-wrap: anywhere;
  word-break: break-word;
}

.listener-card__content span {
  color: #64748b;
}

.empty-state {
  margin: 0;
  color: #64748b;
}

.details-actions,
.details-actions__group {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
}

@media (max-width: 1100px) {
  .details-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 720px) {
  .listener-card,
  .file-row,
  .details-actions,
  .details-actions__group {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
