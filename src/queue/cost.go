package queue

import (
	"driver"
	"fmt"
	"global"
)

func Delegate_order(new_order Order) Order {
	fmt.Println("New order is: ", new_order)
	assigned_elevator := compare_cost(Elevators_online, new_order)
	assigned_elevator_ip := assigned_elevator.Elev_ip
	fmt.Println("The assigned to ip is: ", assigned_elevator_ip)
	new_order.Assigned_to = assigned_elevator_ip
	fmt.Println("New order is assigned to: ", new_order.Assigned_to)
	return new_order
	//Add_new_global_order(new_order)
}

func compare_cost(online_elev_info_list [global.NUM_ELEV]Elev_info, new_order Order) Elev_info {
	fmt.Println("Num elev online: ", global.Num_elev_online)
	lowest_cost := 100
	var best_suited_elevator Elev_info

	for i := 0; i < global.Num_elev_online; i++ {
		cost := calculate_cost(new_order)
		fmt.Println("Online elev info i: ", online_elev_info_list[i])
		fmt.Println("Online elev ip is now: ", online_elev_info_list[i].Elev_ip)

		if cost == -2 {
			best_suited_elevator = online_elev_info_list[i]
			fmt.Println("The cost here is: ", cost)
			break
		} else if cost < lowest_cost {
			best_suited_elevator = online_elev_info_list[i]
			lowest_cost = cost
			fmt.Println("The cost here is: ", cost)
		}
	}

	fmt.Println("The best elev ip is: ", best_suited_elevator.Elev_ip)
	return best_suited_elevator
}

func calculate_cost(new_order Order) int {
	cost := 0

	if elevator_is_idle() {
		cost = -2
		return cost
	} else {
		cost += order_cost()
		cost += direction_cost(new_order.Floor)
		cost += floor_cost(new_order.Floor)
	}
	return cost
}

func elevator_is_idle() bool {
	for i := 0; i < global.NUM_INTERNAL_ORDERS; i++ {
		if Internal_order_list[i].Order_state != Inactive {
			return false
		}
	}
	for i := 0; i < global.NUM_GLOBAL_ORDERS; i++ {
		if External_order_list[i].Order_state != Inactive {
			return false
		}
	}
	return true
}

func order_cost() int {
	order_cost := 0
	for i := 0; i < global.NUM_INTERNAL_ORDERS; i++ {
		if Internal_order_list[i].Order_state != Inactive {
			order_cost += 2
		}
	}
	for i := 0; i < global.NUM_GLOBAL_ORDERS; i++ {
		if External_order_list[i].Order_state != Inactive {
			order_cost += 2
		}
	}
	return order_cost
}

// Add 3 points for wrong direction and -1 for right direction
func direction_cost(order_floor global.Floor_t) int {
	direction_cost := 0
	direction := My_info.Elev_dir

	switch direction {
	case global.DIR_DOWN:
		if order_floor < My_info.Elev_last_floor {
			//Elevator is going down, destination is lower than current floor
			direction_cost = -1
		} else {
			//Elevator going down, destination is higher than current floor
			direction_cost = 3
		}

	case global.DIR_UP:
		if order_floor > My_info.Elev_last_floor {
			//Elevator going up, destination is higher than current floor
			direction_cost = -1
		} else {
			//Elevator going up, destination is lower than current floor
			direction_cost = 3
		}
	}

	return direction_cost
}

// Add 2 points for each floor between the elevator and the order
func floor_cost(order_floor global.Floor_t) int {
	floor_cost := 0

	if My_info.Elev_last_floor < order_floor {
		floor_cost = 2 * (driver.Floor_t_to_floor_int(order_floor) - driver.Floor_t_to_floor_int(My_info.Elev_last_floor) - 1)
	} else {
		floor_cost = (-2) * (driver.Floor_t_to_floor_int(order_floor) - driver.Floor_t_to_floor_int(My_info.Elev_last_floor) + 1)
	}

	return floor_cost
}
