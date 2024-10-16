package rs

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/LeeZXin/remote-sqlite/reqvo"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"io"
	"net/http"
)

const (
	jsonContentType = "application/json;charset=utf-8"
)

type Client struct {
	Host       string
	HttpClient *http.Client
	Secret     string
}

func (c *Client) httpClient() *http.Client {
	if c.HttpClient != nil {
		return c.HttpClient
	}
	return http.DefaultClient
}

func (c *Client) url(path string) string {
	return "http://" + c.Host + path
}

func (c *Client) NewNamespace(ctx context.Context, namespace string) error {
	_, err := c.post(ctx, c.url("/api/v1/newNamespace"), reqvo.NewNamespaceReqVO{
		Namespace: namespace,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteNamespace(ctx context.Context, namespace string) error {
	_, err := c.post(ctx, c.url("/api/v1/deleteNamespace"), reqvo.DeleteNamespaceReqVO{
		Namespace: namespace,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) ShowNamespace(ctx context.Context, namespace string) ([]string, error) {
	body, err := c.post(ctx, c.url("/api/v1/showNamespace"), reqvo.ShowNamespaceReqVO{
		Namespace: namespace,
	})
	if err != nil {
		return nil, err
	}
	ret := make([]string, 0)
	err = json.Unmarshal(body, &ret)
	return ret, err
}

func (c *Client) CreateDB(ctx context.Context, namespace, dbName string) error {
	_, err := c.post(ctx, c.url("/api/v1/createDB"), reqvo.CreateDBReqVO{
		Namespace: namespace,
		DbName:    dbName,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) ExecuteCommand(ctx context.Context, namespace, dbName, cmd string) (int64, error) {
	body, err := c.post(ctx, c.url("/api/v1/executeCommand"), reqvo.ExecuteCommandReqVO{
		Namespace: namespace,
		DbName:    dbName,
		Cmd:       cmd,
	})
	if err != nil {
		return 0, err
	}
	ret := make(gin.H)
	err = json.Unmarshal(body, &ret)
	return cast.ToInt64(ret["affectedRows"]), nil
}

func (c *Client) QueryCommand(ctx context.Context, namespace, dbName, cmd string) ([]map[string]string, error) {
	body, err := c.post(ctx, c.url("/api/v1/queryCommand"), reqvo.QueryCommandReqVO{
		Namespace: namespace,
		DbName:    dbName,
		Cmd:       cmd,
	})
	if err != nil {
		return nil, err
	}
	ret := make([]map[string]string, 0)
	err = json.Unmarshal(body, &ret)
	return ret, nil
}

func (c *Client) DropDB(ctx context.Context, namespace, dbName string) error {
	_, err := c.post(ctx, c.url("/api/v1/dropDB"), reqvo.DropDBReqVO{
		Namespace: namespace,
		DbName:    dbName,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) post(ctx context.Context, url string, req any) ([]byte, error) {
	m, _ := json.Marshal(req)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(m))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", jsonContentType)
	request.Header.Set("Rs-Secret", c.Secret)
	resp, err := c.httpClient().Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	ret, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusBadRequest {
		return nil, errors.New("bad request")
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, errors.New("unauthorized")
	}
	if resp.StatusCode == http.StatusInternalServerError {
		return nil, fmt.Errorf("internal error: %v", string(ret))
	}
	if resp.StatusCode == http.StatusOK {
		return ret, nil
	}
	return nil, errors.New("fail")
}
