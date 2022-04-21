package session

type Server struct{}

type Config struct{}

func New(cfg Config) (*Server, error) {
	return &Server{}, nil
}

func (s *Server) Start() {
}
