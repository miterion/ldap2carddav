package main

import (
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
	card.SetValue(vcard.FieldEmail, entry.GetAttributeValue("mail"))
	card.Set(vcard.FieldTelephone, &vcard.Field{
		Value: entry.GetAttributeValue("mobile"),
		Params: vcard.Params{
			vcard.ParamType: {vcard.TypeCell},
		},
	})
	card.SetValue(vcard.FieldUID, entry.GetAttributeValue("uid"))
	vcard.ToV4(card)
	return &card
}
