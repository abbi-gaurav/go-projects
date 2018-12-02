package unit_testing

import (
	"testing"
	"net/http"
)

const checkMark = "\u2713"
const ballotX = "\u2717"

func TestDownload(t *testing.T) {
	url := "http://www.goinggo.net/feeds/posts/default?alt=rss"
	statusCode := 200

	t.Log("Given the need to test downloading content")
	{
		t.Logf("\tWhen checking [%s] for status code [%d]", url, statusCode)
		{
			resp, err := http.Get(url)
			if err != nil {
				t.Fatal("\t\t Should be able to make Get call", ballotX, err)
			}
			t.Log("\t\t Should be able to make GET call", checkMark)
			defer resp.Body.Close()

			if resp.StatusCode == statusCode {
				t.Logf("\t\tShould receive a [%d] status. %v", statusCode, checkMark)
			} else {
				t.Errorf("Expected [%d] status but got [%v] - %v", statusCode, resp.StatusCode, ballotX)
			}
		}
	}
}
