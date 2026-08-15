package domain

import "fmt"

func Allocate(order Order, stock Stock) ([]Reservation, error) {
	reservations := make([]Reservation, 0, len(order.Packages))
	for _, item := range order.Packages {
		available := stock[item.SKU]
		if available < item.Units {
			return nil, fmt.Errorf("insufficient stock for %s: have %d, need %d", item.SKU, available, item.Units)
		}
		reservations = append(reservations, Reservation{SKU: item.SKU, Units: item.Units})
	}
	return reservations, nil
}
