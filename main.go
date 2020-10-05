// A carddav server which use LDAP as a backend
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/emersion/go-webdav/carddav"
	"github.com/go-ldap/ldap"
	"github.com/spf13/viper"
)

func main() {
	parseConfig()

	// init backend
	backend := CardDAVBackend{StorageRoot: viper.GetStringMapString("carddav")["storagepath"]}

	// init ldapworker
	ldapChannel := make(chan []*ldap.Entry, 10)
	ldapWorker := NewLdapWorker(ldapChannel)
	go ldapWorker.Start()

	// init carddavworker
	carddavWorker := NewCarddavWorker(ldapChannel, &backend)
	go carddavWorker.Start()

	handler := carddav.Handler{
		Backend: backend,
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", viper.GetStringMapString("carddav")["port"]), &handler))
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
