package service

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestClient_Post(t *testing.T) {
	Convey("Sends request", t, func() {
		c := make(chan interface{})
		var request *http.Request
		var body []byte
		ts := httptest.NewTLSServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				request = r
				b, err := ioutil.ReadAll(r.Body)
				if err != nil {
					panic(err)
				}
				body = b
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				io.WriteString(w, `{"success": 1}`)
				close(c)
			}))
		defer ts.Close()
		client := &Client{
			OpenViDuURL: ts.URL,
			Login:       "test login",
			Password:    "test password",
		}

		resp, err := client.Post("test", map[string]interface{}{
			"test": "test_data",
		})
		<-c

		Convey("With correct parameters", func() {

			So(request, ShouldNotBeNil)
			So(request.Method, ShouldEqual, http.MethodPost)
			So(request.URL.Path, ShouldEqual, "/test")
			So(request.Header, ShouldNotBeEmpty)
			So(request.Header.Get("Authorization"), ShouldContainSubstring,
				"Basic dGVzdCBsb2dpbjp0ZXN0IHBhc3N3b3Jk")
			So(request.Header.Get("Content-Type"),
				ShouldContainSubstring, "application/json")
		})

		Convey("with correct body", func() {
			So(string(body), ShouldContainSubstring, "{\"test\":\"test_data\"}")
		})

		Convey("Returns no error", func() {
			So(err, ShouldBeNil)
		})

		Convey("Response is  success", func() {
			So(resp, ShouldNotBeNil)
			So(resp["success"], ShouldEqual, 1)
		})
	})

	Convey("Returns marshall error", t, func() {
		client := &Client{
			OpenViDuURL: "some url",
			Login:       "test login",
			Password:    "test password",
		}

		_, err := client.Post("test", map[string]interface{}{
			"wrong": make(chan interface{}),
		})

		So(err, ShouldNotBeNil)
	})

	Convey("Returns send error", t, func() {
		client := &Client{
			OpenViDuURL: "wrong url",
			Login:       "test login",
			Password:    "test password",
		}

		_, err := client.Post("test", nil)

		So(err, ShouldNotBeNil)
	})

	Convey("Returns unmarshall error", t, func() {
		c := make(chan interface{})
		ts := httptest.NewTLSServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				close(c)
			}))
		defer ts.Close()
		client := &Client{
			OpenViDuURL: ts.URL,
			Login:       "test login",
			Password:    "test password",
		}

		_, err := client.Post("test", map[string]interface{}{
			"test": "test_data",
		})
		<-c

		So(err, ShouldNotBeNil)
	})
}
