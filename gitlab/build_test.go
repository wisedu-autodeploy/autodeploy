package gitlab

import (
	"log"
	"sync"
	"testing"
)

func Test_WatchBuildLog(t *testing.T) {
	t.Skip("skip Test_WatchBuildLog")
	debugger()
	tag, err := NewTag(TestValidParams)
	if err != nil {
		t.Error("fail at NewTag: ", err)
		return
	}

	var wg sync.WaitGroup
	var logChan = make(chan *Logger)
	wg.Add(2)

	go WatchBuildLog(TestValidParams, tag, logChan, &wg)
	go GetBuildLog(func(logger *Logger) {
		log.Println(logger.Status)
	}, logChan, &wg)
	wg.Wait()
	logger := <-logChan
	close(logChan)
	log.Println(logger.Image)
}
