package user_domain

type TokenDetails struct {
	UserId              string `json:"user_id"`
	AccessToken         string `json:"access_token"`
	RefreshToken        string `json:"refresh_token"`
	AccessTokenExpires  int64  `json:"access_token_ttl"`
	RefreshTokenExpires int64  `json:"refresh_token_ttl"`
}
