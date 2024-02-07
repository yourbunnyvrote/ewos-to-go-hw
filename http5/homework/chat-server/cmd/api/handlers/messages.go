package handlers

const UsernameQueryParameter = "username"

type ChattingHandler struct {
	public  *PublicChatHandler
	private *PrivateChatHandler
}

func NewChattingHandler(public PublicChatting, private PrivateChatting) *ChattingHandler {
	return &ChattingHandler{
		public:  NewPublicChatHandler(public),
		private: NewPrivateChatHandler(private),
	}
}
