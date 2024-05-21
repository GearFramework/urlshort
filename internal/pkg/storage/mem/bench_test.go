package mem

import (
	"context"
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func BenchmarkMemory(b *testing.B) {
	inMemStore := NewStorage()
	if err := inMemStore.InitStorage(); err != nil {
		log.Fatalln(err.Error())
	}
	defer inMemStore.Close()
	rand.New(rand.NewSource(time.Now().UnixNano()))
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		code := strconv.FormatInt(rand.Int63n(1_000_000_000_000_000_000), 10)
		if err := inMemStore.Insert(ctx, 1, code, code); err != nil {
			log.Println(err.Error())
		}
	}
}
