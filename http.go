package main

import (
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"
)

type Client struct {
	*resty.Client
}

type Backoff struct {
	minDelay time.Duration
	maxDelay time.Duration
}

func (b *Backoff) next(attempt int) time.Duration {
	if attempt < 0 {
		return b.minDelay
	}

	minf := float64(b.minDelay)
	durf := minf * math.Pow(1.5, float64(attempt))
	durf = durf + rand.Float64()*minf

	delay := time.Duration(durf)
	if delay > b.maxDelay {
		return b.maxDelay
	}

	return delay
}

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Res     string      `json:"res"`
	ErrData string      `json:"errdata"`
}

func NewClient() *Client {
	client := resty.New()
	client.SetTimeout(time.Second * 30)
	client.SetHeaders(map[string]string{
		"content-type": "multipart/form-data; boundary=----WebKitFormBoundaryBoisGGEqBQMlOG7a",
		"user_agent":   "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36",
	})
	return &Client{
		Client: client,
	}
}

func (c *Client) SetHeader(key, value string) {
	c.Client.SetHeader(key, value)
}

func (c *Client) Get(url string, out interface{}) error {
	resp, err := c.R().Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return errors.New("response err")
	}

	if err := decodeResponse(resp.Body(), out); err != nil {
		return err
	}

	return nil
}

func (c *Client) PostForm(url string, params map[string]string, out interface{}) error {
	resp, err := c.R().SetQueryParams(params).Post(url)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return errors.New("response err")
	}

	//fmt.Println(string(resp.Body()))

	if err := decodeResponse(resp.Body(), out); err != nil {
		return err
	}

	return nil
}

func decodeResponse(in []byte, out interface{}) error {
	var ret Response
	if err := json.Unmarshal(in, &ret); err != nil {
		return err
	}

	//if ret.Code != "SUCCESS" {
	//	log.Println(ret.Message)
	//}
	if ret.Res != "succ" {
		log.Println(ret.ErrData)
		return errors.New(ret.ErrData)
	}

	data, err := json.Marshal(ret.Data)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &out); err != nil {
		return err
	}

	return nil
}
