package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

var (
	// DefaultConfigPath : path to the default config file (please do not change it)
	DefaultConfigPath = "config/default_config.yml"
	// ConfigPath : path to the user specific config file (please do not change it)
	ConfigPath = "config/config.yml"
)

// Conf : Modified configuration
var Conf *Config
var privateConf Config

// IsSukebei : Tells if we are on the sukebei website
func IsSukebei() bool {
	return Conf.Models.TorrentsTableName == "sukebei_torrents"
}

// WebAddress : Returns web address for current site
func WebAddress() string {
	if IsSukebei() {
		return Conf.WebAddress.Sukebei
	} else {
		return Conf.WebAddress.Nyaa
	}
}

var allowedDatabaseTypes = map[string]bool{
	"sqlite3":  true,
	"postgres": true,
	"mysql":    true,
	"mssql":    true,
}

var allowedDBLogModes = map[string]bool{
	"default":  true, // errors only
	"detailed": true,
	"silent":   true,
}

func init() {
	Parse()
}

// Parse : Parse config into a config variable
func Parse() {
	getDefaultConfig()
	privateConf = *DefaultConfig
	Conf = &privateConf
	overrideDefaults()
}

func overrideDefaults() {
	data, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		log.Printf("can't read file '%s'", ConfigPath)
	}
	err = yaml.Unmarshal(data, &Conf)
	if err != nil {
		log.Printf("error: %v", err)
	}
}

// BindFlags returns a function which is to be used after
// flag.Parse to check and copy the flags' values to the Config instance.
func (config *Config) BindFlags() func() error {
	confFile := flag.String("conf", "", "path to the configuration file")
	dbType := flag.String("dbtype", Conf.DBType, "database backend")
	host := flag.String("host", Conf.Host, "binding address of the server")
	port := flag.Int("port", Conf.Port, "port of the server")
	dbParams := flag.String("dbparams", Conf.DBParams, "parameters to open the database (see Gorm's doc)")
	dbLogMode := flag.String("dblogmode", Conf.DBLogMode, "database log verbosity (errors only by default)")

	return func() error {
		// You can override fields in the config file with flags.
		config.Host = *host
		config.Port = *port
		config.DBParams = *dbParams
		err := config.SetDBType(*dbType)
		if err != nil {
			return err
		}
		err = config.SetDBLogMode(*dbLogMode)
		if err != nil {
			return err
		}
		err = config.HandleConfFileFlag(*confFile)
		return err
	}
}

// HandleConfFileFlag : Read the config from a file
func (config *Config) HandleConfFileFlag(path string) error {
	if path != "" {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("can't read file '%s'", path)
		}
		err = yaml.Unmarshal(data, config)
		if err != nil {
			return fmt.Errorf("failed to parse file '%s' (%s)", path, err)
		}
	}
	return nil
}

// SetDBType : Set the DataBase type in config
func (config *Config) SetDBType(dbType string) error {
	if !allowedDatabaseTypes[dbType] {
		return fmt.Errorf("unknown database backend '%s'", dbType)
	}
	config.DBType = dbType
	return nil
}

// SetDBLogMode : Set the log mode in config
func (config *Config) SetDBLogMode(dbLogmode string) error {
	if !allowedDBLogModes[dbLogmode] {
		return fmt.Errorf("unknown database log mode '%s'", dbLogmode)
	}
	config.DBLogMode = dbLogmode
	return nil
}

// Read : Decode config from json to config
func (config *Config) Read(input io.Reader) error {
	return json.NewDecoder(input).Decode(config)
}

// Write : Encode config from json to config
func (config *Config) Write(output io.Writer) error {
	return json.NewEncoder(output).Encode(config)
}

// Pretty : Write config json in a file
func (config *Config) Pretty(output io.Writer) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	_, err = output.Write(data)
	return err
}
