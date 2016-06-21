package main

var exampleStartupJSON string = `{
	"bind": "127.0.0.1",
	"port": 80,
	"database_host": "couchdb.local",
	"database_port": 5984,
	"database_user": "waingro",
	"database_pass": "getiton" 
}`

var newConfigurationJSON string = `{
	"_id": "configuration",
	"type": "configuration",
	"web_user": "vincent",
	"web_pass": "hanna",
	"sd_user": "schedules",
	"sd_pass": "direct",
	"sd_countrycode": "US",
	"sd_postalcode": "12345"
}`

type newConfiguration struct {
	id        string `json:_id`
	rev       string `json:_rev`
	typeField string `json:type`

	webUser       string `json:web_user`
	webPass       string `json:web_pass`
	sdUser        string `json:sd_user`
	sdPass        string `json:sd_pass`
	sdCountryCode string `json:sd_countrycode`
	sdPostalCode  string `json:sd_postalcode`
}

func databaseInitialize() {
	// noop
}

func databaseCheckConstitency() {
	// noop
}
