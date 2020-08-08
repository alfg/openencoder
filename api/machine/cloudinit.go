package machine

import (
	"bytes"
	"html/template"

	"github.com/alfg/openencoder/api/config"
)

const userDataTmpl = `
#cloud-config
package_upgrade: false
write_files:
    - path: "/opt/.env"
      content: |
        REDIS_HOST={{.CloudinitRedisHost}}
        REDIS_PORT={{.CloudinitRedisPort}}
        DATABASE_HOST={{.CloudinitDatabaseHost}}
        DATABASE_PORT={{.CloudinitDatabasePort}}
        DATABASE_USER={{.CloudinitDatabaseUser}}
        DATABASE_PASSWORD={{.CloudinitDatabasePassword}}
        DATABASE_NAME={{.CloudinitDatabaseName}}
runcmd:
  - docker run -d --env-file /opt/.env --rm {{.CloudinitWorkerImage}} worker
`

// UserData defines the userdata used for cloud-init.
type UserData struct {
	CloudinitRedisHost        string
	CloudinitRedisPort        int
	CloudinitDatabaseHost     string
	CloudinitDatabasePort     int
	CloudinitDatabaseUser     string
	CloudinitDatabasePassword string
	CloudinitDatabaseName     string
	CloudinitWorkerImage      string
}

func createUserData() string {
	data := &UserData{
		CloudinitRedisHost:        config.Get().CloudinitRedisHost,
		CloudinitRedisPort:        config.Get().CloudinitRedisPort,
		CloudinitDatabaseHost:     config.Get().CloudinitDatabaseHost,
		CloudinitDatabasePort:     config.Get().CloudinitDatabasePort,
		CloudinitDatabaseUser:     config.Get().CloudinitDatabaseUser,
		CloudinitDatabasePassword: config.Get().CloudinitDatabasePassword,
		CloudinitDatabaseName:     config.Get().CloudinitDatabaseName,
		CloudinitWorkerImage:      config.Get().CloudinitWorkerImage,
	}

	var tpl bytes.Buffer
	t := template.Must(template.New("userdata").Parse(userDataTmpl))

	if err := t.Execute(&tpl, data); err != nil {
		log.Println(err)
	}

	result := tpl.String()
	return result
}
