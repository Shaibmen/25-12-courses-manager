import type { UserRole } from '../../types/auth'
import type {
  AdminDashboardResponse,
  DashboardPayload,
  UserDashboardResponse
} from '../../types/dashboard'

export const getWorkerDashboard = async () => {
  const api = useApiClient()

  return await api.core<UserDashboardResponse>('dashboard/user')
}

export const getAdminDashboard = async () => {
  const api = useApiClient()

  return await api.core<AdminDashboardResponse>('dashboard/admin')
}

export const getDashboardByRole = async (role: UserRole): Promise<DashboardPayload> => {
  if (role === 'admin') {
    return {
      kind: 'admin',
      data: await getAdminDashboard()
    }
  }

  return {
    kind: 'user',
    data: await getWorkerDashboard()
  }
}
