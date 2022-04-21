package master

type Server struct {
	storage Storage
}

type Config struct {
	Storage        Storage
	SessionAddress string
}

func New(cfg Config) (*Server, error) {
	if cfg.Storage == nil {
		cfg.Storage = NewUnauthStorage()
	}

	return &Server{
		cfg.Storage,
	}, nil
}

func (s *Server) Start() {
}
