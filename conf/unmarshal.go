package conf

import (
	"encoding/json"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func UnmarshalConf(data string) (conf *Config, err error) {
	conf = new(Config)
	bytes := []byte(data)

	err = json.Unmarshal(bytes, conf)
	if err == nil {
		return conf, nil
	}

	err = yaml.Unmarshal(bytes, conf)
	if err == nil {
		return conf, nil
	}

	err = toml.Unmarshal(bytes, conf)
	if err == nil {
		return conf, nil
	}

	return nil, errors.New("conf unmarshal error: " + data)
}
