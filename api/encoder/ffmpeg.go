package encoder

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
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
	Progress    progress
	cmd         *exec.Cmd
	isCancelled bool
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
	OptionsRaw []string `json:"options_raw"` // Flag options.

	// FFmpeg commander options.
	// TODO: FFmpeg.Run shouold parse and run with struct options.
	// Currently only supports raw options.
	Video videoOptions
	Audio audioOptions
}

type videoOptions struct {
	Input                string
	Output               string
	Container            string
	VideoCodec           string
	VideoSpeed           string
	AudioCodec           string
	HardwareAcceleration string
	Pass                 string
	Crf                  int
	Bitrate              string
	MinRate              string
	MaxRate              string
	BufSize              string
	PixelFormat          string
	FrameRate            string
	Speed                string
	Tune                 string
	Profile              string
	Level                string
}

type audioOptions struct {
	AudioCodec string
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
	options := &ffmpegOptions{}
	if err := json.Unmarshal([]byte(data), &options); err != nil {
		panic(err)
	}

	// Add the list of options from ffmpeg presets.
	for _, v := range options.OptionsRaw {
		args = append(args, strings.Split(v, " ")...)
	}
	args = append(args, output)

	// Execute command.
	log.Info("running FFmpeg with options: ", args)
	f.cmd = exec.Command(ffmpegCmd, args...)
	stdout, _ := f.cmd.StdoutPipe()

	// Capture stderr (if any).
	var stderr bytes.Buffer
	f.cmd.Stderr = &stderr
	f.cmd.Start()

	// Send progress updates.
	go f.trackProgress()

	// Update progress struct.
	f.updateProgress(stdout)

	err := f.cmd.Wait()
	if err != nil {
		if f.isCancelled {
			return errors.New("cancelled")
		}
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

// Cancel stops an FFmpeg job from running.
func (f *FFmpeg) Cancel() {
	log.Warn("killing ffmpeg process")
	f.isCancelled = true
	if err := f.cmd.Process.Kill(); err != nil {
		log.Warn("failed to kill process: ", err)
	}
	log.Warn("killed ffmpeg process")
}

func (f *FFmpeg) finish() {
	close(f.Progress.quit)
}
