package table_test

import (
	"testing"
	"net/http"
)

const checkMark = "\u2713"
const ballotX = "\u2717"

func TestDownload(t *testing.T) {
	var urls = []struct {
		url        string
		statusCode int
	}{
		{
			url:        "http://www.goinggo.net/feeds/posts/default?alt=rss",
			statusCode: http.StatusOK,
		},
		{
			url:        "http://rss.cnn.com/rss/cnn_topstbadurl.rss",
			statusCode: http.StatusNotFound,
		},
	}

	t.Log("Given the need to test downloading different content")
	{
		for _, u := range urls {
			t.Logf("\tWhen checking for status code %d for url %s", u.statusCode, u.url)
			{
				resp, err := http.Get(u.url)
				if err != nil {
					t.Fatal("\t\t Should be able to get the url.", ballotX, err)
				}
				t.Log("\t\tShould be able to get the url", checkMark)
				defer resp.Body.Close()

				if resp.StatusCode == u.statusCode {
					t.Logf("\t\tShould have %d status %v", u.statusCode, checkMark)
				} else {
					t.Errorf("\t\tExpected %d status, Got %v %v", u.statusCode, resp.StatusCode, ballotX)
				}
			}
		}
	}
}
