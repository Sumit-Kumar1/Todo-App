package models

type AuthUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthUserResp struct {
	Msg          string `json:"message,omitempty"`
	Email        string `json:"email,omitempty"`
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}
