package employee

import "europm/internal/hrm/model"

type Employee interface {
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
	InsertRelatives(relatives []model.Relative) (string, error)
	InsertEmergencyContacts(contacts []model.Relative) (string, error)
	InsertCertificates(certificates []model.Certificate) (string, error)
	InsertSalaries(salaries model.Salary) (string, error)
	InsertCareerHistories(careerHistories []model.CareerHistory) (string, error)
	InsertPerformanceEvaluations(performanceEvaluations []model.PerformanceEvaluation) (string, error)
	InsertRewardDiscipline(rds []model.RewardDiscipline) (string, error)
	InsertContractHistories(chs model.ContractHistory) (string, error)
	DelEmployee(empID string) (string, error)
}
