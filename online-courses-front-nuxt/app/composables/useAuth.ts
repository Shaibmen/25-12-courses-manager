import type {
  AuthCredentials,
  AuthState,
  LoginTokenPayload,
  UserRole
} from '../types/auth'

const roleLabelsMap: Record<UserRole, string> = {
  admin: 'Администратор',
  worker: 'Сотрудник',
  accountant: 'Бухгалтер'
}

const buildStateFromPayload = (payload: LoginTokenPayload | null): AuthState => ({
  accessToken: payload?.access_token || null,
  refreshToken: payload?.refresh_token || null,
  username: payload?.user.username || null,
  role: payload?.user.role || null,
  accessTokenExpireAt: payload?.access_token_expire_at || null,
  refreshTokenExpireAt: payload?.refresh_token_expire_at || null
})

export const useAuth = () => {
  const state = useAuthState()

  const accessToken = useCookie<string | null>('courses_access_token', { sameSite: 'lax' })
  const refreshToken = useCookie<string | null>('courses_refresh_token', { sameSite: 'lax' })
  const username = useCookie<string | null>('courses_username', { sameSite: 'lax' })
  const role = useCookie<UserRole | null>('courses_role', { sameSite: 'lax' })
  const accessTokenExpireAt = useCookie<string | null>('courses_access_token_expire_at', { sameSite: 'lax' })
  const refreshTokenExpireAt = useCookie<string | null>('courses_refresh_token_expire_at', { sameSite: 'lax' })

  const hydrateFromCookies = () => {
    state.value = {
      accessToken: accessToken.value || null,
      refreshToken: refreshToken.value || null,
      username: username.value || null,
      role: role.value || null,
      accessTokenExpireAt: accessTokenExpireAt.value || null,
      refreshTokenExpireAt: refreshTokenExpireAt.value || null
    }
  }

  const persist = (payload: LoginTokenPayload) => {
    accessToken.value = payload.access_token
    refreshToken.value = payload.refresh_token
    username.value = payload.user.username
    role.value = payload.user.role
    accessTokenExpireAt.value = payload.access_token_expire_at
    refreshTokenExpireAt.value = payload.refresh_token_expire_at
    state.value = buildStateFromPayload(payload)
  }

  const clear = () => {
    accessToken.value = null
    refreshToken.value = null
    username.value = null
    role.value = null
    accessTokenExpireAt.value = null
    refreshTokenExpireAt.value = null
    state.value = buildStateFromPayload(null)
  }

  const isExpired = computed(() => {
    const expiresAt = state.value.accessTokenExpireAt

    if (!expiresAt) {
      return false
    }

    return new Date(expiresAt).getTime() <= Date.now()
  })

  const isAuthenticated = computed(() =>
    Boolean(state.value.accessToken) && !isExpired.value
  )

  const roleLabel = computed(() =>
    state.value.role ? roleLabelsMap[state.value.role] : null
  )

  const login = async (credentials: AuthCredentials) => {
    const response = await loginRequest(credentials)

    if (!response?.token?.user) {
      throw new Error('Сервер вернул некорректный ответ')
    }

    persist(response.token)

    return response.token.user.role
  }

  const logout = async () => {
    hydrateFromCookies()

    try {
      if (state.value.accessToken) {
        await logoutRequest()
      }
    } catch {
      // Logout should still clear auth state even if the backend is unavailable.
    } finally {
      clear()
      await navigateTo('/login')
    }
  }

  return {
    state,
    isExpired,
    isAuthenticated,
    roleLabel,
    hydrateFromCookies,
    clear,
    login,
    logout
  }
}
