package config

import (
	"bytes"

	tomlv2 "github.com/pelletier/go-toml/v2"
)

func TomlUnmarshaler(p []byte, v interface{}) error {
	return tomlv2.Unmarshal(p, v)
}

func TomlMarshalIndent(cfg interface{}) (string, error) {
	buf := bytes.Buffer{}
	enc := tomlv2.NewEncoder(&buf)
	enc.SetIndentTables(true)
	err := enc.Encode(cfg)
	return buf.String(), err
}
