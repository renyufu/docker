package client

import (
	"io"
	"net/url"

	"golang.org/x/net/context"
)

// ImageSave retrieves one or more images from the docker host as an io.ReadCloser.
// It's up to the caller to store the images and close the stream.
func (cli *Client) ImageReadme(ctx context.Context, imageID string) (io.ReadCloser, error) {
	query := url.Values{
		"names": []string{imageID},
	}

	resp, err := cli.get(ctx, "/images/get", query, nil)
	if err != nil {
		return nil, err
	}
	return resp.body, nil
}
