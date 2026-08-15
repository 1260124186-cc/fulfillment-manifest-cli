package domain

import (
	"fmt"
	"strings"
)

const (
	defaultWindowStart = "09:00"
	defaultWindowEnd   = "18:00"
)

func NormalizeRequest(request ManifestRequest) (Order, error) {
	orderID := strings.TrimSpace(request.OrderID)
	customer := strings.TrimSpace(request.Customer)
	if orderID == "" {
		return Order{}, fmt.Errorf("order id is required")
	}
	if customer == "" {
		return Order{}, fmt.Errorf("customer is required")
	}
	if len(request.Packages) == 0 {
		return Order{}, fmt.Errorf("at least one package is required")
	}

	// 用全新切片承载规范化结果，避免复用调用方 request.Packages 的底层数组而改写原始数据
	packages := make([]PackageRequest, 0, len(request.Packages))
	for _, item := range request.Packages {
		item.SKU = strings.ToLower(strings.TrimSpace(item.SKU))
		if item.SKU == "" || item.Units <= 0 {
			return Order{}, fmt.Errorf("packages must have a sku and positive units")
		}
		packages = append(packages, item)
	}

	window := request.DeliveryWindow
	if window == nil {
		window = &DeliveryWindow{Start: defaultWindowStart, End: defaultWindowEnd}
	} else {
		window = &DeliveryWindow{
			Start: strings.TrimSpace(window.Start),
			End:   strings.TrimSpace(window.End),
		}
		if window.Start == "" || window.End == "" {
			return Order{}, fmt.Errorf("delivery window must include a start and end")
		}
	}

	return Order{
		ID:             orderID,
		Customer:       customer,
		Packages:       packages,
		DeliveryWindow: window,
	}, nil
}
