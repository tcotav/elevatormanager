package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/tcotav/elevatormgr/building"
)

// mock this up for now
var bld *building.Building = building.NewBuilding(1, 10, 3) 

// specific error logging and sets the HTTP response + 400 code
func handleBadRequest(c *gin.Context, source string, err error) {
	log.Error(fmt.Sprintf("%s - %s", source, err.Error()))
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

// in an elevator car, push the button to go to a floor
func PushDestination(c *gin.Context) {
	errloc := "pushdest"	
	elevatorID, err := strconv.Atoi(c.Param("elevator"))
	if err != nil {
		handleBadRequest(c, errloc, err)
		return
	}

	floor, err := strconv.Atoi(c.Param("floor"))
	if err != nil {
		handleBadRequest(c, errloc, err)
		return
	}
	err = bld.PushDestinationButton(elevatorID, floor)
	if err != nil {
		handleBadRequest(c, errloc, err)
		return
	}
	log.Info(fmt.Sprintf("Elevator %d was called to floor %d.", elevatorID, floor))
}

// push the call button, up or down, on a floor
func CallElevator(c *gin.Context) {
	errloc := "callelev"
	floor, err := strconv.Atoi(c.Param("floor"))
	if err != nil {
		handleBadRequest(c, errloc, err)
		return
	}
	direction, err := strconv.Atoi(c.Param("direction"))
	if err != nil {
		handleBadRequest(c, errloc, err)
		return
	}
	elevatorID, err := bld.CallElevator(floor, direction)
	if err != nil {
		handleBadRequest(c, errloc, err)
		return
	}
	log.Info(fmt.Sprintf("Elevator %d was called to floor %d in direction %d.", elevatorID, floor, direction))
	b, err := json.Marshal(map[string]int{"elevator": elevatorID})
	if err != nil {
		handleBadRequest(c, "callelev", err)
		return
	}
	c.Data(http.StatusOK, "application/json", b)
}

// get all elevators' state
func GetAllElevatorState(c *gin.Context) {
	state, err := bld.GetAllElevatorState()
	if err != nil {
		handleBadRequest(c,"getallstate", err)
		return
	}
	c.Data(http.StatusOK, "application/json", state)
}

// reset the elevator -- i.e. call it down to floor 1 and clear its call list
func ResetElevator(c *gin.Context) {
	errloc := "resetelev"
	elevatorID, err := strconv.Atoi(c.Param("elevator"))
	if err != nil {
		handleBadRequest(c, errloc, err)
		return
	}
	callsImpacted, err := bld.ResetElevator(elevatorID)
	if err != nil {
		handleBadRequest(c, errloc, err)
		return
	}
	log.Info(fmt.Sprintf("Elevator %d was reset. %d calls were impacted.", elevatorID, callsImpacted))

}

// have gin log in json format
func jsonLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(
		func(params gin.LogFormatterParams) string {
			logMap := make(map[string]interface{})

			logMap["status_code"] = params.StatusCode
			logMap["method"] = params.Method
			logMap["path"] = params.Path
			logMap["start_time"] = params.TimeStamp.Format("2006/01/02 - 15:04:05")
			logMap["remote_addr"] = params.ClientIP
			logMap["response_time"] = params.Latency.String()
			s, _ := json.Marshal(logMap)
			return string(s) + "\n"
		},
	)
}

// take elevator out of service
func ElevatorOutOfService(c *gin.Context) {
	errloc := "outofservice"
	elevatorID, err := strconv.Atoi(c.Param("elevator"))
	if err != nil {
		handleBadRequest(c, errloc, err)
		return
	}
	err = bld.SetElevatorInServiceStatus(elevatorID, false)
	if err != nil {
		handleBadRequest(c, errloc, err)
		return
	}
	log.Info(fmt.Sprintf("Elevator %d was taken out of service.", elevatorID))
}


// put elevator back in service
func ElevatorBackInService(c *gin.Context) {
	errloc := "backinservice"
	elevatorID, err := strconv.Atoi(c.Param("elevator"))
	if err != nil {
		handleBadRequest(c, errloc, err)
		return
	}
	err = bld.SetElevatorInServiceStatus(elevatorID, true)
	if err != nil {
		handleBadRequest(c, errloc, err)
		return
	}
	log.Info(fmt.Sprintf("Elevator %d was put back in service.", elevatorID))
}

// call elevator to a specific floor and prioritize the call
func MaintenanceCallOverride(c *gin.Context) {
	errloc := "maintoverride"
	elevatorID, err := strconv.Atoi(c.Param("elevator"))
	if err != nil {
		handleBadRequest(c, errloc, err)
		return
	}
	floor, err := strconv.Atoi(c.Param("floor"))
	if err != nil {
		handleBadRequest(c, errloc, err)
		return
	}
	direction, err := strconv.Atoi(c.Param("direction"))
	if err != nil {
		handleBadRequest(c, errloc, err)
		return
	}
	err = bld.MaintenanceCallOverride(elevatorID, floor, direction)
	if err != nil {
		handleBadRequest(c, errloc, err)
		return
	}
	log.Info(fmt.Sprintf("Maintenance override: elevator %d was called to floor %d", elevatorID, floor))
}

// set up the web routes
// and do any other config for gin here (e.g. logging)
func setupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(jsonLogger())

	// the two user-facing routes
	router.POST("/pushDestination/:elevator/:floor", PushDestination)
	router.POST("/callElevator/:floor/:direction", CallElevator)

	// this one is used by both maint and users to see the state
	// I'd tidy it up to share it with users
	router.GET("/getAllElevatorState", GetAllElevatorState)

	// maintenance routes
	router.POST("/maintenanceCallOverride/:elevator/:floor/:direction", MaintenanceCallOverride)
	router.POST("/resetElevator/:elevator", ResetElevator)
	router.POST("/takeElevatorOutOfService/:elevator", ElevatorOutOfService)
	router.POST("/elevatorBackInService/:elevator", ElevatorBackInService)

	return router
}

func main() {
	r := setupRouter()
	log.Info("Starting server on port 8077")
	r.Run(":8077")
}
