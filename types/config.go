package types

import "errors"

type Config struct {
	Regular         string
	Target          string
	ThreadNum       uint
	ThreadWorkRange uint64
}

func (c *Config) Check() error {
	if c.Regular == "" {
		return errors.New("regular cannot be null")
	}
	if c.Target == "" {
		return errors.New("regular cannot be null")
	}
	return nil
}
