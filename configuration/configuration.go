package configuration

import (
	"log"
	"os"

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
		logrus.Fatal(err)
	}

	err = viper.Unmarshal(&c)

	if err != nil {
		log.Fatal(err)
	}

	if c.DownloadDir == "" {
		c.DownloadDir = "/tmp/syncron/"
	}

	logrus.Info("Your configuration file was read succesfully")
	logrus.Info("Reading from bucket: ", c.S3.Bucket)

	return c
}
