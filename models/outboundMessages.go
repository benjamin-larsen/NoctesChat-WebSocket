package models

type AuthErrorOutbound struct {
	MsgType string `json:"type"`
	Message string `json:"error"`
	ErrorCode int  `json:"code"`
}

var AuthErrLoggedOut = AuthErrorOutbound{
	MsgType: "auth_error",
	Message: "You've been logged out. Please log in and try again.",
	ErrorCode: 401,
}