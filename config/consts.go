package config

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		home = os.Getenv("HOME")
		if home == "" {
			home = fmt.Sprintf("/home/%s", os.Getenv("USER"))
		}
	}
	defaultNotifyIconPath = filepath.Join(home, "SleepPls_notify_icon.png")
}

var (
	//go:embed notify_icon.png
	notifyIcon            []byte
	defaultNotifyIconPath string
)

const (
	defaultTimeToShutDown      string        = "23:59:59"
	defaultNotifyDuration      string        = "1h"
	defaultHoursTimeToShutDown time.Duration = 6 * time.Hour
	defaultTGBotToken          string        = ""
	defaultTGChatID            int           = 0
)

const (
	TimeToShutDownFlag string = "t"
	TimeToNotifyFlag   string = "ntf"
	NotifyIconPathFlag string = "nicon"
	TGBotTokenFlag     string = "tg"
	TGChatIDFlag       string = "tgchat"
)

const (
	TimeToShutDownDesc = "The time at which the device is forced to shut down. You can specify the format \"YYYY-MM-DD HH:MM:SS\" or HH:MM:SS if today"
	NotifyDurationDesc = "The time it takes for the app to notify you when the device is turned off."
	NotifyIconPathDesc = "The path to the notification icon."
	TGBotTokenDesc     = "The token of the Telegram bot used to send notifications."
	TGChatIDDesc       = "The ID of the Telegram chat where notifications will be sent."
)
