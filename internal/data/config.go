package data

import (
	"gopkg.in/yaml.v2"
)

var Conf *Config
var db *gorm.DB

type Config struct {
	JibeUser     string `yaml:"jibeUser"`
	JibeHost     string `yaml:"jibeHost"`
	JibeDB       string `yaml:"jibeDB"`
	JibePassword string `yaml:"jibePassword"`
}

func InitConfig() {
	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		log.Panic(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &Conf)
	if err != nil {
		log.Panic(err.Error())
	}
}

func InitDB() {
	dbStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		Conf.JibeHost,
		5432,
		Conf.JibeUser,
		Conf.JibeDB,
		Conf.JibePassword)

	var err error
	db, err = sql.Open("postgres", dbStr)
	if err != nil {
		log.Panic(err)
	}
	if err = db.Ping(); err != nil {
		log.Panic(err)
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
}

func CloseDB() {
	db.Close()
}
