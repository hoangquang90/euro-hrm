package employee

import (
	"europm/internal/hrm/model"
	"time"
)

type Employee interface {
	//employeeImp.go
	GetTotalEmployee(attendanceCode string, fullName string) (int, error)
	GetEmployee(attendanceCode string, fullName string) ([]model.EmployeeResult, error)
	GetDetailEmployeeByID(id string) (model.Employee, error)
	GetCertificatesByEmployeeID(id string) ([]model.Certificate, error)
	GetRelativesByEmployeeID(id string) ([]model.Relative, error)
	GetEmergencyContactsByEmployeeID(id string) ([]model.Relative, error)
	GetContractHistoriesByEmployeeID(id string) ([]model.ContractHistory, error)
	GetSalariesByEmployeeID(id string) ([]model.Salary, error)
	GetCareerHistoriesByEmployeeID(id string) ([]model.CareerHistory, error)
	GetPerformanceEvaluationsByEmployeeID(id string) ([]model.PerformanceEvaluation, error)
	GetRewardDisciplinesByEmployeeID(id string) ([]model.RewardDiscipline, error)
	InstEmployee(emp model.Employee) (string, error)
	InstRelatives(relatives []model.Relative) (string, error)
	InstEmergencyContacts(contacts []model.Relative) (string, error)
	InstCertificates(certificates []model.Certificate) (string, error)
	InstSalaries(salaries model.Salary) (string, error)
	InstCareerHistories(careerHistories []model.CareerHistory) (string, error)
	InstPerformanceEvaluations(performanceEvaluations []model.PerformanceEvaluation) (string, error)
	InstRewardDiscipline(rds []model.RewardDiscipline) (string, error)
	InstContractHistories(chs model.ContractHistory) (string, error)
	DeleteCertificatesByID(id string) (string, error)
	DeleteRelativesByID(id string) (string, error)
	DeleteSalariesByID(id string) (string, error)
	DeleteCareerHistoriesByID(id string) (string, error)
	DeletePerformanceEvaluationsByID(id string) (string, error)
	DeleteRewardDisciplinesByID(id string) (string, error)
	DeleteContractHistoriesByID(id string) (string, error)
	GetTotalEmployeesResign(text string, fromDate time.Time, toDate time.Time) (int, error)
	GetEmployeesResign(text string, fromDate time.Time, toDate time.Time) ([]model.Employee, error)
	GetEmployeeResignByID(id string) (model.Employee, error)
	UpdateEmployeeResign(employ model.Employee) error
	GetHRTraining(year int) ([]model.HRTraining, error)
	GetHRTrainingByID(id string) (model.HRTraining, error)
	InstHRTraining(hrt model.HRTraining) (string, error)
	DeleteHRTraining(id string) error
	GetMedicalHistoryByID(id string) ([]model.MedicalHistory, error)
	InstMedicalHistory(mh model.MedicalHistory) (string, error)

	//EmployeeReportImp.go
	SearchChangesEmployee(fromDate time.Time, toDate time.Time, typeReport string) ([]model.ChangesEmployee, error)
	SearchHRMWorkReport(fromDate time.Time, toDate time.Time, tyeReport string) ([]model.Employee, error)
	SearchHRMResignReport(fromDate time.Time, toDate time.Time, tyeReport string) ([]model.Employee, error)
	GetRecruitmentPlan(year int) ([]model.RecruitmentPlan, error)
	GetRecruitmentPlanByID(id string) (model.RecruitmentPlan, error)
	InstRecruitmentPlan(rp model.RecruitmentPlan) (string, error)
	DeleteRecruitmentPlan(id string) error
}
