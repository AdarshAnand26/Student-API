package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string `yaml:"address" env-required:"true"`
}

type Config struct {
	Env         string     `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string     `yaml:"storage-path" env-required:"true"`
	HTTPServer  HTTPServer `yaml:"http-server"`
}

func MustLoad() *Config {

	var configPath string

	// 1. First check CONFIG_PATH environment variable
	configPath = os.Getenv("CONFIG_PATH")

	// 2. If CONFIG_PATH is empty, take path from command line
	if configPath == "" {

		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("config path is not set")
		}
	}

	// 3. Check whether config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	// 4. Create Config variable
	var cfg Config

	// 5. Read configuration file
	err := cleanenv.ReadConfig(configPath, &cfg)

	if err != nil {
		log.Fatalf("Cannot read config file: %s", err.Error())
	}

	// 6. Return configuration
	return &cfg
}

//working- Find the configuration file → check that it exists → read it → 
// ->convert its values into a Go Config struct → return the configuration.

//contro flow ->
//                  START
//                    │
//                    ↓
//         Look for CONFIG_PATH
//                    │
//           ┌────────┴────────┐
//           │                 │
//        Found             Not Found
//           │                 │
//           │                 ↓
//           │          Check -config argument
//           │                 │
//           │                 ↓
//           │          Still empty?
//           │                 │
//           │                YES
//           │                 ↓
//           │             STOP ❌
//           │
//           ↓
//    Check config file
//        exists?
//           │
//       ┌───┴───┐
//       │       │
//      YES      NO
//       │       │
//       │       ↓
//       │     STOP ❌
//       │
//       ↓
//  Read YAML file
//       │
//       ↓
// Put values into Config struct
//       │
//       ↓
//    Return cfg
//       │
//       ↓
//      DONE ✅