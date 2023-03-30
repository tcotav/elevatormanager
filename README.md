Using Go or Python (feel free to pick another language if you are not strong in either of those) please build a Web API that can act as the elevator control unit for a building with 5 floors and 2 elevators. Users need to be able to call an elevator from any floor, view the current state of all elevators, and set a destination floor for a specific elevator. Think about how the buildings maintenance team (SRE) would monitor and perform maintenance on this service while it is still running.


Abstractions:

We are the building managers.  We will have the base interface of BuildingManager which will have a list of Buildings that we manage.

We manage one or more Buildings.

In each of those buildings there are:
    - one or more floors
    - one or more elevators

Users can call an elevator from any floor, which building will be implicit in the button push (i.e. the button will have the building ID)

Users can request the elevator travel to any floor, and again which building will be implicit in the button push.

Currently we don't have any actual movement.  For this exercise, we've just got a building, elevators, and an ever-growing call stack.  Nothing is working on the call stack and moving the elevator cars.

TODO:

The destination rule we're going to use is that the elevator will continue in one direction until all requests are satisfied or they arrive at floor 1.  If there are no requests in the current direction, the elevator will reverse direction and service the remaining requests.  If there are no requests, the elevator will idle at the current floor.

Cases:

Push call button and elevator is already on that floor -- elevator will ignore, BUT really we should be messaging back "ok, I'm here, let's open the door"