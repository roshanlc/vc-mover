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

func main() {

	// get the download directory's absolute path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	downloadsDir, err := filepath.Abs(homeDir + "/Downloads")

	if err != nil {
		log.Fatal(err)
	}

	// regex to catch all file starting with "vericred_"
	vericredRegex := regexp.MustCompile(`(?i)vericred_.*\d+.*\.*`)

	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
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
					fmt.Println("Create event: ", event.Name)

					// check if the file name matches vericred_pattern

					if vericredRegex.MatchString(event.Name) {
						fmt.Println("--------------------------------")
						fmt.Println("found vericred match: ", event.Name)

						// Get file size info

						stat, err := os.Stat(event.Name)
						if err != nil {
							log.Println(err)
						} else {
							if stat.Size() > 0 {

								// run the (moving, extraction, deletion) process in a separate goroutine
								go func() {
									log.Println(stat.Name(), " : ", stat.Size())
									log.Println("starting file movement for ", stat.Name())
									ext := filepath.Ext(event.Name)

									// get the folder name
									dirName := extractDirName(stat.Name())
									if dirName == "" {
										log.Println("Unable to extract directory name for ", event.Name)
										return
									}

									dirPath := downloadsDir + "/" + dirName

									// create the folder if not exists
									err := createDirIfNotExists(dirPath)
									if err != nil {
										log.Println("Unable to create directory: ", dirName)
										log.Println(err)
										return
									}

									// move the file into the folder
									source := event.Name
									dest, _ := filepath.Abs(dirPath + "/" + stat.Name())

									err = moveToDir(source, dest)
									if err != nil {
										log.Println("Unable to move file to its folder: ", event.Name)
										log.Println(err)
										return
									}

									// unzip the file if required
									if strings.EqualFold(ext, ".zip") {
										log.Println("starting zip extraction")
										fmt.Println("source zip: ", dest)
										destFilePath := strings.TrimSuffix(dest, ".zip")
										fmt.Println("unzipped file: ", destFilePath)
										err = unZipAndRemove(dest, destFilePath)
										if err != nil {
											log.Println("Failed unzip process:", err)
											beeep.Alert("Vericred-mover", "Unzip process failed due to:", err.Error())
											return
										} else {
											log.Println("Unzip successfull for ", destFilePath)
											beeep.Notify("Vericred-mover", destFilePath+"has been moved to its directory", "")
										}
									} else {
										log.Println("File movement successful for: ", dest)
										beeep.Notify("Vericred-mover", "File moved successfuly: "+dest, "")
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
				fmt.Println("error:", err)
			}
		}
	}()

	// Add a path.
	err = watcher.Add(downloadsDir)
	if err != nil {
		log.Fatal(err)
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
