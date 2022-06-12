package helper

var (
	_ ResponseModel = (*RegisterLoginResp)(nil)
)

type ResponseModel interface {
	_Response()
}

type DefaultResp struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

func (d DefaultResp) _Response() {}

type RegisterLoginResp struct {
	DefaultResp
	UserId       int64  `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
