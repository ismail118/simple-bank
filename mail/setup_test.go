package mail

import (
	"github.com/ismail118/simple-bank/util"
	"github.com/rs/zerolog/log"
	"os"
	"testing"
)

var SenderEmailTest SenderEmail

func TestMain(m *testing.M) {
	conf, err := util.LoadConfig("..")
	if err != nil {
		log.Fatal().Msgf("cannot load config error:%s", err)
	}
	SenderEmailTest = NewGmailSender(conf.EmailSenderName, conf.EmailSenderAddress, conf.EmailSenderPassword)

	os.Exit(m.Run())
}
