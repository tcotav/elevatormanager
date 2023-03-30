package elevator

import (
	"testing"
)

func TestElevator(t *testing.T) {
	elevator := NewElevator(1,1,10)
	if elevator.MaxFloor != 10 {
		t.Errorf("Elevator should have a max floor of 10")
	}
	clLen := elevator.CallList.Len()
	if clLen != 0 {
		t.Errorf("Elevator should have an empty call list, got %d", clLen)
	}
}

func TestElevatorTrip(t *testing.T) {
	elevator := NewElevator(1,1,5)

	// adds one to the list
	err := elevator.CallElevator(3, 1)
	if err != nil {
		t.Errorf("Elevator should be called, got %s", err.Error())
	}
	clLen := elevator.CallList.Len()
	if clLen != 1 {
		t.Errorf("Elevator should have 1 call in the call list, got %d", clLen)
	}

	// adds a second one to the call list
	elevator.PushDestinationButton(5)
	if elevator.CallList.Len() != 2 {
		t.Errorf("Elevator should have 2 call in the call list, %d", elevator.CallList.Len())
	}

	call, err := elevator.NextStop()
	if err != nil {
		t.Errorf("Elevator should have a next stop")
	}
	if call.Floor != 3 {
		t.Errorf("Elevator should be going to floor 3, got %d", call.Floor)
	}
	if elevator.CallList.Len() != 1 {
		t.Errorf("Elevator should have 1 call in the call list, %d", elevator.CallList.Len())
	}

	call, err = elevator.NextStop()
	if err != nil {
		t.Errorf("Elevator should have a next stop")
	}
	if call.Floor != 5 {
		t.Errorf("Elevator should be going to floor 5, got %d", call.Floor)
	}

	call, err = elevator.NextStop()
	if err == nil {
		t.Errorf("Elevator should not have a next stop, got %d", call.Floor)
	}
}