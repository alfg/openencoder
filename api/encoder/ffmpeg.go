package encoder

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

// FFMPEG bin
const FFMPEG = "ffmpeg"

// FFmpeg struct.
type FFmpeg struct{}

// Run runs an encoder process with options.
func (f FFmpeg) Run(input string, output string, options []string) {
	args := []string{
		"-hide_banner",
		"-v", "0",
		"-progress", "progress-log.txt",
		"-i", input,
	}

	// Add the list of options from ffmpeg profile.
	for _, v := range options {
		args = append(args, strings.Split(v, " ")...)
	}
	args = append(args, output)

	// Execute command.
	cmd := exec.Command(FFMPEG, args...)

	log.Info("running FFmpeg with options: ", args)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(stdout))
}

// Remux content.
func (f FFmpeg) Remux(video string, audio string, output string) {
	cmd := exec.Command(FFMPEG, "-i", video, "-i", audio, "-c", "copy", "-map", "0:0", "-map", "1:0", output)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(stdout))
}

// Demux demuxes a single MP4 into a video and audio MP4.
func (f FFmpeg) Demux(source string, dest string) {

	createDir(dest)

	// Split video.
	videoOutput := dest + "/video.mp4"
	cmd := exec.Command(FFMPEG, "-i", source, "-an", "-c", "copy", videoOutput)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(stdout))

	audioOutput := dest + "/audio.mp4"
	// Split audio.
	cmd = exec.Command(FFMPEG, "-i", source, "-map", "0:1", "-c", "copy", audioOutput)
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(stdout))
}

// Encode runs an encoder process with progress log.
func (f FFmpeg) Encode(input string, output string) {
	cmd := exec.Command(FFMPEG, "-progress", "progress-log.txt", "-i", input, "-b", "1000000", output)

	fmt.Println("Encoding...")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(stdout))
}

// FFmpegBuildConf outputs the FFmpeg build configuration.
func (f FFmpeg) FFmpegBuildConf() {
	cmd := exec.Command(FFMPEG, "-buildconf")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(stdout))
}

func createDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
}

func buildOptions() {

}
