package main

import (
	"context"
	"fmt"
	"mime"
	"os"
	"os/signal"

	// "strings"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mmaks17/vkvoice"
	"github.com/mmaks17/yavoice"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

var (
	starttime   int64
	client      *whatsmeow.Client
	vktoken     string
	yatoken     string
	VOICE_MODEL string
	// whitechat string
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func init() {
	starttime = time.Now().Unix()
	vktoken = getenv("vktoken", "BLANK")
	yatoken = getenv("yatoken", "BLANK")
	VOICE_MODEL = getenv("VOICE_MODEL", "YANDEX")
	// whitechat = "7XXXXXXXXXX-XXXXXXXXXXXX" // тестовый час с ботом
}

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		if uint64(v.Info.Timestamp.Unix()) > uint64(starttime) { // && strings.Contains(v.Info.Chat.User, whitechat) {
			img := v.Message.GetAudioMessage()
			if img != nil {
				data, err := client.Download(img)
				if err != nil {
					fmt.Printf("Failed to download image: %v", err)
					return
				}
				exts, _ := mime.ExtensionsByType(img.GetMimetype())
				path := fmt.Sprintf("%s%s", v.Info.ID, exts[0])
				err = os.WriteFile(path, data, 0666)
				if err != nil {
					fmt.Printf("Failed to save image: %v", err)
					return
				}
				var rezstr string
				var errv error
				if img.GetSeconds() >= 30 && vktoken != "BLANK" {
					rezstr, errv = vkvoice.Voice2Text(path, vktoken)
				} else {
					if VOICE_MODEL == "YANDEX" {
						rezstr, errv = yavoice.Voice2Text(path, yatoken)
					}
					if VOICE_MODEL == "MAILRU" {
						rezstr, errv = vkvoice.Voice2Text(path, vktoken)

					}

				}

				rezstr = v.Info.Sender.User + "по мнению " + VOICE_MODEL + " сказал: " + rezstr
				_ = os.Remove(path)
				if errv != nil {
					fmt.Println(err)
				} else {
					msg := &waProto.Message{Conversation: proto.String(rezstr)}

					_, err = client.SendMessage(v.Info.Chat, "", msg)
					if err != nil {
						fmt.Errorf("Error sending message: %v", err)
					}
				}
			}

		}

	}

}

func main() {
	dbLog := waLog.Stdout("Database", "INFO", true)
	container, err := sqlstore.New("sqlite3", "file:examplestore.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	clientLog := waLog.Stdout("Client", "INFO", true)
	client = whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(eventHandler)

	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				// Render the QR code here
				// e.g. qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				// or just manually `echo 2@... | qrencode -t ansiutf8` in a terminal
				fmt.Println("QR code:", evt.Code)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}
