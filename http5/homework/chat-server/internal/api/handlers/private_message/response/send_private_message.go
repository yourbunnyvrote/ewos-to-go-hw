package response

type SendPrivateMessageResponse struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}
