package prodservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"route256/cart/config"
	"route256/cart/internal/cart/adapters/prodservice/clientbuilder"
	"route256/cart/internal/cart/adapters/prodservice/roundtripper"
	"route256/cart/internal/cart/models"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

type Client struct {
	baseURL string
	token   string
	client  *http.Client
}

func New(cfg config.ProductProviderCfg, log zerolog.Logger) *Client {
	log.Debug().Str("host", cfg.Address).Msg("creating new product service client")
	clientBuilder := clientbuilder.New(cfg.Timeout).
		Use(roundtripper.Retry(log, cfg.Retries, cfg.MaxDelayForRetry)).
		Use(roundtripper.Tracing()).
		Use(roundtripper.Limit(log, cfg.RateLimit, cfg.BurstLimit)).
		Use(roundtripper.Metrics())

	return &Client{
		baseURL: cfg.Address,
		token:   cfg.AccessToken,
		client:  clientBuilder.Build(),
	}
}

func (c *Client) GetProduct(ctx context.Context, sku int64) (models.ItemDescription, error) {
	localLog := *zerolog.Ctx(ctx)

	url := fmt.Sprintf("%s/get_product", c.baseURL)

	reqBody := getProductReq{
		Token: c.token,
		SKU:   sku,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		localLog.Error().Err(err).Send()
		return models.ItemDescription{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		localLog.Error().Err(err).Send()
		return models.ItemDescription{}, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		localLog.Error().Err(err).Send()
		return models.ItemDescription{}, err
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			localLog.Error().Err(err).Send()
		}
	}()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return models.ItemDescription{}, models.ErrNotFound
		}
		return models.ItemDescription{}, fmt.Errorf("failed to fetch item: status %d", resp.StatusCode)
	}

	return c.validateProductResp(localLog, resp.Body)
}

func (c *Client) validateProductResp(localLog zerolog.Logger, resp io.Reader) (models.ItemDescription, error) {
	body, err := io.ReadAll(resp)
	if err != nil {
		localLog.Error().Err(err).Send()
		return models.ItemDescription{}, err
	}

	var respBody getProductResp
	if err = json.Unmarshal(body, &respBody); err != nil {
		localLog.Error().Err(err).Send()
		return models.ItemDescription{}, err
	}

	if err = validator.New(validator.WithRequiredStructEnabled()).Struct(respBody); err != nil {
		localLog.Error().Err(err).Str("error", models.ErrBadProductInfo.Error()).Send()
		return models.ItemDescription{}, models.ErrBadProductData
	}

	return models.ItemDescription{
		Name:  respBody.Name,
		Price: respBody.Price,
	}, nil
}
