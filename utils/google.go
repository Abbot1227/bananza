package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// Запасной вариант

type UserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
}

func getGoogleUserInfo(token string) (*UserInfo, error) {
	client := http.Client{Timeout: time.Second * 30}
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo?access_token=" + token)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result UserInfo
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
