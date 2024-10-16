package rs

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func newClient() *Client {
	return &Client{
		Host:       "127.0.0.1:15899",
		HttpClient: http.DefaultClient,
		Secret:     "1234",
	}
}

func TestClient_CreateDB(t *testing.T) {
	err := newClient().CreateDB(context.Background(), "lizexin2", "fick")
	if err != nil {
		panic(err)
	}
}

func TestClient_QueryCommand(t *testing.T) {
	command, err := newClient().QueryCommand(context.Background(), "lizexin", "fick", "select * from COMPANY")
	if err != nil {
		panic(err)
	}
	fmt.Println(command)
}

func TestClient_ShowNamespace(t *testing.T) {
	dbs, err := newClient().ShowNamespace(context.Background(), "lizexin2")
	if err != nil {
		panic(err)
	}
	fmt.Println(dbs)
}

func TestClient_ExecuteCommand(t *testing.T) {
	rows, err := newClient().ExecuteCommand(context.Background(), "lizexin", "fick", `
	create index idx_name_2 on COMPANY(name2)
`)
	if err != nil {
		panic(err)
	}
	fmt.Println(rows)
}

func TestClient_GetDBSize(t *testing.T) {
	size, err := newClient().GetDBSize(context.Background(), "lizexin", "fick")
	if err != nil {
		panic(err)
	}
	fmt.Println(size)
}
