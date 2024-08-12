package goserve

// Config exposes the key parameters needed in creating a new server
type Config struct {

	// Port to start the server on, defaults to 8000
	Port int

	// MaxRequestSize sets the buffer size for the bytes array which reads the request.
	MaxRequestSize int

	// Array of domains that are allowed when the CORS middleware inspects the request.
	AllowedOrigins []string
}
