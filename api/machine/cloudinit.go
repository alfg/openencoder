package machine

import (
	"bytes"
	"html/template"
	"log"

	"github.com/alfg/openencoder/api/config"
)

const userDataTmpl = `
#cloud-config
package_upgrade: false
write_files:
    - path: "/opt/.env"
      content: |
        AWS_REGION={{.AWSRegion}}
        AWS_ACCESS_KEY={{.AWSAccessKey}}
        AWS_SECRET_KEY={{.AWSSecretKey}}
        SLACK_WEBHOOK={{.SlackWebhook}}
        REDIS_HOST={{.CloudinitRedisHost}}
        REDIS_PORT={{.CloudinitRedisPort}}
        DATABASE_HOST={{.CloudinitDatabaseHost}}
        DATABASE_PORT={{.CloudinitDatabasePort}}
        DATABASE_USER={{.CloudinitDatabaseUser}}
        DATABASE_PASSWORD={{.CloudinitDatabasePassword}}
        DATABASE_NAME={{.CloudinitDatabaseName}}
runcmd:
  - docker run -d --env-file /opt/.env --rm alfg/openencoder:latest worker
`

// UserData defines the userdata used for cloud-init.
type UserData struct {
	AWSRegion    string
	AWSAccessKey string
	AWSSecretKey string
	SlackWebhook string

	CloudinitRedisHost        string
	CloudinitRedisPort        int
	CloudinitDatabaseHost     string
	CloudinitDatabasePort     int
	CloudinitDatabaseUser     string
	CloudinitDatabasePassword string
	CloudinitDatabaseName     string
}

func createUserData() string {
	data := &UserData{
		AWSRegion:    config.Get().AWSRegion,
		AWSAccessKey: config.Get().AWSAccessKey,
		AWSSecretKey: config.Get().AWSSecretKey,
		SlackWebhook: config.Get().SlackWebhook,

		CloudinitRedisHost:        config.Get().CloudinitRedisHost,
		CloudinitRedisPort:        config.Get().CloudinitRedisPort,
		CloudinitDatabaseHost:     config.Get().CloudinitDatabaseHost,
		CloudinitDatabasePort:     config.Get().CloudinitDatabasePort,
		CloudinitDatabaseUser:     config.Get().CloudinitDatabaseUser,
		CloudinitDatabasePassword: config.Get().CloudinitDatabasePassword,
		CloudinitDatabaseName:     config.Get().CloudinitDatabaseName,
	}

	var tpl bytes.Buffer
	t := template.Must(template.New("userdata").Parse(userDataTmpl))

	if err := t.Execute(&tpl, data); err != nil {
		log.Println(err)
	}

	result := tpl.String()
	return result
}
