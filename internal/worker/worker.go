package worker

import (
	"fmt"
	"sync"
	"time"
)

type Worker struct {
	epoch    int64
	workerId int64

	sequence           int64
	sequenceMask       int64
	workerPart         int64
	timestampLeftShift int64
	lastTimestamp      int64
	maxTimestamp       int64

	mu sync.Mutex
}

func New(workerId int64, epoch int64) (*Worker, error) {
	var workerMaxId int64 = -1 ^ (-1 << 10)
	if workerId < 0 || workerId > workerMaxId {
		return nil, fmt.Errorf("worker ID cannot be smaller than 0 or bigger than %d", workerMaxId)
	}

	if epoch < 0 {
		return nil, fmt.Errorf("epoch time cannot be smaller than 0")
	}

	w := new(Worker)

	w.epoch = epoch
	w.workerId = workerId
	w.sequence = 0
	w.sequenceMask = -1 ^ (-1 << 12)
	w.workerPart = workerId << 12
	w.timestampLeftShift = 22
	w.lastTimestamp = -1
	w.maxTimestamp = -1 ^ (-1 << 41)

	return w, nil
}

func (w *Worker) NextId() (int64, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	var sequenceId int64
	currTimestamp := time.Now().UnixMilli()

	if currTimestamp < w.lastTimestamp {
		return -1, fmt.Errorf("cannot generate ID in %d milliseconds", w.lastTimestamp-currTimestamp)
	}

	if currTimestamp > w.maxTimestamp {
		return -1, fmt.Errorf("this generator cannot work anymore because it has passed maximum timestamp of %d", w.maxTimestamp)
	}

	if w.lastTimestamp == currTimestamp {
		sequenceId = (w.sequence + 1) & w.sequenceMask
		if sequenceId == 0 {
			for currTimestamp <= w.lastTimestamp {
				currTimestamp = time.Now().UnixMilli()
			}
		}
	} else {
		sequenceId = 0
	}

	w.sequence = sequenceId
	w.lastTimestamp = currTimestamp
	sequencePart := sequenceId
	timestampPart := (currTimestamp - w.epoch) << w.timestampLeftShift

	return timestampPart + w.workerPart + sequencePart, nil
}
