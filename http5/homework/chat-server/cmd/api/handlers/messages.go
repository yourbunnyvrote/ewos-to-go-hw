package handlers

import (
	"net/http"
	"strconv"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entities"
)

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

func paginateMessages(r *http.Request, messages []entities.Message) ([]entities.Message, error) {
	limitStr := r.URL.Query().Get(LimitQueryParameter)

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return nil, err
	}

	offsetStr := r.URL.Query().Get(OffsetQueryParameter)

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || limit <= 0 || offset <= 0 {
		return nil, err
	}

	startIndex := offset - 1
	endIndex := offset - 1 + limit

	if startIndex >= len(messages) {
		return nil, ErrEndOfPages
	}

	if endIndex > len(messages) {
		endIndex = len(messages)
	}

	return messages[startIndex:endIndex], nil
}
