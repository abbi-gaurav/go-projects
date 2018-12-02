package mocking

import (
	"net/http/httptest"
	"net/http"
	"fmt"
	"testing"
	"os"
)

var ts *httptest.Server

const checkMark = "\u2713"
const ballotX = "\u2717"

// feed is mocking the XML document we except to receive.
var feed = `<?xml version="1.0" encoding="UTF-8"?>
 <rss>
 <channel>
     <title>Going Go Programming</title>
     <description>Golang : https://github.com/goinggo</description>
     <link>http://www.goinggo.net/</link>
     <item>
         <pubDate>Sun, 15 Mar 2015 15:04:00 +0000</pubDate>
         <title>Object Oriented Programming Mechanics</title>
         <description>Go is an object oriented language.</description>
         <link>http://www.goinggo.net/2015/03/object-oriented</link>
     </item>
 </channel>
 </rss>`

func TestMain(m *testing.M) {
	ts = mockServer()

	ret := m.Run()
	ts.Close()

	os.Exit(ret)
}

func mockServer() *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/xml")
		fmt.Fprint(w, feed)
	}
	return httptest.NewServer(http.HandlerFunc(f))
}

func TestDownload(t *testing.T) {
	statusCode := 200

	t.Log("Given the need to test downloading content.")
	{
		t.Logf("\tWhen checking \"%s\" for status code \"%d\"",
			ts.URL, statusCode)
		{
			resp, err := http.Get(ts.URL)
			if err != nil {
				t.Fatal("\t\tShould be able to make the Get call.",
					ballotX, err)
			}
			t.Log("\t\tShould be able to make the Get call.",
				checkMark)

			defer resp.Body.Close()

			if resp.StatusCode != statusCode {
				t.Fatalf("\t\tShould receive a \"%d\" status. %v %v",
					statusCode, ballotX, resp.StatusCode)
			}
			t.Logf("\t\tShould receive a \"%d\" status. %v",
				statusCode, checkMark)
		}
	}
}
