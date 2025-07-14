package utils

import (
	"context"
	"encoding/json"

	"github.com/sirupsen/logrus"
)

type AuthUser struct {
	UserID         string `json:"userId"`
	MerchantID     string `json:"merchantId"`
	MerchantSecret string `json:"merchantSecret"`
	MerchantApiKey string `json:"merchantApiKey"`
}

func GetUserAuth(ctx context.Context) AuthUser {
	var u AuthUser
	user, ok := ctx.Value("user").(map[string]string)
	if !ok {
		logrus.Errorf("Invalid User")
		return u
	}
	jsonData, _ := json.Marshal(user)

	json.Unmarshal(jsonData, &u)
	return u
}
