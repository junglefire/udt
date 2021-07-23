package main

type Header struct {
	Command 	string
	User		string
	AccessKey	string
}

type ErrMsg struct {
	Errno 	int
	Msg 	string
}