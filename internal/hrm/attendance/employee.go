package attendance

import (
	"europm/internal/hrm/model"
)

type AttendanceCode interface {
	//employeeImp.go
	GetAttendanceDevices(status string, text string) ([]model.AttendanceDevices, error)
	GetAttendanceDevicesById(id string) (model.AttendanceDevices, error)
	InstAttendanceDevices(att model.AttendanceDevices) (string, error)
	DelAttendanceDevices(id string) error
}
