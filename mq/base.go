package mq

type MsgType int

const (
	MSG_TYPE_COMMENT MsgType = iota
)

type CommentEntity struct {
	BookId 		int			`json:"book_id"`
	MemberId 	int			`json:"member_id"`
	Content 	string		`json:"content"`
}

type MsgEntity struct {
	Type 	MsgType 		`json:"msg_type"`
	Data 	string			`json:"data"`
}