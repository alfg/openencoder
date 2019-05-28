package worker

import (
	"path"

	"github.com/alfg/enc/encoder"
	"github.com/alfg/enc/helpers"
	"github.com/alfg/enc/net"
	"github.com/alfg/enc/types"
	log "github.com/sirupsen/logrus"
)

// func runWorkflow(job types.Job) {
// 	wf, err := helpers.GetWorkflow(job.Profile)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	job.LocalSource = helpers.GetLocalSourcePath(job.Source, job.ID)

// 	// Run through tasks.
// 	for _, v := range wf.Tasks {
// 		tasks := strings.Split(v, ".")
// 		name := tasks[0]
// 		task := tasks[1]

// 		switch name {
// 		case "ffmpeg":
// 			runFfmpegTask(task, job)

// 		case "s3":
// 			runS3Task(task, job)
// 		}
// 	}
// }

func runFfmpegTask(task string, job types.Job) {
	log.Info("running ffmpeg task")
	p, err := helpers.GetFFmpegProfile(task)
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

func runS3Task(task string, job types.Job) {
	log.Info("running s3 task")

	p, err := helpers.GetS3Profile(task)
	if err != nil {
		return
	}

	if p.Profile == "download" {
		d := net.GetDownloadFunc()
		err := d(job)
		if err != nil {
			log.Error(err)
		}
	} else if p.Profile == "upload" {
		d := net.GetUploadFunc()
		err := d(job)
		if err != nil {
			log.Error(err)
		}
	}
}
