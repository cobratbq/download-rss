package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type PodcastEntry struct {
	Title string `xml:"title"`
	Link  string `xml:"guid"`
}

func main() {
	path := initialize()
	err := downloadEntriesFromRSSFile(path)
	if err != nil {
		fmt.Println(err)
		fmt.Println()
		flag.Usage()
	}
}

func initialize() *string {
	path := flag.String("f", "rssfeed.xml", "Path to RSS feed xml file.")
	flag.Parse()
	return path
}

func downloadEntriesFromRSSFile(filePath *string) error {

	//Open file for reading.
	file, err := os.Open(*filePath)
	if err != nil {
		log.Printf("Failed to open '%s'. (%s)", *filePath, err.Error())
		return err
	}
	defer file.Close()

	//Parse the xml file for the 'item' tag and read the title and url from each entry.
	decoder := xml.NewDecoder(file)
	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}

		switch t := token.(type) {
		case xml.StartElement:
			name := t.Name.Local
			if name == "item" {
				var entry PodcastEntry
				err := decoder.DecodeElement(&entry, &t)
				if err != nil {
					log.Printf("An error occurred during decoding. (%s)", err.Error())
				}
				fmt.Printf("Downloading '%s'\n", entry.Title)
				err = execute("wget", "-c", entry.Link)
				if err != nil {
					log.Printf("Error occurred during execution: %s", err.Error())
				}
			}
		}
	}

	return nil
}

func execute(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
