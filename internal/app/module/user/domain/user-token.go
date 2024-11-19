package user_domain

type TokenDetails struct {
	UserEmail           string `json:"user_email"`
	AccessToken         string `json:"access_token"`
	RefreshToken        string `json:"refresh_token"`
	AccessTokenExpires  int64  `json:"access_token_ttl"`
	RefreshTokenExpires int64  `json:"refresh_token_ttl"`
}
