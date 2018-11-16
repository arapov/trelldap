package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/arapov/trelldap/env"
)

const (
	configfile = "config.json"
	datafile   = "members.json"
)

// Meta defines the data required to do the synchronization between
// LDAP and Trello Organization
type Meta struct {
	Fullname       string `json:"fullname"`
	TrelloActive   bool   `json:"trello"`
	TrelloID       string `json:"trelloid"`
	TrelloUserName string `json:"trellouser"`
	TrelloMail     string `json:"trellomail"`

	seenInLDAP bool
}

// Members is ssia, map of Meta data
type Members struct {
	Map map[string]*Meta `json:"members"`
}

func (m *Members) Read() error {
	jsonBytes, err := ioutil.ReadFile(datafile)
	json.Unmarshal(jsonBytes, &m)

	return err
}

func (m *Members) Write() error {
	jsonBytes, _ := json.MarshalIndent(m, "", "  ")
	err := ioutil.WriteFile(datafile, jsonBytes, 0644)

	return err
}

func main() {
	c, err := env.LoadConfig(configfile)
	if err != nil {
		log.Fatalln(err)
	}

	// members - TODO: document the importance
	var members Members
	var ldapMembers map[string]string

	if err := members.Read(); err != nil {
		log.Println("no", datafile, "file was found.")

		members.Map = make(map[string]*Meta)
	}

	// trello && ldap connections to work with
	trello := c.Trello.Dial()
	ldap := c.LDAP.Dial()

	// Add newly discovered in LDAP People to 'members'
	ldapMembers = ldap.Query()
	for uid, fullname := range ldapMembers {
		if _, ok := members.Map[uid]; !ok {
			members.Map[uid] = &Meta{
				Fullname:     fullname,
				TrelloActive: true,
			}
		}

		// Mark everyone who is in LDAP, those who end up with
		// false are the material to be removed from Trello.
		members.Map[uid].seenInLDAP = true

		if members.Map[uid].TrelloActive && members.Map[uid].TrelloID == "" {

		}
	}

	// TODO: do the stuff!

	log.Println(trello.Search("aarapov@redhat.com"))

	if err := members.Write(); err != nil {
		log.Fatalln(err)
	}
}
