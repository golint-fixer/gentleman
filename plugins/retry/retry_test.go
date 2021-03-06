package retry

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/lytics/gentleman"
	"github.com/lytics/gentleman/plugins/timeout"
	"github.com/lytics/gentleman/utils"
)

func TestRetryRequest(t *testing.T) {
	calls := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls < RetryTimes {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("foo", r.Header.Get("foo"))
		fmt.Fprintln(w, "Hello, world")
	}))
	defer ts.Close()

	req := gentleman.NewRequest()
	req.SetHeader("foo", "bar")
	req.URL(ts.URL)
	req.Use(New(nil, nil))

	res, err := req.Send()
	utils.Equal(t, err, nil)
	utils.Equal(t, res.Ok, true)
	utils.Equal(t, res.StatusCode, 200)
	utils.Equal(t, res.Header.Get("foo"), "bar")
	utils.Equal(t, calls, 3)
}

func TestRetryRequestWithPayload(t *testing.T) {
	calls := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		buf, _ := ioutil.ReadAll(r.Body)
		fmt.Fprintln(w, string(buf))
	}))
	defer ts.Close()

	req := gentleman.NewRequest()
	req.URL(ts.URL)
	req.Method("POST")
	req.BodyString("Hello, world")
	req.Use(New(nil, nil))

	res, err := req.Send()
	utils.Equal(t, err, nil)
	utils.Equal(t, res.Ok, true)
	utils.Equal(t, res.RawResponse.ContentLength, int64(13))
	utils.Equal(t, res.StatusCode, 200)
	utils.Equal(t, res.String(), "Hello, world\n")
	utils.Equal(t, calls, 3)
}

func TestRetryServerError(t *testing.T) {
	calls := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer ts.Close()

	req := gentleman.NewRequest()
	req.URL(ts.URL)
	req.Use(New(nil, nil))

	res, err := req.Send()
	utils.Equal(t, err, nil)
	utils.Equal(t, res.Ok, false)
	utils.Equal(t, res.StatusCode, 503)
	utils.Equal(t, calls, 4)
}

func TestRetryNetworkError(t *testing.T) {
	req := gentleman.NewRequest()
	req.URL("http://127.0.0.1:9123")
	req.Use(New(nil, nil))

	res, err := req.Send()
	utils.NotEqual(t, err, nil)
	utils.Equal(t, strings.Contains(err.Error(), "connection refused"), true)
	utils.Equal(t, res.Ok, false)
	utils.Equal(t, res.StatusCode, 0)
}

// Timeout retry is not fully supported yet
func TestRetryNetworkTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		w.WriteHeader(200)
	}))
	defer ts.Close()

	req := gentleman.NewRequest()
	req.URL(ts.URL)
	req.Use(timeout.Request(100 * time.Millisecond))
	req.Use(New(nil, nil))

	res, err := req.Send()
	utils.NotEqual(t, err, nil)
	utils.Equal(t, strings.Contains(err.Error(), "request canceled"), true)
	utils.Equal(t, res.Ok, false)
	utils.Equal(t, res.StatusCode, 0)
}

func TestCustomEvaluator(t *testing.T) {
	calls := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if calls < 3 {
			calls++
			w.WriteHeader(418)
			return
		}
		w.WriteHeader(200)
	}))
	defer ts.Close()

	evaluator := func(err error, res *http.Response, req *http.Request) error {
		if res.StatusCode == 418 {
			return errors.New("you are not a teapot")
		}
		return nil
	}

	req := gentleman.NewRequest()
	req.URL(ts.URL)
	req.Use(New(nil, evaluator))

	res, err := req.Send()
	utils.Equal(t, nil, err)
	utils.Equal(t, calls, RetryTimes)
	utils.Equal(t, res.Ok, true)
}

func TestStackEvaluators(t *testing.T) {
	calls := 1
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch calls {
		case 1:
			w.WriteHeader(418)
		case 2:
			w.WriteHeader(420)
		case 3:
			w.WriteHeader(200)
		}
		calls++
	}))
	defer ts.Close()

	evaluations := 0
	evaluator := func(err error, res *http.Response, req *http.Request) error {
		evaluations++
		if res.StatusCode == 418 {
			return errors.New("you are not a teapot")
		}
		return nil
	}
	evaluator2 := func(err error, res *http.Response, req *http.Request) error {
		evaluations++
		if res.StatusCode == 420 {
			return errors.New("enhance your calm")
		}
		return nil
	}

	req := gentleman.NewRequest()
	req.URL(ts.URL)
	req.Use(New(nil, evaluator2))
	req.Use(New(nil, evaluator))

	res, err := req.Send()
	utils.Equal(t, nil, err)
	utils.Equal(t, calls, RetryTimes+1)
	utils.Equal(t, res.Ok, true)
	utils.Equal(t, evaluations, 5)
}
