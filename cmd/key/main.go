package main

import (
	"errors"
	"flag"
	"fmt"
	"keeper/internal/core/encryption"
	"log"
	"os"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	MasterKey string `env:"MASTER_KEY,unset"`
}

func main() {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		log.Fatal("Error parsing environment variables: ", err)
	}

	masterKey, err := encryption.LoadPrivateKey([]byte(cfg.MasterKey))
	if err != nil {
		log.Fatal("Loading master key error", err)
	}

	var version int
	flag.IntVar(&version, "v", 1, "encryption key version")
	flag.Parse()

	path := fmt.Sprintf("encryption_keys/encryption_key_v_%d.pem", version)

	if _, err := os.Stat(path); err == nil {
		log.Fatalf("Key veriosn %d allready exists \n", version)
	} else if !errors.Is(err, os.ErrNotExist) {
		log.Fatal("Checking key version error. Path:", path, err)
	}

	key, err := encryption.NewPrivateKey()
	if err != nil {
		log.Fatal("Creating new encryption key error", err)
	}

	encKey, err := masterKey.Encrypt([]byte(key.Dump()))
	if err != nil {
		log.Fatal("Encripting encryption key error", err)
	}

	log.Println("Creating key file: ", path)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Fatal("Creating key file error", err)
	}

	f.Write(encKey)
	f.Close()
	log.Println("Key file created: ", path)
}
