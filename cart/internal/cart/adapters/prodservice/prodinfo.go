package prodservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"route256/cart/config"
	"route256/cart/internal/cart/models"

	"github.com/rs/zerolog"
)

type Client struct {
	baseURL string
	token   string
	client  *http.Client
	log     zerolog.Logger
}

func New(cfg config.ProductProviderCfg, log zerolog.Logger) *Client {
	return &Client{
		baseURL: cfg.Address,
		token:   cfg.AccessToken,
		client:  &http.Client{},
		log:     log,
	}
}

func (c *Client) GetProduct(ctx context.Context, sku int64) (models.ItemDescription, error) {
	url := fmt.Sprintf("%s/get_product", c.baseURL)

	reqBody := getProductReq{
		Token: c.token,
		SKU:   sku,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return models.ItemDescription{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return models.ItemDescription{}, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return models.ItemDescription{}, err
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			c.log.Error().Err(err).Send()
		}
	}()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return models.ItemDescription{}, models.ErrNotFound
		}
		return models.ItemDescription{}, fmt.Errorf("failed to fetch item: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.ItemDescription{}, err
	}

	var respBody getProductResp
	if err = json.Unmarshal(body, &respBody); err != nil {
		return models.ItemDescription{}, err
	}

	return models.ItemDescription{
		Name:  respBody.Name,
		Price: respBody.Price,
	}, nil
}