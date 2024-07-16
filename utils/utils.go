package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func Format(a []any) string {
	return strings.TrimSuffix(strings.TrimSuffix(fmt.Sprintln(a...), "\n"), "\n")
}

func GetPluralForm(num int, form1 string, form2 string, form3 string) string {
	n := int(math.Abs(float64(num))) % 100
	n1 := n % 10
	if n > 10 && n < 20 {
		return form3
	} else if n1 > 1 && n1 < 5 {
		return form2
	} else if n1 == 1 {
		return form1
	}
	return form3
}

func GetFullPluralForm(num int, form1 string, form2 string, form3 string) string {
	n := int(math.Abs(float64(num))) % 100
	n1 := n % 10
	if n > 10 && n < 20 {
		return strconv.Itoa(num) + " " + form3
	} else if n1 > 1 && n1 < 5 {
		return strconv.Itoa(num) + " " + form2
	} else if n1 == 1 {
		return strconv.Itoa(num) + " " + form1
	}
	return strconv.Itoa(num) + " " + form3
}

func GetTimestamp() int64 {
	return time.Now().Unix()
}

func GetMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func MillisToDate(t int64) string {
	return TimestampToDate(t / 1000)
}

func TimestampToDate(t int64) string {
	return FormatTime(time.Unix(t, 0))
}

func FormatTime(t time.Time) string {
	return t.Format("15:04:05 2006.01.02") // magic numbers https://go.dev/src/time/format.go
}

// Fetches a tableName from the database
func FetchDatabase[T any](tableName string) (map[string]T, error) {
	dbConfig := ReadConfig().Database

	uri := dbConfig.Host + ":" + fmt.Sprint(dbConfig.Port) + "/api/database/" + dbConfig.Name + "/table/" + tableName

	log.Printf("Fetching \"%s\" from: \"%s\"\n", tableName, uri)

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Fatalln("Failed to create new request:", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", dbConfig.Key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		LogErrorToDiscord(err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("failed to fetch data, status code: %d", resp.StatusCode)
		LogErrorToDiscord(err)
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		LogErrorToDiscord(err)
		return nil, err
	}

	var response []struct {
		Data      T      `json:"data"`
		Key       string `json:"_key"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		LogErrorToDiscord(err)
		return nil, err
	}

	obj := make(map[string]T)
	for _, v := range response {
		obj[v.Key] = v.Data
	}
	return obj, nil
}

func LogErrorToDiscord(err error) {
	config := ReadConfig()
	log.Println("Error:", err)
	params := map[string]interface{}{
		"username":   fmt.Sprintf("[%s] Failed to Ping Database", config.Server.Prefix),
		"avatar_url": config.Logging.DiscordErrorLogsIconUrl,
		"content":    "@everyone",
		"embeds": []map[string]interface{}{
			{
				"title":       "Failed to Ping Database",
				"description": fmt.Sprintf("The Proxy has failed to ping the database on port: %s. Please check the database status!", fmt.Sprint(config.Database.Port)),
				"color":       16711680,
				"timestamp":   time.Now().Format(time.RFC3339),
				"fields": []map[string]interface{}{
					{
						"name":   "Error",
						"value":  err.Error(),
						"inline": true,
					},
				},
			},
		},
	}

	jsonParams, _ := json.Marshal(params)
	req, _ := http.NewRequest("POST", os.Getenv("DISCORD_STAFF_ALERTS_WEBHOOK"), io.NopCloser(bytes.NewBuffer(jsonParams)))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to send logs to discord", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Failed to send logs to discord, status code:", resp.StatusCode)
	}
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
