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
        AWS_REGION="{{.AWSRegion}}"
        AWS_ACCESS_KEY="{{.AWSAccessKey}}"
        AWS_SECRET_KEY="{{.AWSSecretKey}}"
        REDIS_HOST="{{.RedisHost}}"
        REDIS_PORT={{.RedisPort}}
        DATABASE_HOST="{{.DatabaseHost}}"
        DATABASE_PORT={{.DatabasePort}}
        DATABASE_USER="{{.DatabaseUser}}"
        DATABASE_PASSWORD="{{.DatabasePassword}}"
        DATABASE_NAME="{{.DatabaseName}}"
runcmd:
  - docker run -d --env-file /opt/.env --rm alfg/openencoder worker
`

// UserData defines the userdata used for cloud-init.
type UserData struct {
	AWSRegion        string
	AWSAccessKey     string
	AWSSecretKey     string
	RedisHost        string
	RedisPort        int
	DatabaseHost     string
	DatabasePort     int
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	SlackWebhook     string
}

func createUserData() string {
	data := &UserData{
		AWSRegion:        config.Get().AWSRegion,
		AWSAccessKey:     config.Get().AWSAccessKey,
		AWSSecretKey:     config.Get().AWSSecretKey,
		RedisHost:        config.Get().RedisHost,
		RedisPort:        config.Get().RedisPort,
		DatabaseHost:     config.Get().DatabaseHost,
		DatabasePort:     config.Get().DatabasePort,
		DatabaseUser:     config.Get().DatabaseUser,
		DatabasePassword: config.Get().DatabasePassword,
		DatabaseName:     config.Get().DatabaseName,
		SlackWebhook:     config.Get().SlackWebhook,
	}

	var tpl bytes.Buffer
	t := template.Must(template.New("userdata").Parse(userDataTmpl))

	if err := t.Execute(&tpl, data); err != nil {
		log.Println(err)
	}

	result := tpl.String()
	return result
}
