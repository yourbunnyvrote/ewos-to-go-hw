package chatutil

type ChatDB struct {
	Users        map[string]User
	PublicChat   []Message
	PrivateChats map[string]map[string][]Message
}
