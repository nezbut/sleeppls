package config

import (
	"flag"
	"os"
	"time"
)

type Config struct {
	TimeToShutDown time.Time
	NotifyDuration time.Duration
	NotifyIconPath string
	TGBotToken     string
	TGChatID       int
	SendToTelegram bool
}

func New(flagSet *flag.FlagSet) (*Config, error) {
	var (
		ttsd     string
		nd       string
		iconPath string
		tgToken  string
		tgChatID int
	)

	flagSet.StringVar(&ttsd, TimeToShutDownFlag, defaultTimeToShutDown, TimeToShutDownDesc)
	flagSet.StringVar(&nd, TimeToNotifyFlag, defaultNotifyDuration, NotifyDurationDesc)
	flagSet.StringVar(&iconPath, NotifyIconPathFlag, defaultNotifyIconPath, NotifyIconPathDesc)
	flagSet.StringVar(&tgToken, TGBotTokenFlag, defaultTGBotToken, TGBotTokenDesc)
	flagSet.IntVar(&tgChatID, TGChatIDFlag, defaultTGChatID, TGChatIDDesc)

	if err := flagSet.Parse(os.Args[1:]); err != nil {
		return &Config{}, err
	}

	if err := createIcon(iconPath); err != nil {
		return &Config{}, err
	}

	ttsdT := parseTimeToShutDown(ttsd)
	ndT := parseNotifyDuration(nd)

	return &Config{
		TimeToShutDown: ttsdT,
		NotifyDuration: ndT,
		NotifyIconPath: iconPath,
		TGBotToken:     tgToken,
		TGChatID:       tgChatID,
		SendToTelegram: tgToken != "" && tgChatID != 0,
	}, nil
}

func parseTimeToShutDown(notParsed string) time.Time {
	res, err := time.ParseInLocation(time.DateTime, notParsed, time.Local)
	now := time.Now()
	if err != nil {
		res, err = time.ParseInLocation(time.TimeOnly, notParsed, time.Local)
		if err != nil {
			res, _ = time.ParseInLocation(time.TimeOnly, defaultTimeToShutDown, time.Local)
		}
		res = time.Date(
			now.Year(), now.Month(), now.Day(), res.Hour(), res.Minute(), res.Second(), res.Nanosecond(), now.Location(),
		)
	}
	if res.Before(now) {
		res = res.Add(now.Sub(res)).Add(defaultHoursTimeToShutDown)
	}
	return res
}

func parseNotifyDuration(notParsed string) time.Duration {
	res, err := time.ParseDuration(notParsed)
	if err != nil {
		res, _ = time.ParseDuration(defaultNotifyDuration)
	}
	return res
}

func createIcon(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.WriteFile(path, notifyIcon, 0644)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}
