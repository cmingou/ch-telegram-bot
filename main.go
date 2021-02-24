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

func getMessageLink(m *tb.Message) string {
	chId := strconv.Itoa(int(m.Chat.ID))[4:]
	link := fmt.Sprintf("https://t.me/c/%v/%d", chId, m.ID)
	return link
}

func writeAuthor(m *tb.Message, docId string) {
	message := fmt.Sprintf("%v %v(%v):\n",
		m.Sender.FirstName, m.Sender.LastName, m.Sender.Username)

	if err := doc.InsertText(docId, message); err != nil {
		fmt.Printf("Insert Author failed, err: %v\n", err)
	}
}

func writeHyperLink(m *tb.Message, docId string) {
	if err := doc.InsertHyperLink(docId, "Telegram-link", getMessageLink(m)); err != nil {
		fmt.Printf("Insert hyper link failed, err: %v\n", err)
	}
}

func getRoomChatId(m *tb.Message) string {
	var docId string
	docId, ok := roomMap[m.Chat.ID]
	if !ok {
		fmt.Printf("The room ID: %v is not in map, using default value\n", m.Chat.ID)
		docId = roomMap[0]
	}
	return docId
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
		docId := getRoomChatId(m)

		lock.Lock()
		defer lock.Unlock()

		writeAuthor(m, docId)

		if err := doc.InsertText(docId, m.Text+"\n"); err != nil {
			fmt.Printf("Insert text failed, err: %v\n", err)
		}

		writeHyperLink(m, docId)
	})

	b.Handle(tb.OnPhoto, func(m *tb.Message) {
		docId := getRoomChatId(m)

		lock.Lock()
		defer lock.Unlock()

		writeAuthor(m, docId)

		photoUrl, err := telegram.GetPhotoUrl(TOKEN, m.Photo.File.FileID)
		if err != nil {
			fmt.Printf("Get photo url failed, err: %v\n", err)
		}

		if err := doc.InsertImage(docId, photoUrl); err != nil {
			fmt.Printf("Insert image to doc failed, err: %v\n", err)
		}

		writeHyperLink(m, docId)
	})

	b.Handle("/help", func(m *tb.Message) {
		b.Reply(m, "Alpha version")
	})

	fmt.Printf("bot started!!\n")
	b.Start()
}
