package model

import "time"

type LoginForm struct {
	FCMToken    string `json:"fcm_token"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginData struct {
	Token string `json:"token"`
}

type JWTResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    time.Time `json:"expires_in"`
}
