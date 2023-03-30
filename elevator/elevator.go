package elevator

import (
	"encoding/json"
	"fmt"
)

type Elevator struct {
	BuildingID int
	ElevatorID int
	CurrentFloor int
	Direction    int
	MaxFloor     int
	InService   bool
	CallList	 *ElevatorCallList	
}

func NewElevator(buildingID int, elevatorID int, maxfloor int) *Elevator {
	callList := NewElevatorCallList()

	return &Elevator{
		ElevatorID: elevatorID,
		CurrentFloor: 1,
		MaxFloor:     maxfloor,
		CallList: callList,
		BuildingID: buildingID,
		InService: true,
	}
}

// these two do the same thing -- move the elevator to a specific floor -- but have different input sources
// which isn't visible at this level
// one is a pull for the elevator car
// the other is a push while IN the elevator car

// PushDestinationButton is called when a user pushes a floor button in the elevator car
func (e *Elevator) PushDestinationButton(floor int) error {
	var direction int
	if floor == e.CurrentFloor || floor > e.MaxFloor{
		return fmt.Errorf("invalid floor: %d for elevator: %d in building: %d", floor, e.ElevatorID, e.BuildingID)
	} 
	
	if floor > e.CurrentFloor {
		direction = 1
	} else {
		direction = -1
	}
	return e.addCall(floor, direction)
}

// CallElevator is called when a user calls an elevator from a floor
// equivalent to pushing the up or down arrow at your floor to summon the elevator
func (e *Elevator) CallElevator(floor int, direction int) error {
	// if at same floor -- we open the door but don't move elevator
	if floor > e.MaxFloor{
		return fmt.Errorf("invalid floor: %d for elevator: %d in building: %d", floor, e.ElevatorID, e.BuildingID)
	}

	// this is a hack -- instead we'd track our direction by which way the car is going for the current call
	// instead we set the direction that we'll be heading to be whatever this FIRST call is
	if e.CallList.Len() == 0 {
		e.Direction = direction
	}
	return e.addCall(floor, direction)
}

// ForceCallElevator is called when a user overrides the existing call stack
// and causes the elevator to go to a specific floor immediately
func (e *Elevator) ForceCallElevator(floor int, direction int) error {
	if floor == e.CurrentFloor || floor > e.MaxFloor{
		return fmt.Errorf("invalid floor: %d for elevator: %d in building: %d on floor: %d", floor, e.ElevatorID, e.BuildingID, e.CurrentFloor)
	}
	if e.CallList.Len() == 0 {
		e.Direction = direction
	}
	call := Call{
		Floor:     floor,
		Direction: direction,
	}
	e.CallList.Prepend(call)
	return nil
}


// GetState returns the current state of the elevator including the call list
func (e Elevator) GetState() ([]byte, error) {
	b, err := json.Marshal(e)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

func (e Elevator) GetCallList() *ElevatorCallList {
	return e.CallList
}

// this is kind of silly -- it just gets the next item in the call list
func (e *Elevator) NextStop() (*Call, error) {
	if e.CallList.Len() == 0 {
		return nil, fmt.Errorf("no calls in call list for elevator: %d in building: %d", e.ElevatorID, e.BuildingID)
	}
	call := e.CallList.PopLeft()
	if call == nil {
		return nil, fmt.Errorf("no calls in call list for elevator: %d in building: %d", e.ElevatorID, e.BuildingID)
	}
	e.CurrentFloor = call.Floor
	e.Direction = call.Direction
	return call, nil
}

// addCall adds a call to the elevator call list - utility method
func (e *Elevator) addCall(floor int, direction int) error {
	call := Call{
		Floor:     floor,
		Direction: direction,
	}
    return e.CallList.Push(call)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (e Elevator) DistanceToFloor(floor int) int {
	return abs(e.CurrentFloor - floor)
}