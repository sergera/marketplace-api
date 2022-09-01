package conf

import (
	"log"
	"sync"

	"github.com/gurkankaymak/hocon"
)

var once sync.Once
var instance *conf

type conf struct {
	hocon           *hocon.Config
	Port            string
	LogPath         string
	DBHost          string
	DBPort          string
	DBName          string
	DBUser          string
	DBPassword      string
	CORSAllowedURLs string
}

func GetConf() *conf {
	once.Do(func() {
		var c *conf = &conf{}
		c.setup()
		instance = c
	})
	return instance
}

func (c *conf) setup() {
	c.parseHOCONConfigFile()
	c.setPort()
	c.setDBHost()
	c.setDBPort()
	c.setDBName()
	c.setDBUser()
	c.setDBPassword()
	c.setCORSAllowedURLs()
}

func (c *conf) parseHOCONConfigFile() {
	hocon, err := hocon.ParseResource("application.conf")
	if err != nil {
		log.Panic("error while parsing configuration file: ", err)
	}

	log.Printf("configurations: %+v", *hocon)

	c.hocon = hocon
}

func (c *conf) setPort() {
	port := c.hocon.GetString("host.port")
	if len(port) == 0 {
		log.Panic("port environment variable not found")
	}

	c.Port = port
}

func (c *conf) setDBHost() {
	dbHost := c.hocon.GetString("db.host")
	if len(dbHost) == 0 {
		log.Panic("database host environment variable not found")
	}

	c.DBHost = dbHost
}

func (c *conf) setDBPort() {
	dbPort := c.hocon.GetString("db.port")
	if len(dbPort) == 0 {
		log.Panic("database port environment variable not found")
	}

	c.DBPort = dbPort
}

func (c *conf) setDBName() {
	dbName := c.hocon.GetString("db.name")

	if len(dbName) == 0 {
		log.Panic("database name environment variable not found")
	}

	c.DBName = dbName
}

func (c *conf) setDBUser() {
	dbUser := c.hocon.GetString("db.user")
	if len(dbUser) == 0 {
		log.Panic("database user environment variable not found")
	}

	c.DBUser = dbUser
}

func (c *conf) setDBPassword() {
	dbPassword := c.hocon.GetString("db.password")
	if len(dbPassword) == 0 {
		log.Panic("database password environment variable not found")
	}

	c.DBPassword = dbPassword
}

func (c *conf) setCORSAllowedURLs() {
	corsAllowedURLs := c.hocon.GetString("cors.urls")
	if len(corsAllowedURLs) == 0 {
		log.Panic("cors allowed urls environment variable not found")
	}

	c.CORSAllowedURLs = corsAllowedURLs
}
