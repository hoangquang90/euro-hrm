package hrm

import (
	attendanceimp "europm/internal/hrm/attendanceImp"
	"europm/internal/hrm/model"
	"europm/internal/util"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SearchAttendanceDevices(c *gin.Context) {
	status := c.Query("status")
	text := c.Query("text")

	attendanceimp := attendanceimp.GetInstance(c.Request.Context())
	lstAttendanceDevice, err := attendanceimp.GetAttendanceDevices(status, text)
	if err != nil {
		log.Printf("Error fetching HR training: %v", err)
		util.NewError(c, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, "get HR training error")
		return
	}

	c.JSON(http.StatusOK, lstAttendanceDevice)
}

func SearchAttendanceDevicesByID(c *gin.Context) {
	id := c.Param("id")

	attendanceimp := attendanceimp.GetInstance(c.Request.Context())
	attendanceDevice, err := attendanceimp.GetAttendanceDevicesById(id)
	if err != nil {
		log.Printf("Error fetching attendance device by ID: %v", err)
		util.NewError(c, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, "get attendance device by ID error")
		return
	}

	c.JSON(http.StatusOK, attendanceDevice)
}

func InsertAttendanceDevices(c *gin.Context) {
	var attendanceDevice model.AttendanceDevices
	if err := c.ShouldBindJSON(&attendanceDevice); err != nil {
		log.Printf("Error binding JSON: %v", err)
		util.NewError(c, http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, "invalid request body")
		return
	}

	attendanceimp := attendanceimp.GetInstance(c.Request.Context())
	id, err := attendanceimp.InstAttendanceDevices(attendanceDevice)
	if err != nil {
		log.Printf("Error inserting attendance device: %v", err)
		util.NewError(c, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, "insert attendance device error")
		return
	}

	c.JSON(http.StatusOK, id)
}

func DeleteAttendanceDevices(c *gin.Context) {
	id := c.Param("id")

	attendanceimp := attendanceimp.GetInstance(c.Request.Context())
	err := attendanceimp.DelAttendanceDevices(id)
	if err != nil {
		log.Printf("Error deleting attendance device: %v", err)
		util.NewError(c, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, "delete attendance device error")
		return
	}

	c.JSON(http.StatusOK, "Attendance device deleted successfully")
}
