package domain

func NewManifest(order Order, reservations []Reservation) Manifest {
	return Manifest{
		OrderID:        order.ID,
		Customer:       order.Customer,
		DeliveryWindow: cloneWindow(order.DeliveryWindow),
		Reservations:   append([]Reservation(nil), reservations...),
		Status:         "planned",
	}
}

func cloneWindow(window *DeliveryWindow) *DeliveryWindow {
	if window == nil {
		return nil
	}
	return &DeliveryWindow{Start: window.Start, End: window.End}
}
