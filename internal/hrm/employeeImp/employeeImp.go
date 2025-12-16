package employeeimp

import (
	"context"
	"europm/internal/db/dbhrm"
	"europm/internal/hrm/employee"
	"europm/internal/hrm/model"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
)

type EmployeeImp struct {
	ctx context.Context
}

func GetInstance(ctx context.Context) employee.Employee {
	return &EmployeeImp{
		ctx: ctx,
	}
}

func (e *EmployeeImp) GetTotalEmployee(attendanceCode string, fullName string) (int, error) {
	var total int
	query := "select hrm.get_employee($1,$2,$3)"
	var cursor string
	tx, _ := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	defer func() {
		if tx != nil {
			tx.Commit(e.ctx)
		}
	}()
	row := tx.QueryRow(e.ctx, query, attendanceCode, fullName, "TOTAL")
	err := row.Scan(&cursor)
	if err != nil {
		log.Printf("Error execute", err)
		return 0, err
	}
	rows, err := tx.Query(e.ctx, "FETCH ALL "+pq.QuoteIdentifier(cursor)+";")
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		log.Printf("Error FETCH refCurser: %v", err)
		return 0, err
	}
	for rows.Next() {
		err := rows.Scan(&total)
		log.Printf("Total employee: %d", total)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return total, err
		}
	}
	return total, nil
}

func (e *EmployeeImp) GetEmployee(attendanceCode string, fullName string) ([]model.EmployeeResult, error) {
	lstEmployee := make([]model.EmployeeResult, 0)
	query := "select hrm.get_employee($1,$2,$3)"
	var cursor string
	tx, _ := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	defer func() {
		if tx != nil {
			tx.Commit(e.ctx)
		}
	}()
	row := tx.QueryRow(e.ctx, query, attendanceCode, fullName, "SELECT")
	err := row.Scan(&cursor)
	if err != nil {
		log.Printf("Error execute", err)
		return nil, err
	}
	rows, err := tx.Query(e.ctx, "FETCH ALL "+pq.QuoteIdentifier(cursor)+";")
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		log.Printf("Error FETCH refCurser: %v", err)
		return nil, err
	}
	for rows.Next() {
		employ := model.EmployeeResult{}
		err := rows.Scan(&employ.ID, &employ.AttendanceCode, &employ.FullName, &employ.BirthDate, &employ.DepartmentName, &employ.PositionTitle, &employ.CompanyPhone, &employ.PersonalPhone, &employ.CompanyEmail)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		lstEmployee = append(lstEmployee, employ)
	}
	return lstEmployee, nil
}

func (e *EmployeeImp) GetDetailEmployeeByID(id string) (model.Employee, error) {
	log.Printf("GetDetailEmployeeByID: %s", id)
	var employ model.Employee
	query := "select hrm.get_employee_by_id($1)"
	var cursor string
	tx, _ := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	defer func() {
		if tx != nil {
			tx.Commit(e.ctx)
		}
	}()
	row := tx.QueryRow(e.ctx, query, id)
	err := row.Scan(&cursor)
	if err != nil {
		log.Printf("Error execute", err)
		return model.Employee{}, err
	}
	rows, err := tx.Query(e.ctx, "FETCH ALL "+pq.QuoteIdentifier(cursor)+";")
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		log.Printf("Error FETCH refCurser: %v", err)
		return model.Employee{}, err
	}
	for rows.Next() {
		err := rows.Scan(&employ.ID,
			&employ.CreatedAt,
			&employ.UpdatedAt,
			&employ.FullName,
			&employ.BirthDate,
			&employ.Gender,
			&employ.IDNumber,
			&employ.IssueDate,
			&employ.IssuePlace,
			&employ.BirthPlace,
			&employ.HomeTown,
			&employ.PermanentAddress,
			&employ.TemporaryAddress,
			&employ.MaritalStatus,
			&employ.PersonalPhone,
			&employ.PersonalEmail,
			&employ.CompanyPhone,
			&employ.CompanyEmail,
			&employ.JoinDate,
			&employ.HighestDegree,
			&employ.Major,
			&employ.SchoolName,
			&employ.GraduationYear,
			&employ.SpecialSkills,
			&employ.AttendanceCode,
			&employ.PositionTitle,
			&employ.DepartmentName,
			&employ.Rank,
			&employ.WorkLocation,
			&employ.OfficialDate,
			&employ.ResignDate,
			&employ.Leader,
			&employ.ManagerID,
			&employ.SocialInsuranceNo,
			&employ.InsuranceStatus,
			&employ.InsuranceDate,
			&employ.InsuranceAmount,
			&employ.KCBPlace,
			&employ.HealthInsurance,
			&employ.HealthInsurExpire,
			&employ.TaxID,
			&employ.BankAccount,
			&employ.BankName,
			&employ.Status,
			&employ.Portrait,
			&employ.DevelopmentPlan,
			&employ.JobObjective,
			&employ.HealthStatus,
			&employ.ImagePath,
			&employ.ImageName,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return model.Employee{}, err
		}
	}
	return employ, nil
}

func (e *EmployeeImp) InstEmployee(emp model.Employee) (string, error) {
	log.Printf("emp.id: %s, emp.FullName: %s, emp.BirthDate: %v, emp.Gender: %v, emp.IDNumber: %v, emp.IssueDate: %v, emp.IssuePlace: %v, emp.BirthPlace: %v, emp.HomeTown: %v, emp.PermanentAddress: %v, emp.TemporaryAddress: %v, emp.MaritalStatus: %v, emp.PersonalPhone: %v, emp.PersonalEmail: %v, emp.CompanyPhone: %v, emp.CompanyEmail: %v, emp.JoinDate: %v, emp.HighestDegree: %v, emp.Major: %v, emp.SchoolName: %v, emp.GraduationYear: %v, emp.SpecialSkills: %v, emp.AttendanceCode: %v, emp.PositionTitle: %v, emp.DepartmentName: %v, emp.Rank: %v, emp.WorkLocation: %v, emp.OfficialDate: %v, emp.ResignDate: %v, emp.Leader: %v, emp.ManagerID: %v, emp.SocialInsuranceNo: %v, emp.InsuranceStatus: %v, emp.InsuranceDate: %v, emp.InsuranceAmount: %v, emp.KCBPlace: %v, emp.HealthInsurance: %v, emp.HealthInsurExpire: %v, emp.TaxID: %v, emp.BankAccount: %v, emp.BankName: %v, emp.Status: %v, emp.Portrait: %v, emp.DevelopmentPlan: %v, emp.JobObjective: %v, emp.HealthStatus: %v, emp.ImagePath: %v, emp.ImageName: %v",
		emp.ID, emp.FullName, emp.BirthDate, emp.Gender, emp.IDNumber, emp.IssueDate, emp.IssuePlace, emp.BirthPlace, emp.HomeTown, emp.PermanentAddress, emp.TemporaryAddress, emp.MaritalStatus, emp.PersonalPhone, emp.PersonalEmail, emp.CompanyPhone, emp.CompanyEmail, emp.JoinDate, emp.HighestDegree, emp.Major, emp.SchoolName, emp.GraduationYear, emp.SpecialSkills, emp.AttendanceCode, emp.PositionTitle, emp.DepartmentName, emp.Rank, emp.WorkLocation, emp.OfficialDate, emp.ResignDate, emp.Leader, emp.ManagerID, emp.SocialInsuranceNo, emp.InsuranceStatus, emp.InsuranceDate, emp.InsuranceAmount, emp.KCBPlace, emp.HealthInsurance, emp.HealthInsurExpire, emp.TaxID, emp.BankAccount, emp.BankName, emp.Status, emp.Portrait, emp.DevelopmentPlan, emp.JobObjective, emp.HealthStatus, emp.ImagePath, emp.ImageName)
	var id string
	query := "select hrm.insert_employee($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35,$36,$37,$38,$39,$40,$41,$42,$43,$44,$45,$46,$47,$48)"

	args := []interface{}{
		emp.ID,
		emp.FullName,
		emp.BirthDate,
		emp.Gender,
		emp.IDNumber,
		emp.IssueDate,
		emp.IssuePlace,
		emp.BirthPlace,
		emp.HomeTown,
		emp.PermanentAddress,
		emp.TemporaryAddress,
		emp.MaritalStatus,
		emp.PersonalPhone,
		emp.PersonalEmail,
		emp.CompanyPhone,
		emp.CompanyEmail,
		emp.JoinDate,
		emp.HighestDegree,
		emp.Major,
		emp.SchoolName,
		emp.GraduationYear,
		emp.SpecialSkills,
		emp.AttendanceCode,
		emp.PositionTitle,
		emp.DepartmentName,
		emp.Rank,
		emp.WorkLocation,
		emp.OfficialDate,
		emp.ResignDate,
		emp.Leader,
		emp.ManagerID,
		emp.SocialInsuranceNo,
		emp.InsuranceStatus,
		emp.InsuranceDate,
		emp.InsuranceAmount,
		emp.KCBPlace,
		emp.HealthInsurance,
		emp.HealthInsurExpire,
		emp.TaxID,
		emp.BankAccount,
		emp.BankName,
		emp.Status,
		emp.Portrait,
		emp.DevelopmentPlan,
		emp.JobObjective,
		emp.HealthStatus,
		emp.ImagePath,
		emp.ImageName,
	}

	tx, err := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	if err != nil {
		return "", err
	}
	defer func() {
		if tx != nil {
			tx.Commit(e.ctx)
		}
	}()

	row := tx.QueryRow(e.ctx, query, args...)
	if err := row.Scan(&id); err != nil {
		log.Printf("Error execute: %v", err)
		return "", err
	}
	return id, nil
}

func (e *EmployeeImp) InsertRelatives(relatives []model.Relative) (string, error) {
	query := "select hrm.insert_relative($1,$2,$3,$4,$5,$6)"
	tx, err := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			tx.Rollback(e.ctx)
		} else {
			tx.Commit(e.ctx)
		}
	}()
	for _, relative := range relatives {
		log.Printf("relative.EmployeeID: %v, relative.ID: %v, relative.FullName: %v, relative.BirthDate: %v, relative.Relationship: %v, relative.Gender: %v",
			relative.EmployeeID, relative.ID, relative.FullName, relative.BirthDate, relative.Relationship, relative.Gender)
		_, err := tx.Exec(e.ctx, query,
			relative.EmployeeID,
			relative.ID,
			relative.FullName,
			relative.BirthDate,
			relative.Relationship,
			relative.Gender)
		if err != nil {
			log.Printf("Error execute: %v", err)
			return "", err
		}
	}
	return "", nil
}

func (e *EmployeeImp) InsertEmergencyContacts(contacts []model.Relative) (string, error) {
	query := "select hrm.insert_emergency($1,$2,$3,$4,$5,$6)"
	tx, err := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			tx.Rollback(e.ctx)
		} else {
			tx.Commit(e.ctx)
		}
	}()
	for _, contact := range contacts {
		log.Printf("contact.EmployeeID: %v, contact.ID: %v, contact.FullName: %v, contact.Relationship: %v, contact.Phone: %v, contact.Address: %v",
			contact.EmployeeID, contact.ID, contact.FullName, contact.Relationship, contact.Phone, contact.Address)
		_, err := tx.Exec(e.ctx, query,
			contact.EmployeeID,
			contact.ID,
			contact.FullName,
			contact.Relationship,
			contact.Phone,
			contact.Address)
		if err != nil {
			log.Printf("Error execute: %v", err)
			return "", err
		}
	}
	return "", nil
}

func (e *EmployeeImp) InsertCertificates(certificates []model.Certificate) (string, error) {
	query := "select hrm.insert_certificate($1,$2,$3,$4,$5,$6,$7)"
	tx, err := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			tx.Rollback(e.ctx)
		} else {
			tx.Commit(e.ctx)
		}
	}()
	for _, certificate := range certificates {
		_, err := tx.Exec(e.ctx, query,
			certificate.EmployeeID,
			certificate.ID,
			certificate.CertificateName,
			certificate.Major,
			certificate.Classification,
			certificate.IssueDate,
			certificate.ExpiryDate)
		if err != nil {
			log.Printf("Error execute: %v", err)
			return "", err
		}
	}
	return "", nil
}

func (e *EmployeeImp) InsertSalaries(salaries model.Salary) (string, error) {
	var id string
	query := "select hrm.insert_salary($1,$2,$3,$4,$5,$6,$7,$8,$9)"
	args := []interface{}{
		salaries.EmployeeID,
		salaries.ID,
		salaries.Type,
		salaries.Description,
		salaries.AmountOld,
		salaries.AmountNew,
		salaries.FileName,
		salaries.FilePath,
		salaries.StartDate,
	}

	tx, err := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	if err != nil {
		return "", err
	}
	defer func() {
		if tx != nil {
			tx.Commit(e.ctx)
		}
	}()

	row := tx.QueryRow(e.ctx, query, args...)
	if err := row.Scan(&id); err != nil {
		log.Printf("Error execute: %v", err)
		return "", err
	}
	return id, nil
}

func (e *EmployeeImp) InsertCareerHistories(careerHistories []model.CareerHistory) (string, error) {
	query := "select hrm.insert_career_history($1,$2,$3,$4,$5,$6)"
	tx, err := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			tx.Rollback(e.ctx)
		} else {
			tx.Commit(e.ctx)
		}
	}()
	for _, careerHistory := range careerHistories {
		_, err := tx.Exec(e.ctx, query,
			careerHistory.EmployeeID,
			careerHistory.ID,
			careerHistory.Position,
			careerHistory.Department,
			careerHistory.Rank,
			careerHistory.StartDate)
		if err != nil {
			log.Printf("Error execute: %v", err)
			return "", err
		}
	}
	return "", nil
}

func (e *EmployeeImp) InsertPerformanceEvaluations(performanceEvaluations []model.PerformanceEvaluation) (string, error) {
	query := "select hrm.insert_performance_evaluation($1,$2,$3,$4,$5,$6,$7)"
	tx, err := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			tx.Rollback(e.ctx)
		} else {
			tx.Commit(e.ctx)
		}
	}()
	for _, performanceEvaluation := range performanceEvaluations {
		_, err := tx.Exec(e.ctx, query,
			performanceEvaluation.EmployeeID,
			performanceEvaluation.ID,
			performanceEvaluation.EvaluationType,
			performanceEvaluation.Purpose,
			performanceEvaluation.Result,
			performanceEvaluation.Score,
			performanceEvaluation.EvaluationDate)
		if err != nil {
			log.Printf("Error execute: %v", err)
			return "", err
		}
	}
	return "", nil
}

func (e *EmployeeImp) InsertRewardDiscipline(rds []model.RewardDiscipline) (string, error) {
	query := "select hrm.insert_reward_discipline($1,$2,$3,$4,$5,$6,$7,$8)"
	tx, err := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			tx.Rollback(e.ctx)
		} else {
			tx.Commit(e.ctx)
		}
	}()
	for _, rd := range rds {
		_, err := tx.Exec(e.ctx, query,
			rd.EmployeeID,
			rd.ID,
			rd.Type,
			rd.Description,
			rd.Title,
			rd.DecisionForm,
			rd.EffectiveDate,
			rd.ExpiryDate)
		if err != nil {
			log.Printf("Error execute: %v", err)
			return "", err
		}
	}
	return "", nil
}

func (e *EmployeeImp) InsertContractHistories(rds model.ContractHistory) (string, error) {
	var id string
	query := "select hrm.insert_contract_history($1,$2,$3,$4,$5,$6,$7,$8)"
	args := []interface{}{
		rds.EmployeeID,
		rds.ID,
		rds.ContractType,
		rds.ContractNo,
		rds.SignDate,
		rds.FilePath,
		rds.FileName,
		rds.DurationMonths,
	}

	tx, err := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	if err != nil {
		return "", err
	}
	defer func() {
		if tx != nil {
			tx.Commit(e.ctx)
		}
	}()

	row := tx.QueryRow(e.ctx, query, args...)
	if err := row.Scan(&id); err != nil {
		log.Printf("Error execute: %v", err)
		return "", err
	}
	return id, nil
}

func (e *EmployeeImp) GetRelativesByEmployeeID(id string) ([]model.Relative, error) {
	lstRelatives := make([]model.Relative, 0)
	query := "select hrm.get_relative_by_id($1)"
	var cursor string
	tx, err := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	if err != nil {
		log.Printf("Error begin tx: %v", err)
		return nil, err
	}
	defer func() {
		if tx != nil {
			tx.Commit(e.ctx)
		}
	}()
	row := tx.QueryRow(e.ctx, query, id)
	err = row.Scan(&cursor)
	if err != nil {
		log.Printf("Error execute: %v", err)
		return nil, err
	}
	rows, err := tx.Query(e.ctx, "FETCH ALL "+pq.QuoteIdentifier(cursor)+";")
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		log.Printf("Error FETCH refCursor: %v", err)
		return nil, err
	}
	for rows.Next() {
		var relative model.Relative
		if err := rows.Scan(&relative.ID, &relative.EmployeeID, &relative.FullName, &relative.Relationship, &relative.BirthDate, &relative.Gender); err != nil {
			log.Printf("Error scan: %v", err)
			return nil, err
		}
		lstRelatives = append(lstRelatives, relative)
	}
	return lstRelatives, nil
}

func (e *EmployeeImp) GetEmergencyContactsByEmployeeID(id string) ([]model.Relative, error) {
	lstRelatives := make([]model.Relative, 0)
	query := "select hrm.get_emergency_contact_by_id($1)"
	var cursor string
	tx, err := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	if err != nil {
		log.Printf("Error begin tx: %v", err)
		return nil, err
	}
	defer func() {
		if tx != nil {
			tx.Commit(e.ctx)
		}
	}()
	row := tx.QueryRow(e.ctx, query, id)
	err = row.Scan(&cursor)
	if err != nil {
		log.Printf("Error execute: %v", err)
		return nil, err
	}
	rows, err := tx.Query(e.ctx, "FETCH ALL "+pq.QuoteIdentifier(cursor)+";")
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		log.Printf("Error FETCH refCursor: %v", err)
		return nil, err
	}
	for rows.Next() {
		var relative model.Relative
		if err := rows.Scan(&relative.ID, &relative.EmployeeID, &relative.FullName, &relative.Relationship, &relative.Phone, &relative.Address); err != nil {
			log.Printf("Error scan: %v", err)
			return nil, err
		}
		lstRelatives = append(lstRelatives, relative)
	}
	return lstRelatives, nil
}

func (e *EmployeeImp) GetCertificatesByEmployeeID(id string) ([]model.Certificate, error) {
	lstCertificates := make([]model.Certificate, 0)
	query := "select hrm.get_certificate_by_id($1)"
	var cursor string
	tx, err := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	if err != nil {
		log.Printf("Error begin tx: %v", err)
		return nil, err
	}
	defer func() {
		if tx != nil {
			tx.Commit(e.ctx)
		}
	}()
	row := tx.QueryRow(e.ctx, query, id)
	err = row.Scan(&cursor)
	if err != nil {
		log.Printf("Error execute: %v", err)
		return nil, err
	}
	rows, err := tx.Query(e.ctx, "FETCH ALL "+pq.QuoteIdentifier(cursor)+";")
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		log.Printf("Error FETCH refCursor: %v", err)
		return nil, err
	}
	for rows.Next() {
		var certificate model.Certificate
		if err := rows.Scan(&certificate.ID, &certificate.EmployeeID, &certificate.CertificateName, &certificate.Major, &certificate.Classification, &certificate.IssueDate, &certificate.ExpiryDate); err != nil {
			log.Printf("Error scan: %v", err)
			return nil, err
		}
		lstCertificates = append(lstCertificates, certificate)
	}
	return lstCertificates, nil
}

func (e *EmployeeImp) GetSalariesByEmployeeID(id string) ([]model.Salary, error) {
	lstSalaries := make([]model.Salary, 0)
	query := "select hrm.get_salary_by_id($1)"
	var cursor string
	tx, err := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	if err != nil {
		log.Printf("Error begin tx: %v", err)
		return nil, err
	}
	defer func() {
		if tx != nil {
			tx.Commit(e.ctx)
		}
	}()
	row := tx.QueryRow(e.ctx, query, id)
	err = row.Scan(&cursor)
	if err != nil {
		log.Printf("Error execute: %v", err)
		return nil, err
	}
	rows, err := tx.Query(e.ctx, "FETCH ALL "+pq.QuoteIdentifier(cursor)+";")
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		log.Printf("Error FETCH refCursor: %v", err)
		return nil, err
	}
	for rows.Next() {
		var salary model.Salary
		if err := rows.Scan(&salary.ID, &salary.EmployeeID, &salary.Type, &salary.Description, &salary.AmountOld, &salary.AmountNew, &salary.StartDate, &salary.EndDate); err != nil {
			log.Printf("Error scan: %v", err)
			return nil, err
		}
		lstSalaries = append(lstSalaries, salary)
	}
	return lstSalaries, nil
}

func (e *EmployeeImp) GetCareerHistoriesByEmployeeID(id string) ([]model.CareerHistory, error) {
	lstCareerHistories := make([]model.CareerHistory, 0)
	query := "select hrm.get_career_history_by_id($1)"
	var cursor string
	tx, err := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	if err != nil {
		log.Printf("Error begin tx: %v", err)
		return nil, err
	}
	defer func() {
		if tx != nil {
			tx.Commit(e.ctx)
		}
	}()
	row := tx.QueryRow(e.ctx, query, id)
	err = row.Scan(&cursor)
	if err != nil {
		log.Printf("Error execute: %v", err)
		return nil, err
	}
	rows, err := tx.Query(e.ctx, "FETCH ALL "+pq.QuoteIdentifier(cursor)+";")
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		log.Printf("Error FETCH refCursor: %v", err)
		return nil, err
	}
	for rows.Next() {
		var careerHistory model.CareerHistory
		if err := rows.Scan(&careerHistory.ID, &careerHistory.EmployeeID, &careerHistory.Position, &careerHistory.Department, &careerHistory.Rank, &careerHistory.EndDate); err != nil {
			log.Printf("Error scan: %v", err)
			return nil, err
		}
		lstCareerHistories = append(lstCareerHistories, careerHistory)
	}
	return lstCareerHistories, nil
}

func (e *EmployeeImp) GetPerformanceEvaluationsByEmployeeID(id string) ([]model.PerformanceEvaluation, error) {
	lstPerformanceEvaluations := make([]model.PerformanceEvaluation, 0)
	query := "select hrm.get_performance_evaluation_by_id($1)"
	var cursor string
	tx, err := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	if err != nil {
		log.Printf("Error begin tx: %v", err)
		return nil, err
	}
	defer func() {
		if tx != nil {
			tx.Commit(e.ctx)
		}
	}()
	row := tx.QueryRow(e.ctx, query, id)
	err = row.Scan(&cursor)
	if err != nil {
		log.Printf("Error execute: %v", err)
		return nil, err
	}
	rows, err := tx.Query(e.ctx, "FETCH ALL "+pq.QuoteIdentifier(cursor)+";")
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		log.Printf("Error FETCH refCursor: %v", err)
		return nil, err
	}
	for rows.Next() {
		var per model.PerformanceEvaluation
		if err := rows.Scan(&per.ID, &per.EmployeeID, &per.EvaluationType, &per.Purpose, &per.Result, &per.Score, &per.EvaluationDate); err != nil {
			log.Printf("Error scan: %v", err)
			return nil, err
		}
		lstPerformanceEvaluations = append(lstPerformanceEvaluations, per)
	}
	return lstPerformanceEvaluations, nil
}

func (e *EmployeeImp) GetRewardDisciplinesByEmployeeID(id string) ([]model.RewardDiscipline, error) {
	lstRewardDisciplines := make([]model.RewardDiscipline, 0)
	query := "select hrm.get_reward_disciptline_by_id($1)"
	var cursor string
	tx, err := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	if err != nil {
		log.Printf("Error begin tx: %v", err)
		return nil, err
	}
	defer func() {
		if tx != nil {
			tx.Commit(e.ctx)
		}
	}()
	row := tx.QueryRow(e.ctx, query, id)
	err = row.Scan(&cursor)
	if err != nil {
		log.Printf("Error execute: %v", err)
		return nil, err
	}
	rows, err := tx.Query(e.ctx, "FETCH ALL "+pq.QuoteIdentifier(cursor)+";")
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		log.Printf("Error FETCH refCursor: %v", err)
		return nil, err
	}
	for rows.Next() {
		var rd model.RewardDiscipline
		if err := rows.Scan(&rd.ID, &rd.EmployeeID, &rd.Type, &rd.Description, &rd.Title, &rd.DecisionForm, &rd.EffectiveDate); err != nil {
			log.Printf("Error scan: %v", err)
			return nil, err
		}
		lstRewardDisciplines = append(lstRewardDisciplines, rd)
	}
	return lstRewardDisciplines, nil
}

func (e *EmployeeImp) GetContractHistoriesByEmployeeID(id string) ([]model.ContractHistory, error) {
	lstContractHistories := make([]model.ContractHistory, 0)
	query := "select hrm.get_contract_history_by_id($1)"
	var cursor string
	tx, err := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	if err != nil {
		log.Printf("Error begin tx: %v", err)
		return nil, err
	}
	defer func() {
		if tx != nil {
			tx.Commit(e.ctx)
		}
	}()
	row := tx.QueryRow(e.ctx, query, id)
	err = row.Scan(&cursor)
	if err != nil {
		log.Printf("Error execute: %v", err)
		return nil, err
	}
	rows, err := tx.Query(e.ctx, "FETCH ALL "+pq.QuoteIdentifier(cursor)+";")
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		log.Printf("Error FETCH refCursor: %v", err)
		return nil, err
	}
	for rows.Next() {
		var ct model.ContractHistory
		if err := rows.Scan(&ct.ID, &ct.EmployeeID, &ct.ContractType, &ct.ContractNo, &ct.SignDate, &ct.FilePath, &ct.FileName, &ct.DurationMonths); err != nil {
			log.Printf("Error scan: %v", err)
			return nil, err
		}
		lstContractHistories = append(lstContractHistories, ct)
	}
	return lstContractHistories, nil
}

func (e *EmployeeImp) DeleteCertificatesByID(id string) (string, error) {
	query := "select hrm.delete_certificates_by_id($1)"
	tx, err := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			tx.Rollback(e.ctx)
		} else {
			tx.Commit(e.ctx)
		}
	}()
	_, err = tx.Exec(e.ctx, query, id)
	if err != nil {
		log.Printf("Error execute: %v", err)
		return "", err
	}
	return id, nil
}

func (e *EmployeeImp) DeleteRelativesByID(id string) (string, error) {
	query := "select hrm.delete_relatives_by_id($1)"
	tx, err := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			tx.Rollback(e.ctx)
		} else {
			tx.Commit(e.ctx)
		}
	}()
	_, err = tx.Exec(e.ctx, query, id)
	if err != nil {
		log.Printf("Error execute: %v", err)
		return "", err
	}
	return id, nil
}

func (e *EmployeeImp) DeleteSalariesByID(id string) (string, error) {
	query := "select hrm.delete_salaries_by_id($1)"
	tx, err := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			tx.Rollback(e.ctx)
		} else {
			tx.Commit(e.ctx)
		}
	}()
	_, err = tx.Exec(e.ctx, query, id)
	if err != nil {
		log.Printf("Error execute: %v", err)
		return "", err
	}
	return id, nil
}

func (e *EmployeeImp) DeleteCareerHistoriesByID(id string) (string, error) {
	query := "select hrm.delete_career_histories_by_id($1)"
	tx, err := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			tx.Rollback(e.ctx)
		} else {
			tx.Commit(e.ctx)
		}
	}()
	_, err = tx.Exec(e.ctx, query, id)
	if err != nil {
		log.Printf("Error execute: %v", err)
		return "", err
	}
	return id, nil
}

func (e *EmployeeImp) DeletePerformanceEvaluationsByID(id string) (string, error) {
	query := "select hrm.delete_performance_evaluations_by_id($1)"
	tx, err := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			tx.Rollback(e.ctx)
		} else {
			tx.Commit(e.ctx)
		}
	}()
	_, err = tx.Exec(e.ctx, query, id)
	if err != nil {
		log.Printf("Error execute: %v", err)
		return "", err
	}
	return id, nil
}

func (e *EmployeeImp) DeleteRewardDisciplinesByID(id string) (string, error) {
	query := "select hrm.delete_reward_discipline_by_id($1)"
	tx, err := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			tx.Rollback(e.ctx)
		} else {
			tx.Commit(e.ctx)
		}
	}()
	_, err = tx.Exec(e.ctx, query, id)
	if err != nil {
		log.Printf("Error execute: %v", err)
		return "", err
	}
	return id, nil
}

func (e *EmployeeImp) DeleteContractHistoriesByID(id string) (string, error) {
	query := "select hrm.delete_contract_histories_by_id($1)"
	tx, err := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			tx.Rollback(e.ctx)
		} else {
			tx.Commit(e.ctx)
		}
	}()
	_, err = tx.Exec(e.ctx, query, id)
	if err != nil {
		log.Printf("Error execute: %v", err)
		return "", err
	}
	return id, nil
}

func (e *EmployeeImp) GetTotalEmployeesResign(text string, fromDate time.Time, toDate time.Time) (int, error) {
	log.Printf("GetTotalEmployeesResign: %s %v %v", text, fromDate, toDate)
	var total int
	query := "select hrm.get_employees_resign($1,$2,$3,$4)"
	var cursor string
	tx, _ := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	defer func() {
		if tx != nil {
			tx.Commit(e.ctx)
		}
	}()
	row := tx.QueryRow(e.ctx, query, fromDate, toDate, text, "TOTAL")
	err := row.Scan(&cursor)
	if err != nil {
		log.Printf("Error execute", err)
		return 0, err
	}
	rows, err := tx.Query(e.ctx, "FETCH ALL "+pq.QuoteIdentifier(cursor)+";")
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		log.Printf("Error FETCH refCurser: %v", err)
		return 0, err
	}
	for rows.Next() {
		err := rows.Scan(&total)
		log.Printf("Total employee: %d", total)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return total, err
		}
	}
	return total, nil
}

func (e *EmployeeImp) GetEmployeesResign(text string, fromDate time.Time, toDate time.Time) ([]model.Employee, error) {
	log.Printf("GetEmployeesResign: %s", text)
	lstEmployee := make([]model.Employee, 0)
	query := "select hrm.get_employees_resign($1,$2,$3,$4)"
	tx, err := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		if tx != nil {
			tx.Commit(e.ctx)
		}
	}()
	row := tx.QueryRow(e.ctx, query, fromDate, toDate, text, "SELECT")
	var cursor string
	err = row.Scan(&cursor)
	if err != nil {
		log.Printf("Error execute", err)
		return nil, err
	}
	rows, err := tx.Query(e.ctx, "FETCH ALL "+pq.QuoteIdentifier(cursor)+";")
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		log.Printf("Error FETCH refCurser: %v", err)
		return nil, err
	}
	for rows.Next() {
		var employ model.Employee
		err := rows.Scan(&employ.ID,
			&employ.AttendanceCode,
			&employ.FullName,
			&employ.BirthDate,
			&employ.ResignCode,
			&employ.ResignRequest,
			&employ.ResignDate,
			&employ.ResignStatus,
			&employ.ResignReason,
			&employ.ResignApprovedBy,
			&employ.ResignApprovedDate,
			&employ.EmploymentType,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		lstEmployee = append(lstEmployee, employ)
	}
	return lstEmployee, nil
}

func (e *EmployeeImp) GetEmployeeResignByID(id string) (model.Employee, error) {
	var employ model.Employee
	query := "select hrm.get_employees_resign_by_id($1)"
	var cursor string
	tx, err := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	if err != nil {
		return model.Employee{}, err
	}
	defer func() {
		if tx != nil {
			tx.Commit(e.ctx)
		}
	}()
	row := tx.QueryRow(e.ctx, query, id)
	err = row.Scan(&cursor)
	if err != nil {
		log.Printf("Error execute: %v", err)
		return model.Employee{}, err
	}

	rows, err := tx.Query(e.ctx, "FETCH ALL "+pq.QuoteIdentifier(cursor)+";")
	if err != nil {
		log.Printf("Error FETCH refCursor: %v", err)
		return model.Employee{}, err
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	for rows.Next() {
		err := rows.Scan(
			&employ.ID,
			&employ.AttendanceCode,
			&employ.FullName,
			&employ.BirthDate,
			&employ.ResignCode,
			&employ.ResignRequest,
			&employ.ResignDate,
			&employ.ResignStatus,
			&employ.ResignReason,
			&employ.ResignApprovedBy,
			&employ.ResignApprovedDate,
			&employ.EmploymentType,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return model.Employee{}, err
		}
	}
	return employ, nil
}

func (e *EmployeeImp) UpdateEmployeeResign(employ model.Employee) error {
	query := "call hrm.update_employee_resign($1,$2,$3,$4,$5,$6,$7,$8)"
	tx, err := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(e.ctx)
		} else {
			tx.Commit(e.ctx)
		}
	}()
	_, err = tx.Exec(e.ctx, query,
		employ.ID,
		employ.ResignCode,
		employ.ResignRequest,
		employ.ResignDate,
		employ.ResignStatus,
		employ.ResignReason,
		employ.ResignApprovedBy,
		employ.ResignApprovedDate,
	)
	if err != nil {
		log.Printf("Error execute: %v", err)
		return err
	}
	return nil
}
