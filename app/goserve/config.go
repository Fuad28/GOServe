package goserve

type Config struct {

	// Port to start the server on, defaults to 4221
	Port int

	// MaxRequestSize sets the buffer size for the bytes array which reads the request.
	MaxRequestSize int

	// Array of domains that are allowed when the CORS middleware inspects the request.
	AllowedOrigins []string
}
