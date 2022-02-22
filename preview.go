package main

import (
	"bytes"
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func videoFirstFrame(name string) []byte {
	return readFrameAsJpeg(name, 1)
}

func imagePreview(name string) []byte {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(name).
		Output("pipe:", ffmpeg.KwArgs{"format": "image2", "vcodec": "mjpeg", "q:v": "4", "vf": "scale=(-1):400"}).
		WithOutput(buf).
		Run()
	if err != nil {
		fmt.Println(err)
	}
	return buf.Bytes()
}

func gifPreview(name string) []byte {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(name).
		Output("pipe:", ffmpeg.KwArgs{"format": "gif", "vf": "fps=8,scale=(-1):200"}).
		WithOutput(buf).
		Run()
	if err != nil {
		fmt.Println(err)
	}
	return buf.Bytes()
}

func videoPreview(name string) []byte {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(name, ffmpeg.KwArgs{"t": "15"}).
		Output("pipe:", ffmpeg.KwArgs{"format": "gif", "vf": "fps=8,scale=(-1):200"}).
		WithOutput(buf).
		Run()
	if err != nil {
		fmt.Println(err)
	}
	return buf.Bytes()
}

func readFrameAsJpeg(inFileName string, frameNum int) []byte {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf).
		Run()
	if err != nil {
		fmt.Println(err)
	}
	return buf.Bytes()
}
