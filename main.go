package main

import (
	"fmt"
	"github.com/cmingou/ch-telegram-bot/internal/google/doc"
	"github.com/cmingou/ch-telegram-bot/internal/telegram"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"strconv"
	"sync"
	"time"
)

var (
	b    *tb.Bot
	lock sync.Mutex
)

func init() {
	doc.DocId = DOC_ID
}

func getMessageLink(m *tb.Message) string {
	chId := strconv.Itoa(int(m.Chat.ID))[4:]
	link := fmt.Sprintf("https://t.me/c/%v/%d", chId, m.ID)
	return link
}

func writeAuthor(m *tb.Message) {
	message := fmt.Sprintf("%v %v(%v):\n",
		m.Sender.FirstName, m.Sender.LastName, m.Sender.Username)

	if err := doc.InsertText(message); err != nil {
		fmt.Printf("%v\n", err)
	}
}

func writeHyperLink(m *tb.Message) {
	if err := doc.InsertHyperLink("Telegram-link", getMessageLink(m)); err != nil {
		fmt.Printf("%v\n", err)
	}
}

func main() {
	var err error
	b, err = tb.NewBot(tb.Settings{
		// Token for bot
		Token:  TOKEN,
		Poller: &tb.LongPoller{Timeout: 5 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle(tb.OnText, func(m *tb.Message) {
		lock.Lock()
		defer lock.Unlock()

		writeAuthor(m)

		if err := doc.InsertText(m.Text + "\n"); err != nil {
			fmt.Printf("%v\n", err)
		}

		writeHyperLink(m)
	})

	b.Handle(tb.OnPhoto, func(m *tb.Message) {
		lock.Lock()
		defer lock.Unlock()

		writeAuthor(m)

		photoUrl, err := telegram.GetPhotoUrl(TOKEN, m.Photo.File.FileID)
		if err != nil {
			fmt.Printf("%v\n", err)
		}

		if err := doc.InsertImage(photoUrl); err != nil {
			fmt.Printf("%v\n", err)
		}

		writeHyperLink(m)
	})

	b.Handle("/help", func(m *tb.Message) {
		b.Reply(m, "Alpha version")
	})

	b.Start()
}
