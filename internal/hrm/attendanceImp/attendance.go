package attendanceimp

import (
	"context"
	"europm/internal/db/dbhrm"
	"europm/internal/hrm/attendance"
	"europm/internal/hrm/model"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
)

type AttendanceImp struct {
	ctx context.Context
}

func GetInstance(ctx context.Context) attendance.AttendanceCode {
	return &AttendanceImp{
		ctx: ctx,
	}
}

func (a *AttendanceImp) GetAttendanceDevices(status string, text string) ([]model.AttendanceDevices, error) {
	var lstAttendanceDevices []model.AttendanceDevices
	query := "select hrm.get_attendance_devices($1, $2)"
	tx, err := dbhrm.Pool.BeginTx(a.ctx, pgx.TxOptions{})
	if err != nil {
		return []model.AttendanceDevices{}, err
	}

	defer func() {
		if err != nil {
			tx.Rollback(a.ctx)
		} else {
			tx.Commit(a.ctx)
		}
	}()
	row := tx.QueryRow(a.ctx, query, status, text)
	var cursor string
	err = row.Scan(&cursor)
	if err != nil {
		log.Printf("Error execute", err)
		return []model.AttendanceDevices{}, err
	}
	rows, err := tx.Query(a.ctx, "FETCH ALL "+pq.QuoteIdentifier(cursor)+";")
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		log.Printf("Error FETCH refCurser: %v", err)
	}
	for rows.Next() {
		var at model.AttendanceDevices
		err := rows.Scan(&at.Id,
			&at.DevicesId,
			&at.DeviceName,
			&at.Location,
			&at.Serial,
			&at.IP,
			&at.Port,
			&at.Status)
		lstAttendanceDevices = append(lstAttendanceDevices, at)
		if err != nil {
			fmt.Println("Error scanning row: ", err)
			return []model.AttendanceDevices{}, err
		}
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating rows in GetAttendanceDevices: ", err)
		return []model.AttendanceDevices{}, err
	}
	return lstAttendanceDevices, nil
}

func (a *AttendanceImp) GetAttendanceDevicesById(id string) (model.AttendanceDevices, error) {
	var attendanceDevice model.AttendanceDevices
	query := "select hrm.get_attendance_devices_by_id($1)"
	tx, err := dbhrm.Pool.BeginTx(a.ctx, pgx.TxOptions{})
	if err != nil {
		return model.AttendanceDevices{}, err
	}

	defer func() {
		if err != nil {
			tx.Rollback(a.ctx)
		} else {
			tx.Commit(a.ctx)
		}
	}()
	row := tx.QueryRow(a.ctx, query, id)
	var cursor string
	err = row.Scan(&cursor)
	if err != nil {
		log.Printf("Error execute", err)
		return model.AttendanceDevices{}, err
	}
	rows, err := tx.Query(a.ctx, "FETCH ALL "+pq.QuoteIdentifier(cursor)+";")
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		log.Printf("Error FETCH refCurser: %v", err)
	}
	for rows.Next() {
		err := rows.Scan(&attendanceDevice.Id, &attendanceDevice.DevicesId, &attendanceDevice.DeviceName, &attendanceDevice.Location, &attendanceDevice.Serial, &attendanceDevice.IP, &attendanceDevice.Port, &attendanceDevice.Status)
		if err != nil {
			fmt.Println("Error scanning row: ", err)
			return model.AttendanceDevices{}, err
		}
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating rows in SearchHRMResignReport: ", err)
		return model.AttendanceDevices{}, err
	}
	return attendanceDevice, nil
}

func (a *AttendanceImp) InstAttendanceDevices(att model.AttendanceDevices) (string, error) {
	var id string
	query := "select hrm.insert_attendance_devices($1, $2, $3, $4, $5, $6, $7, $8)"
	args := []interface{}{
		att.Id,
		att.DevicesId,
		att.DeviceName,
		att.Location,
		att.Serial,
		att.IP,
		att.Port,
		att.Status,
	}
	tx, err := dbhrm.Pool.BeginTx(a.ctx, pgx.TxOptions{})
	if err != nil {
		return "", err
	}
	defer func() {
		if tx != nil {
			tx.Commit(a.ctx)
		}
	}()
	row := tx.QueryRow(a.ctx, query, args...)
	if err := row.Scan(&id); err != nil {
		log.Printf("Error execute: %v", err)
		return "", err
	}
	return id, nil
}

func (a *AttendanceImp) DelAttendanceDevices(id string) error {
	query := "select hrm.delete_attendance_devices($1)"
	tx, err := dbhrm.Pool.BeginTx(a.ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if tx != nil {
			tx.Commit(a.ctx)
		}
	}()
	_, err = tx.Exec(a.ctx, query, id)
	if err != nil {
		log.Printf("Error execute: %v", err)
		return err
	}
	return nil
}
