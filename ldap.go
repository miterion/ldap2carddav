package main

import (
	"log"
	"time"

	"github.com/go-ldap/ldap"

	"github.com/spf13/viper"
)

var (
	ldapattributes = []string{"uid", "givenname", "sn", "mobile", "mail", "jpegPhoto"}
)

type LdapWorkerConfig struct {
	scrapetime time.Duration
	channel    chan []*ldap.Entry
	logger     *log.Logger
}

// NewLdapWorker creates a new LDAPWorker instance
func NewLdapWorker(channel chan []*ldap.Entry) *LdapWorkerConfig {
	duration, err := time.ParseDuration(viper.GetStringMapString("ldap")["scrapetime"])
	if err != nil {
		log.Fatalf("Scrapetime is in an invalid format: %s", err)
	}
	return &LdapWorkerConfig{duration, channel, log.New(log.Writer(), "[LDAPWorker] ", log.Ldate|log.Ltime)}
}

func (config *LdapWorkerConfig) Start() {
	for {
		config.logger.Println("Starting scrape")
		l, err := ldap.DialURL(viper.GetStringMapString("ldap")["url"])
		if err != nil {
			config.logger.Printf("Could not connect to LDAP in this cycle: %s \n", err)
			time.Sleep(config.scrapetime)
			continue
		}

		err = l.Bind(viper.GetStringMapString("ldap")["binddn"], viper.GetStringMapString("ldap")["bindpw"])
		if err != nil {
			config.logger.Fatalf("Error binding to LDAP: %s \n", err)
		}

		sr := ldap.NewSearchRequest(viper.GetStringMapString("ldap")["basedn"], ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, viper.GetStringMapString("ldap")["filter"], ldapattributes, nil)

		res, err := l.Search(sr)
		if err != nil {
			config.logger.Printf("LDAP Search failed: %s \n", err)
		}
		config.logger.Printf("Found %d users\n", len(res.Entries))
		config.channel <- res.Entries

		l.Close()
		time.Sleep(config.scrapetime)
	}
}
