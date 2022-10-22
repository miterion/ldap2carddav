package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/emersion/go-vcard"
	"github.com/go-ldap/ldap/v3"
	"github.com/spf13/viper"
)

type CardDAVWorker struct {
	channel chan []*ldap.Entry
	backend *CardDAVBackend
}

func NewCardDAVWorker(ch chan []*ldap.Entry, backend *CardDAVBackend) *CardDAVWorker {
	return &CardDAVWorker{
		channel: ch,
		backend: backend,
	}
}

func (cw *CardDAVWorker) Start() {
	logger := log.New(log.Writer(), "[CardDAVWorker]	", log.Ldate|log.Ltime)
	for updates := range cw.channel {
		if viper.GetBool("carddav.clear_old_entries") {
			// clearing all vcfs before updating
			logger.Println("Deleting all vCards before sync...")
			err := cw.backend.ClearAddressBook()
			if err != nil {
				log.Fatal(err)
			}
		}
		logger.Printf("Creating %d new vCard files...\n", len(updates))
		for _, update := range updates {
			vCard := createVCardFromLdap(update)
			cw.backend.SaveContact(update.GetAttributeValue(viper.GetString("ldap.unique_id_field")), vCard)
		}
	}
}

func createVCardFromLdap(entry *ldap.Entry) *vcard.Card {
	card := make(vcard.Card)
	card.SetValue(vcard.FieldName, strings.Join([]string{entry.GetAttributeValue("sn"), entry.GetAttributeValue("givenName")}, ";"))
	card.Set(vcard.FieldEmail, &vcard.Field{
		Value: entry.GetAttributeValue("mail"),
		Params: vcard.Params{
			vcard.ParamType: {vcard.TypeWork},
		},
	})
	card.Set(vcard.FieldTelephone, &vcard.Field{
		Value: entry.GetAttributeValue(viper.GetString("ldap.phone_field")),
		Params: vcard.Params{
			vcard.ParamType: {vcard.TypeCell},
		},
	})
	card.SetValue(vcard.FieldUID, entry.GetAttributeValue(viper.GetString("ldap.unique_id_field")))
	card.SetValue(vcard.FieldPhoto, fmt.Sprintf("data:image/jpeg;base64,%s", base64.StdEncoding.EncodeToString([]byte(entry.GetAttributeValue(viper.GetString("ldap.avatar_field"))))))
	if entry.GetAttributeValue("birthyear") != "" || entry.GetAttributeValue("birthmonth") != "" || entry.GetAttributeValue("birthday") != "" {
		card.SetValue(vcard.FieldBirthday, fmt.Sprintf("%s%02s%02s", entry.GetAttributeValue("birthyear"), entry.GetAttributeValue("birthmonth"), entry.GetAttributeValue("birthday")))
	}
	card.SetValue(vcard.FieldVersion, "3.0")
	return &card
}
