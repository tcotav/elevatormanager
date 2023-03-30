package building

import (
	"testing"
	"encoding/json"
	"github.com/tcotav/elevatormgr/elevator"
)

func TestBuilding(t *testing.T) {
	b := NewBuilding(1, 10, 2)

	if len(b.GetElevatorList()) != 2 {
		t.Errorf("Building should have 2 elevators")
	}
	if b.GetElevatorList()[0].MaxFloor != 10 {
		t.Errorf("Elevator should have a max floor of 10")
	}
}

func TestBuildingCallElevator(t *testing.T) {
	b := NewBuilding(1, 10, 2)
	// call elevator to a certain floor
	elID, err := b.CallElevator(3, 1)

	if err != nil {
		t.Errorf("Elevator should be called, got %s", err.Error())
	}

	e := b.GetElevator(elID)
	if e.GetCallList().Len() != 1 {
		t.Errorf("Elevator should have a call list of 1")
	}
}

func TestBuildingElevatorReset(t *testing.T) {
	b := NewBuilding(1, 10, 2)

	elID, err := b.CallElevator(3, 1)
	if err != nil {
		t.Errorf("Elevator should be called, got %s", err.Error())
	}
	b.PushDestinationButton(elID, 5)
	b.ResetElevator(elID)

	e := b.GetElevator(elID)
	if e.GetCallList().Len() != 0 {
		t.Errorf("Elevator should have an empty call list")
	}
}

func TestBuildingElevatorTrip(t *testing.T) {
	b := NewBuilding(1, 10, 2)

	elID, err := b.CallElevator(3, 1)
	if err != nil {
		t.Errorf("Elevator should be called")
	}
	b.PushDestinationButton(elID, 5)

	call, err := b.NextStop(elID)
	if err != nil {
		t.Errorf("Elevator should have a next stop, got %s", err.Error())
	}
	if call.Floor != 3 {
		t.Errorf("Elevator should be going to floor 3, got %d", call.Floor)
	}
	call, err = b.NextStop(elID)
	if err != nil {
		t.Errorf("Elevator should have a next stop, %s", err.Error())
	}
	if call.Floor != 5 {
		t.Errorf("Elevator should be going to floor 5, got %d", call.Floor)
	}
}

func TestGetAllElevatorState(t *testing.T) {
	b := NewBuilding(1, 10, 2)
	elID, err := b.CallElevator(3, 1)
	if err != nil {
		t.Errorf("Elevator should be called")
	}
	b.PushDestinationButton(elID, 5)
	// call stack for elID should be len 2 now

	statJSONb, err := b.GetAllElevatorState()
	if err != nil {
		t.Errorf("Elevator should be called")
	}
	var elList []*elevator.Elevator
	err = json.Unmarshal(statJSONb, &elList)
	if err != nil {
		t.Errorf("JSON unmarshal error, %s", err.Error())
	}
	if len(elList) != 2 {
		t.Errorf("Elevator list should have 2 elevators, got %d", len(elList))
	}
	// now check the call list that we added to above
	var el *elevator.Elevator
	for _, e := range elList {
		if e.ElevatorID == elID {
			el = e 
		}
	}
	if el == nil {
		t.Errorf("Elevator %d should be in list", elID)
	} else {
		callList := el.GetCallList()
		if callList == nil {
			t.Errorf("Elevator should have a call list")
		}

		if callList.Len() != 2 {
			t.Errorf("Elevator should have a call list of 2, instead found %d", callList.Len())
		}

	}
}