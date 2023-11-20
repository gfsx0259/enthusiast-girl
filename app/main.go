package main

import "deployRunner/telegram"

func main() {
	telegram.NewListener().Listen()
}
