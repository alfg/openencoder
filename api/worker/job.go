package worker

import (
	"path"

	"github.com/alfg/openencoder/api/config"
	"github.com/alfg/openencoder/api/data"
	"github.com/alfg/openencoder/api/encoder"
	"github.com/alfg/openencoder/api/helpers"
	"github.com/alfg/openencoder/api/net"
	"github.com/alfg/openencoder/api/types"
	log "github.com/sirupsen/logrus"
)

func download(job types.Job) error {
	log.Info("running download task")

	// Update status.
	data.UpdateJobStatus(job.GUID, types.JobDownloading)

	d := net.GetDownloadFunc()
	err := d(job)
	if err != nil {
		log.Error(err)
	}
	return err
}

func encode(job types.Job) error {
	log.Info("running encode task")

	// Update status.
	data.UpdateJobStatus(job.GUID, types.JobEncoding)

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

	// Update status.
	data.UpdateJobStatus(job.GUID, types.JobUploading)

	d := net.GetUploadFunc()
	err := d(job)
	if err != nil {
		log.Error(err)
	}
	return err
}

func completed(job types.Job) error {
	log.Info("job completed")

	// Update status.
	data.UpdateJobStatus(job.GUID, types.JobCompleted)

	return nil
}

func runEncodeJob(job types.Job) {
	// Set local src path.
	job.LocalSource = helpers.GetLocalSourcePath(job.Source, job.GUID)

	// 1. Download.
	err := download(job)
	if err != nil {
		log.Error(err)
		data.UpdateJobStatus(job.GUID, types.JobError)
	}

	// 2. Encode.
	err = encode(job)
	if err != nil {
		log.Error(err)
		data.UpdateJobStatus(job.GUID, types.JobError)
	}

	// 3. Upload.
	err = upload(job)
	if err != nil {
		log.Error(err)
		data.UpdateJobStatus(job.GUID, types.JobError)
	}

	// 4. Done
	completed(job)
}
