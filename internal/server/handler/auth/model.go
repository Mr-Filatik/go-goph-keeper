package auth

type registerReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerResp struct {
	Token string `json:"token"`
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResp struct {
	Token string `json:"token"`
}
