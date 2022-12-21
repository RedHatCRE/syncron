package configuration

import (
	"os"
	"reflect"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Configuration struct {
	S3          S3Configuration `yaml:"s3"`
	Prefix      string          `yaml:"prefix"`
	DownloadDir string          `yaml:"download_dir"`
	SosReports  []string        `yaml:"sosreports"`
	Insights    []string        `yaml:"insights"`
}

type S3Configuration struct {
	Bucket   string `yaml:"bucket"`
	EndPoint string `yaml:"endpoint"`
	Region   string `yaml:"region"`
}

// Setting up file formatting
// Using Viper
// Reading from file syncron.yaml
func (c *Configuration) GetConfiguration() *Configuration {
	userDirConfig, err := os.UserConfigDir()

	if err != nil {
		logrus.Fatal(err)
	}

	viper.SetConfigName("syncron")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(userDirConfig)
	viper.AddConfigPath(".")

	// Reading from file
	err = viper.ReadInConfig()
	if err != nil {
		logrus.Fatal("Following error reading from config file", err)
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		logrus.Fatal("Error structuring configuration", err)
	}

	if c.DownloadDir == "" {
		c.DownloadDir = "/tmp/syncron/"
	}
	checkConfig(*c)

	logrus.Info("Your configuration file was read succesfully")
	logrus.Info("Reading from bucket: ", c.S3.Bucket)

	return c
}

// This function checks a Configuration struct and its embedded S3 struct for empty fields.
// It does this by using the reflect package to get the values and types.
func checkConfig(c Configuration) {

	values := reflect.ValueOf(c)
	valuesS3 := reflect.ValueOf(c.S3)

	//Check if struct that holds configuration is empty
	if reflect.ValueOf(c).IsZero() {
		logrus.Fatal("Configuration struct appears to be empty")
	}

	typesS3 := valuesS3.Type()
	types := values.Type()

	// Check for configuration struct fields
	for i := 0; i < valuesS3.NumField(); i++ {
		if valuesS3.Field(i).IsZero() {
			logrus.Fatal(typesS3.Field(i).Name, " field in config appears to be empty")
		}
	}

	// Check for S3 struct fields
	for i := 0; i < values.NumField(); i++ {
		if values.Field(i).IsZero() {
			logrus.Fatal(types.Field(i).Name, " field in config appears to be empty")
		}
	}

}
