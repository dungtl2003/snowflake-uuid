package tools

import (
	"errors"
	"math/big"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Configuration struct {
	Port             int
	DatacenterId     *big.Int
	WorkerId         *big.Int
	Epoch            *big.Int
	DatacenterIdBits *big.Int
	WorkerIdBits     *big.Int
	SequenceBits     *big.Int
}

func (config *Configuration) Load(environment string) error {
	var filePath string
	switch environment {
	case "development":
		filePath = "environments/dev/.env"
		break
	case "production":
		filePath = "environments/prod/.env"
		break
	default:
		return errors.New("unsupported environment")
	}

	err := godotenv.Load(filePath)
	if err != nil {
		return err
	}

	port, err := getPort()
	if err != nil {
		return err
	}

	datacenterId, err := getDatacenterId()
	if err != nil {
		return err
	}

	workerId, err := getWorkerId()
	if err != nil {
		return err
	}

	epoch, err := getEpoch()
	if err != nil {
		return err
	}

	datacenterIdBits, err := getDatacenterIdBits()
	if err != nil {
		return err
	}

	workerIdBits, err := getWorkerIdBits()
	if err != nil {
		return err
	}

	sequenceBits, err := getSequenceBits()
	if err != nil {
		return err
	}

	config.Port = port
	config.DatacenterId = datacenterId
	config.WorkerId = workerId
	config.Epoch = epoch
	config.DatacenterIdBits = datacenterIdBits
	config.WorkerIdBits = workerIdBits
	config.SequenceBits = sequenceBits

	return nil
}

func getSequenceBits() (*big.Int, error) {
	raw := os.Getenv("SEQUENCE_BITS")
	if raw == "" {
		return nil, nil
	}

	return convertToBigint(raw)
}

func getWorkerIdBits() (*big.Int, error) {
	raw := os.Getenv("WORKER_ID_BITS")
	if raw == "" {
		return nil, nil
	}

	return convertToBigint(raw)
}

func getDatacenterIdBits() (*big.Int, error) {
	raw := os.Getenv("DATACENTER_ID_BITS")
	if raw == "" {
		return nil, nil
	}

	return convertToBigint(raw)
}

func getEpoch() (*big.Int, error) {
	raw := os.Getenv("EPOCH")
	if raw == "" {
		return nil, nil
	}

	return convertToBigint(raw)
}

func getWorkerId() (*big.Int, error) {
	raw := os.Getenv("WORKER_ID")
	if raw == "" {
		return nil, errors.New("missing WORKER_ID")
	}

	return convertToBigint(raw)
}

func getDatacenterId() (*big.Int, error) {
	raw := os.Getenv("DATACENTER_ID")
	if raw == "" {
		return nil, errors.New("missing DATACENTER_ID")
	}

	return convertToBigint(raw)
}

func getPort() (int, error) {
	raw := os.Getenv("PORT")
	if raw == "" {
		return 9000, nil
	}

	return strconv.Atoi(raw)
}

func convertToBigint(str string) (*big.Int, error) {
	num, err := strconv.Atoi(str)
	if err != nil {
		return nil, err
	}

	return big.NewInt(int64(num)), nil
}
