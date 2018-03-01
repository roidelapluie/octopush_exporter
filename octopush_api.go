package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type apiResponse struct {
	Balance map[string]string
}

func getBalances(r route) map[string]string {
	q := url.Values{}
	q.Add("user_login", r.Login)
	q.Add("api_key", r.Key)
	f := q.Encode()

	req, err := http.NewRequest("POST", "https://www.octopush-dm.com/api/balance/json", strings.NewReader(f))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	timeout := time.Duration(10 * time.Second)
	client := &http.Client{Timeout: timeout}
	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("err while fetching balance: %v\n", err)
		return nil
	}
	resp := apiResponse{}
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		fmt.Printf("err while decoding balance: %v\n", err)
		return nil
	}
	return resp.Balance
}
