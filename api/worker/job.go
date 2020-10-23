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
)

var progressCh chan struct{}

func generatePresignedURL(job types.Job) (string, error) {
	log.Info("generating a presigned URL")

	// Update status.
	db := data.New()
	db.Jobs.UpdateJobStatusByGUID(job.GUID, types.JobDownloading)

	// Get presigned URL.
	str, err := net.GetPresignedURL(job)
	if err != nil {
		log.Error(err)
	}
	return str, nil
}

func download(job types.Job, storageDriver string) error {
	log.Info("running download task for: ", storageDriver)

	// Update status.
	db := data.New()
	db.Jobs.UpdateJobStatusByGUID(job.GUID, types.JobDownloading)

	// Get job data.
	j, err := db.Jobs.GetJobByGUID(job.GUID)
	if err != nil {
		log.Error(err)
		return err
	}
	encodeID := j.EncodeID

	if err := net.Download(job); err != nil {
		log.Error(err)
		return err
	}

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
	probeData := f.Run(job.Source)

	// Add probe data to DB.
	b, err := json.Marshal(probeData)
	if err != nil {
		log.Error(err)
	}
	j, _ := db.Jobs.GetJobByGUID(job.GUID)
	db.Jobs.UpdateEncodeProbeByID(j.EncodeID, string(b))

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
	db.Jobs.UpdateEncodeOptionsByID(j.EncodeID, p.Data)

	// Run FFmpeg.
	f := &encoder.FFmpeg{}
	go trackEncodeProgress(j.GUID, j.EncodeID, probeData, f)
	err = f.Run(job.Source, dest, p.Data)
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
	j, err := db.Jobs.GetJobByGUID(job.GUID)
	if err != nil {
		log.Error(err)
		return err
	}
	encodeID := j.EncodeID

	if err := net.Upload(job); err != nil {
		log.Error(err)
		return err
	}

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
	webhook, err := db.Settings.GetSetting(types.SlackWebhook)
	if err != nil {
		return err
	}

	// Only send Slack alert if configured.
	if webhook.Value != "" {
		message := fmt.Sprintf(AlertMessageFormat, job.GUID, job.Preset, job.Source, job.Destination)
		err := notify.SendSlackMessage(webhook.Value, message)
		if err != nil {
			return err
		}
	}
	return nil
}

func runEncodeJob(job types.Job) {
	// Set local src path.
	job.LocalSource = helpers.CreateLocalSourcePath(
		config.Get().WorkDirectory, job.Source, job.GUID)

	db := data.New()
	storageDriver, err := db.Settings.GetSetting(types.StorageDriver)
	if err != nil {
		log.Error(err)
		return
	}

	// If STREAMING setting is enabled, get a presigned URL and update
	// the job.Source.
	s3Streaming, err := db.Settings.GetSetting(types.S3Streaming)
	if err != nil {
		log.Error(err)
		return
	}

	if s3Streaming.Value == "enabled" && storageDriver.Value == "s3" {
		// 1a. Get presigned URL.
		presigned, err := generatePresignedURL(job)
		if err != nil {
			log.Error(err)
			db.Jobs.UpdateJobStatusByGUID(job.GUID, types.JobError)
			return
		}

		// Update source with presigned URL.
		job.Source = presigned

	} else {
		// 1b. Download.
		err := download(job, storageDriver.Value)
		if err != nil {
			log.Error(err)
			db.Jobs.UpdateJobStatusByGUID(job.GUID, types.JobError)
			return
		}

		job.Source = job.LocalSource
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
	ticker := time.NewTicker(ProgressInterval)

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
				log.Infof("progress: %d / %d - %0.2f%%", currentFrame, totalFrames, pct)
				db.Jobs.UpdateEncodeProgressByID(encodeID, pct, speed, fps)
			}
		}
	}
}
