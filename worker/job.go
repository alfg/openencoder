package worker

import (
	"path"

	"github.com/alfg/enc/encoder"
	"github.com/alfg/enc/helpers"
	"github.com/alfg/enc/net"
	"github.com/alfg/enc/types"
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

	p, err := helpers.GetFFmpegProfile(job.Profile)
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
	job.LocalSource = helpers.GetLocalSourcePath(job.Source, job.ID)

	// 1. Download.
	download(job)

	// 2. Encode.
	encode(job)

	// 3. Upload.
	upload(job)
}
