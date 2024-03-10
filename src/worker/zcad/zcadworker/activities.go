package zcadworker

import (
	"context"
	"transform2/worker/zcad/libzcad"
)

type ZCAD_LoadFileResult struct {
	File   string
	Status string
}

func ZCAD_LoadFile(ctx context.Context, file string) (*ZCAD_LoadFileResult, error) {
	log.Infof("ZCAD_LoadFile %s", file)
	libzcad.Hello(file)
	return &ZCAD_LoadFileResult{file, "Success"}, nil
}
