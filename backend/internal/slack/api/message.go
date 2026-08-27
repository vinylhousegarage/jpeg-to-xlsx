package api

type Message struct {
	Text   string
	Button *Button
}

type Button struct {
	Text string
	URL  string
}
