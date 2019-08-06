package machine

import (
	"context"

	"github.com/alfg/openencoder/api/config"
	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

const (
	providerName    = "digitalocean"
	workerName      = "openencoder-worker"
	tagName         = "openencoder"
	dockerImageName = "docker-18-04"
)

var (
	sizesLimiter = []string{"s-1vcpu-1gb", "s-1vcpu-2gb"}
)

const userData = `
#cloud-config
package_upgrade: true
write-files:
    - path: "/etc/profile.env"
      content: |
        export MY_VAR="foo"
runcmd:
  - touch /test.txt
`

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
func NewDigitalOceanClient() (*DigitalOcean, error) {
	tokenSource := &TokenSource{
		AccessToken: config.Get().DigitalOceanAccessToken,
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
				Provider: providerName,
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
func (do *DigitalOcean) CreateDroplets(ctx context.Context, count int) ([]MachineCreated, error) {

	var names []string
	for i := 0; i < count; i++ {
		names = append(names, workerName)
	}

	createRequest := &godo.DropletMultiCreateRequest{
		Names:  names,
		Region: "sfo2",
		Size:   "s-1vcpu-1gb",
		Image: godo.DropletCreateImage{
			Slug: dockerImageName,
		},
		// SSHKeys: []godo.DropletCreateSSHKey{
		// 	godo.DropletCreateSSHKey{ID: 107149},
		// },
		IPv6:       true,
		Tags:       []string{tagName},
		Monitoring: true,
		UserData:   userData,
	}

	droplets, _, err := do.client.Droplets.CreateMultiple(ctx, createRequest)
	if err != nil {
		return nil, err
	}

	list := []MachineCreated{}
	for _, d := range droplets {
		list = append(list, MachineCreated{
			ID:       d.ID,
			Provider: providerName,
		})
	}
	return list, nil
}

// DeleteDropletByID deletes a DigitalOcean droplet by ID.
func (do *DigitalOcean) DeleteDropletByID(ctx context.Context, ID int) (*MachineDeleted, error) {
	_, err := do.client.Droplets.Delete(ctx, ID)
	if err != nil {
		return nil, err
	}

	deleted := &MachineDeleted{
		ID:       ID,
		Provider: providerName,
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
		if contains(sizesLimiter, d.Slug) {
			list = append(list, Size{
				Slug:         d.Slug,
				Available:    d.Available,
				PriceMonthly: d.PriceMonthly,
				PriceHourly:  d.PriceHourly,
			})
		}
	}

	return list, err
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
