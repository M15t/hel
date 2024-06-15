package main

import (
	"fmt"
	"hel/fshare"
	"hel/receiver"
	"hel/tools"
	"strings"
	"sync"
	"time"
)

func main() {
	fs := fshare.New()

	// step 1: receive input from command line
	for {
		var input string
		fmt.Print("Enter something: ")
		fmt.Scan(&input)

		if input == "exit" {
			break // Exit the loop if the user types 'exit'
		}

		q := strings.TrimSpace(input)

		// step 2: classify the input: free text or url
		// step 2.1: get creds in redis
		creds, err := tools.GetCreds(fs)
		if err != nil {
			checkErr(fmt.Errorf("token or sessionID is empty: %s", err.Error()))
		}

		// step 3: get location
		receiverSvc := receiver.URL{FS: fs, Creds: creds}

		switch {
		case strings.Contains(q, "https://www.fshare.vn/file"): // for single file
			if err := receiverSvc.GetLocation(q); err != nil {
				checkErr(err)
			}
		case strings.Contains(q, "https://www.fshare.vn/folder"): // for folder
			results, err := fs.GetFilesInFolder(creds.Token, creds.SessionID, q)
			if err != nil {
				checkErr(err)
			}

			enqueued := tools.StartQueue(results)
			errors := make(chan string, 0)

			wg := sync.WaitGroup{}
			wg.Add(5)
			for i := 0; i < 5; i++ {
				go func(workerID int) {
					defer wg.Done()
					for q := range enqueued {
						// fmt.Printf("Worker %d is working on %s\n", workerID, q)
						if err := receiverSvc.GetLocation(q); err != nil {
							fmt.Printf("Error in getLocation: %v\n", err)
							errors <- err.Error()
						}
						time.Sleep(time.Second * 2)
					}
				}(i)
			}

			wg.Wait()

		default:
			receiverSvc := receiver.FreeText{FS: fs, Creds: creds}
			if err := receiverSvc.GetLocation(q); err != nil {
				checkErr(err)
			}
		}
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
