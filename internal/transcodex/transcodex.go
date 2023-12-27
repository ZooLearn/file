package transcodex

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func CreateHLS(inputFile string, fileName, outputDir string, segmentDuration int) error {
	// Create the output directory if it does not exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}
	// TODO: optimize ffmpeg command line
	// TODO: add a check to see if the output file already exists
	// TODO: add a check to see if the input file exists
	// TODO: make resolutions configurable
	// Create the HLS playlist and segment the video using ffmpeg

	ffmpegCmd := exec.Command(
		"ffmpeg",
		"-i", inputFile,
		"-map", "0:v:0", "-map", "0:a:0", "-map", "0:v:0", "-map", "0:a:0", "-map", "0:v:0", "-map", "0:a:0",
		"-c:v", "libx264", "-crf", "22", "-c:a", "aac", "-ar", "48000",
		"-filter:v:0", "scale=w=480:h=360", "-maxrate:v:0", "600k", "-b:a:0", "64k",
		"-filter:v:1", "scale=w=640:h=480", "-maxrate:v:1", "900k", "-b:a:1", "128k",
		"-filter:v:2", "scale=w=1280:h=720", "-maxrate:v:2", "900k", "-b:a:2", "128k",
		"-var_stream_map", "v:0,a:0,name:480p v:1,a:1,name:640p v:2,a:2,name:1280p",
		"-preset", "slow", "-hls_list_size", "0", "-threads", "0", "-f", "hls",
		"-hls_playlist_type", "event", "-hls_time", strconv.Itoa(segmentDuration),
		"-hls_flags", "independent_segments", "-master_pl_name", fmt.Sprintf("%s-master-playlist.m3u8", fileName),
		fmt.Sprintf("%s/%s-%%v.m3u8", outputDir, fileName),
	)
	output, err := ffmpegCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create HLS: %v\nOutput: %s", err, string(output))
	}
	return nil
}
