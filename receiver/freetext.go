package receiver

import (
	"fmt"
	"hel/fshare"
	"hel/tools"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
)

// FreeText struct
type FreeText struct {
	FS    *fshare.Fshare
	Creds *tools.FSCredentials
}

// GetLocation retrieves the location and information of a file based on the input string.
// It processes the results to find the file with the largest size, fetches additional information and locations,
// converts the file size to int64, replaces "http://" with "https://" in the URLs, and writes the play URL to the clipboard.
// It returns an error if any operation fails during the process.
func (f *FreeText) GetLocation(input string) error {
	results, err := f.FS.Query(input)
	if err != nil {
		return err
	}

	ls := make(map[int64]*fshare.QueryDetailResponse, 0)
	for _, v := range results.Data {
		s := strings.ToLower(v.Name)
		// remove all extension of file name
		s = strings.TrimSuffix(filepath.Base(s), filepath.Ext(s))

		// Check if the element starts with "ss1" or contains "[ss1"
		if strings.HasPrefix(strings.ToLower(s), input) || strings.Contains(s, "]"+input) {
			// Add the element to the new slice
			ls[v.Size] = v
		}
	}

	if len(ls) == 0 {
		fmt.Println("Nothing else matter")
		return nil
	}

	// sorting to get the biggest size
	keys := make([]int64, 0, len(ls))

	for k := range ls {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})

	capturedURL := prettyURL(ls[keys[0]].URL)

	info, err := f.FS.GetFileInfo(f.Creds.Token, f.Creds.SessionID, capturedURL)
	if err != nil {
		return err
	}

	location, err := f.FS.GetFileLocation(f.Creds.Token, f.Creds.SessionID, capturedURL)
	if err != nil {
		return err
	}

	pLocation, err := f.FS.GetFileLocation(f.Creds.Token, f.Creds.SessionID, capturedURL)
	if err != nil {
		return err
	}

	// replace http:// to https://
	downloadURL := strings.Replace(location.Location, "http://", "https://", 1)
	playURL := strings.Replace(pLocation.Location, "http://", "https://", 1)

	// convert info.Size to int64
	size, err := strconv.ParseInt(info.Size, 10, 64)
	if err != nil {
		return err
	}

	fmt.Println("File size: ", prettyByteSize(size))
	fmt.Println("Download URL: ", downloadURL)

	return clipboard.WriteAll(playURL)
}
