package main

import (
	"fmt"
	"os"
	pathlib "path"
	"path/filepath"

	"github.com/emersion/go-vcard"
	"github.com/emersion/go-webdav/carddav"
)

type CardDAVBackend struct {
	StorageRoot string
	Subdirectory string
}

func (cb CardDAVBackend) AddressBook() (*carddav.AddressBook, error) {
	return &carddav.AddressBook{
		Path:            cb.Subdirectory,
		Name:            "LDAP Adressbook",
		Description:     "Adressbook for LDAP Contacts",
		MaxResourceSize: 100 * 1024,
	}, nil
}

func (cb CardDAVBackend) GetAddressObject(path string, req *carddav.AddressDataRequest) (*carddav.AddressObject, error) {
	dirname, filename := pathlib.Split(path)
	ext := pathlib.Ext(filename)
	if dirname != cb.Subdirectory || ext != ".vcf" {
		return nil, fmt.Errorf("Contact not found: %s%s", dirname, filename)
	}
	return cb.getContact(filename)
}

func (cb CardDAVBackend) ListAddressObjects(req *carddav.AddressDataRequest) ([]carddav.AddressObject, error) {
	vcards, err := filepath.Glob(pathlib.Join(cb.StorageRoot, "*.vcf"))
	if err != nil {
		return nil, err
	}
	contacts := []carddav.AddressObject{}
	for _, cardpath := range vcards {
		_, card := pathlib.Split(cardpath)
		contact, err := cb.getContact(card)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, *contact)
	}
	return contacts, nil
}

func (cb CardDAVBackend) QueryAddressObjects(query *carddav.AddressBookQuery) ([]carddav.AddressObject, error) {
	panic("not implemented") // TODO: Implement
}

func (cb CardDAVBackend) PutAddressObject(path string, card vcard.Card) (loc string, err error) {
	return "Not supported", nil
}

func (cb CardDAVBackend) DeleteAddressObject(path string) error {
	return nil
}

func (cb CardDAVBackend) getContact(filename string) (*carddav.AddressObject, error) {
	f, err := os.Open(pathlib.Join(cb.StorageRoot, filename))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec := vcard.NewDecoder(f)
	card, err := dec.Decode()
	if err != nil {
		return nil, err
	}

	filestats, err := f.Stat()
	if err != nil {
		return nil, err
	}

	return &carddav.AddressObject{
		Path:    filename,
		ModTime: filestats.ModTime(),
		ETag:    fmt.Sprintf("%x%x", filestats.ModTime(), filestats.Size()),
		Card:    card,
	}, nil
}

func (cb CardDAVBackend) SaveContact(name string, card *vcard.Card) error {
	dest, err := os.Create(pathlib.Join(cb.StorageRoot, fmt.Sprintf("%s.vcf", name)))
	if err != nil {
		return err
	}
	defer dest.Close()

	enc := vcard.NewEncoder(dest)
	return enc.Encode(*card)
}
