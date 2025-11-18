<template>
  <div class="page-container">
    <Header />
    <nav class="nav-bar">
      <button class="btn btn-primary btn-sm" @click="$router.push('/listeners')">Слушатели</button>
      <button class="btn btn-primary btn-sm" @click="$router.push('/divisions')">Подразделения</button>
      <button class="btn btn-primary btn-sm" @click="$router.push('/types')">Типы обучения</button>
      <button class="btn btn-primary btn-sm" @click="$router.push('/levels')">Уровни</button>
      <button class="btn btn-primary btn-sm" @click="$router.push('/programs')">Программы</button>
      <button class="btn btn-primary btn-sm" @click="$router.push('/enrollments')">Записи на курс</button>
    </nav>
    <main>
      <slot />
    </main>
    <Footer />
  </div>
</template>

<script setup>
import Header from '../components/Header.vue'
import Footer from '../components/Footer.vue'
</script>

<style scoped>
.page-container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  width: 100%;
  margin: 0;
}

main {
  flex: 1 0 auto;
  width: 100%;
  padding: 16px;
  box-sizing: border-box;
}

.nav-bar {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  background-color: #f4f5f7;
  border-bottom: 1px solid #ccc;
  padding: 8px 0;
  width: 100%;
}
</style>
