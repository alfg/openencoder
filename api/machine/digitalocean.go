package machine

import (
	"context"

	"github.com/alfg/openencoder/api/config"
	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

const providerName = "digitalocean"
const tagName = "openencoder"

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
		names = append(names, "openencoder-worker")
	}

	createRequest := &godo.DropletMultiCreateRequest{
		Names:  names,
		Region: "sfo2",
		Size:   "s-1vcpu-1gb",
		Image: godo.DropletCreateImage{
			Slug: "docker-18-04",
		},
		// SSHKeys: []godo.DropletCreateSSHKey{
		// 	godo.DropletCreateSSHKey{ID: 107149},
		// },
		IPv6: true,
		Tags: []string{tagName},
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
