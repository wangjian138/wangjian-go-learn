package prometheus

// A Config is a prometheus config.
type Config struct {
	Host string `json:",default=localhost"`
	Port int    `json:",default=9101"`
	Path string `json:",default=/metrics"`
}
