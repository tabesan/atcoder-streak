package commit

import (
	td "atcoder-streak/commit/test_data"
	"fmt"
	"net/http"
	"net/http/httptest"
)

const (
	allCommitURL  = "/repos/tabe/Atcoder/commits"
	lastCommitURL = "/repos/tabe/Atcoder/commits&per_page=1"
	header        = "Accept"
	value         = "application/vnd.github.v3+json"
)

func NewMockServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc(allCommitURL, func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set(header, value)
		fmt.Fprint(w, td.ResultAll)
	})

	GithubMockServer := httptest.NewServer(mux)
	defer GithubMockServer.Close()

	return GithubMockServer
}
