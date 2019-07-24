package encoder

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

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

	// Start reading logs.
	go readLog()

	log.Info("running FFmpeg with options: ", args)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(stdout))
}

func readLog() {
	c := time.Tick(10 * time.Second)
	for _ = range c {
		file, err := os.Open("progress-log.txt")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		buf := make([]byte, 62)
		stat, err := os.Stat("progress-log.txt")
		start := stat.Size() - 62
		_, err = file.ReadAt(buf, start)
		if err == nil {
			fmt.Printf("%s\n", buf)
		}
	}
}
