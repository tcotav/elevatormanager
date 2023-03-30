package building

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/tcotav/elevatormgr/elevator"
)

type Building struct {
	mu           sync.Mutex
	ID           int
	ElevatorList []*elevator.Elevator
	NumFloors    int
}

func NewBuilding(buildingID int, maxfloors int, numberElevators int) *Building {
	elevatorList := make([]*elevator.Elevator,0)
	for i := 0; i < numberElevators; i++ {
		// for simplicity sake, we use the count as elevatorID
		elevatorList = append(elevatorList, elevator.NewElevator(buildingID, i, maxfloors))
	}

	return &Building{
		ID:           buildingID,
		ElevatorList: elevatorList,
		NumFloors:    maxfloors,
	}
}

// maintenance function -- resets elevator to the ground floor and clears the call list
func (b *Building) ResetElevator(elevatorID int) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	e := b.GetElevator(elevatorID)
	if e == nil {
		return -1, fmt.Errorf("elevator with ID: %d does not exist in building: %d", elevatorID, b.ID)
	}
	if !e.InService {
		return -1, fmt.Errorf("elevator with ID: %d is not in service in building: %d", elevatorID, b.ID)
	}
	callsImpacted := e.CallList.Len()
	b.ElevatorList[elevatorID] = elevator.NewElevator(b.ID, elevatorID, b.NumFloors)

	// we want to return error because in a real system, something co
	// go wrong with resetting an elevator
	return callsImpacted, nil
}

func (b *Building) GetElevatorList() []*elevator.Elevator {
	return b.ElevatorList
}

func (b *Building) GetAllElevatorState() ([]byte, error) {
	retbytes, err := json.Marshal(b.ElevatorList)
	if err != nil {
		return []byte{}, err
	}
	return retbytes, nil
}

func (b *Building) SetElevatorInServiceStatus(elevatorID int, inService bool) error {
	e := b.GetElevator(elevatorID)
	if e == nil {
		return fmt.Errorf("elevator with ID: %d does not exist in building: %d", elevatorID, b.ID)
	}
	if e.InService == inService {
		// NOOP
		// do we message out that the elevator is already in the state requested?
		return nil
	}
	// set the inservice flag
	e.InService = inService
	// if it is being taken out of service, reset the elevator
	// clearing the call stack and bringing the elevator back to the ground floor
	if !inService {
		b.ResetElevator(elevatorID)
	}
	return nil
}

func (b *Building) CallElevator(floor int, direction int) (int, error) {
	if direction != 1 && direction != -1 {
		return -1, fmt.Errorf("invalid direction: %d in building: %d", direction, b.ID)
	}

	// we want to use the CLOSEST elevator to the floor
	var el *elevator.Elevator
	for _, e := range b.ElevatorList {
		if e.InService {
			if el == nil {
				el = e
			} else {
				if el.DistanceToFloor(floor) > e.DistanceToFloor(floor) {
					el = e
				}
			}
		}
	}

	// if we didn't find an elevator, return an error
	if el == nil {
		return -1, fmt.Errorf("no elevators in service in building: %d", b.ID)
	}
	// then do the actual call
	err := el.CallElevator(floor, direction)
	// and return which elevator we called
	return el.ElevatorID, err
}

func (b *Building) PushDestinationButton(elevatorID int, floor int) error {
	e := b.GetElevator(elevatorID)
	if e == nil {
		return fmt.Errorf("elevator with ID: %d does not exist in building: %d", elevatorID, b.ID)
	}
	if !e.InService {
		return fmt.Errorf("elevator with ID: %d is not in service in building: %d", elevatorID, b.ID)
	}
	return e.PushDestinationButton(floor)
}

func (b *Building) NextStop(elevatorID int) (*elevator.Call, error) {
	e := b.GetElevator(elevatorID)
	if e == nil {
		return nil, fmt.Errorf("elevator with ID: %d does not exist in building: %d", elevatorID, b.ID)
	}
	if !e.InService {
		return nil, fmt.Errorf("elevator with ID: %d is not in service in building: %d", elevatorID, b.ID)
	}
	return e.NextStop()
}

func (b *Building) MaintenanceCallOverride(elevatorID int, floor int, direction int) error {
	if direction != 1 && direction != -1 {
		return fmt.Errorf("invalid direction: %d for elevator: %d in building: %d", direction, elevatorID, b.ID)
	}
	e := b.GetElevator(elevatorID)
	if e == nil {
		return fmt.Errorf("elevator with ID: %d does not exist in building: %d", elevatorID, b.ID)
	}
	if !e.InService {
		return fmt.Errorf("elevator with ID: %d is not in service in building: %d", elevatorID, b.ID)
	}
	return e.ForceCallElevator(floor, direction)
}

// make sure elevator exists in our list
func (b *Building) doesElevatorExist(elevatorID int) bool {
	if len(b.ElevatorList) == 0 {
		return false
	}
	for _, e := range b.ElevatorList {
		if e.ElevatorID == elevatorID {
			return true
		}
	}
	return false
}

// GetElevator returns a pointer to the elevator object, checks if the elevator exists first
func (b *Building) GetElevator(elevatorID int) *elevator.Elevator {
	if b.doesElevatorExist(elevatorID) {
		return b.ElevatorList[elevatorID]
	}
	return nil
}
