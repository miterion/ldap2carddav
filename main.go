// A CardDAV server which uses LDAP as backend
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/emersion/go-webdav/carddav"
	"github.com/go-ldap/ldap/v3"
	"github.com/spf13/viper"
)

func main() {
	setDefaultConfig()
	parseConfig()

	// init backend
	backend := NewCardDAVBackend(
		viper.GetString("carddav.storage_path"),
		viper.GetString("carddav.subdirectory"),
		viper.GetString("carddav.address_book_name"),
	)

	// init LDAP worker
	ldapChannel := make(chan []*ldap.Entry, 10)
	ldapWorker := NewLdapWorker(ldapChannel)
	go ldapWorker.Start()

	// init CardDAV worker
	carddavWorker := NewCardDAVWorker(ldapChannel, &backend)
	go carddavWorker.Start()

	handler := carddav.Handler{
		Backend: backend,
	}

	addr := fmt.Sprintf("%s:%s", viper.GetString("carddav.address"), viper.GetString("carddav.port"))
	log.Printf("Starting carddav server on %s", addr)

	log.Fatal(http.ListenAndServe(addr, &handler))
}

func setDefaultConfig() {
	viper.SetDefault("carddav", map[string]interface{}{
		"storage_path":      "/srv/ldap2carddav",
		"subdirectory":      "cards",
		"address_book_name": "LDAP address book",
		"address":           "127.0.0.1",
		"port":              "8000",
		"clear_old_entries": true,
	})
	viper.SetDefault("ldap", map[string]interface{}{
		"unique_id_field": "uid",
		"phone_field":     "mobile",
		"avatar_field":    "jpegPhoto",
		"filter":          "(objectClass=Person)",
		"scrape_time":     "6000s",
	})
}

func parseConfig() {
	viper.SetConfigName("ldap2carddav")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Problem loading config file: %s \n", err)
	}
}
