package hrm

import (
	employeeimp "europm/internal/hrm/employeeImp"
	"europm/internal/hrm/model"
	"europm/internal/util"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SearchChangesEmployee(c *gin.Context) {
	// Bind into DTO with string dates to support both JSON body and query params
	var dto model.ChangesEmployeeFilterDTO
	if err := c.ShouldBind(&dto); err != nil {
		log.Printf("Error binding JSON: %v", err)
		util.NewError(c, http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, "bind reward discipline error")
		return
	}

	if dto.FromDate == "" || dto.ToDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from_date and to_date are required"})
		return
	}

	// Parse dates supporting RFC3339 and YYYY-MM-DD
	parseDate := func(s string) (time.Time, error) {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			return t, nil
		}
		return time.Parse("2006-01-02", s)
	}
	fromTime, err := parseDate(dto.FromDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from_date must be RFC3339 or YYYY-MM-DD"})
		return
	}
	toTime, err := parseDate(dto.ToDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to_date must be RFC3339 or YYYY-MM-DD"})
		return
	}
	if toTime.Before(fromTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to_date must be greater than or equal to from_date"})
		return
	}

	filter := model.ChangesEmployeeFilter{
		FromDate:   fromTime,
		ToDate:     toTime,
		TypeReport: dto.TypeReport,
	}

	employeeDao := employeeimp.GetInstance(c.Request.Context())
	changes, err := employeeDao.SearchChangesEmployee(&filter)

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
