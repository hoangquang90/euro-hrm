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

	// Aggregate totals by WorkLocation
	type Totals struct {
		WorkLocation string `json:"work_location"`
		TotalLeave   int64  `json:"total_leave"`
		TotalJoin    int64  `json:"total_join"`
		TotalEmp     int64  `json:"total_emp"`
	}

	totalsMap := make(map[string]*Totals)
	for _, ch := range changes {
		key := ch.WorkLocation
		if _, ok := totalsMap[key]; !ok {
			totalsMap[key] = &Totals{WorkLocation: key}
		}
		t := totalsMap[key]
		t.TotalLeave += ch.CountLeave
		t.TotalJoin += ch.CountJoin
		t.TotalEmp += ch.CountEmployee
	}

	totals := make([]Totals, 0, len(totalsMap))
	for _, v := range totalsMap {
		totals = append(totals, *v)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":                    changes,
		"totals_by_work_location": totals,
	})
}
