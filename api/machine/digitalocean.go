package machine

import (
	"context"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

// TokenSource defines an access token for oauth2.TokenSource.
type TokenSource struct {
	AccessToken string
}

// Token creates a token for TokenSource.
func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

// DigitalOcean client.
type DigitalOcean struct {
	client *godo.Client
}

// NewDigitalOceanClient creates a Digital Ocean client.
func NewDigitalOceanClient(token string) (*DigitalOcean, error) {
	tokenSource := &TokenSource{
		AccessToken: token,
	}

	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	client := godo.NewClient(oauthClient)

	return &DigitalOcean{
		client: client,
	}, nil
}

// ListDropletByTag lists the droplets for a digitalocean account.
func (do *DigitalOcean) ListDropletByTag(ctx context.Context, tag string) ([]Machine, error) {
	list := []Machine{}

	opt := &godo.ListOptions{}
	for {
		droplets, resp, err := do.client.Droplets.ListByTag(ctx, tag, opt)
		if err != nil {
			return nil, err
		}

		for _, d := range droplets {
			list = append(list, Machine{
				ID:       d.ID,
				Name:     d.Name,
				Status:   d.Status,
				SizeSlug: d.SizeSlug,
				Created:  d.Created,
				Region:   d.Region.Name,
				Tags:     d.Tags,
				Provider: digitalOceanProviderName,
			})
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}
		opt.Page = page + 1
	}
	return list, nil
}

// CreateDroplets creates a new DigitalOcean droplet.
func (do *DigitalOcean) CreateDroplets(ctx context.Context, region, size, vpc string, count int) ([]CreatedResponse, error) {

	var (
		ipv6              = true
		tags              = []string{tagName, workerTagName}
		monitoring        = true
		privateNetworking = true
	)

	var names []string
	for i := 0; i < count; i++ {
		names = append(names, workerTagName)
	}

	createRequest := &godo.DropletMultiCreateRequest{
		Names:  names,
		Region: region,
		Size:   size,
		Image: godo.DropletCreateImage{
			Slug: dockerImageName,
		},
		// SSHKeys: []godo.DropletCreateSSHKey{
		// 	godo.DropletCreateSSHKey{ID: 107149},
		// },
		UserData:          createUserData(),
		Tags:              tags,
		Monitoring:        monitoring,
		IPv6:              ipv6,
		PrivateNetworking: privateNetworking,
		VPCUUID:           vpc,
	}

	droplets, _, err := do.client.Droplets.CreateMultiple(ctx, createRequest)
	if err != nil {
		return nil, err
	}

	list := []CreatedResponse{}
	for _, d := range droplets {
		list = append(list, CreatedResponse{
			ID:       d.ID,
			Provider: digitalOceanProviderName,
		})
	}
	return list, nil
}

// DeleteDropletByID deletes a DigitalOcean droplet by ID.
func (do *DigitalOcean) DeleteDropletByID(ctx context.Context, ID int) (*DeletedResponse, error) {
	_, err := do.client.Droplets.Delete(ctx, ID)
	if err != nil {
		return nil, err
	}

	deleted := &DeletedResponse{
		ID:       ID,
		Provider: digitalOceanProviderName,
	}
	return deleted, nil
}

// DeleteDropletByTag deletes a DigitalOcean droplet by Tag.
func (do *DigitalOcean) DeleteDropletByTag(ctx context.Context, tag string) error {
	_, err := do.client.Droplets.DeleteByTag(ctx, tag)
	return err
}

// ListRegions gets a list of DigitalOcean regions.
func (do *DigitalOcean) ListRegions(ctx context.Context) ([]Region, error) {
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	regions, _, err := do.client.Regions.List(ctx, opt)
	if err != nil {
		return nil, err
	}

	list := []Region{}
	for _, d := range regions {
		list = append(list, Region{
			Name:      d.Name,
			Slug:      d.Slug,
			Sizes:     d.Sizes,
			Available: d.Available,
		})
	}

	return list, err
}

// ListSizes gets a list of DigitalOcean sizes.
func (do *DigitalOcean) ListSizes(ctx context.Context) ([]Size, error) {
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	sizes, _, err := do.client.Sizes.List(ctx, opt)
	if err != nil {
		return nil, err
	}

	list := []Size{}
	for _, d := range sizes {
		// if contains(sizesLimiter, d.Slug) {
		list = append(list, Size{
			Slug:         d.Slug,
			Available:    d.Available,
			PriceMonthly: d.PriceMonthly,
			PriceHourly:  d.PriceHourly,
		})
		// }
	}
	return list, err
}

// ListVPCs gets a list of DigitalOcean VPCs.
func (do *DigitalOcean) ListVPCs(ctx context.Context) ([]VPC, error) {
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	vpcs, _, err := do.client.VPCs.List(ctx, opt)
	if err != nil {
		return nil, err
	}

	list := []VPC{}
	for _, d := range vpcs {
		list = append(list, VPC{
			ID:   d.ID,
			Name: d.Name,
		})
	}
	return list, err
}

// GetCurrentPricing gets the current pricing data of running machines.
func (do *DigitalOcean) GetCurrentPricing(ctx context.Context, tag string) (*Pricing, error) {

	// Get sizes first.
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	sizes, _, err := do.client.Sizes.List(ctx, opt)
	if err != nil {
		return nil, err
	}

	sizeList := []Size{}
	for _, d := range sizes {
		sizeList = append(sizeList, Size{
			Slug:         d.Slug,
			Available:    d.Available,
			PriceMonthly: d.PriceMonthly,
			PriceHourly:  d.PriceHourly,
		})
	}

	// Get current running machines.
	machines := []Machine{}

	opt = &godo.ListOptions{}
	for {
		droplets, resp, err := do.client.Droplets.ListByTag(ctx, tag, opt)
		if err != nil {
			return nil, err
		}

		for _, d := range droplets {
			machines = append(machines, Machine{
				ID:       d.ID,
				Name:     d.Name,
				Status:   d.Status,
				SizeSlug: d.SizeSlug,
				Created:  d.Created,
				Region:   d.Region.Name,
				Tags:     d.Tags,
				Provider: digitalOceanProviderName,
			})
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}
		opt.Page = page + 1
	}

	// Calculate pricing.
	var running = len(machines)
	var priceHourly float64
	var priceMonthly float64

	for _, m := range machines {
		slug := m.SizeSlug
		for _, sl := range sizeList {
			if sl.Slug == slug {
				priceHourly += sl.PriceHourly
				priceMonthly += sl.PriceMonthly
			}
		}
	}

	p := &Pricing{
		Count:        running,
		PriceHourly:  priceHourly,
		PriceMonthly: priceMonthly,
	}

	return p, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
