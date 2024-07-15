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
	Bot struct {
		Enabled     bool
		XUID        string
		DisplayName string
	}
	Rcon struct {
		Enabled  bool
		Port     int
		Password string
	}
	Server struct {
		ViewDistance    int32
		Whitelist       bool
		DisableXboxAuth bool
		Ops             []string
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
		DiscordErrorLogsWebhook   string
		DiscordErrorLogsIconUrl   string
		DiscordStaffAlertsWebhook string
	}
}

func ReadConfig() Config {
	c := Config{
		Connection: struct {
			ProxyAddress  string
			RemoteAddress string
		}{"0.0.0.0:19134", "0.0.0.0:19132"},
	}

	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		f, err := os.Create("config.toml")
		if err != nil {
			log.Fatalf("error creating config: %v", err)
		}
		data, err := toml.Marshal(c)
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

	if err := toml.Unmarshal(data, &c); err != nil {
		log.Fatalf("error decoding config: %v", err)
	}

	if c.Connection.ProxyAddress == "" {
		panic("ProxyAddress is not assigned in config!")
	}

	if c.WorldBorder.Enabled && c.WorldBorder.MaxX == 0 && c.WorldBorder.MaxZ == 0 && c.WorldBorder.MinX == 0 && c.WorldBorder.MinZ == 0 {
		c.WorldBorder.MaxX = 1200
		c.WorldBorder.MaxZ = 1200
		c.WorldBorder.MinX = -1200
		c.WorldBorder.MinZ = -1200
	}

	if c.Server.ViewDistance == 0 {
		c.Server.ViewDistance = 10
	}

	if c.Rcon.Enabled && (c.Rcon.Port == 0 || c.Rcon.Password == "") {
		panic("Rcon is enabled and not configured in config!")
	}

	data, _ = toml.Marshal(c)
	if err := os.WriteFile("config.toml", data, 0644); err != nil {
		log.Fatalf("error writing config file: %v", err)
	}

	return c
}
