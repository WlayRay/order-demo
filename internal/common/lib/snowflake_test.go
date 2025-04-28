package lib_test

import (
	"fmt"
	"hash/fnv"
	"sync"
	"testing"
	"time"

	"github.com/WlayRay/order-demo/common/lib"
)

func HashString(str string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(str))
	return h.Sum64()
}

func TestSnowflake(t *testing.T) {
	fmt.Printf("The current timestamp %d\n", time.Now().UnixMilli())

	ip, err := lib.GetLocalIP()
	if err != nil {
		t.Fatalf("Failed to get local IP: %v", err)
	}
	workerID := HashString(ip)
	fmt.Printf("ip: %s\tWorker ID: %d\n", ip, workerID)

	instance, err := lib.GetSnowflakeInstance(workerID%1024, 10*time.Millisecond)
	if err != nil {
		t.Fatalf("Failed to create instance: %v", err)
	}

	idSet := make(map[uint64]struct{}, 5000)

	p, c := 20, 10
	ch := make(chan uint64, 250)
	pwg := &sync.WaitGroup{}
	cwg := &sync.WaitGroup{}
	mtx := &sync.Mutex{}

	pwg.Add(p)
	for range p {
		go func() {
			defer pwg.Done()
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Panic occurred: %v", r)
				}
			}()

			for range 250 {
				id, err := instance.GetID()
				if err != nil {
					t.Errorf("Failed to get ID: %v", err)
				}
				ch <- id
			}
		}()
	}

	cwg.Add(c)
	for range c {
		go func() {
			defer cwg.Done()
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Panic occurred: %v", r)
				}
			}()

			for id := range ch {
				mtx.Lock()
				if _, ok := idSet[id]; ok {
					t.Errorf("Duplicate ID found: %d", id)
				}
				idSet[id] = struct{}{}
				mtx.Unlock()
			}
		}()
	}

	pwg.Wait()
	close(ch)
	cwg.Wait()

	fmt.Printf("Total IDs generated: %d\n", len(idSet))

	i := 0
	for k := range idSet {
		fmt.Printf("ID: %d\n", k)
		i++
		if i > 5 {
			break
		}
	}
}
