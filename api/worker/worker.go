package worker

import "time"

// Worker constants.
const (
	ProgressInterval = time.Second * 5
)

// Worker variables.
var (
	AlertMessageFormat = `
*Encode Successful!* :tada:\n
"*Job ID*: %s:\n"
"*Preset*: %s\n"
"*Source*: %s\n"
"*Destination*: %s\n\n"
`
)
