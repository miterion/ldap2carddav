package main

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/emersion/go-vcard"
	"github.com/go-ldap/ldap"
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
	for updates := range cw.channel {
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
	card.SetValue(vcard.FieldBirthday, fmt.Sprintf("%s%02s%02s", entry.GetAttributeValue("birthyear"), entry.GetAttributeValue("birthmonth"), entry.GetAttributeValue("birthday")))
	vcard.ToV4(card)
	return &card
}
