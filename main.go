//	@title			SMSGate Twilio Fallback API
//	@version		{{VERSION}}
//	@description	Provides a fallback for Twilio SMS messages via SMSGate

//	@contact.name	SMSGate Support
//	@contact.email	support@sms-gate.app

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		[YOUR_SERVICE_URL]
//	@schemes	https

package main

import (
	"runtime"
	"strconv"

	"github.com/android-sms-gateway/twilio-fallback/internal"
	"github.com/go-core-fx/healthfx"
	"github.com/samber/lo"
)

//go:generate swag init --parseDependency --outputTypes go -g ./main.go -o ./internal/server/docs

//nolint:gochecknoglobals // build metadata
var (
	appVersion   = "dev"
	appReleaseID = "0"
	appBuildDate = "unknown"
	appGitCommit = "unknown"
	appGoVersion = runtime.Version()
)

func main() {
	internal.Run(healthfx.Version{
		Version:   appVersion,
		ReleaseID: lo.Must1(strconv.Atoi(appReleaseID)),
		BuildDate: appBuildDate,
		GitCommit: appGitCommit,
		GoVersion: appGoVersion,
	})
}
