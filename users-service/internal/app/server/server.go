package server

import "net/http"

type Server struct {
	Router http.Handler
}

func CreateNewServer(appRouter http.Handler) *Server {
	s := &Server{}
	s.Router = appRouter
	return s
}
