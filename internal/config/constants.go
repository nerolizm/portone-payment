package config

import (
	"fmt"

	"github.com/Netflix/go-env"
)

type Environment struct {
	Port      string `env:"PORT,default=:8080"`
	BaseURL   string `env:"BASE_URL,default=https://api.iamport.kr"`
	ImpKey    string `env:"IMP_KEY,required=true"`
	ImpSecret string `env:"IMP_SECRET,required=true"`
}

var Env Environment

func Init() error {
	_, err := env.UnmarshalFromEnviron(&Env)
	if err != nil {
		return fmt.Errorf("failed to unmarshal env: %v", err)
	}

	return nil
}
