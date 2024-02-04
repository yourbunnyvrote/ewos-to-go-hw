package chatutil

type ChatDB struct {
	Users        map[string]User
	PublicChat   []Message
	PrivateChats map[Chat][]Message
}

type Chat struct {
	User1 string
	User2 string
}
