package utils

import (
	"log"
	"os"

	"github.com/pelletier/go-toml"
)

type Config struct {
	Connection struct {
		ProxyAddress  string
		RemoteAddress string
	}
	Server struct {
		ViewDistance    int32
		Whitelist       bool
		DisableXboxAuth bool
		Prefix          string
	}
	WorldBorder struct {
		Enabled bool
		MinX    int32
		MinZ    int32
		MaxX    int32
		MaxZ    int32
	}
	Api struct {
		ApiHost string
		ApiKey  string
	}
	Resources struct {
		PackURLs []string
	}
	Database struct {
		Host string
		Port int
		Key  string
		Name string
	}
	Logging struct {
		DiscordCommandLogsWebhook string
		DiscordChatLogsWebhook    string
		DiscordSignLogsWebhook    string
		DiscordSignLogsIconURL    string
		DiscordStaffAlertsWebhook string
	}
}

func ReadConfig() Config {
	// Initialize with default values
	defaultConfig := Config{
		Connection: struct {
			ProxyAddress  string
			RemoteAddress string
		}{
			ProxyAddress:  "0.0.0.0:19132",
			RemoteAddress: "0.0.0.0:19134",
		},
	}

	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		log.Println("config.toml not found, creating default config...")
		f, err := os.Create("config.toml")
		if err != nil {
			log.Fatalf("error creating config: %v", err)
		}
		data, err := toml.Marshal(defaultConfig)
		if err != nil {
			log.Fatalf("error encoding default config: %v", err)
		}
		if _, err := f.Write(data); err != nil {
			log.Fatalf("error writing encoded default config: %v", err)
		}
		_ = f.Close()
	}

	data, err := os.ReadFile("config.toml")
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	c := Config{}
	if err := toml.Unmarshal(data, &c); err != nil {
		log.Fatalf("error decoding config: %v", err)
	}

	// Validate required fields and set defaults if necessary
	if c.Connection.ProxyAddress == "" {
		panic("ProxyAddress is not assigned in config!")
	}

	if c.Connection.RemoteAddress == "" {
		panic("RemoteAddress is not assigned in config!")
	}

	if c.Server.ViewDistance <= 0 {
		panic("ViewDistance must be a value greater than 0!")
	}

	if c.Database.Host == "" {
		panic("Database Host must be a valid address!")
	}

	if c.Database.Port == 0 {
		panic("Database Port must be a valid port number!")
	}

	if c.Api.ApiHost == "" {
		panic("API Host must be a valid address!")
	}

	if c.Logging.DiscordCommandLogsWebhook == "" {
		panic("Discord Command Logs Webhook must be provided!")
	}

	if c.Logging.DiscordChatLogsWebhook == "" {
		panic("Discord Chat Logs Webhook must be provided!")
	}

	if c.Logging.DiscordStaffAlertsWebhook == "" {
		panic("Discord Staff Alerts Webhook must be provided!")
	}

	return c
}
