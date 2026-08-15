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
	if order.DeliveryWindow == nil || order.DeliveryWindow.Start != defaultWindowStart {
		t.Fatalf("NormalizeRequest() window = %#v", order.DeliveryWindow)
	}
	if order.Packages[0].SKU != "book" {
		t.Fatalf("NormalizeRequest() sku = %q", order.Packages[0].SKU)
	}
}
