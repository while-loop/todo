package parser

import (
	"fmt"
	"io"
	"net/http"
	"github.com/pkg/errors"
)

// DownloadFile downloads repository files. Http clients should have access
// tokens already set. With this, parser does not need knowledge of private/public repos
func DownloadFile(client *http.Client, u string) (io.ReadCloser, error) {
	resp, err := client.Get(u)
	if err != nil {
		return nil, errors.Wrap(err, "failed to download file to parse")
	}

	if resp.StatusCode != http.StatusOK {
		if resp.Body != nil {
			resp.Body.Close()
		}
		return nil, fmt.Errorf("failed to retrieve download file: status code %d", resp.StatusCode)
	}

	return resp.Body, nil
}
