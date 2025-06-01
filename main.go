package main

import (
	"embed"
	"os"

	"github.com/android-sms-gateway/twilio-fallback/internal"
)

//go:generate swag init --parseDependency -g ./main.go -o ./api

//	@title			SMS Gate Twilio Fallback API
//	@version		{{VERSION}}
//	@description	Provides a fallback for Twilio SMS messages via SMS Gate

//	@contact.name	SMSGate Support
//	@contact.email	support@sms-gate.app

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		twilio.sms-gate.app
//	@schemes	https

//go:embed migrations
var migrationsFS embed.FS

func main() {
	cmd := ""
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "":
		fallthrough
	case "run":
		internal.Run()

	case "db:auto-migrate":
		internal.RunORMMigrations()

	case "db:migrate":
		internal.RunMigrations(migrationsFS)

	default:
		panic("unknown command: " + cmd)
	}

}
