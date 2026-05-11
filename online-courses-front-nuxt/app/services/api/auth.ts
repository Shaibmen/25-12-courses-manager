import type {
  AuthCredentials,
  LoginResponse,
  RegisterUserPayload,
  RegisterUserResponse
} from '../../types/auth'

export const loginRequest = async (credentials: AuthCredentials) => {
  const api = useApiClient()

  try {
    return await api.auth<LoginResponse>('login', {
      method: 'POST',
      body: credentials,
      auth: false
    })
  } catch (error) {
    const statusCode = typeof error === 'object' && error !== null && 'statusCode' in error
      ? Number(error.statusCode)
      : null

    if (statusCode === 401) {
      throw new Error('Неверный логин или пароль')
    }

    throw new Error('Не удалось выполнить вход')
  }
}

export const logoutRequest = async () => {
  const api = useApiClient()

  return await api.auth('protected/logout', {
    method: 'POST'
  })
}

export const registerUserRequest = async (payload: RegisterUserPayload) => {
  const api = useApiClient()

  try {
    return await api.auth<RegisterUserResponse>('register', {
      method: 'POST',
      body: payload
    })
  } catch (error) {
    const statusCode = typeof error === 'object' && error !== null && 'statusCode' in error
      ? Number(error.statusCode)
      : null

    const responseMessage =
      typeof error === 'object' &&
      error !== null &&
      'data' in error &&
      typeof error.data === 'object' &&
      error.data !== null &&
      'message' in error.data &&
      typeof error.data.message === 'string'
        ? error.data.message
        : null

    if (statusCode === 401) {
      throw new Error(responseMessage || 'Регистрация доступна только администратору')
    }

    if (responseMessage) {
      throw new Error(responseMessage)
    }

    throw new Error('Не удалось зарегистрировать пользователя')
  }
}
