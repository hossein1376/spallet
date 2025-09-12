package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand/v2"
	"net/http"
	"time"
)

func main() {
	var usersNum int
	var baseURL string
	flag.IntVar(&usersNum, "users", 10, "Number of users")
	flag.StringVar(&baseURL, "url", "http://127.0.0.1:4004", "Base URL")
	flag.Parse()

	client := &http.Client{Timeout: 5 * time.Second}

	users, err := createUsers(client, baseURL, usersNum)
	if err != nil {
		fmt.Println("create users:", err)
		return
	}

	for _, user := range users {
		go sendRequest(client, baseURL, user)
	}

	select {}
}

type TopUpRequest struct {
	Amount      int        `json:"amount"`
	ReleaseDate *time.Time `json:"release_date,omitempty"`
}

type WithdrawRequest struct {
	Amount int `json:"amount"`
}

func sendRequest(client *http.Client, baseURL string, userID int) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C

		// Launch 3 top-ups
		for range 3 {
			go func() {
				delay := time.Duration(rand.IntN(600)) * time.Millisecond
				time.Sleep(delay)

				req := TopUpRequest{
					Amount: rand.IntN(90) + 10, // between 10 and 100
				}

				// ~20% chance to include a release date
				if rand.Float64() < 0.2 {
					release := time.Now().Add(time.Duration(rand.IntN(10)+1) * time.Second)
					req.ReleaseDate = &release
				}
				body, _ := json.Marshal(req)
				url := fmt.Sprintf(baseURL+"/wallets/%d/topup", userID)

				resp, err := client.Post(url, "application/json", bytes.NewBuffer(body))
				if err != nil {
					fmt.Println("topup error:", err)
					return
				}
				resp.Body.Close()
			}()
		}

		// Launch 1 withdraw
		go func() {
			delay := time.Duration(rand.IntN(3)) * time.Second
			time.Sleep(delay)

			req := WithdrawRequest{
				Amount: rand.IntN(270) + 30, // between 30 and 300
			}

			body, _ := json.Marshal(req)
			url := fmt.Sprintf(baseURL+"/wallets/%d/withdraw", userID)

			resp, err := client.Post(url, "application/json", bytes.NewBuffer(body))
			if err != nil {
				fmt.Println("withdraw error:", err)
				return
			}
			resp.Body.Close()
		}()
	}
}

type UserResponse struct {
	Data struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
	} `json:"data"`
}

func createUsers(client *http.Client, baseURL string, n int) ([]int, error) {
	userIDs := make([]int, n)

	for i := range n {
		userID, err := createUser(client, baseURL)
		if err != nil {
			return nil, fmt.Errorf("creating user: %w", err)
		}
		userIDs[i] = userID
	}

	return userIDs, nil
}

func createUser(client *http.Client, baseURL string) (int, error) {
	username := fmt.Sprintf("user_%d", rand.Uint64())
	body := []byte(fmt.Sprintf(`{"username": "%s"}`, username))

	resp, err := client.Post(
		baseURL+"/users", "application/json", bytes.NewBuffer(body),
	)
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return 0, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var userResp UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return 0, fmt.Errorf("decode failed: %w", err)
	}

	return userResp.Data.ID, nil
}
