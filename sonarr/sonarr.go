package sonarr

import (
	"crypto/tls"
	"net/http"

	"golift.io/starr"
)

// Sonarr contains all the methods to interact with a Sonarr server.
type Sonarr struct {
	starr.APIer
}

// Filter values are integers. Given names for ease of discovery.
//nolint:lll
// https://github.com/Sonarr/Sonarr/blob/0cb8d93069d6310abd39ee2fe73219e17aa83fe6/src/NzbDrone.Core/History/EpisodeHistory.cs#L34-L41
const (
	FilterUnknown starr.Filtering = iota
	FilterGrabbed
	FilterSeriesFolderImported
	FilterDownloadFolderImported
	FilterDownloadFailed
	FilterDeleted
	FilterRenamed
	FilterImportFailed
)

// New returns a Sonarr object used to interact with the Sonarr API.
func New(config *starr.Config) *Sonarr {
	if config.Client == nil {
		//nolint:exhaustivestruct,gosec
		config.Client = &http.Client{
			Timeout: config.Timeout.Duration,
			CheckRedirect: func(r *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: !config.ValidSSL},
			},
		}
	}

	if config.Debugf == nil {
		config.Debugf = func(string, ...interface{}) {}
	}

	return &Sonarr{APIer: config}
}
