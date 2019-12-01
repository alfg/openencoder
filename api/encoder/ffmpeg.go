package encoder

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	ffmpegCmd      = "ffmpeg"
	updateInterval = time.Second * 5
)

// FFmpeg struct.
type FFmpeg struct {
	Progress progress
}

type progress struct {
	quit       chan struct{}
	Frame      int
	FPS        float64
	Bitrate    float64
	TotalSize  int
	OutTimeMS  int
	OutTime    string
	DupFrames  int
	DropFrames int
	Speed      string
	Progress   float64
}

// ffmpegOptions struct passed into Ffmpeg.Run.
type ffmpegOptions struct {
	Options []string `json:"options"`
}

// Run runs the ffmpeg encoder with options.
func (f *FFmpeg) Run(input string, output string, data string) error {
	args := []string{
		"-hide_banner",
		"-loglevel", "error", // Set loglevel to fail job on errors.
		"-progress", "pipe:1",
		"-i", input,
	}

	// Decode JSON get options list from data.
	dat := &ffmpegOptions{}
	if err := json.Unmarshal([]byte(data), &dat); err != nil {
		panic(err)
	}

	// Add the list of options from ffmpeg presets.
	for _, v := range dat.Options {
		args = append(args, strings.Split(v, " ")...)
	}
	args = append(args, output)

	// Execute command.
	log.Info("running FFmpeg with options: ", args)
	cmd := exec.Command(ffmpegCmd, args...)
	stdout, _ := cmd.StdoutPipe()

	// Capture stderr (if any).
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Start()

	// Send progress updates.
	go f.trackProgress()

	// Update progress struct.
	f.updateProgress(stdout)

	err := cmd.Wait()
	if err != nil {
		fmt.Println(stderr.String())
		f.finish()
		return err
	}
	f.finish()
	return nil
}

func (f *FFmpeg) updateProgress(stdout io.ReadCloser) {
	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		line := scanner.Text()
		str := strings.Replace(line, " ", "", -1)

		parts := strings.Split(str, " ")
		f.setProgressParts(parts)
	}
}

func (f *FFmpeg) setProgressParts(parts []string) {
	for i := 0; i < len(parts); i++ {
		progressSplit := strings.Split(parts[i], "=")
		k := progressSplit[0]
		v := progressSplit[1]

		switch k {
		case "frame":
			frame, _ := strconv.Atoi(v)
			f.Progress.Frame = frame
		case "fps":
			fps, _ := strconv.ParseFloat(v, 64)
			f.Progress.FPS = fps
		case "bitrate":
			v = strings.Replace(v, "kbits/s", "", -1)
			bitrate, _ := strconv.ParseFloat(v, 64)
			f.Progress.Bitrate = bitrate
		case "total_size":
			size, _ := strconv.Atoi(v)
			f.Progress.TotalSize = size
		case "out_time_ms":
			outTimeMS, _ := strconv.Atoi(v)
			f.Progress.OutTimeMS = outTimeMS
		case "out_time":
			f.Progress.OutTime = v
		case "dup_frames":
			frames, _ := strconv.Atoi(v)
			f.Progress.DupFrames = frames
		case "drop_frames":
			frames, _ := strconv.Atoi(v)
			f.Progress.DropFrames = frames
		case "speed":
			f.Progress.Speed = v
		case "progress":
			progress, _ := strconv.ParseFloat(v, 64)
			f.Progress.Progress = progress
		}
	}
}

func (f *FFmpeg) trackProgress() {
	f.Progress.quit = make(chan struct{})
	ticker := time.NewTicker(updateInterval)

	for {
		select {
		case <-f.Progress.quit:
			ticker.Stop()
			return
		case <-ticker.C:
			// fmt.Println(f.Progress)
		}
	}
}

func (f *FFmpeg) finish() {
	close(f.Progress.quit)
}
