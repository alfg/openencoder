package worker

import (
	"path"

	"github.com/alfg/openencoder/api/config"
	"github.com/alfg/openencoder/api/encoder"
	"github.com/alfg/openencoder/api/helpers"
	"github.com/alfg/openencoder/api/net"
	"github.com/alfg/openencoder/api/types"
	log "github.com/sirupsen/logrus"
)

func download(job types.Job) error {
	log.Info("running download task")

	d := net.GetDownloadFunc()
	err := d(job)
	if err != nil {
		log.Error(err)
	}
	return err
}

func encode(job types.Job) error {
	log.Info("running encode task")

	p, err := config.GetFFmpegProfile(job.Profile)
	if err != nil {
		return err
	}
	dest := path.Dir(job.LocalSource) + "/dst/" + p.Output

	// Run FFmpeg.
	f := encoder.FFmpeg{}
	f.Run(job.LocalSource, dest, p.Options)

	return err
}

func upload(job types.Job) error {
	log.Info("running upload task")

	d := net.GetUploadFunc()
	err := d(job)
	if err != nil {
		log.Error(err)
	}
	return err
}

func runEncodeJob(job types.Job) {
	// Set local src path.
	job.LocalSource = helpers.GetLocalSourcePath(job.Source, job.GUID)

	// 1. Download.
	download(job)

	// 2. Encode.
	encode(job)

	// 3. Upload.
	upload(job)
}
