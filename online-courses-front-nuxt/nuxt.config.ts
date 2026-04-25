const env = (
  (globalThis as { process?: { env?: Record<string, string | undefined> } }).process?.env ??
  {}
)

export default defineNuxtConfig({
  compatibilityDate: '2026-04-11',
  devtools: { enabled: true },
  modules: ['@pinia/nuxt'],
  css: ['~/assets/styles/main.css'],
  app: {
    head: {
      title: 'Система курсов 25-12',
      meta: [
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'theme-color', content: '#0f172a' }
      ]
    }
  },
  runtimeConfig: {
    public: {
      apiUrlCore: env.NUXT_PUBLIC_API_URL_CORE || 'http://localhost:8080/api/v1',
      apiUrlAuth: env.NUXT_PUBLIC_API_URL_AUTH || 'http://localhost:8081/auth/v1',
      apiUrlDoc: env.NUXT_PUBLIC_API_URL_DOC || 'http://localhost:8082/v1/doc'
    }
  },
  imports: {
    dirs: ['services/**', 'stores/**', 'types/**']
  }
})
