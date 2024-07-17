package config

type Config struct {
	ServerHost string
	ServerPort string
	MongoURI   string
}

// TODO: Implement the GetConfig function
func GetConfig() *Config {
	return &Config{
		ServerHost: "localhost",
		ServerPort: "3001",
		MongoURI:   "mongodb://localhost:27017",
	}
}
