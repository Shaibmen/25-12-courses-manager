package service

import (
	"context"
	"online-courses/internal/domain/repository"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

var (
	pathReports = "reports/excel/"
)

type ReportService struct {
	service repository.ReportRepository
}

func NewRepostService(repo repository.ReportRepository) *ReportService {
	return &ReportService{service: repo}
}

func (r *ReportService) ReportPeriod(ctx context.Context, startDate, endDate time.Time) (string, error) {
	data, err := r.service.ReportPeriod(ctx, startDate, endDate)
	if err != nil {
		return "", err
	}

	f := excelize.NewFile()
	index, err := f.NewSheet("Sheet1")
	if err != nil {
		return "", err
	}

	f.SetCellValue("Sheet1", "A1", "Отчётная ведомость за выбранный период: "+startDate.Format("2006-01-02")+"/"+endDate.Format("2006-01-02"))
	f.SetCellValue("Sheet1", "A3", "№")
	f.SetCellValue("Sheet1", "B3", "ФИО")
	f.SetCellValue("Sheet1", "C3", "Код дохода")
	f.SetCellValue("Sheet1", "D3", "Программа")
	f.SetCellValue("Sheet1", "E3", "Период курса")
	f.SetCellValue("Sheet1", "F3", "Доход за период (₽)")
	f.SetCellValue("Sheet1", "G3", "Удержано НДФЛ (13%)")
	f.SetCellValue("Sheet1", "H3", "")
	f.SetCellValue("Sheet1", "I3", "К перечислению (₽)")

	rowA, rowB, rowC, rowD, rowE, rowF, rowG, rowH, rowI := 4, 4, 4, 4, 4, 4, 4, 4, 4
	num := 1

	var totalPrice float32

	for _, l := range data {
		ndfl := l.Payment * 13 / 100
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(rowA), num)
		rowA++
		num++

		f.SetCellValue("Sheet1", "B"+strconv.Itoa(rowB), l.FullName)
		rowB++

		f.SetCellValue("Sheet1", "C"+strconv.Itoa(rowC), "2000")
		rowC++

		f.SetCellValue("Sheet1", "D"+strconv.Itoa(rowD), l.ProgramName)
		rowD++

		f.SetCellValue("Sheet1", "E"+strconv.Itoa(rowE), l.StartDate.Format("2006-01-02")+"-"+l.EndDate.Format("2006-01-02"))
		rowE++

		f.SetCellValue("Sheet1", "F"+strconv.Itoa(rowF), l.Payment)
		rowF++

		f.SetCellValue("Sheet1", "G"+strconv.Itoa(rowG), ndfl)
		rowG++

		f.SetCellValue("Sheet1", "H"+strconv.Itoa(rowH), "")
		rowH++

		f.SetCellValue("Sheet1", "I"+strconv.Itoa(rowI), l.Payment-ndfl)
		rowI++

		totalPrice = totalPrice + (l.Payment - ndfl)
	}

	f.SetCellValue("Sheet1", "J3", "Итого")
	f.SetCellValue("Sheet1", "J4", totalPrice)

	f.SetColWidth("Sheet1", "A", "A", 5)
	f.SetColWidth("Sheet1", "B", "B", 35)
	f.SetColWidth("Sheet1", "C", "C", 12)
	f.SetColWidth("Sheet1", "D", "D", 55)
	f.SetColWidth("Sheet1", "E", "K", 22)

	f.SetActiveSheet(index)
	nameFile := pathReports + "отчёт_за_" + startDate.Format("2006-01-02") + "_" + endDate.Format("2006-01-02") + ".xlsx"
	if err := f.SaveAs(nameFile); err != nil {
		return "", err
	}

	return nameFile, nil

}
func (r *ReportService) MostExpensiveProgram(ctx context.Context, startDate, endDate time.Time) (string, error) {

	data, err := r.service.MostExpensiveProgram(ctx, startDate, endDate)
	if err != nil {
		return "", err
	}

	f := excelize.NewFile()
	index, err := f.NewSheet("Sheet1")
	if err != nil {
		return "", err
	}

	f.SetCellValue("Sheet1", "A1", "Отчётная ведомость за выбранный период: "+startDate.Format("2006-01-02")+"/"+endDate.Format("2006-01-02"))
	f.SetCellValue("Sheet1", "A3", "№")
	f.SetCellValue("Sheet1", "B3", "Курс")
	f.SetCellValue("Sheet1", "C3", "Кол-во слушателей")
	f.SetCellValue("Sheet1", "D3", "Общая сумма(₽)")

	rowA, rowB, rowC, rowD := 4, 4, 4, 4
	num := 1

	var totalPrice float32
	var bestRevenue float32
	var bestProgram string

	for _, l := range data {

		f.SetCellValue("Sheet1", "A"+strconv.Itoa(rowA), num)
		rowA++
		num++

		f.SetCellValue("Sheet1", "B"+strconv.Itoa(rowB), l.ProgramName)
		rowB++

		f.SetCellValue("Sheet1", "C"+strconv.Itoa(rowC), l.TotalListeners)
		rowC++

		f.SetCellValue("Sheet1", "D"+strconv.Itoa(rowD), l.TotalRevenue)
		rowD++

		if l.TotalRevenue > bestRevenue {
			bestRevenue = l.TotalRevenue
			bestProgram = l.ProgramName
		}

		totalPrice = totalPrice + l.TotalRevenue
	}
	f.SetCellValue("Sheet1", "E3", "Самый дорогой курс")
	f.SetCellValue("Sheet1", "E4", bestProgram)
	f.SetCellValue("Sheet1", "E5", bestRevenue)

	f.SetCellValue("Sheet1", "F3", "Итого")
	f.SetCellValue("Sheet1", "F4", totalPrice)

	f.SetColWidth("Sheet1", "A", "A", 5)
	f.SetColWidth("Sheet1", "B", "B", 55)
	f.SetColWidth("Sheet1", "C", "C", 20)
	f.SetColWidth("Sheet1", "D", "D", 55)
	f.SetColWidth("Sheet1", "E", "E", 55)
	f.SetColWidth("Sheet1", "F", "F", 12)

	f.SetActiveSheet(index)
	nameFile := pathReports + "отчёт_программ_за_" + startDate.Format("2006-01-02") + "_" + endDate.Format("2006-01-02") + ".xlsx"
	if err := f.SaveAs(nameFile); err != nil {
		return "", err
	}

	return nameFile, nil

}
