package pingpong

type AliveInfoReply struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type Message struct {
	UUID    string `json:"uuid"`
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
	Time    string `json:"time"`
}
