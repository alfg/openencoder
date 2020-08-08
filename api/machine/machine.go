package machine

import "github.com/alfg/openencoder/api/logging"

var log = logging.Log

const (
	digitalOceanProviderName = "digitalocean"
	workerTagName            = "openencoder-worker"
	tagName                  = "openencoder"
	dockerImageName          = "docker-18-04"
)

// var (
// 	sizesLimiter = []string{"s-1vcpu-1gb", "s-1vcpu-2gb"}
// )

// Machine defines a machine struct from a provider.
type Machine struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Status   string   `json:"status"`
	SizeSlug string   `json:"size_slug"`
	Created  string   `json:"created_at"`
	Region   string   `json:"region"`
	Tags     []string `json:"tags"`

	Provider string `json:"provider"`
}

// CreatedResponse defines the response for creating a machine.
type CreatedResponse struct {
	ID       int    `json:"id"`
	Provider string `json:"provider"`
}

// DeletedResponse defines the response for deleted a machine.
type DeletedResponse struct {
	ID       int    `json:"id"`
	Provider string `json:"provider"`
}

// Region defines the response for listing regions.
type Region struct {
	Name      string   `json:"name"`
	Slug      string   `json:"slug"`
	Sizes     []string `json:"sizes"`
	Available bool     `json:"available"`
}

// Size defines the response for listing sizes.
type Size struct {
	Slug         string  `json:"slug"`
	Available    bool    `json:"available"`
	PriceMonthly float64 `json:"price_monthly"`
	PriceHourly  float64 `json:"price_hourly"`
}

// Pricing defines the response for listing pricing.
type Pricing struct {
	Count        int     `json:"count"`
	PriceHourly  float64 `json:"price_hourly"`
	PriceMonthly float64 `json:"price_monthly"`
}

// VPC defines the response for listing VPCs.
type VPC struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
