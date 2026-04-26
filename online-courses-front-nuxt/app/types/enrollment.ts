import type { ExecuterItem } from './executer'
import type { LegalEntityDetailsResponse } from './legalentity'
import type { ProgramListItem } from './program'

export type ContractOption = {
  id_contract: string
  name: string
  type: 'bilateral' | 'trilateral'
}

export const singleEnrollmentContracts: ContractOption[] = [
  { id_contract: 'DO_3_FIZ', name: 'ДО с оплатой физическим лицом', type: 'trilateral' },
  { id_contract: 'PK_2_FIZ', name: 'ПК с оплатой физическим лицом', type: 'bilateral' },
  { id_contract: 'PK_3_FIZ', name: 'ПК с оплатой физическим лицом', type: 'trilateral' },
  { id_contract: 'PK_3_YUR', name: 'ПК с оплатой юридическим лицом', type: 'trilateral' },
  { id_contract: 'PP_2_FIZ', name: 'ПП с оплатой физическим лицом', type: 'bilateral' },
  { id_contract: 'PP_3_FIZ', name: 'ПП с оплатой физическим лицом', type: 'trilateral' },
  { id_contract: 'PP_3_YUR', name: 'ПП с оплатой юридическим лицом', type: 'trilateral' }
]

export const legalEntityEnrollmentContracts: ContractOption[] = [
  { id_contract: 'PK_3_YUR', name: 'ПК с оплатой юридическим лицом', type: 'trilateral' },
  { id_contract: 'PP_3_YUR', name: 'ПП с оплатой юридическим лицом', type: 'trilateral' }
]

export const loadVariantsDO = {
  1: 'с пониженной недельной учебной нагрузкой (1 акад. час в неделю)',
  2: 'с умеренной недельной учебной нагрузкой (2 акад. часа в неделю)',
  3: 'со стандартной недельной учебной нагрузкой (3 акад. часа в неделю)',
  4: 'с высокой недельной учебной нагрузкой (4 акад. часа в неделю)',
  5: 'с повышенной недельной учебной нагрузкой (6 акад. часов в неделю)'
} as const

export const loadVariantsNotDO = {
  1: 'с пониженной недельной учебной нагрузкой (3 акад. часа в неделю)',
  2: 'с умеренной недельной учебной нагрузкой (6 акад. часов в неделю)',
  3: 'со стандартной недельной учебной нагрузкой (12 акад. часов в неделю)',
  4: 'с высокой недельной учебной нагрузкой (15 акад. часов в неделю)',
  5: 'с повышенной недельной учебной нагрузкой (30 акад. часов в неделю)',
  6: 'с интенсивной недельной учебной нагрузкой (36 акад. часов в неделю)'
} as const

export const ageCategories = {
  BELOW_EIGHTEEN: 'Меньше восемнадцати',
  FOURTEEN: 'Меньше четырнадцати',
  EIGHTEEN: 'Восемнадцать'
} as const

export const optDocumentOptions = {
  1: 'удостоверение о повышении квалификации вручается по окончании',
  2: 'удостоверение выдаётся одновременно с дипломом СПО/ВО. До этого момента хранится у Исполнителя.'
} as const

export type ContractorCore = {
  id_contractor?: string
  first_name: string
  second_name: string
  middle_name: string
  contact_phone: string
  email: string
}

export type ContractorPassport = {
  place_birth: string
  citizenship: string
  gender: string
  seria: string
  number: string
  passport_given: string
  date_given: string
  code: string
}

export type ContractorRegistrationAddress = {
  mail_index: string
  region: string
  city: string
  street: string
  house: string
  building: string
  apartment: string
}

export type ContractorPayload = {
  contractor: ContractorCore
  passport: ContractorPassport
  reg_address: ContractorRegistrationAddress
}

export type ListenerEnrollmentContext = {
  listener: {
    id_listener: string
    first_name: string
    second_name: string
    middle_name?: string
  }
  contractor?: ContractorPayload | null
}

export type EnrollmentItem = {
  id_listener: string
  first_name: string
  second_name: string
  middle_name?: string
  name_prof_education: string
  start_date: string
  end_date: string
  id_group: string
  type_of_retraining: string
}

export type AccurateEnrollmentItem = {
  id_listener: string
  name_prof_education: string
  time_education: number
  price: number
  education_type: string
  division: string
}

export type EnrollmentPayload = {
  id_listener: string
  id_program: string
  start_date: string
  end_date: string
  id_group: string
  type_of_retraining: string
  is_active: boolean
}

export type DocumentPayload = {
  id_listener: string
  id_program: string
  id_executor: string | null
  front_data: Record<string, unknown>
}

export type EnrollmentPageContext = {
  programs: ProgramListItem[]
  executers: ExecuterItem[]
}

export type LegalEntityEnrollmentContext = LegalEntityDetailsResponse
