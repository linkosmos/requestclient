package main

import (
	"fmt"
	"net/url"

	"github.com/ernestas-poskus/requestclient"
)

func main() {

	client := requestclient.New(nil)

	u, _ := url.Parse("http://www.example.com")

	responseClient, _ := client.Do(client.GET(u)) // For higher level API

	fmt.Println(responseClient.Status)
}
