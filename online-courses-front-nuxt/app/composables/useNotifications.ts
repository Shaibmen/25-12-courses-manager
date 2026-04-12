type NotificationVariant = 'success' | 'error' | 'info'

type AppNotification = {
  id: string
  message: string
  title?: string
  variant: NotificationVariant
  duration: number
}

const DEFAULT_DURATION_BY_VARIANT: Record<NotificationVariant, number> = {
  success: 9000,
  error: 12000,
  info: 9000
}

export const useNotifications = () => {
  const notifications = useState<AppNotification[]>('app-notifications', () => [])

  const remove = (id: string) => {
    notifications.value = notifications.value.filter((item) => item.id !== id)
  }

  const push = (
    message: string,
    options?: {
      title?: string
      variant?: NotificationVariant
      duration?: number
    }
  ) => {
    const variant = options?.variant || 'info'
    const notification: AppNotification = {
      id: `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
      title: options?.title,
      message,
      variant,
      duration: options?.duration ?? DEFAULT_DURATION_BY_VARIANT[variant]
    }

    notifications.value = [...notifications.value, notification]

    if (import.meta.client && notification.duration > 0) {
      window.setTimeout(() => {
        remove(notification.id)
      }, notification.duration)
    }

    return notification.id
  }

  return {
    notifications: readonly(notifications),
    push,
    remove,
    success: (message: string, title?: string, duration?: number) =>
      push(message, { title, duration, variant: 'success' }),
    error: (message: string, title?: string, duration?: number) =>
      push(message, { title, duration, variant: 'error' }),
    info: (message: string, title?: string, duration?: number) =>
      push(message, { title, duration, variant: 'info' })
  }
}
