package response

type ShowPrivateMessagesResponse struct {
	Messages []Message `json:"messages"`
}
