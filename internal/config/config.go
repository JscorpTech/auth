package config

import "os"

type Config struct {
	PrivateKey []byte
	PublicKey  []byte
	Addr       string
	AccessExp  int64
	RefreshExp int64
}

func NewConfig() *Config {
	privKey, err := os.ReadFile("keys/private.pem")
	if err != nil {
		panic("private.pem faylini o'qishda xatolik yuz berdi: " + err.Error())
	}
	pubKey, err := os.ReadFile("keys/public.pem")
	if err != nil {
		panic("public.pem faylini o'qishda xatolik yuz berdi: " + err.Error())
	}

	return &Config{
		PrivateKey: privKey,
		PublicKey:  pubKey,
		Addr:       os.Getenv("ADDR"),
		AccessExp:  60,
		RefreshExp: 43200,
	}
}
