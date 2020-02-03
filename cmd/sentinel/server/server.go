package server

import (
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

type Server struct {
	router      *mux.Router
	redisClient *redis.Client
}

func NewServer(rc *redis.Client) Server {
	return Server{
		router:      mux.NewRouter(),
		redisClient: rc,
	}
}

func (s Server) Router() *mux.Router {
	return s.router
}
