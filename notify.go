package sleeppls

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gen2brain/beeep"
)

var (
	ErrTGResponseNotOK = fmt.Errorf("telegram response not ok")
)

type baseTGResponse struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description,omitempty"`
}

type getMeResponse struct {
	baseTGResponse
	Result struct {
		ID        int    `json:"id"`
		FirstName string `json:"first_name,omitempty"`
		Username  string `json:"username"`
	} `json:"result,omitempty"`
}

// Notifier interface
type Notifier interface {
	Notify(msg string) error
	String() string
}

// Desktop Notifier
type DesktopNotifier struct {
	IconPath string
}

func NewDesktopNotifier(IconPath string) Notifier {
	return DesktopNotifier{
		IconPath: IconPath,
	}
}

func (d DesktopNotifier) Notify(msg string) error {
	err := beeep.Notify("SleepPls", msg, d.IconPath)
	if err != nil {
		return err
	}
	return nil
}

func (d DesktopNotifier) String() string {
	return "Desktop"
}

// TelegramBot Notifier
type TelegramBotNotifier struct {
	client      *http.Client
	token       string
	baseURL     string
	chatID      int
	botUsername string
	botID       int
}

func NewTelegramBotNotifier(client *http.Client, token string, chatID int) (Notifier, error) {
	if token == "" || chatID == 0 {
		return nil, fmt.Errorf("token or chatID is empty")
	}

	if client == nil {
		client = http.DefaultClient
	}

	notifier := &TelegramBotNotifier{
		client:  client,
		token:   token,
		baseURL: "https://api.telegram.org/bot" + token,
		chatID:  chatID,
	}
	if err := notifier.setUsernameAndID(); err != nil {
		return nil, err
	}
	return notifier, nil
}

func (tg *TelegramBotNotifier) Notify(msg string) error {
	var response baseTGResponse
	body := fmt.Sprintf(`{
		"chat_id": %d,
		"text": "%s"
	}`, tg.chatID, msg)
	err := tg.request(http.MethodPost, tg.baseURL+"/sendMessage", bytes.NewBuffer([]byte(body)), &response)
	if err != nil {
		return err
	}
	return handleBaseTGResponse(response)
}

func (tg *TelegramBotNotifier) String() string {
	return fmt.Sprintf(
		"TelegramBot: bot=%s, botID=%d, chatID=%d", tg.botUsername, tg.botID, tg.chatID,
	)
}

func (tg *TelegramBotNotifier) setUsernameAndID() error {
	var response getMeResponse
	err := tg.request(http.MethodGet, tg.baseURL+"/getMe", nil, &response)
	if err != nil {
		return err
	}

	if err := handleBaseTGResponse(response.baseTGResponse); err != nil {
		return err
	}

	if response.Result.FirstName != "" {
		tg.botUsername = response.Result.FirstName
	} else {
		tg.botUsername = response.Result.Username
	}
	tg.botID = response.Result.ID
	return nil
}

func (tg *TelegramBotNotifier) request(method string, url string, body io.Reader, responseModel any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return failedTGRequest(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := tg.client.Do(req)
	if err != nil {
		return failedTGRequest(err)
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(responseModel); err != nil {
		return failedTGRequest(err)
	}
	return nil
}

func handleBaseTGResponse(response baseTGResponse) error {
	if !response.Ok {
		return failedTGRequest(
			fmt.Errorf(
				"%w: %s, error code is %d",
				ErrTGResponseNotOK,
				response.Description,
				response.ErrorCode,
			),
		)
	}
	return nil
}

func failedTGRequest(err error) error {
	return fmt.Errorf("failed send request to TelegramBot: %w", err)
}
