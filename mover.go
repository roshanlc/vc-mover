package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// autoMoveAtStart goes through content of downloads folder at startup
// and moves them  to respective directory
func autoMoveAtStart(downloadsDir string) {
	files, err := os.ReadDir(downloadsDir)

	if err != nil {
		log.Println(err)
		return
	}

	var wg sync.WaitGroup

	for _, file := range files {
		if vericredRegex.MatchString(file.Name()) && !strings.HasSuffix(file.Name(), ".part") {
			details, err := file.Info()
			if err != nil {
				log.Println(err)
				continue
			}

			if details.Size() > 0 {
				wg.Add(1)

				go func(file fs.DirEntry, downloadsDir string, wg *sync.WaitGroup) {

					log.Println("found vericred match: ", file.Name())
					// run the (moving, extraction, deletion) process in a separate goroutine

					log.Println("starting file movement for ", file.Name())
					ext := filepath.Ext(file.Name())

					// get the folder name
					dirName := extractDirName(file.Name())
					if dirName == "" {
						log.Println("Unable to extract directory name for ", file.Name())

						wg.Done()
						return
					}

					dirPath := downloadsDir + "/" + dirName

					// create the folder if not exists
					err := createDirIfNotExists(dirPath)
					if err != nil {
						log.Println("Unable to create directory: ", dirName)
						log.Println(err)

						wg.Done()
						return
					}

					// move the file into the folder
					source := downloadsDir + "/" + file.Name()
					dest, _ := filepath.Abs(dirPath + "/" + file.Name())

					err = moveToDir(source, dest)
					if err != nil {
						log.Println("Unable to move file to its folder: ", file.Name())
						log.Println(err)

						wg.Done()
						return
					}

					// unzip the file if required
					if strings.EqualFold(ext, ".zip") {
						log.Println("starting zip extraction")
						destFilePath := strings.TrimSuffix(dest, ".zip")

						err = unZipAndRemove(dest, destFilePath)
						if err != nil {
							log.Println("Failed unzip process:", err.Error())

							wg.Done()
							return
						} else {
							log.Println("Unzip successfull for ", destFilePath)
						}
					} else {
						log.Println("File movement successful for: ", file.Name())
					}

					wg.Done()
				}(file, downloadsDir, &wg)

			}

		}
	}

	// wait for go routines to complete
	wg.Wait()
}
