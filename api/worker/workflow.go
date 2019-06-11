package worker

import (
	"path"

	"github.com/alfg/enc/api/config"
	"github.com/alfg/enc/api/encoder"
	"github.com/alfg/enc/api/types"
	log "github.com/sirupsen/logrus"
)

func runFfmpegTask(task string, job types.Job) {
	log.Info("running ffmpeg task")
	p, err := config.GetFFmpegProfile(task)
	if err != nil {
		return
	}
	dest := path.Dir(job.LocalSource) + "/"

	// Set output to publish or not.
	var out string
	if !p.Publish {
		out = dest + "/src/"
	} else {
		out = dest + "/dst/"
	}

	f := encoder.FFmpeg{}
	f.Run(job.LocalSource, out+p.Output, p.Options)
}
