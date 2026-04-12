export default defineNuxtRouteMiddleware((to) => {
  const auth = useAuth()

  auth.hydrateFromCookies()

  if (auth.isExpired.value) {
    auth.clear()
  }

  if (to.path === '/login' && auth.isAuthenticated.value) {
    return navigateTo('/dashboard')
  }

  if (to.path === '/login') {
    return
  }

  if (!auth.isAuthenticated.value) {
    return navigateTo('/login')
  }
})
