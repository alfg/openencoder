package helpers

import (
	"errors"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

// C is a config instance available as a public config object.
var C config

type config struct {
	Settings  settings
	WorkFlows []workflow
	Ffmpeg    []ffmpegProfile
	S3        []s3Profile
}

type settings struct {
	S3Bucket string `yaml:"s3_bucket"`
	S3Region string `yaml:"s3_region"`
}

type workflow struct {
	Name  string
	Tasks []string
}

type ffmpegProfile struct {
	Profile string
	Output  string
	Publish bool
	Options []string
}

type s3Profile struct {
	Profile string
	Output  string
	Options []string
}

func LoadConfig(file string) {
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println("yamlfile err: ", err)
	}

	err = yaml.Unmarshal(yamlFile, &C)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

// func LoadConfig() {
// 	err := viper.Unmarshal(&C)
// 	if err != nil {
// 		panic(fmt.Errorf("unable to decode config into struct, %v", err))
// 	}
// }

// func GetSettings() (t *settings, err error) {

// }

func GetWorkflow(name string) (t *workflow, err error) {
	for _, v := range C.WorkFlows {
		if v.Name == name {
			return &v, nil
		}
	}
	return nil, errors.New("No task")
}

func GetFFmpegProfile(profile string) (t *ffmpegProfile, err error) {
	for _, v := range C.Ffmpeg {
		if v.Profile == profile {
			return &v, nil
		}
	}
	return nil, errors.New("No task")
}

func GetS3Profile(profile string) (t *s3Profile, err error) {
	for _, v := range C.S3 {
		if v.Profile == profile {
			return &v, nil
		}
	}
	return nil, errors.New("No task")
}
