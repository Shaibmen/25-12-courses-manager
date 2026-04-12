export type LegalEntityListItem = {
  id_legalentity: string
  name_company: string
  inn: string
  kpp: string
  ogrn: string
  phone: string
  email: string
  first_name: string
  second_name: string
  middle_name?: string
  status: string
}

export type LegalEntityCore = {
  id_legalentity?: string
  name_company: string
  inn: string
  kpp: string
  ogrn: string
  phone: string
  email: string
  first_name: string
  second_name: string
  middle_name: string
  status: string
}

export type LegalEntityRegistrationAddress = {
  mail_index: string
  region: string
  city: string
  street: string
  house: string
  building: string
  apartment: string
}

export type LegalEntityListener = {
  id_listener: string
  first_name: string
  second_name: string
  middle_name?: string
  snils: string
}

export type LegalEntityDetailsResponse = {
  legal_entity: LegalEntityCore & {
    listeners?: LegalEntityListener[]
    id_regaddress?: string
  }
  reg_address: LegalEntityRegistrationAddress
}

export type LegalEntityPayload = {
  legal_entity: LegalEntityCore
  reg_address: LegalEntityRegistrationAddress
}

export type LegalEntityFormState = {
  legal: LegalEntityCore
  regAddress: LegalEntityRegistrationAddress
}

export const createEmptyLegalEntityFormState = (): LegalEntityFormState => ({
  legal: {
    name_company: '',
    inn: '',
    kpp: '',
    ogrn: '',
    phone: '',
    email: '',
    first_name: '',
    second_name: '',
    middle_name: '',
    status: ''
  },
  regAddress: {
    mail_index: '',
    region: '',
    city: '',
    street: '',
    house: '',
    building: '',
    apartment: ''
  }
})
