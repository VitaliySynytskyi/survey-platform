package config

// Config represents the application configuration
type Config struct {
	DB   DBConfig
	Port string
}

// DBConfig represents the database configuration
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}
