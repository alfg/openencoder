package encoder

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

const ffprobeCmd = "ffprobe"

// FFProbe struct.
type FFProbe struct{}

// Run runs an FFProbe command.
func (f FFProbe) Run(input string) *FFProbeResponse {
	args := []string{
		"-i", input,
		"-show_streams",
		"-print_format", "json",
		"-v", "quiet",
	}

	fmt.Println(args)

	// Execute command.
	cmd := exec.Command(ffprobeCmd, args...)
	log.Info("Running FFprobe...")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err.Error())
	}
	// log.Info((string(stdout))

	// j, _ := json.Marshal(string(stdout))
	dat := &FFProbeResponse{}
	if err := json.Unmarshal([]byte(stdout), &dat); err != nil {
		panic(err)
	}
	return dat
}

type FFProbeResponse struct {
	Streams []Stream `json:"streams"`
}

type Stream struct {
	Index              int         `json:"index"`
	CodecName          string      `json:"codec_name"`
	CodecLongName      string      `json:"codec_long_name"`
	Profile            string      `json:"profile"`
	CodecType          string      `json:"codec_type"`
	CodecTimeBase      string      `json:"codec_time_base"`
	CodecTagString     string      `json:"codec_tag_string"`
	CodecTag           string      `json:"codec_tag"`
	Width              int         `json:"width"`
	Height             int         `json:"height"`
	CodedWidth         int         `json:"coded_width"`
	CodedHeight        int         `json:"coded_height"`
	HasBFrames         int         `json:"has_b_frames"`
	SampleAspectRatio  string      `json:"sample_aspect_ratio"`
	DisplayAspectRatio string      `json:"display_aspect_ratio"`
	PixFmt             string      `json:"pix_fmt"`
	Level              int         `json:"level"`
	ChromaLocation     string      `json:"chroma_location"`
	Refs               int         `json:"refs"`
	IsAVC              string      `json:"is_avc"`
	NalLengthSize      string      `json:"nal_length_size"`
	RFrameRate         string      `json:"r_frame_rate"`
	AvgFrameRate       string      `json:"avg_frame_rate"`
	TimeBase           string      `json:"time_base"`
	StartPts           int         `json:"start_pts"`
	StartTime          string      `json:"start_time"`
	DurationTS         int         `json:"duration_ts"`
	Duration           string      `json:"duration"`
	BitRate            string      `json:"bit_rate"`
	BitsPerRawSample   string      `json:"bits_per_raw_sample"`
	NbFrames           string      `json:"nb_frames"`
	Disposition        Disposition `json:"disposition"`
	Tags               Tags        `json:"tags"`
}

type Disposition struct {
	Default         int `json:"default"`
	Dub             int `json:"dub"`
	Original        int `json:"original"`
	Comment         int `json:"comment"`
	Lyrics          int `json:"lyrics"`
	Karoake         int `json:"karaoke"`
	Forced          int `json:"forced"`
	HearingImpaired int `json:"hearing_impaired"`
	VisualImpaired  int `json:"visual_empaired"`
	CleanEffects    int `json:"clean_effects"`
	AttachedPic     int `json:"attached_pic"`
	TimedThumbnails int `json:"timed_thumbnails"`
}

type Tags struct {
	Language    string `json:"language"`
	HandlerName string `json:"handler_name"`
}
