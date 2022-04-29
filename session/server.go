package session

type Server struct {
	config Config
}

type Config struct {
	MasterAddress string // MasterAddress is address of master.Server
	Token         []byte // Token is auth token.
}

func New(cfg Config) (*Server, error) {
	return &Server{
		config: cfg,
	}, nil
}

func (s *Server) Serve() error {
	panic("implement")
}
