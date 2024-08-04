package tools

import (
	"errors"
	"os"
	"strconv"
)

type Configuration struct {
	WorkerId int64
	Epoch    int64
}

func (config *Configuration) Load() error {
	workerId, err := getWorkerId()
	if err != nil {
		return err
	}

	epoch, err := getEpoch()
	if err != nil {
		return err
	}

	config.WorkerId = workerId
	config.Epoch = epoch

	return nil
}

func getEpoch() (int64, error) {
	raw := os.Getenv("EPOCH")
	if raw == "" {
		return 1704067200, nil
	}

	return strconv.ParseInt(raw, 10, 64)
}

func getWorkerId() (int64, error) {
	raw := os.Getenv("WORKER_ID")
	if raw == "" {
		return -1, errors.New("missing WORKER_ID")
	}

	return strconv.ParseInt(raw, 10, 64)
}
