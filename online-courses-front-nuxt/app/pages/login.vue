<script setup lang="ts">
import LoginForm from '../components/features/auth/LoginForm.vue'

definePageMeta({
  layout: 'auth'
})
</script>

<template>
  <section class="page-grid">
    <LoginForm />
  </section>
</template>
