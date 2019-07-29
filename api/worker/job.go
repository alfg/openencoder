package worker

import (
	"encoding/json"
	"math"
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

	// Add probe data to DB.
	b, err := json.Marshal(probeData)
	if err != nil {
		log.Error(err)
	}
	j, _ := data.GetJobByGUID(job.GUID)
	data.UpdateEncodeDataByID(j.EncodeDataID, string(b))

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

	// Get job data.
	j, _ := data.GetJobByGUID(job.GUID)
	encodeID := j.EncodeDataID

	// Run FFmpeg.
	f := &encoder.FFmpeg{}
	go trackProgress(encodeID, probeData, f)
	f.Run(job.LocalSource, dest, p.Options)
	close(quit)

	// Set encode progress to 100.
	data.UpdateEncodeProgressByID(encodeID, 100)

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

func trackProgress(encodeID int64, p *encoder.FFProbeResponse, f *encoder.FFmpeg) {
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

			// Update DB with progress.
			pct = math.Round(pct*100) / 100
			log.Infof("progress: %d / %d - %0.2f%%\n", currentFrame, totalFrames, pct)
			data.UpdateEncodeProgressByID(encodeID, pct)
		}
	}
}
