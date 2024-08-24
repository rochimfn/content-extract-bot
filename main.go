package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

const RETRY_GET_LIMIT = 3
const RETRY_DOWNLOAD_LIMIT = 3

func main() {
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatalln("TELEGRAM_TOKEN not found in environment")
	}

	if os.Getenv("TIKA_HOST") == "" {
		log.Fatalln("TIKA_HOST not found in environment")
	}

	b, err := gotgbot.NewBot(token, nil)
	if err != nil {
		log.Fatalln("failed to create new bot with error: ", err)
	}

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})

	dispatcher.AddHandler(handlers.NewCommand("start", start))
	dispatcher.AddHandler(handlers.NewCommand("help", help))
	dispatcher.AddHandler(handlers.NewMessage(message.Photo, handlePhoto))
	dispatcher.AddHandler(handlers.NewMessage(message.Document, handleDocument))

	updater := ext.NewUpdater(dispatcher, nil)
	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		log.Fatalln("failed to start polling: " + err.Error())
	}
	log.Printf("%s has been started...\n", b.User.Username)

	updater.Idle()
}

func start(b *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("got start command from", dumpStructAsJSON(ctx.Message.From))

	_, err := ctx.EffectiveMessage.Reply(
		b, "Welcome to Content Extractor Bot"+
			"\nThis bot will help you extract text from many kinds of formats, including images!."+
			"\n"+
			"\nAvailable Command"+
			"\n/start - Getting started"+
			"\n/help  - Help",
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed reply start command: %w", err)
	}
	return nil
}

func help(b *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("got help command from", dumpStructAsJSON(ctx.Message.From))

	_, err := ctx.EffectiveMessage.Reply(
		b, "This bot will help you extract text from many kinds of formats (including Images)."+
			"\nSend a file to this chat room, and the bot will start working."+
			"\n"+
			"\nThe list of supported formats can be found in https://tika.apache.org/1.28.1/formats.html",
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed reply help command: %w", err)
	}
	return nil
}

func handlePhoto(b *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("got photo from", dumpStructAsJSON(ctx.Message.From))

	photos := ctx.EffectiveMessage.Photo
	err := sendAccepted(b, ctx)
	if err != nil {
		return err
	}

	var contents []string
	for _, photo := range photos {
		file, err := getFile(b, photo.FileId, 0)
		if err != nil {
			return fmt.Errorf("failed get file: %w", err)
		}

		path, err := downloadFile(file.URL(b, nil), photo.FileUniqueId, 0)
		if err != nil {
			return fmt.Errorf("failed to download file: %w", err)
		}

		content, err := parseFile(path, "")
		if err != nil {
			return fmt.Errorf("failed to parse file: %w", err)
		}

		contents = append(contents, content)
	}

	err = sendContent(b, ctx, strings.Join(contents, "\n\n"))
	if err != nil {
		return err
	}

	return nil
}

func handleDocument(b *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("got document from", dumpStructAsJSON(ctx.Message.From))

	doc := ctx.EffectiveMessage.Document
	err := sendAccepted(b, ctx)
	if err != nil {
		return err
	}

	file, err := getFile(b, doc.FileId, 0)
	if err != nil {
		return fmt.Errorf("failed get file: %w", err)
	}

	path, err := downloadFile(file.URL(b, nil), doc.FileName, 0)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}

	content, err := parseFile(path, doc.MimeType)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	err = sendContent(b, ctx, content)
	if err != nil {
		return err
	}

	return nil
}

func sendAccepted(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b, "Please wait, your file is being extracted!", nil)
	if err != nil {
		return fmt.Errorf("failed reply: %w", err)
	}

	return nil
}

func sendContent(b *gotgbot.Bot, ctx *ext.Context, content string) error {
	if len(content) == 0 {
		_, err := ctx.EffectiveMessage.Reply(b, "Cannot find text", nil)
		if err != nil {
			return fmt.Errorf("failed send content: %w", err)
		}

	} else if len(content) > 4096 {
		_, err := ctx.EffectiveMessage.Reply(
			b, "Text length exceeds telegram message limit :("+
				"\nThe text will be sent as a text file", nil)
		if err != nil {
			return fmt.Errorf("failed reply: %w", err)
		}

		filename := fmt.Sprint(time.Now().UnixNano(), ".txt")
		_, err = b.SendDocument(ctx.EffectiveChat.Id,
			gotgbot.InputFileByReader(filename, strings.NewReader(content)), nil)
		if err != nil {
			return fmt.Errorf("failed send content in file: %w", err)
		}

	} else {
		_, err := ctx.EffectiveMessage.Reply(b, content, nil)
		if err != nil {
			return fmt.Errorf("failed send content: %w", err)
		}
	}

	return nil
}

func getFile(b *gotgbot.Bot, fileid string, retry uint) (*gotgbot.File, error) {
	if retry > RETRY_GET_LIMIT {
		return nil, errors.New("get file retry limit exceeded")
	}

	file, err := b.GetFile(fileid, nil)
	if err != nil {
		return getFile(b, fileid, retry+1)
	}

	return file, nil
}

func downloadFile(url, fileName string, retry uint) (string, error) {
	if retry > RETRY_DOWNLOAD_LIMIT {
		return "", errors.New("download retry limit exceeded")
	}

	path := "/tmp/" + fileName

	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return downloadFile(url, fileName, retry+1)
	}

	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)

	return path, err
}

func parseFile(path, mime string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPut, os.Getenv("TIKA_HOST")+"/tika", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "text/plain")

	if mime != "" {
		req.Header.Set("Content-Type", mime)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func dumpStructAsJSON(structure any) string {
	b, err := json.Marshal(structure)
	if err != nil {
		log.Println("error marshal: ", err)
	}

	return string(b)
}
