package main

const (
	fullCommand = "/full"
)

func main() {
	tgbot, err := NewTgBot()
	if err != nil {
		panic(err)
	}

	for update := range tgbot.GetUpdates() {
		switch {
		case update.Message == nil:
			continue
		case update.Message.Document != nil && update.Message.Document.MimeType == "image/heic":
			tgbot.Heic2JpgCompress(update.Message)
		case update.Message.Text == fullCommand && update.Message.ReplyToMessage != nil &&
			update.Message.ReplyToMessage.Document != nil && update.Message.ReplyToMessage.Document.MimeType == "image/heic":
			tgbot.Heic2JpgDoc(update.Message.ReplyToMessage)
		}
	}
}
