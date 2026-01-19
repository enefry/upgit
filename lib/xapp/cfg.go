package xapp

type Config struct {
	DefaultUploader string            `toml:"default_uploader,omitempty"`
	Rename          string            `toml:"rename,omitempty"`
	Replacements    map[string]string `toml:"replacements,omitempty"`
	OutputFormats   map[string]string `toml:"output_formats,omitempty"`
	HmacKey         string            `toml:"hmac_key,omitempty"`
	HmacFormat      string            `toml:"hmac_format,omitempty"`
	HmacLen         int               `toml:"hmac_len,omitempty"`
}

var AppCfg Config
