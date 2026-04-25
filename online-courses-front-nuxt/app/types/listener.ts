export type ListenerListItem = {
  id_listener: string
  first_name: string
  second_name: string
  middle_name: string
  date_of_birth: string
  snils: string
  contact_phone: string
  email: string
}

export type ListenerCore = {
  first_name: string
  second_name: string
  middle_name: string
  date_of_birth: string
  snils: string
  contact_phone: string
  email: string
  looting_education: boolean
  id_legalentity?: string
  id_contractor?: string
}

export type PassportData = {
  place_birth: string
  citizenship: string
  gender: string
  seria: string
  number: string
  passport_given: string
  date_given: string
  code: string
}

export type RegistrationAddressData = {
  mail_index: string
  region: string
  city: string
  street: string
  house: string
  building: string
  apartment: string
}

export type EducationData = {
  diplom_seria: string
  diplom_number: string
  date_given: string
  city: string
  region: string
  educational_institution: string
  speciality: string
  level_education: string
}

export type PlaceWorkData = {
  name_company: string
  job_title: string
  all_experience: number
  job_title_expirience: number
}

export type PlaceWorkFormData = {
  name_company: string
  job_title: string
  all_experience: string | number
  job_title_expirience: string | number
}

export type ListenerFormPayload = {
  listener: ListenerCore
  registration_address: RegistrationAddressData
  passport?: PassportData
  education?: EducationData
  placeWork?: PlaceWorkData
}

export type ListenerDetailsResponse = {
  listener: ListenerCore & {
    id_listener: string
    id_passport?: string | null
    id_reg_address?: string
    id_education_listener?: string | null
    id_placework?: string | null
  }
  passport?: PassportData | null
  regaddress: RegistrationAddressData
  education_listener?: EducationData | null
  placework?: {
    name_company: string
    job_title: string
    all_experience: number | string
    job_title_experience?: number | string
    job_title_expirience?: number | string
  } | null
}

export type LevelEducationItem = {
  id_level_education: string
  education: string
}

export type EnrollmentProgramDetails = {
  id_listener: string
  id_program_education: string
  name_prof_education: string
  time_education: number
  price: number
  education_type: string
  division_education: string
  start_date: string
  end_date: string
  id_group: string
  type_of_retraining: string
}

export type ListenerFormState = {
  listener: ListenerCore
  passport: PassportData
  registrationAddress: RegistrationAddressData
  education: EducationData
  placeWork: PlaceWorkFormData
}

export const createEmptyListenerFormState = (): ListenerFormState => ({
  listener: {
    first_name: '',
    second_name: '',
    middle_name: '',
    date_of_birth: '',
    snils: '',
    contact_phone: '',
    email: '',
    looting_education: false
  },
  passport: {
    place_birth: '',
    citizenship: '',
    gender: '',
    seria: '',
    number: '',
    passport_given: '',
    date_given: '',
    code: ''
  },
  registrationAddress: {
    mail_index: '',
    region: '',
    city: '',
    street: '',
    house: '',
    building: '',
    apartment: ''
  },
  education: {
    diplom_seria: '',
    diplom_number: '',
    date_given: '',
    city: '',
    region: '',
    educational_institution: '',
    speciality: '',
    level_education: ''
  },
  placeWork: {
    name_company: '',
    job_title: '',
    all_experience: '',
    job_title_expirience: ''
  }
})
