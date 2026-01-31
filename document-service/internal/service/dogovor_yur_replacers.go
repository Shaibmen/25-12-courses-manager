package service

import (
	"document-service/internal/domain/dto"
	"fmt"
	"strings"

	"github.com/nguyenthenguyen/docx"
)

const PRIMECHANIE_ODIN = "Подписанный Акт является основанием для выставления счёта."
const PRIMECHANIE_DVA = "подтверждает согласие Заказчика с условиями настоящего договора и исполнение обязанности по оплате."
const PRIMECHAINE_TRI = "подтверждает надлежащее исполнение Заказчиком обязанности по оплате услуг, оказанных Исполнителем."

func replacePP3YUR(doc *docx.Docx, model *dto.DogovorDTO) {

	doc.Replace("ZAKAZCHIKSTATUS", model.Zakazchik.Status, -1)
	doc.Replace("ZAKAZCHIKFIO", model.Zakazchik.FIO, -1)
	doc.Replace("INN", model.Zakazchik.INN, -1)
	doc.Replace("KPP", model.Zakazchik.KPP, -1)
	doc.Replace("OGRN", model.Zakazchik.OGRN, -1)
	doc.Replace("ZAKAZCHIKPHONE", model.Zakazchik.Phone, -1)
	doc.Replace("ZSTATUS", model.Zakazchik.Status, -1)
	doc.Replace("ZFIO", model.Zakazchik.FIO, -1)

	cityStreetHouseBuildingApartment := fmt.Sprintf("%s, %s, %s, %s, %s. ",
		model.Zakazchik.Address.City,
		model.Zakazchik.Address.Street,
		model.Zakazchik.Address.House,
		model.Zakazchik.Address.Building,
		model.Zakazchik.Address.Apartment)

	doc.Replace("ZAKAZCHIKADDRESS", cityStreetHouseBuildingApartment, -1)
	doc.Replace("ZAKAZCHIKPHONE", model.Zakazchik.Phone, -1)
	doc.Replace("ZAKAZCHIKEMAIL", model.Zakazchik.Email, -1)

	doc.Replace("STATUS", model.Executor.Status, -1)

	executorFio := model.Executor.ExecutorSurname + " " + model.Executor.ExecutorName + " " + model.Executor.ExecutorMiddlename
	doc.Replace("EXECUTORFIO", executorFio, -1)

	doc.Replace("NAMEPROFEDUCATION", model.ProgramEducation.NameProfEducation, -1)
	doc.Replace("SINCE", model.Enrollment.StartDate.Format("02.01.2006"), -1)
	doc.Replace("FOR", model.Enrollment.EndDate.Format("02.01.2006"), -1)

	doc.Replace("ZAKAZCHI", model.Zakazchik.CompanyName, -1)
	doc.Replace("K", "", -1)

	if strings.Contains(model.OptionPrice, "Оплата подлежит перечислению на расчётный счёт Исполнителя в срок до 5 (пяти) рабочих дней") { // третья опция оплаты
		doc.Replace("PRIMECHANIEODIN", PRIMECHANIE_ODIN, -1)
		doc.Replace("PRIMECHANIEDVA", PRIMECHAINE_TRI, -1)
	} else {
		doc.Replace("PRIMECHANIEODIN", "", -1)
		doc.Replace("PRIMECHANIEDVA", PRIMECHANIE_DVA, -1)
	}
}

func replacePK3YUR(doc *docx.Docx, model *dto.DogovorDTO) {
	doc.Replace("ZAKAZCHIKSTATUS", model.Zakazchik.Status, -1)
	doc.Replace("ZAKAZCHIKFIO", model.Zakazchik.FIO, -1)
	doc.Replace("ZFIO", model.Zakazchik.FIO, -1)
	doc.Replace("INN", model.Zakazchik.INN, -1)
	doc.Replace("KPP", model.Zakazchik.KPP, -1)
	doc.Replace("OGRN", model.Zakazchik.OGRN, -1)
	doc.Replace("ZAKAZCHIKPHONE", model.Zakazchik.Phone, -1)
	doc.Replace("ZSTATUS", model.Zakazchik.Status, -1)

	cityStreetHouseBuildingApartment := fmt.Sprintf("%s, %s, %s, %s, %s. ",
		model.Zakazchik.Address.City,
		model.Zakazchik.Address.Street,
		model.Zakazchik.Address.House,
		model.Zakazchik.Address.Building,
		model.Zakazchik.Address.Apartment)

	doc.Replace("ZAKAZCHIKADDRESS", cityStreetHouseBuildingApartment, -1)
	doc.Replace("ZAKAZCHIKPHONE", model.Zakazchik.Phone, -1)
	doc.Replace("ZAKAEMAIL", model.Zakazchik.Email, -1)

	doc.Replace("STATUS", model.Executor.Status, -1)

	doc.Replace("NAMEPROFEDUCATION", model.ProgramEducation.NameProfEducation, -1)
	doc.Replace("SINCE", model.Enrollment.StartDate.Format("02.01.2006"), -1)
	doc.Replace("FOR", model.Enrollment.EndDate.Format("02.01.2006"), -1)

	executorFio := model.Executor.ExecutorSurname + " " + model.Executor.ExecutorName + " " + model.Executor.ExecutorMiddlename
	doc.Replace("EXECUTORFIO", executorFio, -1)

	doc.Replace("NAMEPROFEDUCATION", model.ProgramEducation.NameProfEducation, -1)
	doc.Replace("SINCE", model.Enrollment.StartDate.Format("02.01.2006"), -1)
	doc.Replace("FOR", model.Enrollment.EndDate.Format("02.01.2006"), -1)

	doc.Replace("ZAKAZCHIK", model.Zakazchik.CompanyName, -1)
	doc.Replace("ZAKA", model.Zakazchik.CompanyName, -1)

	if strings.Contains(model.OptionPrice, "Оплата подлежит перечислению на расчётный счёт Исполнителя в срок до 5 (пяти) рабочих дней") { // третья опция оплаты
		doc.Replace("PRIMECHANIEODIN", PRIMECHANIE_ODIN, -1)
		doc.Replace("PRIMECHANIEDVA", PRIMECHAINE_TRI, -1)
	} else {
		doc.Replace("PRIMECHANIEODIN", "", -1)
		doc.Replace("PRIMECHANIEDVA", PRIMECHANIE_DVA, -1)
	}
}
