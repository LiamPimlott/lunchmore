package mail

// ClientConfig describes configuration for connecting to an smtp server
type ClientConfig struct {
	Username string
	Password string
	Host     string
	Port     int
}
