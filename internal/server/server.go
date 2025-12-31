package server

import (
	"context"
	"europm/internal/config"
	"europm/internal/hrm"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var srv *http.Server

func Start() (err error) {
	hostName, _ := os.Hostname()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(otelgin.Middleware(hostName))

	// Chuẩn hoá: bỏ "/" cuối nếu có → "/api/v1"
	prefix := strings.TrimSuffix(config.GetString("api_prefix"), "/")

	// Health cho chính prefix → "/api/v1/"
	r.GET(prefix+"/", func(c *gin.Context) { c.Status(http.StatusOK) })
	r.HEAD(prefix+"/", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Group các API con → "/api/v1/..."
	rg := r.Group(prefix + "/")
	{
		rg.GET("employee", hrm.SearchEmployee)
		rg.GET("employee/:id", hrm.SearchEmployeeByID)
		rg.GET("employee/certificates/:id", hrm.SearchCertificatesByID)
		rg.GET("employee/relatives/:id", hrm.SearchRelativesByID)
		rg.GET("employee/emergency_contacts/:id", hrm.SearchEmergencyContactsByID)
		rg.GET("employee/salaries/:id", hrm.SearchSalariesByID)
		rg.GET("employee/career_histories/:id", hrm.SearchCareerHistoriesByID)
		rg.GET("employee/performance_evaluations/:id", hrm.SearchPerformanceEvaluationsByID)
		rg.GET("employee/reward_disciplines/:id", hrm.SearchRewardDisciplinesByID)
		rg.GET("employee/contract_histories/:id", hrm.SearchContractHistoriesByID)
		rg.POST("employee", hrm.SaveEmployee)
		rg.PUT("employee/relatives", hrm.UpdateRelatives)
		rg.PUT("employee/emergency_contacts", hrm.UpdateEmergencyContacts)
		rg.PUT("employee/salaries", hrm.UpdateSalaries)
		rg.PUT("employee/certificates", hrm.UpdateCertificates)
		rg.PUT("employee/career_histories", hrm.UpdateCareerHistories)
		rg.PUT("employee/performance_evaluations", hrm.UpdatePerformanceEvaluations)
		rg.PUT("employee/reward_disciplines", hrm.UpdateRewardDisciplines)
		rg.PUT("employee/contract_histories", hrm.UpdateContractHistories)
		rg.DELETE("employee/certificates/:id", hrm.DeleteCertificatesByID)
		rg.DELETE("employee/relatives/:id", hrm.DeleteRelativesByID)
		rg.DELETE("employee/emergency_contacts/:id", hrm.DeleteRelativesByID)
		rg.DELETE("employee/salaries/:id", hrm.DeleteSalariesByID)
		rg.DELETE("employee/career_histories/:id", hrm.DeleteCareerHistoriesByID)
		rg.DELETE("employee/performance_evaluations/:id", hrm.DeletePerformanceEvaluationsByID)
		rg.DELETE("employee/reward_disciplines/:id", hrm.DeleteRewardDisciplinesByID)
		rg.DELETE("employee/contract_histories/:id", hrm.DeleteContractHistoriesByID)
		rg.GET("employee/employees_resign", hrm.SearchEmployeesResign)
		rg.GET("employee/employees_resign/:id", hrm.SearchEmployeeResignByID)
		rg.PUT("employee/update_employee_resign", hrm.UpdateEmployeeResign)
		rg.GET("employee/changes_employee", hrm.SearchChangesEmployee)
		rg.GET("employee/hrm_work_report", hrm.SearchHRMWorkReport)
		rg.GET("employee/employees_resign_report", hrm.SearhResignReport)
		rg.GET("employee/recruitment_plan", hrm.SearchRecruitmentPlan)
		rg.GET("employee/recruitment_plan/:id", hrm.SearchRecruitmentPlanByID)
		rg.POST("employee/recruitment_plan", hrm.InsertRecruitmentPlan)
		rg.DELETE("employee/recruitment_plan/:id", hrm.DeleteRecruitmentPlan)
		rg.GET("employee/hr_training/", hrm.SearchHRTraining)
		rg.GET("employee/hr_training/:id", hrm.SearchHRTrainingByID)
		rg.POST("employee/hr_training", hrm.InsertHRTraining)
		rg.DELETE("employee/hr_training/:id", hrm.DeleteHRTraining)
		rg.GET("employee/medical_history/:id", hrm.SearchMedicalHistoryByID)
		rg.POST("employee/medical_history", hrm.InsertMedicalHistory)
	}

	srv = &http.Server{
		Addr:    ":" + config.GetString("server.port"), // 8433
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err)
		}
	}()
	return nil
}

func Stop() {
	// The context is used to inform the server it has 90 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	fmt.Println("defer cancel()")
	defer cancel()
	fmt.Println("srv.Shutdown(ctx)")
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown:")
		log.Fatal("Server forced to shutdown: ", err)
	}

	fmt.Println("Server exiting")
}
