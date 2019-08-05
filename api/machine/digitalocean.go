package machine

import (
	"context"

	"github.com/alfg/openencoder/api/config"
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

// DropletListByTag lists the droplets for a digitalocean account.
func (do *DigitalOcean) DropletListByTag(ctx context.Context, tag string) ([]Machine, error) {
	// create a list to hold our droplets
	// list := []godo.Droplet{}
	list := []Machine{}

	// create options. initially, these will be blank
	opt := &godo.ListOptions{}
	for {
		droplets, resp, err := do.client.Droplets.ListByTag(ctx, tag, opt)
		if err != nil {
			return nil, err
		}

		// append the current page's droplets to our list
		for _, d := range droplets {
			list = append(list, Machine{
				ID:       d.ID,
				Name:     d.Name,
				Status:   d.Status,
				SizeSlug: d.SizeSlug,
				Created:  d.Created,
				Region:   d.Region.Name,
				Tags:     d.Tags,
				Provider: "digitalocean",
			})
		}

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

	return list, nil
}
