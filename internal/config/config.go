package config

import (
	"os"

	"go.uber.org/zap"
)

type Config struct {
	PrivateKey     []byte
	PublicKey      []byte
	Addr           string
	AccessExp      int64
	RefreshExp     int64
	GoogleClientID string
	DatabaseDsn    string
	DatabaseType   string
}

func NewConfig(logger *zap.Logger) *Config {
	var privKey []byte
	var pubKey []byte
	var err error
	if os.Getenv("PRIVATE_KEY") != "" {
		privKey = []byte(os.Getenv("PRIVATE_KEY"))
	} else {
		privKey, err = os.ReadFile("keys/private.pem")
		if err != nil {
			logger.Info("private.pem faylini o'qishda xatolik yuz berdi: ", zap.Error(err))
		}
	}
	if os.Getenv("PUBLIC_KEY") != "" {
		pubKey = []byte(os.Getenv("PUBLIC_KEY"))
	} else {
		pubKey, err = os.ReadFile("keys/public.pem")
		if err != nil {
			logger.Info("public.pem faylini o'qishda xatolik yuz berdi: ", zap.Error(err))
		}
	}

	return &Config{
		PrivateKey:     privKey,
		PublicKey:      pubKey,
		Addr:           os.Getenv("ADDR"),
		AccessExp:      60,
		RefreshExp:     43200,
		GoogleClientID: os.Getenv("GOOGLE_CLIENT_ID"),
		DatabaseType:   os.Getenv("DATABASE_TYPE"),
		DatabaseDsn:    os.Getenv("DATABASE_DSN"),
	}
}
