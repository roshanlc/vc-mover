package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gen2brain/beeep"
)

// setupLog prepares folder for log output
func setupLogFolder(homeDir string) error {
	return createDirIfNotExists(homeDir + "/.cache/vericred-mover")
}

// vericred regexp
var vericredRegex *regexp.Regexp

func main() {

	// regex to catch all file starting with "vericred_"
	vericredRegex = regexp.MustCompile(`(?i)vericred_.*\d+.*\.*`)

	// get the download directory's absolute path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	// setup for logs
	err = setupLogFolder(homeDir)
	if err != nil {
		log.Fatal(err)
	}

	// get current date
	t := time.Now()
	year, month, day := t.Date()

	logFile, err := os.OpenFile(fmt.Sprintf(homeDir+"/.cache/vericred-mover/vericred-mover_%d_%d_%d.log", year, month, day), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		log.Println("Error creating log file")
		log.Fatal(err)
	}

	defer logFile.Close()

	logger := log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
	downloadsDir, err := filepath.Abs(homeDir + "/Downloads")

	if err != nil {
		logger.Fatal(err)
	}

	// auto manage the downloads folder at program startup
	go autoMoveAtStart(downloadsDir)

	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Has(fsnotify.Create) {
					// check if the file name matches vericred_pattern and is not a part file

					if vericredRegex.MatchString(event.Name) && !strings.HasSuffix(event.Name, ".part") {

						// Get file size info

						stat, err := os.Stat(event.Name)
						if err != nil {
							logger.Println(err)
						} else {
							if stat.Size() > 0 {

								logger.Println("found vericred match: ", event.Name)
								// run the (moving, extraction, deletion) process in a separate goroutine
								go func() {
									logger.Println("starting file movement for ", stat.Name())
									ext := filepath.Ext(event.Name)

									// get the folder name
									dirName := extractDirName(stat.Name())
									if dirName == "" {
										logger.Println("Unable to extract directory name for ", event.Name)
										beeep.Alert("Vericred-mover", "Unable to extract directory name for "+event.Name, "")
										return
									}

									dirPath := downloadsDir + "/" + dirName

									// create the folder if not exists
									err := createDirIfNotExists(dirPath)
									if err != nil {
										logger.Println("Unable to create directory: ", dirName)
										logger.Println(err)
										beeep.Alert("Vericred-mover", "Unable to create directory: "+dirName, "")
										return
									}

									// move the file into the folder
									source := event.Name
									dest, _ := filepath.Abs(dirPath + "/" + stat.Name())

									err = moveToDir(source, dest)
									if err != nil {
										logger.Println("Unable to move file to its folder: ", event.Name)
										logger.Println(err)
										beeep.Alert("Vericred-mover", "Unable to move file to its folder: "+event.Name, "")

										return
									}

									// unzip the file if required
									if strings.EqualFold(ext, ".zip") {
										logger.Println("starting zip extraction")
										destFilePath := strings.TrimSuffix(dest, ".zip")

										err = unZipAndRemove(dest, destFilePath)
										if err != nil {
											logger.Println("Failed unzip process:", err.Error())
											beeep.Alert("Vericred-mover", "Unzip process failed due to:", "")
											return
										} else {
											logger.Println("Unzip successfull for ", destFilePath)
											beeep.Notify("Vericred-mover", destFilePath+"has been moved to its directory", "")
										}
									} else {
										logger.Println("File movement successful for: ", stat.Name())
										beeep.Notify("Vericred-mover", "File moved successfuly to: "+dest, "")
									}

								}()

							}
						}

					}

				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logger.Println("error:", err)
			}
		}
	}()

	// Add a path.
	err = watcher.Add(downloadsDir)
	if err != nil {
		logger.Fatal(err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}

// extracts directory name from the vericred files
func extractDirName(fileName string) string {
	var re = regexp.MustCompile(`(?i)(vericred_\w+)_\d+.*\.*`)

	matches := re.FindStringSubmatch(fileName)
	if len(matches) == 2 {
		return matches[1]
	}
	return ""
}

// createDir creates a directory if it does not exist
func createDirIfNotExists(dir string) error {
	_, err := os.Stat(dir)

	if err != nil && os.IsNotExist(err) {
		// create the directory
		return os.Mkdir(dir, 0755)
	}

	return nil

}

// moveToDir moves the provided file from source to destination
func moveToDir(source, dest string) error {
	return os.Rename(source, dest)
}

// unZipIfRequired will operate on zip files only
func unZipAndRemove(zipFilePath, destFilePath string) error {

	zipReader, err := zip.OpenReader(zipFilePath)
	if err != nil && err != zip.ErrInsecurePath {
		fmt.Println("openreader err: ", err)
		return err
	}

	var fileInArchive zip.File

	for _, file := range zipReader.File {
		if filepath.Ext(file.Name) == ".json" {
			fileInArchive = *file
		}
	}

	destFile, err := os.Create(destFilePath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	archiveReader, err := fileInArchive.Open()
	if err != nil {
		return err
	}

	if _, err := io.Copy(destFile, archiveReader); err != nil {
		return err
	}

	// close the opened archive file manually
	zipReader.Close()
	archiveReader.Close()

	time.Sleep(200 * time.Millisecond)
	// remove the zip file
	err = os.Remove(zipFilePath)
	if err != nil {
		return err
	}

	return nil
}
