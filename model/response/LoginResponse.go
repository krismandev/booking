package response

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiryTime   int64  `json:"expiryTime"`
}
