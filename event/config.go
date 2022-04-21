package event

type Config struct {
	Transport  Transport // Transport is used for io
	PublicAddr string    // PublicAddr is used for listen (Server) / connect (Client)
}
