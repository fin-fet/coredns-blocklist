package blocklist

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"strings"
)

func buildBlocklist(options BlocklistOptions) (Blocklist, error) {
	if options.MatchSubdomains {
		return NewRadixBlocklist(), nil
	} else {
		return NewBasicBlocklist(), nil
	}
}

func (blp *BlocklistPlugin) loadDomains() {
	if blp.options.SourceType == SourceTypeFile {
		file, err := os.Open(blp.options.Url)
		if err != nil {
			log.Errorf("failed to open '%s': %v", blp.options.Url, err)
			return
		}
		defer file.Close()

		// Collect stats, so we can compare previous file access

		stat, err := file.Stat()
		if err != nil {
			return
		}

		size := blp.lastFileSize

		if size == stat.Size() && blp.lastFileTime.Equal(stat.ModTime()) {
			log.Debug("No file changes")
			return
		}

		blp.lastFileSize = stat.Size()
		blp.lastFileTime = stat.ModTime()
	}

	// Empty blocklist
	blp.blocklist, _ = buildBlocklist(*blp.options)

	for name := range loadFromSource(*blp.options) {
		//log.Infof("added domain '%s'", name)
		blp.blocklist.Add(name)
	}
}

func loadFromSource(options BlocklistOptions) chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)

		// Load the source, depending on type
		var sourceData io.Reader
		{
			if options.SourceType == SourceTypeFile {
				log.Infof("loading blocklist file '%s'", options.Url)

				file, err := os.Open(options.Url)
				if err != nil {
					log.Errorf("failed to load file '%s': %v", options.Url, err)
				}
				defer file.Close()
				sourceData = file
			} else if options.SourceType == SourceTypeHttp {
				log.Infof("loading blocklist URL '%s'", options.Url)

				response, err := http.Get(options.Url)
				if err != nil {
					log.Errorf("failed loading HTTP '%s': %v", options.Url, err)
				}
				defer response.Body.Close()
				sourceData = response.Body
			}
		}

		scanner := bufio.NewScanner(sourceData)
		for scanner.Scan() {
			name := scanner.Text()

			if name == "" || strings.HasPrefix(name, "#") {
				// Skip comments and empty strings
				continue
			}

			// Remove preceeding IP, if applicable (ex Hostfile)
			fields := strings.Fields(name)
			if len(fields) == 1 {
				name = fields[0]
			} else {
				name = fields[1]
			}

			// Add root, if not in name
			if !strings.HasSuffix(name, ".") {
				name += "."
			}

			ch <- name
		}
	}()

	return ch
}
