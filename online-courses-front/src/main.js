import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import router from './router'
import 'bootstrap/dist/css/bootstrap.min.css'
import 'bootstrap'
import 'bootstrap-icons/font/bootstrap-icons.css'
import Vue3Toastify from 'vue3-toastify'
import 'vue3-toastify/dist/index.css'
import './assets/theme.css'

const app = createApp(App)

app.use(router)
app.use(Vue3Toastify, { autoClose: 3000 })

app.mount('#app')
