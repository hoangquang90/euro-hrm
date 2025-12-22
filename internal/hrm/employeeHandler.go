package hrm

import (
	"encoding/json"
	"europm/internal/config"
	employeeimp "europm/internal/hrm/employeeImp"
	"europm/internal/hrm/model"
	"europm/internal/util"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func SearchEmployee(c *gin.Context) {
	lstEmployee := make([]model.EmployeeResult, 0)
	attendanceCode := c.Query("attendance_code")
	fullName := c.Query("full_name")

	employeeDao := employeeimp.GetInstance(c.Request.Context())
	total, err := employeeDao.GetTotalEmployee(attendanceCode, fullName)
	if err != nil {
		log.Printf("Error fetching total employee: %v", err)
		util.NewError(c, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, "get total employee error")
		return
	}
	if total > 0 {
		lstEmployee, err = employeeDao.GetEmployee(attendanceCode, fullName)
		if err != nil {
			log.Printf("Error fetching employee: %v", err)
			util.NewError(c, http.StatusInternalServerError, err)
			c.JSON(http.StatusInternalServerError, "get employee error")
			return
		}
	}
	c.JSON(http.StatusOK, lstEmployee)
}

func SearchEmployeeByID(c *gin.Context) {
	id := c.Param("id")
	employeeDao := employeeimp.GetInstance(c.Request.Context())
	lstEmployee, err := employeeDao.GetDetailEmployeeByID(id)
	if err != nil {
		log.Printf("Error fetching employee by ID: %v", err)
		util.NewError(c, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, "get employee by ID error")
		return
	}
	log.Printf("employee list: %v", lstEmployee)
	c.JSON(http.StatusOK, lstEmployee)
}

func SaveEmployee(c *gin.Context) {
	// Parse multipart form (optional)
	c.Request.ParseMultipartForm(32 << 20)

	// Try to get file, but do not fail if not present
	file, err := c.FormFile("file")
	var savedFilePath, fileName string
	if err == nil && file != nil {
		savePath := config.GetString("file_pathhrm") + file.Filename
		log.Printf("Saving file to: %s", savePath)
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			log.Printf("Error saving uploaded file: %v", err)
			c.JSON(http.StatusInternalServerError, "save file error")
			return
		}
		savedFilePath = savePath
		fileName = file.Filename
	}

	// Lấy employee JSON string từ form-data
	employeeJSON := c.PostForm("employee")
	if employeeJSON == "" {
		log.Printf("employee form field is empty")
		c.JSON(http.StatusBadRequest, "employee field empty")
		return
	}

	var emp model.Employee
	if err := json.Unmarshal([]byte(employeeJSON), &emp); err != nil {
		log.Printf("Error parsing employee JSON: %v", err)
		c.JSON(http.StatusBadRequest, "invalid employee json")
		return
	}

	// Only set ImagePath/ImageName if file was uploaded
	if savedFilePath != "" {
		emp.ImagePath = &savedFilePath
	} else {
		emp.ImagePath = nil
	}
	if fileName != "" {
		emp.ImageName = &fileName
	} else {
		emp.ImageName = nil
	}

	// Lưu vào DB
	employeeDao := employeeimp.GetInstance(c.Request.Context())
	id, err := employeeDao.InstEmployee(emp)
	if err != nil {
		log.Printf("Error inserting employee: %v", err)
		c.JSON(http.StatusInternalServerError, "insert employee error")
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func UpdateRelatives(c *gin.Context) {
	var emp []model.Relative
	if err := c.ShouldBindJSON(&emp); err != nil {
		log.Printf("Error binding JSON: %v", err)
		util.NewError(c, http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, "bind relatives error")
		return
	}
	employeeDao := employeeimp.GetInstance(c.Request.Context())
	if emp != nil {
		_, err := employeeDao.InsertRelatives(emp)
		if err != nil {
			log.Printf("Error inserting relatives: %v", err)
			util.NewError(c, http.StatusInternalServerError, err)
			c.JSON(http.StatusInternalServerError, "insert relatives error")
			return
		}
	}
	c.JSON(200, "Success")
}

func UpdateEmergencyContacts(c *gin.Context) {
	var emp []model.Relative
	if err := c.ShouldBindJSON(&emp); err != nil {
		log.Printf("Error binding JSON: %v", err)
		util.NewError(c, http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, "bind relatives error")
		return
	}
	employeeDao := employeeimp.GetInstance(c.Request.Context())
	if emp != nil {
		_, err := employeeDao.InsertEmergencyContacts(emp)
		if err != nil {
			log.Printf("Error inserting emergency contacts: %v", err)
			util.NewError(c, http.StatusInternalServerError, err)
			c.JSON(http.StatusInternalServerError, "insert emergency contacts error")
			return
		}
	}
	c.JSON(200, "Success")
}

func UpdateSalaries(c *gin.Context) {
	c.Request.ParseMultipartForm(32 << 20)
	file, err := c.FormFile("file")
	var savedFilePath, fileName string
	if err == nil && file != nil {
		savePath := config.GetString("file_pathhrm") + file.Filename
		log.Printf("Saving file to: %s", savePath)
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			log.Printf("Error saving uploaded file: %v", err)
			c.JSON(http.StatusInternalServerError, "save file error")
			return
		}
		savedFilePath = savePath
		fileName = file.Filename
	}

	// Lấy employee JSON string từ form-data
	salary := c.PostForm("salary")
	if salary == "" {
		log.Printf("salary form field is empty")
		c.JSON(http.StatusBadRequest, "salary field empty")
		return
	}

	var emp model.Salary
	if err := json.Unmarshal([]byte(salary), &emp); err != nil {
		log.Printf("Error parsing salary JSON: %v", err)
		c.JSON(http.StatusBadRequest, "invalid salary json")
		return
	}

	if savedFilePath != "" {
		emp.FilePath = &savedFilePath
	} else {
		emp.FilePath = nil
	}
	if fileName != "" {
		emp.FileName = &fileName
	} else {
		emp.FileName = nil
	}
	employeeDao := employeeimp.GetInstance(c.Request.Context())
	id, err := employeeDao.InsertSalaries(emp)
	if err != nil {
		log.Printf("Error inserting salaries: %v", err)
		util.NewError(c, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, "insert salaries error")
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func UpdateCertificates(c *gin.Context) {
	var emp []model.Certificate
	if err := c.ShouldBindJSON(&emp); err != nil {
		log.Printf("Error binding JSON: %v", err)
		util.NewError(c, http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, "bind certificates error")
		return
	}
	employeeDao := employeeimp.GetInstance(c.Request.Context())
	if emp != nil {
		_, err := employeeDao.InsertCertificates(emp)
		if err != nil {
			log.Printf("Error inserting certificates: %v", err)
			util.NewError(c, http.StatusInternalServerError, err)
			c.JSON(http.StatusInternalServerError, "insert certificates error")
			return
		}
	}
	c.JSON(200, "Success")
}
func UpdateCareerHistories(c *gin.Context) {
	var emp []model.CareerHistory
	if err := c.ShouldBindJSON(&emp); err != nil {
		log.Printf("Error binding JSON: %v", err)
		util.NewError(c, http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, "bind career histories error")
		return
	}
	employeeDao := employeeimp.GetInstance(c.Request.Context())
	if emp != nil {
		_, err := employeeDao.InsertCareerHistories(emp)
		if err != nil {
			log.Printf("Error inserting career histories: %v", err)
			util.NewError(c, http.StatusInternalServerError, err)
			c.JSON(http.StatusInternalServerError, "insert career histories error")
			return
		}
	}
	c.JSON(200, "Success")
}

func UpdatePerformanceEvaluations(c *gin.Context) {
	var emp []model.PerformanceEvaluation
	if err := c.ShouldBindJSON(&emp); err != nil {
		log.Printf("Error binding JSON: %v", err)
		util.NewError(c, http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, "bind performance evaluations error")
		return
	}
	employeeDao := employeeimp.GetInstance(c.Request.Context())
	if emp != nil {
		_, err := employeeDao.InsertPerformanceEvaluations(emp)
		if err != nil {
			log.Printf("Error inserting performance evaluations: %v", err)
			util.NewError(c, http.StatusInternalServerError, err)
			c.JSON(http.StatusInternalServerError, "insert performance evaluations error")
			return
		}
	}
	c.JSON(200, "Success")
}

func UpdateRewardDisciplines(c *gin.Context) {
	var emp []model.RewardDiscipline
	if err := c.ShouldBindJSON(&emp); err != nil {
		log.Printf("Error binding JSON: %v", err)
		util.NewError(c, http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, "bind reward discipline error")
		return
	}
	employeeDao := employeeimp.GetInstance(c.Request.Context())
	if emp != nil {
		_, err := employeeDao.InsertRewardDiscipline(emp)
		if err != nil {
			log.Printf("Error inserting reward disciplines: %v", err)
			util.NewError(c, http.StatusInternalServerError, err)
			c.JSON(http.StatusInternalServerError, "insert reward disciplines error")
			return
		}
	}
	c.JSON(200, "Success")
}

func UpdateContractHistories(c *gin.Context) {
	c.Request.ParseMultipartForm(32 << 20)
	file, err := c.FormFile("file")
	var savedFilePath, fileName string
	if err == nil && file != nil {
		savePath := config.GetString("file_pathhrm") + file.Filename
		log.Printf("Saving file to: %s", savePath)
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			log.Printf("Error saving uploaded file: %v", err)
			c.JSON(http.StatusInternalServerError, "save file error")
			return
		}
		savedFilePath = savePath
		fileName = file.Filename
	}

	// Lấy employee JSON string từ form-data
	ContractHistory := c.PostForm("contractHistory")
	if ContractHistory == "" {
		log.Printf("contractHistory form field is empty")
		c.JSON(http.StatusBadRequest, "contractHistory field empty")
		return
	}

	var emp model.ContractHistory
	if err := json.Unmarshal([]byte(ContractHistory), &emp); err != nil {
		log.Printf("Error parsing contractHistory JSON: %v", err)
		c.JSON(http.StatusBadRequest, "invalid contractHistory json")
		return
	}

	if savedFilePath != "" {
		emp.FilePath = &savedFilePath
	} else {
		emp.FilePath = nil
	}
	if fileName != "" {
		emp.FileName = &fileName
	} else {
		emp.FileName = nil
	}

	employeeDao := employeeimp.GetInstance(c.Request.Context())
	id, err := employeeDao.InsertContractHistories(emp)
	if err != nil {
		log.Printf("Error inserting contract histories: %v", err)
		util.NewError(c, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, "insert contract histories error")
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func SearchCertificatesByID(c *gin.Context) {
	id := c.Param("id")
	employeeDao := employeeimp.GetInstance(c.Request.Context())
	lstCertificates, err := employeeDao.GetCertificatesByEmployeeID(id)
	if err != nil {
		log.Printf("Error fetching certificates by employee ID: %v", err)
		util.NewError(c, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, "get certificates by employee ID error")
		return
	}

	c.JSON(http.StatusOK, lstCertificates)
}

func SearchRelativesByID(c *gin.Context) {
	id := c.Param("id")
	employeeDao := employeeimp.GetInstance(c.Request.Context())
	lstRelatives, err := employeeDao.GetRelativesByEmployeeID(id)
	if err != nil {
		log.Printf("Error fetching relatives by employee ID: %v", err)
		util.NewError(c, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, "get relatives by employee ID error")
		return
	}

	c.JSON(http.StatusOK, lstRelatives)
}

func SearchEmergencyContactsByID(c *gin.Context) {
	id := c.Param("id")
	employeeDao := employeeimp.GetInstance(c.Request.Context())
	lstEmergencyContacts, err := employeeDao.GetEmergencyContactsByEmployeeID(id)
	if err != nil {
		log.Printf("Error fetching emergency contacts by employee ID: %v", err)
		util.NewError(c, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, "get emergency contacts by employee ID error")
		return
	}

	c.JSON(http.StatusOK, lstEmergencyContacts)
}

func SearchSalariesByID(c *gin.Context) {
	id := c.Param("id")
	employeeDao := employeeimp.GetInstance(c.Request.Context())
	lstSalaries, err := employeeDao.GetSalariesByEmployeeID(id)
	if err != nil {
		log.Printf("Error fetching salaries by employee ID: %v", err)
		util.NewError(c, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, "get salaries by employee ID error")
		return
	}

	c.JSON(http.StatusOK, lstSalaries)
}

func SearchCareerHistoriesByID(c *gin.Context) {
	id := c.Param("id")
	employeeDao := employeeimp.GetInstance(c.Request.Context())
	lstCareerHistories, err := employeeDao.GetCareerHistoriesByEmployeeID(id)
	if err != nil {
		log.Printf("Error fetching career histories by employee ID: %v", err)
		util.NewError(c, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, "get career histories by employee ID error")
		return
	}

	c.JSON(http.StatusOK, lstCareerHistories)
}

func SearchPerformanceEvaluationsByID(c *gin.Context) {
	id := c.Param("id")
	employeeDao := employeeimp.GetInstance(c.Request.Context())
	lstPerformanceEvaluations, err := employeeDao.GetPerformanceEvaluationsByEmployeeID(id)
	if err != nil {
		log.Printf("Error fetching performance evaluations by employee ID: %v", err)
		util.NewError(c, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, "get performance evaluations by employee ID error")
		return
	}

	c.JSON(http.StatusOK, lstPerformanceEvaluations)
}

func SearchRewardDisciplinesByID(c *gin.Context) {
	id := c.Param("id")
	employeeDao := employeeimp.GetInstance(c.Request.Context())
	lstRewardDisciplines, err := employeeDao.GetRewardDisciplinesByEmployeeID(id)
	if err != nil {
		log.Printf("Error fetching reward disciplines by employee ID: %v", err)
		util.NewError(c, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, "get reward disciplines by employee ID error")
		return
	}

	c.JSON(http.StatusOK, lstRewardDisciplines)
}

func SearchContractHistoriesByID(c *gin.Context) {
	id := c.Param("id")
	employeeDao := employeeimp.GetInstance(c.Request.Context())
	lstContractHistories, err := employeeDao.GetContractHistoriesByEmployeeID(id)
	if err != nil {
		log.Printf("Error fetching contract histories by employee ID: %v", err)
		util.NewError(c, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, "get contract histories by employee ID error")
		return
	}

	c.JSON(http.StatusOK, lstContractHistories)
}

func DeleteCertificatesByID(c *gin.Context) {
	id := c.Param("id")
	employeeDao := employeeimp.GetInstance(c.Request.Context())
	_, err := employeeDao.DeleteCertificatesByID(id)
	if err != nil {
		log.Printf("Error deleting certificates by ID: %v", err)
		util.NewError(c, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, "delete certificates by ID error")
		return
	}
	c.JSON(http.StatusOK, "Success")
}

func DeleteRelativesByID(c *gin.Context) {
	id := c.Param("id")
	employeeDao := employeeimp.GetInstance(c.Request.Context())
	_, err := employeeDao.DeleteRelativesByID(id)
	if err != nil {
		log.Printf("Error deleting relatives by ID: %v", err)
		util.NewError(c, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, "delete relatives by ID error")
		return
	}
	c.JSON(http.StatusOK, "Success")
}

func DeleteSalariesByID(c *gin.Context) {
	id := c.Param("id")
	employeeDao := employeeimp.GetInstance(c.Request.Context())
	_, err := employeeDao.DeleteSalariesByID(id)
	if err != nil {
		log.Printf("Error deleting salaries by ID: %v", err)
		util.NewError(c, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, "delete salaries by ID error")
		return
	}
	c.JSON(http.StatusOK, "Success")
}

func DeleteCareerHistoriesByID(c *gin.Context) {
	id := c.Param("id")
	employeeDao := employeeimp.GetInstance(c.Request.Context())
	_, err := employeeDao.DeleteCareerHistoriesByID(id)
	if err != nil {
		log.Printf("Error deleting career histories by ID: %v", err)
		util.NewError(c, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, "delete career histories by ID error")
		return
	}
	c.JSON(http.StatusOK, "Success")
}

func DeletePerformanceEvaluationsByID(c *gin.Context) {
	id := c.Param("id")
	employeeDao := employeeimp.GetInstance(c.Request.Context())
	_, err := employeeDao.DeletePerformanceEvaluationsByID(id)
	if err != nil {
		log.Printf("Error deleting performance evaluations by ID: %v", err)
		util.NewError(c, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, "delete performance evaluations by ID error")
		return
	}
	c.JSON(http.StatusOK, "Success")
}

func DeleteRewardDisciplinesByID(c *gin.Context) {
	id := c.Param("id")
	employeeDao := employeeimp.GetInstance(c.Request.Context())
	_, err := employeeDao.DeleteRewardDisciplinesByID(id)
	if err != nil {
		log.Printf("Error deleting reward disciplines by ID: %v", err)
		util.NewError(c, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, "delete reward disciplines by ID error")
		return
	}
	c.JSON(http.StatusOK, "Success")
}

func DeleteContractHistoriesByID(c *gin.Context) {
	id := c.Param("id")
	employeeDao := employeeimp.GetInstance(c.Request.Context())
	_, err := employeeDao.DeleteContractHistoriesByID(id)
	if err != nil {
		log.Printf("Error deleting contract histories by ID: %v", err)
		util.NewError(c, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, "delete contract histories by ID error")
		return
	}
	c.JSON(http.StatusOK, "Success")
}

func SearchEmployeesResign(c *gin.Context) {
	lstEmployee := make([]model.Employee, 0)
	text := c.Query("text")
	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")

	var fromDate, toDate time.Time
	var err error
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
	total, err := employeeDao.GetTotalEmployeesResign(text, fromDate, toDate)
	if err != nil {
		log.Printf("Error fetching total employees resign: %v", err)
		util.NewError(c, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, "get total employees resign error")
		return
	}
	if total > 0 {
		lstEmployee, err = employeeDao.GetEmployeesResign(text, fromDate, toDate)
		if err != nil {
			log.Printf("Error fetching employees resign: %v", err)
			util.NewError(c, http.StatusInternalServerError, err)
			c.JSON(http.StatusInternalServerError, "get employees resign error")
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"total": total, "employees": lstEmployee})
}

func SearchEmployeeResignByID(c *gin.Context) {
	id := c.Param("id")
	employeeDao := employeeimp.GetInstance(c.Request.Context())
	employeeResign, err := employeeDao.GetEmployeeResignByID(id)
	if err != nil {
		log.Printf("Error fetching employee resign by ID: %v", err)
		util.NewError(c, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, "get employee resign by ID error")
		return
	}

	c.JSON(http.StatusOK, employeeResign)
}

func UpdateEmployeeResign(c *gin.Context) {
	var employee model.Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		log.Printf("Error binding JSON: %v", err)
		util.NewError(c, http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, "bind reward discipline error")
		return
	}
	employeeDao := employeeimp.GetInstance(c.Request.Context())
	err := employeeDao.UpdateEmployeeResign(employee)
	if err != nil {
		log.Printf("Error updating employee resign: %v", err)
		util.NewError(c, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, "update employee resign error")
		return
	}
	c.JSON(200, "Success")
}

func SearchRecruitmentPlan(c *gin.Context) {
	lstRecruitmentPlan := make([]model.RecruitmentPlan, 0)
	yearStr := c.Query("year")
	var year int
	var err error
	if yearStr != "" {
		year, err = strconv.Atoi(yearStr)
		if err != nil {
			log.Printf("Error parsing year: %v", err)
			util.NewError(c, http.StatusBadRequest, err)
			c.JSON(http.StatusBadRequest, "invalid year format")
			return
		}
	}
	employeeDao := employeeimp.GetInstance(c.Request.Context())
	lstRecruitmentPlan, err = employeeDao.GetRecruitmentPlan(year)
	if err != nil {
		log.Printf("Error fetching recruitment plan: %v", err)
		util.NewError(c, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, "get recruitment plan error")
		return
	}

	c.JSON(http.StatusOK, gin.H{"recruitment_plans": lstRecruitmentPlan})
}

func SearchRecruitmentPlanByID(c *gin.Context) {
	ID := c.Query("id")
	employeeDao := employeeimp.GetInstance(c.Request.Context())
	recruitmentPlan, err := employeeDao.GetRecruitmentPlanByID(ID)
	if err != nil {
		log.Printf("Error fetching recruitment plan: %v", err)
		util.NewError(c, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, "get recruitment plan error")
		return
	}

	c.JSON(http.StatusOK, recruitmentPlan)
}

func InsertRecruitmentPlan(c *gin.Context) {
	var rp model.RecruitmentPlan
	if err := c.ShouldBindJSON(&rp); err != nil {
		log.Printf("Error binding JSON: %v", err)
		util.NewError(c, http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, "bind recruitment plan error")
		return
	}
	employeeDao := employeeimp.GetInstance(c.Request.Context())
	id, err := employeeDao.InstRecruitmentPlan(rp)
	if err != nil {
		log.Printf("Error inserting recruitment plan: %v", err)
		util.NewError(c, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, "insert recruitment plan error")
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func DeleteRecruitmentPlan(c *gin.Context) {
	id := c.Param("id")
	employeeDao := employeeimp.GetInstance(c.Request.Context())
	err := employeeDao.DeleteRecruitmentPlan(id)
	if err != nil {
		log.Printf("Error deleting recruitment plan by ID: %v", err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "Success")
}
