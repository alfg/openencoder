package worker

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/alfg/openencoder/api/config"
	"github.com/alfg/openencoder/api/data"
	"github.com/alfg/openencoder/api/encoder"
	"github.com/alfg/openencoder/api/helpers"
	"github.com/alfg/openencoder/api/net"
	"github.com/alfg/openencoder/api/notify"
	"github.com/alfg/openencoder/api/types"
	log "github.com/sirupsen/logrus"
)

var progressCh chan struct{}

const progressInterval = time.Second * 5
const slackWebhookKey = "SLACK_WEBHOOK"

func download(job types.Job) error {
	log.Info("running download task")

	// Update status.
	db := data.New()
	db.Jobs.UpdateJobStatusByGUID(job.GUID, types.JobDownloading)

	// Get job data.
	j, _ := db.Jobs.GetJobByGUID(job.GUID)
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
	db.Jobs.UpdateTransferProgressByID(encodeID, 100)
	return err
}

func probe(job types.Job) (*encoder.FFProbeResponse, error) {
	log.Info("running probe task")

	// Update status.
	db := data.New()
	db.Jobs.UpdateJobStatusByGUID(job.GUID, types.JobProbing)

	// Run FFProbe.
	f := encoder.FFProbe{}
	probeData := f.Run(job.LocalSource)

	// Add probe data to DB.
	b, err := json.Marshal(probeData)
	if err != nil {
		log.Error(err)
	}
	j, _ := db.Jobs.GetJobByGUID(job.GUID)
	db.Jobs.UpdateEncodeProbeByID(j.EncodeDataID, string(b))

	return probeData, nil
}

func encode(job types.Job, probeData *encoder.FFProbeResponse) error {
	log.Info("running encode task")

	// Update status.
	db := data.New()
	db.Jobs.UpdateJobStatusByGUID(job.GUID, types.JobEncoding)

	p, err := db.Presets.GetPresetByName(job.Preset)
	if err != nil {
		return err
	}
	dest := path.Dir(job.LocalSource) + "/dst/" + p.Output

	// Get job data.
	j, _ := db.Jobs.GetJobByGUID(job.GUID)

	// Update encode options in DB.
	db.Jobs.UpdateEncodeOptionsByID(j.EncodeDataID, p.Data)

	// Run FFmpeg.
	f := &encoder.FFmpeg{}
	go trackEncodeProgress(j.GUID, j.EncodeDataID, probeData, f)
	err = f.Run(job.LocalSource, dest, p.Data)
	if err != nil {
		close(progressCh)
		return err
	}
	close(progressCh)
	return err
}

func upload(job types.Job) error {
	log.Info("running upload task")

	// Update status.
	db := data.New()
	db.Jobs.UpdateJobStatusByGUID(job.GUID, types.JobUploading)

	// Get job data.
	j, _ := db.Jobs.GetJobByGUID(job.GUID)
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
	db.Jobs.UpdateTransferProgressByID(encodeID, 100)
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
	db := data.New()
	db.Jobs.UpdateJobStatusByGUID(job.GUID, types.JobCompleted)
	return nil
}

func sendAlert(job types.Job) error {
	log.Info("sending alert")

	db := data.New()
	webhook := db.Settings.GetSetting(slackWebhookKey).Value
	message := fmt.Sprintf(
		"*Encode Successful!* :tada:\n"+
			"*Job ID*: %s:\n"+
			"*Preset*: %s\n"+
			"*Source*: %s\n"+
			"*Destination*: %s\n\n",
		job.GUID, job.Preset, job.Source.String, job.Destination.String)
	err := notify.SendSlackMessage(webhook, message)
	if err != nil {
		return err
	}
	return nil
}

func runEncodeJob(job types.Job) {
	// Set local src path.
	job.LocalSource = helpers.CreateLocalSourcePath(
		config.Get().WorkDirectory, job.Source.String, job.GUID)

	db := data.New()

	// 1. Download.
	err := download(job)
	if err != nil {
		log.Error(err)
		db.Jobs.UpdateJobStatusByGUID(job.GUID, types.JobError)
		return
	}

	// 2. Probe data.
	probeData, err := probe(job)
	if err != nil {
		log.Error(err)
		db.Jobs.UpdateJobStatusByGUID(job.GUID, types.JobError)
		return
	}

	// 3. Encode.
	err = encode(job, probeData)
	if err != nil {
		log.Error(err)
		if err := cleanup(job); err != nil {
			log.Error("cleanup err", err)
		}

		// Set job to 'cancelled' if it was cancelled.
		if err.Error() == types.JobCancelled {
			db.Jobs.UpdateJobStatusByGUID(job.GUID, types.JobCancelled)
			return
		}
		db.Jobs.UpdateJobStatusByGUID(job.GUID, types.JobError)
		return
	}

	// 4. Upload.
	err = upload(job)
	if err != nil {
		log.Error(err)
		db.Jobs.UpdateJobStatusByGUID(job.GUID, types.JobError)
		return
	}

	// 5. Cleanup.
	err = cleanup(job)
	if err != nil {
		log.Error(err)
		db.Jobs.UpdateJobStatusByGUID(job.GUID, types.JobError)
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

func trackEncodeProgress(guid string, encodeID int64, p *encoder.FFProbeResponse, f *encoder.FFmpeg) {
	db := data.New()
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
			speed := f.Progress.Speed
			fps := f.Progress.FPS

			// Check cancel.
			status, _ := db.Jobs.GetJobStatusByGUID(guid)
			if status == types.JobCancelled {
				f.Cancel()
			}

			// Only track progress if we know the total frames.
			if totalFrames != 0 {
				pct := (float64(currentFrame) / float64(totalFrames)) * 100

				// Update DB with progress.
				pct = math.Round(pct*100) / 100
				fmt.Printf("progress: %d / %d - %0.2f%%\r", currentFrame, totalFrames, pct)
				db.Jobs.UpdateEncodeProgressByID(encodeID, pct, speed, fps)
			}
		}
	}
}

func trackTransferProgress(encodeID int64, d *net.S3) {
	db := data.New()
	progressCh = make(chan struct{})
	ticker := time.NewTicker(progressInterval)

	for {
		select {
		case <-progressCh:
			ticker.Stop()
			return
		case <-ticker.C:
			fmt.Println("transfer progress: ", d.Progress.Progress)
			db.Jobs.UpdateTransferProgressByID(encodeID, float64(d.Progress.Progress))
		}
	}
}
