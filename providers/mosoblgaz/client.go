package mosoblgaz

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/hasura/go-graphql-client"
)

const (
	BaseURL = "https://lkk.mosoblgaz.ru/graphql/"
)

type Client struct {
	token  string
	client *graphql.Client
}

func New() *Client {
	c := &Client{}
	c.client = graphql.NewClient(BaseURL, nil).WithRequestModifier(c.requestModifier)

	return c
}

func (c *Client) WithToken(token string) *Client {
	c.token = token

	return c
}

func (c *Client) WithDebug(flag bool) *Client {
	c.client = c.client.WithDebug(flag)

	return c
}

func (c *Client) requestModifier(r *http.Request) {
	if c.token != "" {
		r.Header.Set("token", c.token)
		r.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36")
	}
}

func (c *Client) Balance(ctx context.Context) (_ map[string]float64, err error) {
	var q struct {
		Me struct {
			Contracts []struct {
				Number      string
				LiveBalance struct {
					LiveBalance string
				}
			}
		}
	}

	if err = c.client.Query(ctx, &q, nil); err != nil {
		return nil, err
	}

	if len(q.Me.Contracts) == 0 {
		return nil, errors.New("contract not found")
	}

	balances := make(map[string]float64, len(q.Me.Contracts))
	for _, item := range q.Me.Contracts {
		balances[item.Number], err = strconv.ParseFloat(item.LiveBalance.LiveBalance, 64)

		if err != nil {
			return nil, err
		}
	}

	return balances, nil
}
