package main

import (
	"github.com/eaigner/hood"
)

type FfprobeFormatRecord struct {
	Id         hood.Id `json:"id"`
	FileId     int32   `json:"file_id"` // FOREIGN KEY
	Duration   float64 `json:"duration"`
	FormatName string  `json:"format_name"` // separate table too?
	Size       int64   `json:"size"`
	// add tags?
	// add streams?
}

type FfprobeStreamRecord struct {
	Id        hood.Id `json:"id"`
	FileId    int32   `json:"file_id"` // FOREIGN KEY
	CodecName string  `json:"codec_name"`
	CodecType string  `json:"codec_type"`
	Width     int32   `json:"width"`
	Height    int32   `json:"height"`
	Duration  float64 `json:"duration"`
}
