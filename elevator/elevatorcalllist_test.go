package elevator

import (
	"testing"
)

func TestElevatorCallList(t *testing.T){
	elevatorCallList := NewElevatorCallList()
	elevatorCallList.Push(Call{Floor: 3, Direction: 1})
	elevatorCallList.Push(Call{Floor: 5, Direction: 1})
	elevatorCallList.Push(Call{Floor: 1, Direction: -1})
	elevatorCallList.Push(Call{Floor: 2, Direction: 1})
	elevatorCallList.Push(Call{Floor: 4, Direction: -1})

	clLen := elevatorCallList.Len()
	if clLen != 5 {
		t.Errorf("Elevator call list should have 5 calls, got %d", clLen)
	}
	call := elevatorCallList.PopLeft()
	if call.Floor != 3 {
		t.Errorf("Elevator call list should have 3 as the floor value, got %d", call.Floor)
	}
	if call.Direction != 1 {
		t.Errorf("Elevator call list should have 1 as the direction value, got %d", call.Direction)
	}
	if elevatorCallList.Len() != 4 {
		t.Errorf("Elevator call list should have 4 calls, got %d", elevatorCallList.Len())
	}
	// try to add a duplicate
	elevatorCallList.Push(Call{Floor: 5, Direction: 1})
	if elevatorCallList.Len() != 4 {
		t.Errorf("Elevator call list should have 4 calls after dupe push, got %d", elevatorCallList.Len())
	}
	elevatorCallList.Push(Call{Floor: 5, Direction: -1})
	if elevatorCallList.Len() != 5 {
		t.Errorf("Elevator call list should have 5 calls, got %d", elevatorCallList.Len())
	}

}

func TestElevatorCallListPrepend(t *testing.T){
	elevatorCallList := NewElevatorCallList()
	elevatorCallList.Push(Call{Floor: 3, Direction: 1})
	elevatorCallList.Push(Call{Floor: 5, Direction: 1})

	elevatorCallList.Prepend(Call{Floor: 1, Direction: -1})
	if elevatorCallList.Len() != 3 {
		t.Errorf("Elevator call list should have 3 calls, got %d", elevatorCallList.Len())
	}
	call := elevatorCallList.PopLeft()
	if call.Floor != 1 {
		t.Errorf("Elevator call list should have 1 as the floor value, got %d", call.Floor)
	}
	if call.Direction != -1 {
		t.Errorf("Elevator call list should have -1 as the direction value, got %d", call.Direction)
	}
}


