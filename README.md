# SleepPls

**SleepPls** is a CLI application that helps you fall asleep. This application turns off the device at a certain time and notifies the shutdown in some time.

## Installation

```bash
go install github.com/nezbut/sleeppls/cmd/sleeppls-cli
```

or

Clone the repository and install local:

```bash
make install
```

After installation and the first run you will have a picture that you need for `DesktopNotifier`, it will appear in the path you specified in `-nicon`, default is `$HOME/SleepPls_notify_icon.png`.

## Usage

```bash
sleeppls-cli --help

2025/05/07 19:15:47 Usage of SleepPls:
  -nicon string
        The path to the notification icon. (default "$HOME/SleepPls_notify_icon.png")
  -ntf string
        The time it takes for the app to notify you when the device is turned off. (default "1h")
  -t string
        The time at which the device is forced to shut down. You can specify the format "YYYY-MM-DD HH:MM:SS" or HH:MM:SS if today (default "23:59:59")
  -tg string
        The token of the Telegram bot used to send notifications.
  -tgchat int
        The ID of the Telegram chat where notifications will be sent.
```

## How it works

The app can turn off your device at 23:30:30 and will notify you 1 hour before shutting it down

```bash
sleeppls-cli -t "23:30:30" -ntf "1h"
```

By default the notification comes only via `DesktopNotifier`, but you can specify the token of the Telegram bot and the chat where the notification will also be sent via `TelegramBotNotifier`.

```bash
sleeppls-cli -t "23:30:30" -ntf "1h" -tg "<token>" -tgchat 123
```

You can also specify a custom path to the picture to be used in `DesktopNotifier`

```bash
sleeppls-cli -t "23:30:30" -ntf "1h" -nicon "/path/to/icon.png"
```

## License

[MIT License](https://github.com/nezbut/sleeppls/blob/main/LICENSE)
