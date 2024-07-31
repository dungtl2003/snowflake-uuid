package internal

import (
	"fmt"
	"log"
	"math/big"
	"time"
)

type Worker struct {
	epoch        *big.Int
	datacenterId *big.Int
	workerId     *big.Int

	sequence              *big.Int
	sequenceMask          *big.Int
	timestampLeftShift    *big.Int
	datacenterIdLeftShift *big.Int
	workerIdLeftShift     *big.Int
	lastTimestamp         *big.Int
}

func New(datacenterId *big.Int, workerId *big.Int, epoch *big.Int, datacenterIdBits *big.Int, workerIdBits *big.Int, sequenceBits *big.Int) *Worker {
	if datacenterId == nil || workerId == nil {
		log.Fatalf("data center ID and worker ID cannot be nil")
	}

	w := new(Worker)

	if epoch == nil {
		epoch = big.NewInt(1609459200000)
	}
	if sequenceBits == nil {
		sequenceBits = big.NewInt(12)
	}
	if datacenterIdBits == nil {
		datacenterIdBits = big.NewInt(5)
	}
	if workerIdBits == nil {
		workerIdBits = big.NewInt(5)
	}

	if epoch.Cmp(big.NewInt(0)) == -1 {
		log.Fatalf("epoch time cannot be smaller than 0")
	}

	if sequenceBits.Cmp(big.NewInt(0)) == -1 {
		log.Fatalf("sequence bits cannot be smaller than 0")
	}

	workerIdMaxBits := big.NewInt(0).Xor(big.NewInt(-1), big.NewInt(0).Lsh(big.NewInt(-1), uint(workerIdBits.Uint64())))
	if workerId.Cmp(big.NewInt(0)) == -1 || workerId.Cmp(workerIdMaxBits) == 1 {
		log.Fatalf("worker ID cannot be greater than %s or smaller than 0", workerIdMaxBits.String())
	}

	datacenterIdMaxBits := big.NewInt(0).Xor(big.NewInt(-1), big.NewInt(0).Lsh(big.NewInt(-1), uint(datacenterIdBits.Uint64())))
	if datacenterId.Cmp(big.NewInt(0)) == -1 || datacenterId.Cmp(datacenterIdMaxBits) == 1 {
		log.Fatalf("datacenter ID cannot be greater than %s or smaller than 0", datacenterIdMaxBits.String())
	}

	sequenceMask := big.NewInt(0).Xor(big.NewInt(-1), big.NewInt(0).Lsh(big.NewInt(-1), uint(sequenceBits.Uint64())))
	workerIdLeftShift := sequenceBits
	datacenterIdLeftShift := big.NewInt(0).Add(workerIdBits, sequenceBits)
	timestampLeftShift := big.NewInt(0).Add(datacenterIdBits, big.NewInt(0).Add(workerIdBits, sequenceBits))

	log.Printf("worker configuration:")
	log.Printf("epoch time in milliseconds: %s", epoch.String())
	log.Printf("sequence bits: %s", sequenceBits.String())
	log.Printf("worker ID bits: %s", workerIdMaxBits.String())
	log.Printf("datacenter ID bits: %s", datacenterIdMaxBits.String())

	w.epoch = epoch
	w.datacenterId = datacenterId
	w.workerId = workerId
	w.sequence = big.NewInt(0)
	w.sequenceMask = sequenceMask
	w.workerIdLeftShift = workerIdLeftShift
	w.datacenterIdLeftShift = datacenterIdLeftShift
	w.timestampLeftShift = timestampLeftShift
	w.lastTimestamp = big.NewInt(-1)

	return w
}

func (w *Worker) NextId() (*big.Int, error) {
	var sequenceId *big.Int
	currTimestamp := w.now()

	if currTimestamp.Cmp(w.lastTimestamp) == -1 {
		return nil, fmt.Errorf("cannot generate ID in %v milliseconds", big.NewInt(0).Sub(w.lastTimestamp, currTimestamp).String())
	}

	if w.lastTimestamp == currTimestamp {
		sequenceId = big.NewInt(0).And(big.NewInt(0).Add(w.sequence, big.NewInt(1)), w.sequenceMask)
		if sequenceId.Cmp(big.NewInt(0)) == 0 {
			currTimestamp = w.tilNextMillis() // we cannot just +1 millisecond because it will not syncronize with the clock
		}
	} else {
		sequenceId = big.NewInt(0)
	}

	w.lastTimestamp = currTimestamp
	sequencePart := sequenceId
	workerPart := big.NewInt(0).Lsh(w.workerId, uint(w.workerIdLeftShift.Uint64()))
	datacenterPart := big.NewInt(0).Lsh(w.datacenterId, uint(w.datacenterIdLeftShift.Uint64()))
	timestampPart := big.NewInt(0).Lsh(big.NewInt(0).Sub(currTimestamp, w.epoch), uint(w.timestampLeftShift.Uint64()))

	return big.NewInt(0).Add(big.NewInt(0).Add(timestampPart, datacenterPart), big.NewInt(0).Add(workerPart, sequencePart)), nil
}

func (w *Worker) tilNextMillis() *big.Int {
	var timestamp *big.Int

	for timestamp := w.now(); timestamp.Cmp(w.lastTimestamp) == -1; {
	}

	return timestamp
}

func (w *Worker) now() *big.Int {
	return big.NewInt(time.Now().UnixMilli())
}
