// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultFileName = "config"

	defaultAppWorkingDir     = "~/.platform.d"
	defaultAppDefaultProfile = ""
	defaultAppDefaultOutput  = "human"
	defaultAppPager          = true

	defaultLogLevel             = "INFO"
	defaultLogFileJson          = false
	defaultLogConsoleJson       = false
	defaultLogFileEnabled       = false
	defaultLogTimestampTimezone = "utc"

	defaultTerminalNoColor           = false
	defaultTerminalTimestampTimezone = "utc"
)

type Config struct {
	// Application settings
	WorkingDir     string `json:"working_dir"`
	DefaultProfile string `json:"default_profile"`
	DefaultOutput  string `json:"default_output"`
	Pager          bool   `json:"pager"`

	// Profiles
	profileName string
	profiles    map[string]*Profile

	// Repositories
	repositories map[string]*Repository

	// Log settings
	LogLevel             string         `json:"log_level"`
	LogFileJSON          bool           `json:"log_file_json"`
	LogConsoleJSON       bool           `json:"log_console_json"`
	LogFileEnabled       bool           `json:"log_file_enabled"`
	LogTimestampTimezone *time.Location `json:"log_timestamp_timezone"`

	// Terminal settings
	TerminalNoColor           bool           `json:"terminal_no_color"`
	TerminalTimestampTimezone *time.Location `json:"terminal_timestamp_timezone"`

	// Mongo settings
	MongoUri string `json:"mongo_uri"`
}

func NewConfig(defaults map[string]interface{}, envBinding map[string]string, appWorkingDir, sysConfigPath, fileName string) *Config {
	ac := new(Config)

	if defaults == nil {
		defaults = defaultValues
	}

	if envBinding == nil {
		envBinding = defaultEnvVarBindings
	}

	if appWorkingDir == "" {
		appWorkingDir = defaultAppWorkingDir
	}

	if fileName == "" {
		fileName = defaultFileName
	}

	ac.initConfig(defaults, envBinding, appWorkingDir, sysConfigPath, fileName)

	return ac
}

func (ac *Config) DumpConfig() string {
	bs, _ := json.Marshal(ac)
	return string(bs)
}

func (ac *Config) initConfig(defaultsVariables map[string]interface{}, environmentBindings map[string]string, appWorkingDir, sysConfigPath, fileName string) {
	// Set the default values within the application.
	for k, v := range defaultsVariables {
		viper.SetDefault(k, v)
	}

	// Set the EnvVar binding.
	for k, v := range environmentBindings {
		err := viper.BindEnv(k, v)
		if err != nil {
			// todo: once cleaned up we can come back and get this situated
			fmt.Println("binding error")
		}
	}
	// Now we will set the config file
	setConfigFile(appWorkingDir, sysConfigPath, fileName)
	ac.populateFields()
	var err error

	// @kevin told me about this issue and I didn't copy it over. :facepalm:
	ac.WorkingDir, err = homedir.Expand(ac.WorkingDir)
	if err != nil {
		handleError("", err)
	}

	ac.profiles = map[string]*Profile{}
	ac.repositories = map[string]*Repository{}

	var defaults map[string]interface{}

	if value, exists := viper.AllSettings()["profile default"]; exists {
		defaults = value.(map[string]interface{})
	}

	ac.profiles["default"] = loadProfile(defaults, defaults, map[string]interface{}{})

	for key, value := range viper.AllSettings() {
		if strings.HasPrefix(key, "repository ") {
			parts := strings.Split(key, " ")

			if len(parts) > 2 {
				handleError("repository names cannot contain spaces", nil)
			}

			var overrides = map[string]interface{}{}

			for _, ele := range getRepositoryFields() {
				if val, exists := os.LookupEnv(fmt.Sprintf("IPCTL_REPOSITORY_%s_%s", strings.ToUpper(parts[1]), strings.ToUpper(ele))); exists {
					overrides[ele] = val
				}
			}

			ac.repositories[parts[1]] = loadRepository(value.(map[string]any), overrides)

		} else if strings.HasPrefix(key, "profile ") {
			parts := strings.Split(key, " ")

			if len(parts) > 2 {
				handleError("profile names cannot contain spaces", nil)
			} else if parts[1] == "default" {
				continue
			}

			var overrides = map[string]interface{}{}

			for _, ele := range getProfileFields() {
				if val, exists := os.LookupEnv(fmt.Sprintf("IPCTL_PROFILE_%s_%s", strings.ToUpper(parts[1]), strings.ToUpper(ele))); exists {
					overrides[ele] = val
				}
			}

			ac.profiles[parts[1]] = loadProfile(value.(map[string]any), defaults, overrides)
		}
	}

	ac.profileName = getProfileFromFlag()
	if ac.profileName == "" {
		ac.profileName = ac.DefaultProfile
	}
}

func (ac *Config) populateFields() {
	ac.WorkingDir = GetAndExpandDirectory("application.working_dir")
	ac.DefaultProfile = viper.GetString("application.default_profile")
	ac.Pager = viper.GetBool("application.pager")
	ac.DefaultOutput = viper.GetString("application.default_output")

	ac.LogLevel = viper.GetString("log.level")
	ac.LogFileJSON = viper.GetBool("log.file_json")
	ac.LogConsoleJSON = viper.GetBool("log.console_json")
	ac.LogFileEnabled = viper.GetBool("log.file_enabled")
	ac.LogTimestampTimezone = getTzLocation("log.timestamp_timezone")

	ac.TerminalNoColor = viper.GetBool("terminal.no_color")
	ac.TerminalTimestampTimezone = getTzLocation("terminal.timestamp_timezone")

	ac.MongoUri = viper.GetString("mongo.uri")
}

var defaultValues = map[string]interface{}{
	"application.working_dir":     defaultAppWorkingDir,
	"application.default_profile": defaultAppDefaultProfile,
	"application.default_output":  defaultAppDefaultOutput,
	"application.pager":           defaultAppPager,

	"log.level":              defaultLogLevel,
	"log.file_json":          defaultLogFileJson,
	"log.console_json":       defaultLogConsoleJson,
	"log.file_enabled":       defaultLogFileEnabled,
	"log.timestamp_timezone": defaultLogTimestampTimezone,

	"terminal.no_color":           defaultTerminalNoColor,
	"terminal.timestamp_timezone": defaultTerminalTimestampTimezone,
}

var defaultEnvVarBindings = map[string]string{
	"application.working_dir":     "IPCTL_APPLICATION_WORKING_DIR",
	"application.default_profile": "IPCTL_APPLICATION_DEFAULT_PROFILE",
	"applicationd.default_output": "IPCTL_APPLICATION_DEFAULT_OUTPUT",
	"application.pager":           "IPCTL_APPLICATION_PAGER",

	"log.level":              "IPCTL_LOG_LEVEL",
	"log.file_json":          "IPCTL_LOG_FILE_JSON",
	"log.console_json":       "IPCTL_LOG_CONSOLE_JSON",
	"log.file_enabled":       "IPCTL_LOG_FILE_ENABLED",
	"log.timestamp_timezone": "IPCTL_LOG_TIMESTAMP_TIMEZONE",

	"terminal.no_color":           "IPCTL_TERMINAL_NO_COLOR",
	"terminal.timestamp_timezone": "IPCTL_TERMINAL_TIMESTAMP_TIMEZONE",

	"mongo.uri": "IPCTL_MONGO_URI",
}

// getConfigFileFromFlag reads in the file passed in using the --config flag on the cli
func getConfigFileFromFlag() string {
	var configPath string

	cfgFlags := pflag.NewFlagSet("configFlag", pflag.ContinueOnError)
	cfgFlags.StringVar(&configPath, "config", "", "Path to config file")
	cfgFlags.ParseErrorsWhitelist.UnknownFlags = true // Ignore unknown flags
	cfgFlags.Usage = func() {}                        // Suppress default usage message

	// Non Nil ignore due to result empty slice
	if err := cfgFlags.Parse(os.Args[1:]); err != nil && err != pflag.ErrHelp {
		handleError("failed to parse config command line argument", err)
	}

	if value := os.Getenv("IPCTL_CONFIG_FILE"); value != "" {
		configPath = value
	}

	if configPath == "" {
		return "" // No config path provided
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		handleError(fmt.Sprintf("The path '%s' provided by the --config flag does not exist", configPath), err)
	}

	expanded, err := homedir.Expand(configPath)
	if err != nil {
		handleError("failed to expand the config path", err)
	}

	return expanded
}

func setConfigFile(appWorkingDir, sysConfigPath, fileName string) {
	viper.SetConfigName(fileName)
	viper.SetConfigType("ini")

	envConfFile := os.Getenv("IPCTL_CONFIG")
	if envConfFile != "" {
		viper.SetConfigFile(envConfFile)
	}

	confFile := getConfigFileFromFlag()
	if confFile != "" {
		viper.SetConfigFile(confFile)
	}

	expandedWorkingDir, err := homedir.Expand(appWorkingDir)
	if err != nil {
		handleError(fmt.Sprintf("An error occurred while attempting to find the configuration directory at %s", appWorkingDir), err)
	}
	viper.AddConfigPath(expandedWorkingDir)

	viper.AddConfigPath(sysConfigPath) // defaults to /etc/torero

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			handleError("An error occurred while attempting to read the configuration file", err)
		}
	}
}

// GetAndExpandDirectory will grab a directory from viper's store and expand its value if the directory contains a
// ~ to denote a home directory
func GetAndExpandDirectory(configName string) string {
	cfgVal := viper.GetString(configName)
	expandedVal, err := homedir.Expand(cfgVal)
	if err != nil {
		handleError(fmt.Sprintf("failed to expand directory for configuration variable %s with value %s", defaultEnvVarBindings[configName], cfgVal), err)
	}
	return expandedVal
}

// parseTz correctly formats the capitalization of a timezone so that it can be ingested by the time package's
// time.LoadLocation function. If there are any issues, an error log is displayed and timestamps will default to UTC
func getTzLocation(configName string) *time.Location {
	tz := viper.GetString(configName)
	var parsedTz string
	switch strings.ToLower(tz) {
	case "utc":
		parsedTz = "UTC"
	case "local":
		parsedTz = "Local"
	default:
		parsedTz = tz
	}

	location, err := time.LoadLocation(parsedTz)
	if err != nil { // Log a comment to not cause errors when running `torero completion xyz` in a shell sourcing script
		fmt.Fprintf(os.Stderr, "# Error: failed to load the config variable %s. "+
			"The provided timezone of '%s' does not match either utc, local, or a valid tz identifier. "+
			"Defautling to utc: %s\n\n", configName, tz, err)

		location = time.UTC
	}
	return location
}

func handleError(message string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s: %s\n", message, err)
	} else {
		fmt.Fprintf(os.Stderr, "Error: %s\n", message)
	}
	os.Exit(1)
}

type Option func()
