package worker_test

import (
	"dungtl2003/snowflake-uuid/internal/worker"
	"testing"
)

func TestCreateWorker(t *testing.T) {
	_, err := worker.New(1, 1722613500)
	if err != nil {
		t.Fatalf("failed to create worker: %v", err)
	}

	_, err = worker.New(300, 1532213569)
	if err != nil {
		t.Fatalf("failed to create worker: %v", err)
	}
}

func TestGenerateUniqueId(t *testing.T) {
	w, _ := worker.New(1, 1722613500)
	idsSet := make(map[int64]int)

	for i := 0; i < 1000000; i++ {
		id, err := w.NextId()
		if err != nil {
			t.Fatalf("error: %v", err)
		}

		if _, ok := idsSet[id]; ok {
			t.Fatalf("error: generated duplicated ID %d", id)
		}

		idsSet[id] = 1
	}
}

func TestSequentialExecution(t *testing.T) {
	w, _ := worker.New(1, 1722613500)
	idsSet := make(map[int64]int)
	c := make(chan int64)
	errors := make(chan error)

	go func() {
		defer close(c)
		defer close(errors)

		for i := 0; i < 1000000; i++ {
			id, err := w.NextId()

			if err != nil {
				errors <- err
				return
			} else {
				c <- id
			}
		}
	}()

	for id := range c {
		if _, ok := idsSet[id]; ok {
			t.Fatalf("error: generated duplicated ID %d", id)
		}

		idsSet[id] = 1
	}

	if err := <-errors; err != nil {
		t.Fatalf("error: %v", err)
	}
}

// check if there is any problems of creating multiple workers when a worker is generating IDs
func TestRace(t *testing.T) {
	w, _ := worker.New(1, 1722613500)

	go func() {
		for i := 0; i < 10000000; i++ {
			worker.New(1, 1722613500)
		}
	}()

	for i := 0; i < 100000; i++ {
		w.NextId()
	}
}

// BENCHMARK

func BenchmarkNextId(b *testing.B) {
	w, _ := worker.New(1, 1722613500)

	b.ReportAllocs()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		w.NextId()
	}
}
