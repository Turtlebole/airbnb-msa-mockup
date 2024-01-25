package client

import (
	"backend/domain"
	"backend/models"
	"bytes"
	"context"
	"github.com/goccy/go-json"
	"github.com/sony/gobreaker"
	"log"
	"net/http"
	"time"
)

type ProfileClient struct {
	client  *http.Client
	address string
	cb      *gobreaker.CircuitBreaker
}

func NewProfileClient(client *http.Client, address string, cb *gobreaker.CircuitBreaker) ProfileClient {
	return ProfileClient{
		client:  client,
		address: address,
		cb:      cb,
	}

}

func (c *ProfileClient) SendUserData(ctx context.Context, user models.User) (interface{}, error) {
	req := convertUser(user)

	reqBytes, err := json.Marshal(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var timeout time.Duration
	deadline, reqHasDeadline := ctx.Deadline()
	if reqHasDeadline {
		timeout = time.Until(deadline)
	}

	cbResp, err := c.cb.Execute(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.address+"/profiles/create", bytes.NewBuffer(reqBytes))
		if err != nil {
			return nil, err
		}
		return c.client.Do(req)
	})
	if err != nil {
		return nil, handleHttpReqErr(err, c.address, http.MethodPost, timeout)
	}
	resp := cbResp.(*http.Response)
	if resp.StatusCode != http.StatusCreated {
		return nil, domain.ErrResp{
			URL:        resp.Request.URL.String(),
			Method:     resp.Request.Method,
			StatusCode: resp.StatusCode,
		}
	}

	return true, nil
}

func convertUser(user models.User) UserData {
	userData := UserData{
		First_name: user.First_name,
		Last_Name:  user.Last_name,
		Email:      user.Email,
		Address:    user.Address,
		Role:       user.User_type,
	}
	return userData
}
