package domain

import "testing"

func TestNormalizeRequestDefaultsDeliveryWindow(t *testing.T) {
	order, err := NormalizeRequest(ManifestRequest{
		OrderID:  " order-1 ",
		Customer: " Ada ",
		Packages: []PackageRequest{{SKU: " BOOK ", Units: 2}},
	})
	if err != nil {
		t.Fatalf("NormalizeRequest() error = %v", err)
	}
	if order.DeliveryWindow == nil ||
		order.DeliveryWindow.Start != defaultWindowStart ||
		order.DeliveryWindow.End != defaultWindowEnd {
		t.Fatalf("NormalizeRequest() window = %#v", order.DeliveryWindow)
	}
	if order.Packages[0].SKU != "book" {
		t.Fatalf("NormalizeRequest() sku = %q", order.Packages[0].SKU)
	}
}

func TestNormalizeRequestPreservesSpecifiedDeliveryWindow(t *testing.T) {
	order, err := NormalizeRequest(ManifestRequest{
		OrderID:  "order-2",
		Customer: "Ada",
		Packages: []PackageRequest{{SKU: "book", Units: 1}},
		DeliveryWindow: &DeliveryWindow{
			Start: " 10:00 ",
			End:   " 12:00 ",
		},
	})
	if err != nil {
		t.Fatalf("NormalizeRequest() error = %v", err)
	}
	if order.DeliveryWindow == nil ||
		order.DeliveryWindow.Start != "10:00" ||
		order.DeliveryWindow.End != "12:00" {
		t.Fatalf("NormalizeRequest() window = %#v", order.DeliveryWindow)
	}
}
