package model

import (
	"encoding/json"
	"fmt"
	"time"
)

type Employee struct {
	ID                *string    `json:"id"`
	CreatedAt         *time.Time `json:"created_at,omitempty"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty"`
	FullName          *string    `json:"full_name"`
	BirthDate         *time.Time `json:"birth_date,omitempty"`
	Gender            *string    `json:"gender,omitempty"`
	IDNumber          *string    `json:"id_number,omitempty"`
	IssueDate         *time.Time `json:"issue_date,omitempty"`
	IssuePlace        *string    `json:"issue_place,omitempty"`
	BirthPlace        *string    `json:"birth_place,omitempty"`
	HomeTown          *string    `json:"home_town,omitempty"`
	PermanentAddress  *string    `json:"permanent_address,omitempty"`
	TemporaryAddress  *string    `json:"temporary_address,omitempty"`
	MaritalStatus     *string    `json:"marital_status,omitempty"`
	PersonalPhone     *string    `json:"personal_phone,omitempty"`
	PersonalEmail     *string    `json:"personal_email,omitempty"`
	CompanyPhone      *string    `json:"company_phone,omitempty"`
	CompanyEmail      *string    `json:"company_email,omitempty"`
	JoinDate          *time.Time `json:"join_date,omitempty"`
	HighestDegree     *string    `json:"highest_degree,omitempty"`
	Major             *string    `json:"major,omitempty"`
	SchoolName        *string    `json:"school_name,omitempty"`
	GraduationYear    *int64     `json:"graduation_year,omitempty"`
	SpecialSkills     *string    `json:"special_skills,omitempty"`
	AttendanceCode    *string    `json:"attendance_code,omitempty"`
	PositionTitle     *string    `json:"position_title,omitempty"`
	DepartmentName    *string    `json:"department_name,omitempty"`
	Rank              *string    `json:"rank,omitempty"`
	WorkLocation      *string    `json:"work_location,omitempty"`
	OfficialDate      *time.Time `json:"official_date,omitempty"`
	Leader            *string    `json:"leader,omitempty"`
	ManagerID         *string    `json:"manager_id,omitempty"`
	SocialInsuranceNo *string    `json:"social_insurance_no,omitempty"`
	InsuranceStatus   *string    `json:"insurance_status,omitempty"`
	InsuranceDate     *time.Time `json:"insurance_date,omitempty"`
	InsuranceAmount   *int64     `json:"insurance_amount,omitempty"`
	KCBPlace          *string    `json:"kcb_place,omitempty"`
	HealthInsurance   *string    `json:"health_insurance,omitempty"`
	HealthInsurExpire *time.Time `json:"health_insur_expire,omitempty"`
	TaxID             *string    `json:"tax_id,omitempty"`
	BankAccount       *string    `json:"bank_account,omitempty"`
	BankName          *string    `json:"bank_name,omitempty"`
	Status            *string    `json:"status,omitempty"`
	Portrait          *string    `json:"portrait,omitempty"`
	DevelopmentPlan   *string    `json:"development_plan,omitempty"`
	JobObjective      *string    `json:"job_objective,omitempty"`
	HealthStatus      *string    `json:"health_status,omitempty"`
	ImagePath         *string    `json:"image_path,omitempty"`
	ImageName         *string    `json:"image_name,omitempty"`
	// Additional fields resign infomation
	ResignCode         *string    `json:"resign_code,omitempty"`
	ResignRequest      *time.Time `json:"resign_request,omitempty"`
	ResignDate         *time.Time `json:"resign_date,omitempty"`
	ResignStatus       *string    `json:"resign_status,omitempty"`
	ResignReason       *string    `json:"resign_reason,omitempty"`
	ResignApprovedBy   *string    `json:"resign_approved_by,omitempty"`
	ResignApprovedDate *time.Time `json:"resign_approved_date,omitempty"`
	EmploymentType     *string    `json:"employment_type,omitempty"`
}

// MarshalJSON customizes JSON output: pointer string fields nil -> ""

func (e Employee) MarshalJSON() ([]byte, error) {
	type Alias Employee
	aux := struct {
		Alias
		CreatedAt          string `json:"created_at,omitempty"`
		UpdatedAt          string `json:"updated_at,omitempty"`
		IssueDate          string `json:"issue_date,omitempty"`
		JoinDate           string `json:"join_date,omitempty"`
		OfficialDate       string `json:"official_date,omitempty"`
		ResignDate         string `json:"resign_date,omitempty"`
		InsuranceDate      string `json:"insurance_date,omitempty"`
		HealthInsurExpire  string `json:"health_insur_expire,omitempty"`
		Age                string `json:"age,omitempty"`
		GraduationYear     string `json:"graduation_year,omitempty"`
		InsuranceAmount    string `json:"insurance_amount,omitempty"`
		HealthInsurance    string `json:"health_insurance,omitempty"`
		TaxID              string `json:"tax_id,omitempty"`
		BankAccount        string `json:"bank_account,omitempty"`
		BankName           string `json:"bank_name,omitempty"`
		Status             string `json:"status,omitempty"`
		Portrait           string `json:"portrait,omitempty"`
		DevelopmentPlan    string `json:"development_plan,omitempty"`
		JobObjective       string `json:"job_objective,omitempty"`
		HealthStatus       string `json:"health_status,omitempty"`
		ImagePath          string `json:"image_path,omitempty"`
		ImageName          string `json:"image_name,omitempty"`
		ResignCode         string `json:"resign_code,omitempty"`
		ResignRequest      string `json:"resign_request,omitempty"`
		ResignStatus       string `json:"resign_status,omitempty"`
		ResignReason       string `json:"resign_reason,omitempty"`
		ResignApprovedBy   string `json:"resign_approved_by,omitempty"`
		ResignApprovedDate string `json:"resign_approved_date,omitempty"`
		EmploymentType     string `json:"employment_type,omitempty"`
	}{
		Alias:              (Alias)(e),
		CreatedAt:          derefTime(e.CreatedAt),
		UpdatedAt:          derefTime(e.UpdatedAt),
		IssueDate:          derefTime(e.IssueDate),
		JoinDate:           derefTime(e.JoinDate),
		OfficialDate:       derefTime(e.OfficialDate),
		ResignDate:         derefTime(e.ResignDate),
		InsuranceDate:      derefTime(e.InsuranceDate),
		HealthInsurExpire:  derefTime(e.HealthInsurExpire),
		GraduationYear:     derefInt64(e.GraduationYear),
		InsuranceAmount:    derefInt64(e.InsuranceAmount),
		HealthInsurance:    derefStr(e.HealthInsurance),
		TaxID:              derefStr(e.TaxID),
		BankAccount:        derefStr(e.BankAccount),
		BankName:           derefStr(e.BankName),
		Status:             derefStr(e.Status),
		Portrait:           derefStr(e.Portrait),
		DevelopmentPlan:    derefStr(e.DevelopmentPlan),
		JobObjective:       derefStr(e.JobObjective),
		HealthStatus:       derefStr(e.HealthStatus),
		ImagePath:          derefStr(e.ImagePath),
		ImageName:          derefStr(e.ImageName),
		ResignCode:         derefStr(e.ResignCode),
		ResignRequest:      derefTime(e.ResignRequest),
		ResignStatus:       derefStr(e.ResignStatus),
		ResignReason:       derefStr(e.ResignReason),
		ResignApprovedBy:   derefStr(e.ResignApprovedBy),
		ResignApprovedDate: derefTime(e.ResignApprovedDate),
		EmploymentType:     derefStr(e.EmploymentType),
	}
	return json.Marshal(aux)
}

func derefTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

func derefInt64(i *int64) string {
	if i == nil {
		return ""
	}
	return fmt.Sprintf("%d", *i)
}

func derefInt(i *int) string {
	if i == nil {
		return ""
	}
	return fmt.Sprintf("%d", *i)
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

type Relative struct {
	ID           string     `json:"id"`
	EmployeeID   string     `json:"employee_id,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	FullName     *string    `json:"full_name"`
	BirthDate    *time.Time `json:"birth_date,omitempty"`
	Phone        *string    `json:"phone,omitempty"`
	Relationship *string    `json:"relationship,omitempty"`
	IsEmergency  *string    `json:"is_emergency,omitempty"`
	IsRelative   *string    `json:"is_relative,omitempty"`
	Address      *string    `json:"address,omitempty"`
	Gender       *string    `json:"gender,omitempty"`
	Age          *int64     `json:"age,omitempty"`
}

func (r Relative) MarshalJSON() ([]byte, error) {
	type Alias Relative
	aux := struct {
		Alias
		FullName     string `json:"full_name"`
		Phone        string `json:"phone,omitempty"`
		Relationship string `json:"relationship,omitempty"`
		IsEmergency  string `json:"is_emergency,omitempty"`
		IsRelative   string `json:"is_relative,omitempty"`
		Address      string `json:"address,omitempty"`
		Gender       string `json:"gender,omitempty"`
		Age          int64  `json:"age,omitempty"`
	}{
		Alias:        (Alias)(r),
		FullName:     derefStr(r.FullName),
		Phone:        derefStr(r.Phone),
		Relationship: derefStr(r.Relationship),
		IsEmergency:  derefStr(r.IsEmergency),
		IsRelative:   derefStr(r.IsRelative),
		Address:      derefStr(r.Address),
		Gender:       derefStr(r.Gender),
		Age:          derefInt64Default(r.Age),
	}
	return json.Marshal(aux)
}

func derefInt64Default(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

type Certificate struct {
	ID              string     `json:"id"`
	EmployeeID      string     `json:"employee_id,omitempty"`
	CreatedAt       *time.Time `json:"created_at,omitempty"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
	CertificateName string     `json:"certificate_name,omitempty"`
	Major           string     `json:"major,omitempty"`
	Classification  string     `json:"classification,omitempty"`
	IssueDate       *time.Time `json:"issue_date,omitempty"`
	ExpiryDate      *time.Time `json:"expiry_date,omitempty"`
}

type ContractHistory struct {
	ID             string     `json:"id"`
	EmployeeID     string     `json:"employee_id,omitempty"`
	ContractNo     string     `json:"contract_no,omitempty"`
	ContractType   string     `json:"contract_type,omitempty"`
	SignDate       *time.Time `json:"sign_date,omitempty"`
	DurationMonths *int64     `json:"duration_months,omitempty"`
	SignedBy       *string    `json:"signed_by,omitempty"`
	FilePath       *string    `json:"file_path,omitempty"`
	FileName       *string    `json:"file_name,omitempty"`
}

type Salary struct {
	ID          string     `json:"id"`
	EmployeeID  string     `json:"employee_id,omitempty"`
	Type        string     `json:"type"`
	Description string     `json:"description,omitempty"`
	AmountOld   float64    `json:"amount_old"`
	AmountNew   float64    `json:"amount_new"`
	StartDate   time.Time  `json:"start_date"`
	FilePath    *string    `json:"file_path,omitempty"`
	FileName    *string    `json:"file_name,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	CreatedBy   string     `json:"created_by,omitempty"`
}

type CareerHistory struct {
	ID         string     `json:"id"`
	EmployeeID string     `json:"employee_id,omitempty"`
	Position   *string    `json:"position"`
	Department *string    `json:"department"`
	Rank       *string    `json:"rank,omitempty"`
	StartDate  time.Time  `json:"start_date"`
	EndDate    *time.Time `json:"end_date,omitempty"`
}

type SalaryHistory struct {
	ID           string    `json:"id"`
	EmployeeID   string    `json:"employee_id,omitempty"`
	OldSalary    float64   `json:"old_salary"`
	NewSalary    float64   `json:"new_salary"`
	AdjustDate   time.Time `json:"adjust_date"`
	Reason       string    `json:"reason,omitempty"`
	DocumentPath string    `json:"document_path,omitempty"`
}

type PerformanceEvaluation struct {
	ID             string     `json:"id"`
	EmployeeID     string     `json:"employee_id,omitempty"`
	EvaluationType string     `json:"evaluation_type"`
	Purpose        string     `json:"purpose,omitempty"`
	Result         string     `json:"result,omitempty"`
	Score          *int64     `json:"score,omitempty"`
	EvaluationDate time.Time  `json:"evaluation_date"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}

type RewardDiscipline struct {
	ID            string     `json:"id"`
	EmployeeID    string     `json:"employee_id,omitempty"`
	Type          string     `json:"type"`
	Description   string     `json:"description,omitempty"`
	Title         string     `json:"title,omitempty"`         //danh hiệu
	DecisionForm  string     `json:"decision_form,omitempty"` // hình thức khen thưởng/kỷ luật
	EffectiveDate time.Time  `json:"effective_date"`
	ExpiryDate    time.Time  `json:"expiry_date"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}

type EmployeeResult struct {
	ID             string     `json:"id"`
	AttendanceCode string     `json:"attendance_code"`
	FullName       string     `json:"full_name"`
	BirthDate      *time.Time `json:"birth_date,omitempty"`
	DepartmentName string     `json:"department_name"`
	PositionTitle  string     `json:"position_title"`
	CompanyPhone   *string    `json:"company_phone"`
	PersonalPhone  *string    `json:"personal_phone"`
	CompanyEmail   *string    `json:"company_email"`
}

func (r EmployeeResult) MarshalJSON() ([]byte, error) {
	type empRs EmployeeResult
	aux := struct {
		empRs
		CompanyPhone  string `json:"company_phone"`
		PersonalPhone string `json:"personal_phone"`
		CompanyEmail  string `json:"company_email"`
	}{
		empRs: empRs{
			ID:             r.ID,
			AttendanceCode: r.AttendanceCode,
			FullName:       r.FullName,
			BirthDate:      r.BirthDate,
			DepartmentName: r.DepartmentName,
			PositionTitle:  r.PositionTitle,
		},
	}
	return json.Marshal(aux)
}

type ChangesEmployee struct {
	TypeReport     string `json:"type_report"`
	CountLeave     int64  `json:"count_leave"`
	CountJoin      int64  `json:"count_join"`
	CountEmployee  int64  `json:"count_employee"`
	DepartmentName string `json:"department_name"`
	WorkLocation   string `json:"work_location"`
}

// ChangesEmployeeFilterDTO is the transport object for binding from JSON or query params.
// Dates are strings here for flexible client input; they will be parsed to time.Time in the handler.
type ChangesEmployeeFilterDTO struct {
	FromDate   string `json:"from_date" form:"from_date"`
	ToDate     string `json:"to_date" form:"to_date"`
	TypeReport string `json:"type_report" form:"type_report"`
}

// ChangesEmployeeFilter is the internal filter used by DAO with parsed time values
type ChangesEmployeeFilter struct {
	FromDate   time.Time `json:"from_date,omitempty"`
	ToDate     time.Time `json:"to_date,omitempty"`
	TypeReport string    `json:"type_report,omitempty"`
}

type RecruitmentPlan struct {
	ID           string     `json:"id"`
	PlanYear     *int       `json:"plan_year,omitempty"`
	WorkLocation *string    `json:"work_location,omitempty"`
	Department   *string    `json:"department,omitempty"`
	Position     *string    `json:"position,omitempty"`
	QuantityPlan *int       `json:"quantity_plan,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	Difficulty   *string    `json:"difficulty,omitempty"` // mức độ khó tuyển
	Solution     *string    `json:"reason,omitempty"`     // nguyên nhân
	Proposal     *string    `json:"proposal,omitempty"`   // đề xuất
}

func (r RecruitmentPlan) MarshalJSON() ([]byte, error) {
	type Alias RecruitmentPlan
	aux := struct {
		Alias
		PlanYear     string `json:"plan_year,omitempty"`
		WorkLocation string `json:"work_location,omitempty"`
		Department   string `json:"department,omitempty"`
		Position     string `json:"position,omitempty"`
		QuantityPlan string `json:"quantity_plan,omitempty"`
		CreatedAt    string `json:"created_at,omitempty"`
		Difficulty   string `json:"difficulty,omitempty"` // mức độ khó tuyển
		Solution     string `json:"reason,omitempty"`     // nguyên nhân
		Proposal     string `json:"proposal,omitempty"`   // đề xuất
	}{
		Alias:        (Alias)(r),
		PlanYear:     derefInt(r.PlanYear),
		WorkLocation: derefStr(r.WorkLocation),
		Department:   derefStr(r.Department),
		Position:     derefStr(r.Position),
		QuantityPlan: derefInt(r.QuantityPlan),
		CreatedAt:    derefTime(r.CreatedAt),
		Difficulty:   derefStr(r.Difficulty),
		Solution:     derefStr(r.Solution),
		Proposal:     derefStr(r.Proposal),
	}
	return json.Marshal(aux)
}
