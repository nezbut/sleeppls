package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"runtime"
	"time"

	"github.com/nezbut/sleeppls"
	"github.com/nezbut/sleeppls/cli"
	"github.com/nezbut/sleeppls/config"
)

func main() {
	// Configure CLI application
	confErrBuf := bytes.NewBuffer([]byte{})
	confFlagSet := flag.NewFlagSet("SleepPls", flag.ContinueOnError)
	confFlagSet.SetOutput(confErrBuf)

	conf, err := config.New(confFlagSet)
	if err != nil {
		if confErrBuf.Len() != 0 {
			log.Fatal(confErrBuf.String())
		}
		log.Fatal(err)
	}

	// Initialize notifiers
	var notifiers []sleeppls.Notifier
	if conf.SendToTelegram {
		tg, err := sleeppls.NewTelegramBotNotifier(nil, conf.TGBotToken, conf.TGChatID)
		if err != nil {
			slog.Warn("Failed to create telegram bot notifier", "error", err)
			notifiers = []sleeppls.Notifier{
				sleeppls.NewDesktopNotifier(conf.NotifyIconPath),
			}
		} else {
			notifiers = []sleeppls.Notifier{
				sleeppls.NewDesktopNotifier(conf.NotifyIconPath),
				tg,
			}
		}
	} else {
		notifiers = []sleeppls.Notifier{
			sleeppls.NewDesktopNotifier(conf.NotifyIconPath),
		}
	}

	// Print welcome message and notifiers
	art := `
   _____ __                ____  __
  / ___// /__  ___  ____  / __ \/ /____
  \__ \/ / _ \/ _ \/ __ \/ /_/ / / ___/
 ___/ / /  __/  __/ /_/ / ____/ (__  )
/____/_/\___/\___/ .___/_/   /_/____/
                /_/
	`
	fmt.Println(art)
	fmt.Printf("Welcome %s to the \"SleepPls\" app\nThis app is needed to help you go to bed on time.\nShutdown is scheduled in %s, The app will remind you %s before shutdown.\n",
		os.Getenv("USER"), conf.TimeToShutDown.Format(time.DateTime), conf.NotifyDuration)

	fmt.Println("Im using the following notifiers:")
	for _, notifier := range notifiers {
		fmt.Printf("Notifier: %s\n", notifier)
	}

	// Start CLI app with computer shut downer
	shutDowner := sleeppls.NewComputerShutDowner(runtime.GOOS)
	app := cli.New(conf, shutDowner, notifiers...)
	if err = app.Start(); err != nil {
		log.Fatal(err)
	}
}
