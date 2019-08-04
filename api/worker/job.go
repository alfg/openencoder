package worker

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/alfg/openencoder/api/alert"
	"github.com/alfg/openencoder/api/config"
	"github.com/alfg/openencoder/api/data"
	"github.com/alfg/openencoder/api/encoder"
	"github.com/alfg/openencoder/api/helpers"
	"github.com/alfg/openencoder/api/net"
	"github.com/alfg/openencoder/api/types"
	log "github.com/sirupsen/logrus"
)

var progressCh chan struct{}

const progressInterval = time.Second * 2

func download(job types.Job) error {
	log.Info("running download task")

	// Update status.
	data.UpdateJobStatus(job.GUID, types.JobDownloading)

	// Get job data.
	j, _ := data.GetJobByGUID(job.GUID)
	encodeID := j.EncodeDataID

	// Get downloader type.
	d := net.GetDownloader()

	// Do download and track progress.
	go trackTransferProgress(encodeID, d)
	err := d.S3Download(job)
	if err != nil {
		log.Error(err)
	}

	// Close channel to stop progress updates.
	close(progressCh)

	// Set progress to 100.
	data.UpdateEncodeProgressByID(encodeID, 100)
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
	go trackEncodeProgress(encodeID, probeData, f)
	f.Run(job.LocalSource, dest, p.Options)
	close(progressCh)

	// Set encode progress to 100.
	data.UpdateEncodeProgressByID(encodeID, 100)
	return err
}

func upload(job types.Job) error {
	log.Info("running upload task")

	// Update status.
	data.UpdateJobStatus(job.GUID, types.JobUploading)

	// Get job data.
	j, _ := data.GetJobByGUID(job.GUID)
	encodeID := j.EncodeDataID

	d := net.GetUploader()

	// Do download and track progress.
	go trackTransferProgress(encodeID, d)
	err := d.S3Upload(job)
	if err != nil {
		log.Error(err)
	}

	// Close channel to stop progress updates.
	close(progressCh)

	// Set progress to 100.
	data.UpdateEncodeProgressByID(encodeID, 100)
	return err
}

func cleanup(job types.Job) error {
	log.Info("running cleanup task")

	tmpPath := helpers.GetTmpPath(config.Get().WorkDirectory, job.GUID)
	err := os.RemoveAll(tmpPath)
	if err != nil {
		return err
	}
	return nil
}

func completed(job types.Job) error {
	log.Info("job completed")

	// Update status.
	data.UpdateJobStatus(job.GUID, types.JobCompleted)
	return nil
}

func sendAlert(job types.Job) error {
	log.Info("sending alert")

	message := fmt.Sprintf(
		"*Encode Successful!* :tada:\n"+
			"*Job ID*: %s:\n"+
			"*Profile*: %s\n"+
			"*Source*: %s\n"+
			"*Destination*: %s\n\n",
		job.GUID, job.Profile, job.Source, job.Destination)
	err := alert.SendSlackMessage(config.Get().SlackWebhook, message)
	if err != nil {
		return err
	}
	return nil
}

func runEncodeJob(job types.Job) {
	// Set local src path.
	job.LocalSource = helpers.CreateLocalSourcePath(
		config.Get().WorkDirectory, job.Source, job.GUID)

	// 1. Download.
	err := download(job)
	if err != nil {
		log.Error(err)
		data.UpdateJobStatus(job.GUID, types.JobError)
		return
	}

	// 2. Probe data.
	probeData, err := probe(job)
	if err != nil {
		log.Error(err)
		data.UpdateJobStatus(job.GUID, types.JobError)
		return
	}

	// 3. Encode.
	err = encode(job, probeData)
	if err != nil {
		log.Error(err)
		data.UpdateJobStatus(job.GUID, types.JobError)
		return
	}

	// 4. Upload.
	err = upload(job)
	if err != nil {
		log.Error(err)
		data.UpdateJobStatus(job.GUID, types.JobError)
		return
	}

	// 5. Cleanup.
	err = cleanup(job)
	if err != nil {
		log.Error(err)
		data.UpdateJobStatus(job.GUID, types.JobError)
		return
	}

	// 6. Done
	completed(job)
	if err != nil {
		log.Error(err)
	}

	// 7. Alert
	sendAlert(job)
	if err != nil {
		log.Error(err)
	}
}

func trackEncodeProgress(encodeID int64, p *encoder.FFProbeResponse, f *encoder.FFmpeg) {
	progressCh = make(chan struct{})
	ticker := time.NewTicker(progressInterval)

	for {
		select {
		case <-progressCh:
			ticker.Stop()
			return
		case <-ticker.C:
			currentFrame := f.Progress.Frame
			totalFrames, _ := strconv.Atoi(p.Streams[0].NbFrames)

			pct := (float64(currentFrame) / float64(totalFrames)) * 100

			// Update DB with progress.
			pct = math.Round(pct*100) / 100
			fmt.Printf("progress: %d / %d - %0.2f%%\r", currentFrame, totalFrames, pct)
			data.UpdateEncodeProgressByID(encodeID, pct)
		}
	}
}

func trackTransferProgress(encodeID int64, d *net.S3) {
	progressCh = make(chan struct{})
	ticker := time.NewTicker(progressInterval)

	for {
		select {
		case <-progressCh:
			ticker.Stop()
			return
		case <-ticker.C:
			fmt.Println("transfer progress: ", d.Progress.Progress)
			data.UpdateEncodeProgressByID(encodeID, float64(d.Progress.Progress))
		}
	}
}
