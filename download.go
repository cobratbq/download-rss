package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

// PodcastEntry containing all relevant fields for downloading.
type PodcastEntry struct {
	Title     string `xml:"title"`
	Enclosure struct {
		URL    string `xml:"url,attr"`
		Length uint64 `xml:"length,attr"`
		Type   string `xml:"type,attr"`
	} `xml:"enclosure"`
}

// Application for downloading podcasts from RSS feeds.
//
// First argument is the file path to the RSS file. If no argument was
// provided, Stdin is used as the RSS feed source.
func main() {
	flag.Parse()
	source, err := getFeedSource()
	if err != nil {
		os.Stderr.WriteString("failed to acquire feed: " + err.Error() + "\n")
		return
	}
	defer source.Close()
	numFailedDownloads, err := downloadEntriesFromRSSFile(source)
	if err != nil {
		os.Stderr.WriteString("error while downloading podcasts: " + err.Error() + "\n")
	}
	fmt.Printf("Finished downloading podcasts. (%d downloads failed)\n", numFailedDownloads)
	return
}

func getFeedSource() (io.ReadCloser, error) {
	if flag.NArg() < 1 {
		return os.Stdin, nil
	}
	// Open file for reading.
	filePath := flag.Arg(0)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func downloadEntriesFromRSSFile(source io.Reader) (int, error) {
	var numberFailed = 0
	// Parse the xml file for the 'item' tag and read the title and url from each entry.
	decoder := xml.NewDecoder(source)
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
					return numberFailed, err
				}
				fmt.Printf("Downloading '%s'\n", entry.Title)
				err = execute("wget", "-c", entry.Enclosure.URL)
				if err != nil {
					numberFailed += 1
					log.Printf("Error occurred during execution: %s", err.Error())
				}
			}
		}
	}
	return numberFailed, nil
}

func execute(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
