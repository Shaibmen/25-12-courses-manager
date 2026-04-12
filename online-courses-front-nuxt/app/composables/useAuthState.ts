import type { UserRole } from '../types/auth'

type StoredAuthState = {
  accessToken: string | null
  refreshToken: string | null
  username: string | null
  role: UserRole | null
  accessTokenExpireAt: string | null
  refreshTokenExpireAt: string | null
}

const authStateDefault = (): StoredAuthState => ({
  accessToken: null,
  refreshToken: null,
  username: null,
  role: null,
  accessTokenExpireAt: null,
  refreshTokenExpireAt: null
})

export const useAuthState = () => useState<StoredAuthState>('auth-state', authStateDefault)
