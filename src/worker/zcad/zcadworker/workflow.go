package zcadworker

import (
	"encoding/base64"
	"time"

	"github.com/bytedance/sonic"
	"go.temporal.io/sdk/workflow"
)

type ZCAD_LoadFileParams struct {
	Files []string `json:"files"`
}

func ScheduleWorkflow(ctx workflow.Context, token string, parameters string) error {
	// Apply the options.
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Second * 5,
		ScheduleToCloseTimeout: time.Second * 20,
		StartToCloseTimeout:    time.Second * 5,
	})

	futures := []workflow.Future{}

	dec, err := base64.StdEncoding.DecodeString(parameters)
	if err != nil {
		return err
	}

	parameters = string(dec)

	LoadFileParams := ZCAD_LoadFileParams{}
	if err = sonic.Unmarshal([]byte(parameters), &LoadFileParams); err != nil {
		return err
	}

	// get all activity futures into futures variant
	for _, file := range LoadFileParams.Files {
		futures = append(futures, workflow.ExecuteActivity(ctx, ZCAD_LoadFile, file))
	}

	// wait all futures return.
	for _, future := range futures {
		var res ZCAD_LoadFileResult
		if err = future.Get(ctx, &res); err != nil || res.Status != "Success" {
			log.Error("ZCAD_LoadFile failed.", err)
		} else {
			log.Infof("ZCAD_LoadFile %s success.", res.File)
		}
	}

	return nil
}
