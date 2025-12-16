package employeeimp

import (
	"europm/internal/db/dbhrm"
	"europm/internal/hrm/model"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
)

func (e *EmployeeImp) SearchChangesEmployee(fromDate time.Time, toDate time.Time, typeReport string) ([]model.ChangesEmployee, error) {
	var changes []model.ChangesEmployee
	query := "select hrm.get_employee_report_monthly($1, $2, $3)"
	tx, err := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		if tx != nil {
			tx.Commit(e.ctx)
		}
	}()
	row := tx.QueryRow(e.ctx, query, fromDate, toDate, typeReport)
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
		employ := model.ChangesEmployee{}
		err := rows.Scan(&employ.TypeReport, &employ.CountLeave, &employ.CountJoin, &employ.CountEmployee, &employ.DepartmentName, &employ.WorkLocation)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		changes = append(changes, employ)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows in SearchChangesEmployee: %v", err)
		return nil, err
	}

	return changes, nil
}

func (e *EmployeeImp) SearchHRMWorkReport(fromDate time.Time, toDate time.Time, typeReport string) ([]model.Employee, error) {
	var employees []model.Employee
	query := "select hrm.get_hrm_work_report($1, $2, $3)"
	tx, err := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		if tx != nil {
			tx.Commit(e.ctx)
		}
	}()
	row := tx.QueryRow(e.ctx, query, fromDate, toDate, typeReport)
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
		employ := model.Employee{}
		err := rows.Scan(&employ.ID, &employ.DepartmentName, &employ.WorkLocation, &employ.Status, &employ.EmploymentType)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		employees = append(employees, employ)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows in SearchHRMWorkReport: %v", err)
		return nil, err
	}

	return employees, nil
}
func (e *EmployeeImp) SearchHRMResignReport(fromDate time.Time, toDate time.Time, typeReport string) ([]model.Employee, error) {
	var employees []model.Employee
	fmt.Println("SearchHRMResignReport called with fromDate:", fromDate, "toDate:", toDate, "typeReport:", typeReport)
	query := "select hrm.get_hr_resigns_report($1, $2, $3)"
	tx, err := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback(e.ctx)
		} else {
			tx.Commit(e.ctx)
		}
	}()
	row := tx.QueryRow(e.ctx, query, fromDate, toDate, typeReport)
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
		employ := model.Employee{}
		err := rows.Scan(&employ.ID, &employ.DepartmentName, &employ.WorkLocation, &employ.Status, &employ.EmploymentType, &employ.ResignDate)
		if err != nil {
			fmt.Println("Error scanning row: ", err)
			return nil, err
		}
		employees = append(employees, employ)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating rows in SearchHRMResignReport: ", err)
		return nil, err
	}
	return employees, nil
}
