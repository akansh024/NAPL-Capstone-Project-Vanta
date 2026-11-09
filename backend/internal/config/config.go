package config

import "os"

type Config struct{ DatabaseURL, FrontendOrigin, Port string }

func Load() Config {
	c := Config{DatabaseURL: os.Getenv("DATABASE_URL"), FrontendOrigin: os.Getenv("FRONTEND_ORIGIN"), Port: os.Getenv("PORT")}
	if c.Port == "" {
		c.Port = "8080"
	}
	if c.FrontendOrigin == "" {
		c.FrontendOrigin = "http://localhost:5173"
	}
	return c
}
