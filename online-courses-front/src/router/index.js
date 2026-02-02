import { createRouter, createWebHistory } from 'vue-router'
import Login from '../components/Login.vue'
import AdminDashboard from '../components/AdminDashboard.vue'
import UserDashboard from '../components/UserDashboard.vue'
import AccountantDashboard from '../components/AccountantDashboard.vue'
import Listeners from '../components/Listeners.vue'
import Executers from '../components/Executers.vue'


const routes = [
  { path: '/', redirect: '/login' },
  { path: '/login', component: Login },

  { path: '/dashboard/admin', component: AdminDashboard },
  { path: '/dashboard/worker', component: UserDashboard },
  { path: '/dashboard/accountant', component: AccountantDashboard },

  { path: '/listeners', component: Listeners },
{ path: '/listeners/create', component: () => import('../components/CreateListener.vue') },
{ path: '/listeners/edit/:id', component: () => import('../components/EditListener.vue') },
{ path: '/listeners/:id', component: () => import('../components/ListenerDetails.vue') },

  { path: '/executers', component: Executers },
{ path: '/executers/create', component: () => import('../components/CreateExecuters.vue') },


  { path: '/divisions', component: () => import('../components/Divisions.vue') },
  { path: '/divisions/create', component: () => import('../components/CreateDivision.vue') },
  { path: '/divisions/edit/:id', component: () => import('../components/EditDivisions.vue') },

  { path: '/types', component: () => import('../components/EducationTypes.vue') },
  { path: '/types/create', component: () => import('../components/CreateEducationType.vue') },
  { path: '/types/edit/:id', component: () => import('../components/EditEducationType.vue') },

  { path: '/levels', component: () => import('../components/EducationLevels.vue') },

  { path: '/programs', component: () => import('../components/ProgramsEducation.vue') },
  { path: '/programs/create', component: () => import('../components/ProgramEducationCreate.vue') },
  { path: '/programs/edit/:id', component: () => import('../components/ProgramEducationEdit.vue') },

  { path: '/enrollments', component: () => import('../components/EnrollmentList.vue') },
  { path: '/enrollment/details/:listenerId', component: () => import('../components/EnrollmentDetails.vue') },
  { path: '/enrollment/edit/:idStudent/:idProgram', component: () => import('../components/EnrollmentEdit.vue') },
  { path: '/enrollment/create/:listenerId', component: () => import('../components/EnrollmentCreate.vue') },
  { path: '/enrollment/by-course', component: () => import('../components/EnrollmentByCourse.vue') },

  { path: '/legalentities', component: () => import('../components/LegalEntity.vue') },
   { path: '/legalentity/create', component: () => import('../components/CreateLegalEntity.vue') },
   { path: '/legalentity/edit/:id', component: () => import('../components/EditLegalEntity.vue') },
{ path: '/legalentity/:id', component: () => import('../components/LegalEntityDetails.vue') },

{ path: '/enrollment/accurate', component: () => import('../components/AccurateEnrollments.vue') },
{ path: '/enrollment/yur/:id', component: () => import('../components/EnrollmentYUR.vue') },

]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('access_token')
  const role = localStorage.getItem('role')
  const expiresAt = localStorage.getItem('access_token_expire_at')

  if (expiresAt && new Date(expiresAt) < new Date()) {
    localStorage.clear()
    return next('/login')
  }

  if (!token && to.path !== '/login') {
    return next('/login')
  }

  if (token && to.path === '/login') {
    switch (role) {
      case 'admin':
        return next('/dashboard/admin')
      case 'worker':
        return next('/dashboard/worker')
      case 'accountant':
        return next('/dashboard/accountant')
    }
  }

  if (to.path.startsWith('/dashboard')) {
    if (role === 'admin' && !to.path.includes('/admin')) return next('/dashboard/admin')
    if (role === 'worker' && !to.path.includes('/worker')) return next('/dashboard/worker')
    if (role === 'accountant' && !to.path.includes('/accountant')) return next('/dashboard/accountant')
  }

  next()
})

export default router
