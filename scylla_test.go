package gocql

import (
	"runtime"
	"sync"
	"testing"
)

func TestScyllaConnPickerPickNilToken(t *testing.T) {
	t.Parallel()

	s := scyllaConnPicker{
		nrShards:  4,
		msbIgnore: 12,
	}

	t.Run("no conns", func(t *testing.T) {
		s.conns = []*Conn{{}}
		if s.Pick(token(nil)) != s.conns[0] {
			t.Fatal("expected connection")
		}
	})

	t.Run("one shard", func(t *testing.T) {
		s.conns = []*Conn{{}}
		if s.Pick(token(nil)) != s.conns[0] {
			t.Fatal("expected connection")
		}
	})

	t.Run("multiple shards", func(t *testing.T) {
		s.conns = []*Conn{nil, {}}
		if s.Pick(token(nil)) != s.conns[1] {
			t.Fatal("expected connection")
		}
		if s.Pick(token(nil)) != s.conns[1] {
			t.Fatal("expected connection")
		}
	})

	t.Run("multiple shards no conns", func(t *testing.T) {
		s.conns = []*Conn{nil, nil}
		if s.Pick(token(nil)) != nil {
			t.Fatal("expected nil")
		}
		if s.Pick(token(nil)) != nil {
			t.Fatal("expected nil")
		}
	})
}

func hammerConnPicker(t *testing.T, wg *sync.WaitGroup, s *scyllaConnPicker, loops int) {
	t.Helper()
	for i := 0; i < loops; i++ {
		if c := s.Pick(nil); c == nil {
			t.Error("unexpected nil")
		}
	}
	wg.Done()
}

func TestScyllaConnPickerHammerPickNilToken(t *testing.T) {
	t.Parallel()

	s := scyllaConnPicker{
		nrShards:  4,
		msbIgnore: 12,
	}
	s.conns = make([]*Conn, 100)
	for i := range s.conns {
		if i%7 == 0 {
			continue
		}
		s.conns[i] = &Conn{}
	}

	n := runtime.GOMAXPROCS(0)
	loops := 10000 / n

	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go hammerConnPicker(t, &wg, &s, loops)
	}
	wg.Wait()
}

func TestScyllaConnPickerShardOf(t *testing.T) {
	t.Parallel()

	s := scyllaConnPicker{
		nrShards:  4,
		msbIgnore: 12,
	}
	for _, test := range scyllaShardOfTests {
		if shard := s.shardOf(murmur3Token(test.token)); shard != test.shard {
			t.Errorf("wrong scylla shard calculated for token %d, expected %d, got %d", test.token, test.shard, shard)
		}
	}
}
