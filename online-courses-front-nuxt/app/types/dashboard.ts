export type ProgramEndingSoon = {
  name_prof_education: string
  end_date: string
  total_listeners: number
}

export type UserDashboardResponse = {
  total_listener: number
  total_program: number
  active_enrollments: number
  program_ending_soon: ProgramEndingSoon[]
}

export type AdminStat = {
  database_name?: string
  total_size?: string
  active_connection?: number
  committed_tx?: number
  rolledback_tx?: number
  disk_block_read?: number
  buffer_hits?: number
}

export type AdminProcessStat = {
  pid: number
  state?: string
  client_addr?: string
  client_port?: number
  query_start?: string
}

export type AdminDashboardResponse = {
  admin_stat: AdminStat
  pg_stat: AdminProcessStat[]
  role?: Array<{
    id: string | number
    role: string
  }>
}

export type DashboardPayload =
  | {
      kind: 'admin'
      data: AdminDashboardResponse
    }
  | {
      kind: 'user'
      data: UserDashboardResponse
    }
