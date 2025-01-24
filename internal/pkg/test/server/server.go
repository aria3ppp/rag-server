package server

import (
	"sync"
)

type TestServerFunc func(Fatalizer) (port int, cleanup func())

func SetupServers(f Fatalizer, testcontainers map[string]TestServerFunc) (ports map[string]int, cleanupFunc func()) {
	var wg sync.WaitGroup

	ports = make(map[string]int, len(testcontainers))
	cleanups := make([]func(), 0, len(testcontainers))

	for k, fn := range testcontainers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var c func()
			ports[k], c = fn(f)
			cleanups = append(cleanups, c)
		}()
	}

	wg.Wait()

	cleanupFunc = func() {
		for _, c := range cleanups {
			wg.Add(1)
			go func() {
				defer wg.Done()
				c()
			}()
		}
		wg.Wait()
	}

	return ports, cleanupFunc
}
