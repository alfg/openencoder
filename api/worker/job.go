package worker

import (
	"fmt"
	"path"
	"strconv"
	"time"

	"github.com/alfg/openencoder/api/config"
	"github.com/alfg/openencoder/api/data"
	"github.com/alfg/openencoder/api/encoder"
	"github.com/alfg/openencoder/api/helpers"
	"github.com/alfg/openencoder/api/net"
	"github.com/alfg/openencoder/api/types"
	log "github.com/sirupsen/logrus"
)

var quit chan struct{}

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

func probe(job types.Job) (*encoder.FFProbeResponse, error) {
	log.Info("running probe task")

	// Update status.
	data.UpdateJobStatus(job.GUID, types.JobProbing)

	// Run FFProbe.
	f := encoder.FFProbe{}
	probeData := f.Run(job.LocalSource)
	return probeData, nil
}

func encode(job types.Job, probeData *encoder.FFProbeResponse) error {
	log.Info("running encode task")

	// Update status.
	data.UpdateJobStatus(job.GUID, types.JobEncoding)

	p, err := config.GetFFmpegProfile(job.Profile)
	if err != nil {
		return err
	}
	dest := path.Dir(job.LocalSource) + "/dst/" + p.Output

	// Run FFmpeg.
	f := &encoder.FFmpeg{}
	go trackProgress(probeData, f)
	f.Run(job.LocalSource, dest, p.Options)
	close(quit)

	// Set encode progress to 100.

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

	// 2. Probe data.
	probeData, err := probe(job)
	if err != nil {
		log.Error(err)
		data.UpdateJobStatus(job.GUID, types.JobError)
	}

	// 3. Encode.
	err = encode(job, probeData)
	if err != nil {
		log.Error(err)
		data.UpdateJobStatus(job.GUID, types.JobError)
	}

	// 4. Upload.
	err = upload(job)
	if err != nil {
		log.Error(err)
		data.UpdateJobStatus(job.GUID, types.JobError)
	}

	// 5. Done
	completed(job)
}

func trackProgress(p *encoder.FFProbeResponse, f *encoder.FFmpeg) {
	quit = make(chan struct{})
	ticker := time.NewTicker(time.Second * 2)

	for {
		select {
		case <-quit:
			ticker.Stop()
			return
		case <-ticker.C:
			currentFrame := f.Progress.Frame
			totalFrames, _ := strconv.Atoi(p.Streams[0].NbFrames)

			pct := (float64(currentFrame) / float64(totalFrames)) * 100

			fmt.Println("progress", currentFrame, totalFrames)
			fmt.Printf("%0.2f\n", pct)

			// TODO: Update DB with progress.
		}
	}
}
