package employeeimp

import (
	"europm/internal/db/dbhrm"
	"europm/internal/hrm/model"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
)

func (e *EmployeeImp) SearchChangesEmployee(filter *model.ChangesEmployeeFilter) ([]model.ChangesEmployee, error) {
	var changes []model.ChangesEmployee
	query := `SELECT hrm.get_employee_report_monthly($1, $2, $3);`
	var cursor string
	tx, _ := dbhrm.Pool.BeginTx(e.ctx, pgx.TxOptions{})
	defer func() {
		if tx != nil {
			tx.Commit(e.ctx)
		}
	}()

	row := tx.QueryRow(e.ctx, query, filter.FromDate, filter.ToDate, filter.TypeReport)
	err := row.Scan(&cursor)
	if err != nil {
		log.Printf("Error execute", err)
		return changes, err
	}
	rows, err := tx.Query(e.ctx, "FETCH ALL "+pq.QuoteIdentifier(cursor)+";")
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		log.Printf("Error FETCH refCurser: %v", err)
		return changes, err
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
