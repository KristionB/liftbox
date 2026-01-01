package sync

import "sync"

func UploadFiles(files []string) {
	var wg sync.WaitGroup

	for _, f := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			encryptAndSend(file)
		}(f)
	}

	wg.Wait()
}

