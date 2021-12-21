package main

import (
	"log"
	"time"

	"github.com/go-ldap/ldap"
	"github.com/spf13/viper"
)

var (
	ldapAttributes = []string{
		viper.GetString("ldap.unique_id_field"),
		"givenname",
		"sn",
		viper.GetString("ldap.phone_field"),
		"mail",
		viper.GetString("ldap.avatar_field"),
		"birthday", "birthmonth", "birthyear",
	}
)

type LdapWorkerConfig struct {
	scrapeTime time.Duration
	channel    chan []*ldap.Entry
	logger     *log.Logger
}

// NewLdapWorker creates a new LdapWorker instance
func NewLdapWorker(channel chan []*ldap.Entry) *LdapWorkerConfig {
	duration, err := time.ParseDuration(viper.GetString("ldap.scrapetime"))
	if err != nil {
		log.Fatalf("Scrapetime is in an invalid format: %s", err)
	}
	return &LdapWorkerConfig{duration, channel, log.New(log.Writer(), "[LDAPWorker]	", log.Ldate|log.Ltime)}
}

func (config *LdapWorkerConfig) Start() {
	for {
		config.logger.Println("Starting scrape...")
		l, err := ldap.DialURL(viper.GetString("ldap.url"))
		if err != nil {
			config.logger.Printf("Could not connect to LDAP in this cycle: %s \n", err)
			time.Sleep(config.scrapeTime)
			continue
		}

		bindDn := viper.GetString("ldap.binddn")
		if pw := viper.GetString("ldap.bindpw"); pw != "" {
			config.logger.Println("LDAP password set, using authenticated bind.")
			err = l.Bind(bindDn, pw)
		} else {
			config.logger.Println("LDAP password empty, using unauthenticated bind.")
			err = l.UnauthenticatedBind(bindDn)
		}

		if err != nil {
			config.logger.Fatalf("Error binding to LDAP: %s \n", err)
		}

		sr := ldap.NewSearchRequest(viper.GetString("ldap.basedn"), ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, viper.GetString("ldap.filter"), ldapAttributes, nil)

		res, err := l.Search(sr)
		if err != nil {
			config.logger.Printf("LDAP search failed: %s \n", err)
		}
		config.logger.Printf("Found %d users.\n", len(res.Entries))
		config.channel <- res.Entries

		l.Close()
		time.Sleep(config.scrapeTime)
	}
}
