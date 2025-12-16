package hrm

import (
	employeeimp "europm/internal/hrm/employeeImp"
	"europm/internal/hrm/model"
	"europm/internal/util"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SearchChangesEmployee(c *gin.Context) {
	// Bind into DTO with string dates to support both JSON body and query params
	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")
	typeReport := c.Query("type_report")

	var fromDate, toDate time.Time
	var err error
	// Parse dates supporting RFC3339 and YYYY-MM-DD
	if fromDateStr != "" {
		fromDate, err = time.Parse("2006-01-02", fromDateStr)
		if err != nil {
			log.Printf("Error parsing from_date: %v", err)
			util.NewError(c, http.StatusBadRequest, err)
			c.JSON(http.StatusBadRequest, "invalid from_date format")
			return
		}
	}
	if toDateStr != "" {
		toDate, err = time.Parse("2006-01-02", toDateStr)
		if err != nil {
			log.Printf("Error parsing to_date: %v", err)
			util.NewError(c, http.StatusBadRequest, err)
			c.JSON(http.StatusBadRequest, "invalid to_date format")
			return
		}
	}

	employeeDao := employeeimp.GetInstance(c.Request.Context())
	changes, err := employeeDao.SearchChangesEmployee(fromDate, toDate, typeReport)

	if err != nil {
		log.Printf("Error searching changes employee: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search changes employee"})
		return
	}

	// 1) Department summaries by WorkLocation

	// 2) Build work_location keyed response with department summaries
	type DeptSummary struct {
		FirstPeriod    int64  `json:"first_period"`
		CountLeave     int64  `json:"count_leave"`
		CountJoin      int64  `json:"count_join"`
		EndPeriod      int64  `json:"end_period"`
		DepartmentName string `json:"department_name"`
	}

	// workLocation -> departmentName -> summary accumulator
	grouped := make(map[string]map[string]*DeptSummary)
	for _, ch := range changes {
		wl := ch.WorkLocation
		dept := ch.DepartmentName
		if _, ok := grouped[wl]; !ok {
			grouped[wl] = make(map[string]*DeptSummary)
		}
		if _, ok := grouped[wl][dept]; !ok {
			grouped[wl][dept] = &DeptSummary{DepartmentName: dept}
		}
		rec := grouped[wl][dept]
		// Sum join/leave across entries in the period
		rec.CountJoin += ch.CountJoin
		rec.CountLeave += ch.CountLeave
		// Capture first period from 'Đầu kỳ' snapshot if present
		if ch.TypeReport == "Đầu kỳ" {
			rec.FirstPeriod = ch.CountEmployee
		}
	}

	// After aggregation, compute end_period = first_period - count_leave + count_join
	// and compute totals per work_location: total_first_period, total_leave, total_join, total_end_period
	type Totals struct {
		WorkLocation     string `json:"work_location"`
		TotalFirstPeriod int64  `json:"total_first_period"`
		TotalLeave       int64  `json:"total_leave"`
		TotalJoin        int64  `json:"total_join"`
		TotalEndPeriod   int64  `json:"total_end_period"`
	}

	totals := make([]Totals, 0, len(grouped))
	resp := gin.H{}
	for wl, deptMap := range grouped {
		items := make([]DeptSummary, 0, len(deptMap))
		var sumFirst, sumLeave, sumJoin, sumEnd int64
		for _, v := range deptMap {
			v.EndPeriod = v.FirstPeriod - v.CountLeave + v.CountJoin
			items = append(items, *v)
			sumFirst += v.FirstPeriod
			sumLeave += v.CountLeave
			sumJoin += v.CountJoin
			sumEnd += v.EndPeriod
		}
		resp[wl] = items
		totals = append(totals, Totals{
			WorkLocation:     wl,
			TotalFirstPeriod: sumFirst,
			TotalLeave:       sumLeave,
			TotalJoin:        sumJoin,
			TotalEndPeriod:   sumEnd,
		})
	}
	resp["totals_by_work_location"] = totals

	c.JSON(http.StatusOK, resp)
}

func SearchHRMWorkReport(c *gin.Context) {
	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")
	typeReport := c.Query("type_report")

	var fromDate, toDate time.Time
	var err error
	// Parse dates supporting RFC3339 and YYYY-MM-DD
	if fromDateStr != "" {
		fromDate, err = time.Parse("2006-01-02", fromDateStr)
		if err != nil {
			log.Printf("Error parsing from_date: %v", err)
			util.NewError(c, http.StatusBadRequest, err)
			c.JSON(http.StatusBadRequest, "invalid from_date format")
			return
		}
	}
	if toDateStr != "" {
		toDate, err = time.Parse("2006-01-02", toDateStr)
		if err != nil {
			log.Printf("Error parsing to_date: %v", err)
			util.NewError(c, http.StatusBadRequest, err)
			c.JSON(http.StatusBadRequest, "invalid to_date format")
			return
		}
	}

	employeeDao := employeeimp.GetInstance(c.Request.Context())
	employees, err := employeeDao.SearchHRMWorkReport(fromDate, toDate, typeReport)

	if err != nil {
		log.Printf("Error searching HRM work report: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search HRM work report"})
		return
	}

	// Summary calculations per user request:
	// 1. Nhân sự đầu kỳ = 1.1 + 1.2
	// 1.1 Số nhân sự thuộc định biên Nhà máy = Chính thức + Thuê khoán
	//     Conditions: Status == "Đang làm việc", WorkLocation == "Nhà máy"
	//     Breakdown by EmploymentType: "Chính thức" and "Thuê khoán"
	// 1.2 Số nhân sự thuộc định biên Công ty (Văn phòng)
	//     Conditions: Status == "Đang làm việc", WorkLocation == "Văn phòng"

	var factoryOfficial, factoryContract, factoryTotal int64
	var officeTotal int64

	for _, e := range employees {
		// Dereference pointers safely
		status := ""
		location := ""
		employmentType := ""
		if e.Status != nil {
			status = *e.Status
		}
		if e.WorkLocation != nil {
			location = *e.WorkLocation
		}
		if e.EmploymentType != nil {
			employmentType = *e.EmploymentType
		}

		if status == "Đang làm việc" {
			// Factory staff
			if location == "Nhà máy" {
				if employmentType == "Chính thức" {
					factoryOfficial++
				} else if employmentType == "Thuê khoán" {
					factoryContract++
				}
			}
			// Office staff
			if location == "Văn phòng" {
				officeTotal++
			}
		}
	}

	factoryTotal = factoryOfficial + factoryContract
	nhanSuDauKy := factoryTotal + officeTotal

	// Section 2: Tuyển mới
	// Assumption: When type_report == "Tuyển mới" the returned employees are new hires in the period.
	// Counts split by WorkLocation (Nhà máy vs Văn phòng) with Status == "Đang làm việc".
	var newFactory, newOffice int64
	if typeReport == "Tuyển mới" {
		for _, e := range employees {
			status := ""
			location := ""
			if e.Status != nil {
				status = *e.Status
			}
			if e.WorkLocation != nil {
				location = *e.WorkLocation
			}
			if status == "Đang làm việc" {
				if location == "Nhà máy" {
					newFactory++
				} else if location == "Văn phòng" {
					newOffice++
				}
			}
		}
	}
	newTotal := newFactory + newOffice

	summary := gin.H{
		"1_nhan_su_dau_ky":                 nhanSuDauKy,
		"1_1_dinh_bien_nha_may":            factoryTotal,
		"1_1_dinh_bien_nha_may_chinh_thuc": factoryOfficial,
		"1_1_dinh_bien_nha_may_thue_khoan": factoryContract,
		"1_2_dinh_bien_cong_ty":            officeTotal,
	}
	if typeReport == "Tuyển mới" {
		summary["2_tuyen_moi"] = newTotal
		summary["2_1_tuyen_moi_nha_may"] = newFactory
		summary["2_2_tuyen_moi_cong_ty"] = newOffice
	}

	// Section 3: Nghỉ việc
	var resignFactoryOfficial, resignFactoryContract, resignFactoryTotal int64
	var resignOfficeTotal int64
	var forecastFactory, forecastOffice int64

	for _, e := range employees {
		status := ""
		location := ""
		employmentType := ""
		var resignDate *time.Time
		if e.Status != nil {
			status = *e.Status
		}
		if e.WorkLocation != nil {
			location = *e.WorkLocation
		}
		if e.EmploymentType != nil {
			employmentType = *e.EmploymentType
		}
		if e.ResignDate != nil {
			resignDate = e.ResignDate
		}

		// Actual resigns within reporting period
		if resignDate != nil && (resignDate.Equal(fromDate) || resignDate.After(fromDate)) && (resignDate.Equal(toDate) || resignDate.Before(toDate)) {
			if location == "Nhà máy" {
				if employmentType == "Chính thức" {
					resignFactoryOfficial++
				} else if employmentType == "Thuê khoán" {
					resignFactoryContract++
				}
			} else if location == "Văn phòng" {
				resignOfficeTotal++
			}
		}

		// Forecast resigns for next period: still working and resign date after current period
		if status == "Đang làm việc" && resignDate != nil && resignDate.After(toDate) {
			if location == "Nhà máy" {
				forecastFactory++
			} else if location == "Văn phòng" {
				forecastOffice++
			}
		}
	}

	resignFactoryTotal = resignFactoryOfficial + resignFactoryContract

	summary["3_nghi_viec"] = resignFactoryTotal + resignOfficeTotal + forecastFactory + forecastOffice
	summary["3_1_nghi_viec_nha_may"] = resignFactoryTotal
	summary["3_1_nha_may_chinh_thuc"] = resignFactoryOfficial
	summary["3_1_nha_may_thue_khoan"] = resignFactoryContract
	summary["3_2_nghi_viec_cong_ty"] = resignOfficeTotal
	summary["3_3_nha_may_du_kien_nghi_ky_toi"] = forecastFactory
	summary["3_4_cong_ty_du_kien_nghi_ky_toi"] = forecastOffice
	// End of Section 3
	// Section 4: Điều chuyển đến
	// Assumption: When type_report == "Điều chuyển đến" the returned employees are transfers-in in the period.
	// Counts split by WorkLocation (Nhà máy vs Văn phòng).
	var transferFactory, transferOffice int64
	if typeReport == "Điều chuyển đến" {
		for _, e := range employees {
			location := ""
			if e.WorkLocation != nil {
				location = *e.WorkLocation
			}
			if location == "Nhà máy" {
				transferFactory++
			} else if location == "Văn phòng" {
				transferOffice++
			}
		}
	}
	transferTotal := transferFactory + transferOffice
	if typeReport == "Điều chuyển đến" {
		summary["4_dieu_chuyen_den"] = transferTotal
		summary["4_1_dieu_chuyen_den_nha_may"] = transferFactory
		summary["4_2_dieu_chuyen_den_cong_ty"] = transferOffice
	}

	// Section 5: Điều chuyển đi (transfer-out)
	// Conditions: Resign within period AND ResignReason == "Điều chuyển đi"
	// Count by WorkLocation: Nhà máy vs Văn phòng
	var transferOutFactory, transferOutOffice int64
	for _, e := range employees {
		location := ""
		var resignDate *time.Time
		resignReason := ""
		if e.WorkLocation != nil {
			location = *e.WorkLocation
		}
		if e.ResignDate != nil {
			resignDate = e.ResignDate
		}
		if e.ResignReason != nil {
			resignReason = *e.ResignReason
		}
		if resignDate != nil && resignReason == "Điều chuyển đi" &&
			(resignDate.Equal(fromDate) || resignDate.After(fromDate)) && (resignDate.Equal(toDate) || resignDate.Before(toDate)) {
			if location == "Nhà máy" {
				transferOutFactory++
			} else if location == "Văn phòng" {
				transferOutOffice++
			}
		}
	}
	summary["5_dieu_chuyen_di"] = transferOutFactory + transferOutOffice
	summary["5_1_nha_may_dieu_chuyen_di"] = transferOutFactory
	summary["5_2_cong_ty_dieu_chuyen_di"] = transferOutOffice

	// Section 6: Điều chuyển nội bộ
	// Assumption: When type_report == "Điều chuyển nội bộ" the returned employees are internal transfers in the period.
	// This does not affect company totals; we provide total transfers and per-department counts.
	if typeReport == "Điều chuyển nội bộ" {
		var internalTotal int64
		deptCounts := map[string]int64{}
		for _, e := range employees {
			dept := ""
			if e.DepartmentName != nil {
				dept = *e.DepartmentName
			}
			internalTotal++
			deptCounts[dept] = deptCounts[dept] + 1
		}
		// Flatten department counts to an array for consistent JSON
		deptItems := make([]gin.H, 0, len(deptCounts))
		for d, cnt := range deptCounts {
			deptItems = append(deptItems, gin.H{"department_name": d, "count": cnt})
		}
		summary["6_dieu_chuyen_noi_bo"] = internalTotal
		summary["6_dieu_chuyen_noi_bo_theo_phong_ban"] = deptItems
	}

	// Section 7: Nhân sự cuối kỳ
	// Sections 7, 8, 9 require combining multiple components regardless of type_report
	// We'll fetch datasets for required categories and derive the metrics.
	countLocType := func(list []model.Employee) (factoryOfficial, factoryContract, factoryTotal, officeOfficial, officeContract, officeTotal int64) {
		for _, e := range list {
			loc := ""
			typ := ""
			if e.WorkLocation != nil {
				loc = *e.WorkLocation
			}
			if e.EmploymentType != nil {
				typ = *e.EmploymentType
			}
			if loc == "Nhà máy" {
				factoryTotal++
				if typ == "Chính thức" {
					factoryOfficial++
				} else if typ == "Thuê khoán" {
					factoryContract++
				}
			} else if loc == "Văn phòng" {
				officeTotal++
				if typ == "Chính thức" {
					officeOfficial++
				} else if typ == "Thuê khoán" {
					officeContract++
				}
			}
		}
		return
	}

	fetch := func(tp string) []model.Employee {
		lst, err := employeeDao.SearchHRMWorkReport(fromDate, toDate, tp)
		if err != nil {
			log.Printf("SearchHRMWorkReport(%s) error: %v", tp, err)
			return nil
		}
		return lst
	}

	beginList := fetch("Đầu kỳ")
	newList := fetch("Tuyển mới")
	resignList := fetch("Nghỉ việc")
	transferInList := fetch("Điều chuyển đến")
	transferOutList := fetch("Điều chuyển đi")

	bfO, bfC, bFT, _, _, bOT := countLocType(beginList)
	nfO, nfC, nFT, _, _, nOT := countLocType(newList)
	rfO, rfC, rFT, _, _, rOT := countLocType(resignList)
	tifO, tifC, tifT, _, _, tioT := countLocType(transferInList)
	tofO, tofC, tofT, _, _, tooT := countLocType(transferOutList)

	// 7. End-of-period headcount per location
	endFactory := bFT + nFT - rFT + tifT - tofT
	endOffice := bOT + nOT - rOT + tioT - tooT
	summary["7_nhan_su_cuoi_ky"] = endFactory + endOffice
	summary["7_1_dinh_bien_nha_may"] = endFactory
	summary["7_2_dinh_bien_cong_ty"] = endOffice

	// 8. Turnover rate (%)
	calcRate := func(numerator, begin, end int64) float64 {
		den := float64(begin+end) / 2.0
		if den <= 0 {
			return 0
		}
		pct := (float64(numerator) / den) * 100.0
		return math.Round(pct*100) / 100 // round to 2 decimals
	}

	summary["8_1_ty_le_nghi_viec_nha_may"] = calcRate(rFT, bFT, endFactory)
	summary["8_1_chinh_thuc"] = calcRate(rfO, bfO, bfO+nfO-rfO+tifO-tofO)
	summary["8_1_thue_khoan"] = calcRate(rfC, bfC, bfC+nfC-rfC+tifC-tofC)
	summary["8_2_ty_le_nghi_viec_cong_ty"] = calcRate(rOT, bOT, endOffice)

	// 9. Hiring needs (based on resigns and forecast already computed above)
	needFactory := rFT + forecastFactory
	needOffice := rOT + forecastOffice
	summary["9_nhan_su_can_tuyen_them"] = needFactory + needOffice
	summary["9_1_can_tuyen_nha_may"] = needFactory
	summary["9_2_can_tuyen_cong_ty"] = needOffice

	response := gin.H{"employees": employees, "summary": summary}

	c.JSON(http.StatusOK, response)
}

func SearhResignReport(c *gin.Context) {
	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")
	typeReport := c.Query("type_report")

	if fromDateStr == "" || toDateStr == "" || typeReport == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from_date, to_date and type_report are required"})
		return
	}

	var fromDate, toDate time.Time
	var err error
	// Parse dates supporting RFC3339 and YYYY-MM-DD
	if fromDateStr != "" {
		fromDate, err = time.Parse("2006-01-02", fromDateStr)
		fmt.Println("Parsed fromDate:", fromDate)
		if err != nil {
			log.Printf("Error parsing from_date: %v", err)
			util.NewError(c, http.StatusBadRequest, err)
			c.JSON(http.StatusBadRequest, "invalid from_date format")
			return
		}
	}
	if toDateStr != "" {
		toDate, err = time.Parse("2006-01-02", toDateStr)
		fmt.Println("Parsed toDate:", toDate)
		if err != nil {
			log.Printf("Error parsing to_date: %v", err)
			util.NewError(c, http.StatusBadRequest, err)
			c.JSON(http.StatusBadRequest, "invalid to_date format")
			return
		}
	}

	employeeDao := employeeimp.GetInstance(c.Request.Context())
	employees, err := employeeDao.SearchHRMResignReport(fromDate, toDate, typeReport)
	fmt.Println("SearhResignReport: found employees" + fmt.Sprint(len(employees)))
	if err != nil {
		log.Printf("Error searching resign report: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search resign report"})
		return
	}

	// Helper to safely deref
	getStr := func(p *string) string {
		if p == nil {
			return ""
		}
		return *p
	}

	// 1) Tổng số nghỉ việc: Status == "nghỉ việc" within period
	var totalResign int64
	var factoryResign, officeResign int64
	var monthlyCounts = map[string]int64{} // key: YYYY-MM
	var reasonCounts = map[string]int64{}
	var contractCounts = map[string]int64{}
	var deptCountsOverall = map[string]int64{}
	var deptCountsFactory = map[string]int64{}
	var deptCountsOffice = map[string]int64{}

	for _, e := range employees {
		status := getStr(e.Status)
		loc := getStr(e.WorkLocation)
		dept := getStr(e.DepartmentName)
		reason := getStr(e.ResignReason)
		contract := getStr(e.EmploymentType)
		var rdate *time.Time
		if e.ResignDate != nil {
			rdate = e.ResignDate
		}

		// Count only resigns within period
		if status == "Nghỉ việc" && rdate != nil && (rdate.Equal(fromDate) || rdate.After(fromDate)) && (rdate.Equal(toDate) || rdate.Before(toDate)) {
			totalResign++
			if loc == "Nhà máy" {
				factoryResign++
				deptCountsFactory[dept] = deptCountsFactory[dept] + 1
			} else if loc == "Văn phòng" {
				officeResign++
				deptCountsOffice[dept] = deptCountsOffice[dept] + 1
			}
			deptCountsOverall[dept] = deptCountsOverall[dept] + 1

			// Monthly bucket YYYY-MM
			monthKey := rdate.Format("2006-01")
			monthlyCounts[monthKey] = monthlyCounts[monthKey] + 1

			// Reason
			reasonKey := reason
			if reasonKey == "" {
				reasonKey = "Khác/Không rõ"
			}
			reasonCounts[reasonKey] = reasonCounts[reasonKey] + 1

			// Contract type
			contractKey := contract
			if contractKey == "" {
				contractKey = "Khác/Không rõ"
			}
			contractCounts[contractKey] = contractCounts[contractKey] + 1
		}
	}

	// 2) Previous period comparisons
	periodMonths := int64(1)
	// estimate number of months between fromDate and toDate inclusive
	{
		y := toDate.Year()*12 + int(toDate.Month())
		x := fromDate.Year()*12 + int(fromDate.Month())
		delta := y - x + 1
		if delta > 0 {
			periodMonths = int64(delta)
		}
	}
	prevFrom := fromDate.AddDate(0, -int(periodMonths), 0)
	prevTo := toDate.AddDate(0, -int(periodMonths), 0)
	prevEmployees, err := employeeDao.SearchHRMResignReport(prevFrom, prevTo, typeReport)
	if err != nil {
		log.Printf("Error previous resign report: %v", err)
	}
	var prevTotalResign int64
	var prevMonthlyAvg float64
	if prevEmployees != nil {
		for _, e := range prevEmployees {
			status := getStr(e.Status)
			var rdate *time.Time
			if e.ResignDate != nil {
				rdate = e.ResignDate
			}
			if status == "nghỉ việc" && rdate != nil && (rdate.Equal(prevFrom) || rdate.After(prevFrom)) && (rdate.Equal(prevTo) || rdate.Before(prevTo)) {
				prevTotalResign++
			}
		}
		prevMonthlyAvg = float64(prevTotalResign) / float64(periodMonths)
	}

	// % change total
	pctChangeTotal := 0.0
	if prevTotalResign > 0 {
		pctChangeTotal = (float64(totalResign-prevTotalResign) / float64(prevTotalResign)) * 100.0
	}

	// Factory/Office percentages
	factoryPct := 0.0
	officePct := 0.0
	if totalResign > 0 {
		factoryPct = (float64(factoryResign) / float64(totalResign)) * 100.0
		officePct = (float64(officeResign) / float64(totalResign)) * 100.0
	}

	// Monthly average and % change vs previous
	monthlyAvg := 0.0
	if periodMonths > 0 {
		monthlyAvg = float64(totalResign) / float64(periodMonths)
	}
	pctChangeMonthly := 0.0
	if prevMonthlyAvg > 0 {
		pctChangeMonthly = ((monthlyAvg - prevMonthlyAvg) / prevMonthlyAvg) * 100.0
	}

	// Build chart data
	// Fill months sequence between fromDate and toDate
	monthsSeq := []string{}
	cur := time.Date(fromDate.Year(), fromDate.Month(), 1, 0, 0, 0, 0, fromDate.Location())
	end := time.Date(toDate.Year(), toDate.Month(), 1, 0, 0, 0, 0, toDate.Location())
	for !cur.After(end) {
		monthsSeq = append(monthsSeq, cur.Format("2006-01"))
		cur = cur.AddDate(0, 1, 0)
	}
	monthlySeries := make([]gin.H, 0, len(monthsSeq))
	for _, m := range monthsSeq {
		monthlySeries = append(monthlySeries, gin.H{"month": m, "count": monthlyCounts[m]})
	}

	// Convert maps to arrays for JSON
	toItems := func(m map[string]int64, total int64) []gin.H {
		items := make([]gin.H, 0, len(m))
		for k, v := range m {
			pct := 0.0
			if total > 0 {
				pct = (float64(v) / float64(total)) * 100.0
			}
			items = append(items, gin.H{"label": k, "count": v, "percent": math.Round(pct*100) / 100})
		}
		return items
	}

	// Top 5 departments overall
	topN := func(m map[string]int64, n int) []gin.H {
		type kv struct {
			K string
			V int64
		}
		arr := make([]kv, 0, len(m))
		for k, v := range m {
			arr = append(arr, kv{k, v})
		}
		// simple selection sort up to n
		for i := 0; i < len(arr); i++ {
			maxIdx := i
			for j := i + 1; j < len(arr); j++ {
				if arr[j].V > arr[maxIdx].V {
					maxIdx = j
				}
			}
			arr[i], arr[maxIdx] = arr[maxIdx], arr[i]
		}
		limit := n
		if len(arr) < n {
			limit = len(arr)
		}
		res := make([]gin.H, 0, limit)
		for i := 0; i < limit; i++ {
			res = append(res, gin.H{"department_name": arr[i].K, "count": arr[i].V})
		}
		return res
	}

	response := gin.H{
		"summary": gin.H{
			"tong_so_nghi_viec":                     totalResign,
			"ty_le_thay_doi_vs_ky_truoc":            math.Round(pctChangeTotal*100) / 100,
			"nha_may_so_nghi":                       factoryResign,
			"nha_may_percent":                       math.Round(factoryPct*100) / 100,
			"van_phong_so_nghi":                     officeResign,
			"van_phong_percent":                     math.Round(officePct*100) / 100,
			"trung_binh_thang":                      math.Round(monthlyAvg*100) / 100,
			"ty_le_thay_doi_trung_binh_vs_ky_truoc": math.Round(pctChangeMonthly*100) / 100,
		},
		"charts": gin.H{
			"nghi_viec_theo_thoi_gian":           monthlySeries,
			"ty_le_theo_ly_do":                   toItems(reasonCounts, totalResign),
			"ty_le_theo_loai_hdld":               toItems(contractCounts, totalResign),
			"top5_phong_ban_nghi_nhieu":          topN(deptCountsOverall, 5),
			"nghi_viec_theo_phong_ban_nha_may":   toItems(deptCountsFactory, factoryResign),
			"nghi_viec_theo_phong_ban_van_phong": toItems(deptCountsOffice, officeResign),
		},
	}

	c.JSON(http.StatusOK, response)
}
