type ApiMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'

type ApiRequestOptions = {
  method?: ApiMethod
  query?: Record<string, string | number | boolean | undefined | null>
  body?: unknown
  headers?: HeadersInit
  auth?: boolean
}

const buildUrl = (
  baseUrl: string,
  path: string,
  query?: ApiRequestOptions['query']
) => {
  const url = new URL(path, `${baseUrl}/`)

  if (!query) {
    return url.toString()
  }

  for (const [key, value] of Object.entries(query)) {
    if (value === undefined || value === null || value === '') {
      continue
    }

    url.searchParams.set(key, String(value))
  }

  return url.toString()
}

export const useApiClient = () => {
  const config = useRuntimeConfig()
  const auth = useAuthState()
  const coreBaseUrl = import.meta.server ? config.apiUrlCoreInternal : config.public.apiUrlCore
  const authBaseUrl = import.meta.server ? config.apiUrlAuthInternal : config.public.apiUrlAuth
  const docBaseUrl = import.meta.server ? config.apiUrlDocInternal : config.public.apiUrlDoc

  const request = async <T>(baseUrl: string, path: string, options: ApiRequestOptions = {}) => {
    const headers = new Headers(options.headers)

    if (options.body !== undefined && !headers.has('Content-Type')) {
      headers.set('Content-Type', 'application/json')
    }

    if (options.auth !== false && auth.value.accessToken) {
      headers.set('Authorization', `Bearer ${auth.value.accessToken}`)
    }

    return await $fetch<T>(buildUrl(baseUrl, path, options.query), {
      method: options.method || 'GET',
      headers,
      body: options.body as BodyInit | Record<string, any> | null | undefined
    })
  }

  return {
    core: <T>(path: string, options?: ApiRequestOptions) =>
      request<T>(coreBaseUrl, path, options),
    auth: <T>(path: string, options?: ApiRequestOptions) =>
      request<T>(authBaseUrl, path, options),
    doc: <T>(path: string, options?: ApiRequestOptions) =>
      request<T>(docBaseUrl, path, options)
  }
}
