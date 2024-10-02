package main

import (
	"rent_service/server"
	"rent_service/server/controllers/dummy"
)

func main() {
	var s = server.New(server.WithPort(42069))

	s.Extend(dummy.New())

	s.Run()
}

