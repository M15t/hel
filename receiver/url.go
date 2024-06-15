package receiver

import (
	"fmt"
	"hel/fshare"
	"hel/tools"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
)

// URL struct
type URL struct {
	FS    *fshare.Fshare
	Creds *tools.FSCredentials
}

// GetLocation retrieves the location and information of a file based on the input URL.
// It replaces "http://" with "https://" in the download URL, converts the file size to int64,
// and writes the download URL to the clipboard using the clipboard package.
// It returns an error if any operation fails during the process.
func (u *URL) GetLocation(input string) error {
	location, err := u.FS.GetFileLocation(u.Creds.Token, u.Creds.SessionID, input)
	if err != nil {
		return err
	}

	info, err := u.FS.GetFileInfo(u.Creds.Token, u.Creds.SessionID, input)
	if err != nil {
		return err
	}

	// replace http:// to https://
	downloadURL := strings.Replace(location.Location, "http://", "https://", 1)

	// convert info.Size to int64
	size, err := strconv.ParseInt(info.Size, 10, 64)
	if err != nil {
		return err
	}

	fmt.Println("File size: ", prettyByteSize(size))
	fmt.Println("Download URL: ", downloadURL)

	return clipboard.WriteAll(downloadURL)
}
