package domain

type PackageRequest struct {
	SKU   string `json:"sku"`
	Units int    `json:"units"`
}

type DeliveryWindow struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type ManifestRequest struct {
	OrderID        string           `json:"order_id"`
	Customer       string           `json:"customer"`
	Packages       []PackageRequest `json:"packages"`
	DeliveryWindow *DeliveryWindow  `json:"delivery_window,omitempty"`
}

type Order struct {
	ID             string
	Customer       string
	Packages       []PackageRequest
	DeliveryWindow *DeliveryWindow
}

type Reservation struct {
	SKU   string `json:"sku"`
	Units int    `json:"units"`
}

type Manifest struct {
	OrderID        string          `json:"order_id"`
	Customer       string          `json:"customer"`
	DeliveryWindow *DeliveryWindow `json:"delivery_window,omitempty"`
	Reservations   []Reservation   `json:"reservations"`
	Status         string          `json:"status"`
}

type Stock map[string]int
