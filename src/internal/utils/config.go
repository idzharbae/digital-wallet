package utils

import (
	"fmt"

	"github.com/idzharbae/digital-wallet/src/internal/constants"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(constants.ENV_FILE)
	viper.AddConfigPath(constants.ENV_FILE_DIRECTORY)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error occurred while reading env file, might fallback to OS env config: %s", err.Error())
	}
	viper.AutomaticEnv()
}

// This function can be used to get ENV Var in our App
// Modify this if you change the library to read ENV
func GetEnvVar(name string) string {
	if !viper.IsSet(name) {
		fmt.Printf("Environment variable %s is not set\n", name)
		return ""
	}
	value := viper.GetString(name)
	return value
}
