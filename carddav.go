package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/emersion/go-vcard"
	"github.com/go-ldap/ldap"
	"github.com/spf13/viper"
)

type CarddavWorker struct {
	channel chan []*ldap.Entry
	backend *CardDAVBackend
}

func NewCarddavWorker(ch chan []*ldap.Entry, backend *CardDAVBackend) *CarddavWorker {
	return &CarddavWorker{
		channel: ch,
		backend: backend,
	}
}

func (cw *CarddavWorker) Start() {
	logger := log.New(log.Writer(), "[CarddavWorker]	", log.Ldate|log.Ltime)
	for updates := range cw.channel {
		if viper.GetBool("carddav.clear_old_entries") {
			// clearing all vcfs before updating
			logger.Println("Deleting all vcards before sync")
			err := cw.backend.ClearAddressbook()
			if err != nil {
				log.Fatal(err)
			}
		}
		logger.Printf("Creating %d new vcard files\n", len(updates))
		for _, update := range updates {
			vcard := createVcardFromLdap(update)
			cw.backend.SaveContact(update.GetAttributeValue("uid"), vcard)
		}
	}
}

func createVcardFromLdap(entry *ldap.Entry) *vcard.Card {
	card := make(vcard.Card)
	card.SetValue(vcard.FieldName, strings.Join([]string{entry.GetAttributeValue("sn"), entry.GetAttributeValue("givenName")}, ";"))
	card.Set(vcard.FieldEmail, &vcard.Field{
		Value: entry.GetAttributeValue("mail"),
		Params: vcard.Params{
			vcard.ParamType: {vcard.TypeWork},
		},
	})
	//card.SetValue(vcard.FieldEmail, entry.GetAttributeValue("mail"))
	card.Set(vcard.FieldTelephone, &vcard.Field{
		Value: entry.GetAttributeValue("mobile"),
		Params: vcard.Params{
			vcard.ParamType: {vcard.TypeCell},
		},
	})
	card.SetValue(vcard.FieldUID, entry.GetAttributeValue("uid"))
	card.SetValue(vcard.FieldPhoto, fmt.Sprintf("data:image/jpeg;base64,%s", base64.StdEncoding.EncodeToString([]byte(entry.GetAttributeValue("jpegPhoto")))))
	if entry.GetAttributeValue("birthyear") != "" || entry.GetAttributeValue("birthmonth") != "" || entry.GetAttributeValue("birthday") != "" {
		card.SetValue(vcard.FieldBirthday, fmt.Sprintf("%s%02s%02s", entry.GetAttributeValue("birthyear"), entry.GetAttributeValue("birthmonth"), entry.GetAttributeValue("birthday")))
	}
	vcard.ToV4(card)
	return &card
}
