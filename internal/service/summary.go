package service

import "github.com/1260124186-cc/fulfillment-manifest-cli/internal/domain"

func Summary(manifest domain.Manifest) string {
	return manifest.OrderID + ":" + manifest.Status
}
