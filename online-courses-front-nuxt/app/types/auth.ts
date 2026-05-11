export type UserRole = 'admin' | 'worker' | 'accountant'

export type AuthState = {
  accessToken: string | null
  refreshToken: string | null
  username: string | null
  role: UserRole | null
  accessTokenExpireAt: string | null
  refreshTokenExpireAt: string | null
}

export type AuthCredentials = {
  username: string
  password: string
}

export type RegisterUserPayload = {
  username: string
  password: string
  role: string
}

export type LoginTokenPayload = {
  access_token: string
  refresh_token: string
  access_token_expire_at: string
  refresh_token_expire_at: string
  user: {
    username: string
    role: UserRole
  }
}

export type LoginResponse = {
  token: LoginTokenPayload
}

export type RegisterUserResponse = {
  message?: string
  error?: string | Record<string, unknown>
}
