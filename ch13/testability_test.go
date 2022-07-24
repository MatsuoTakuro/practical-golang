package ch13

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type User2 struct {
	Name    string
	Address string
}

func DumpUser(u *User2) {
	DumpUserTo(os.Stdout, u)
}

// DumpUserTo has the abstract interface, io.Writer as argument.
func DumpUserTo(w io.Writer, u *User2) {
	if u.Address == "" {
		fmt.Fprintf(w, "%s(address unknown)", u.Name)
	} else {
		fmt.Fprintf(w, "%s@%s", u.Name, u.Address)
	}
}

func TestConsoleOut(t *testing.T) {
	var b bytes.Buffer
	DumpUserTo(&b, &User2{
		Name:    "Takuro Matsuo",
		Address: "",
	})
	if b.String() != "Takuro Matsuo(address unknown)" {
		t.Errorf("error (expected: 'Takuro Matsuo(address unknown)', actual='%s')", b.String())
	}
}

// helper for testing
type testTransport struct {
	req **http.Request
	res *http.Response
	err error
}

func (tr *testTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	*(tr.req) = req
	return tr.res, tr.err
}

var _ http.RoundTripper = &testTransport{}

func newTransport(req **http.Request, res *http.Response, err error) http.RoundTripper {
	return &testTransport{
		req: req,
		res: res,
		err: err,
	}
}

func TestHttpRequest(t *testing.T) {
	var req *http.Request
	res := httptest.NewRecorder()
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	_, err := res.WriteString(`{"ranking": ["Back to the Future, "Rambo"]}`)
	if err != nil {
		log.Fatal(err)
	}

	clt := http.Client{
		Transport: newTransport(&req, res.Result(), nil),
	}
	resp, err := clt.Get("http://example.com/movies/1985")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

type ctxTimeKey string

const timeKey ctxTimeKey = "timeKey"

func CurrentTime(ctx context.Context) time.Time {
	// for test: value for timeKey is set
	v := ctx.Value(timeKey)
	if t, ok := v.(time.Time); ok {
		return t
	}
	// default: no value for timeKey exists
	return time.Now()
}

// call this func only for test
func SetFixTime(ctx context.Context, t time.Time) context.Context {
	// set fix time as value of timekey
	return context.WithValue(ctx, timeKey, t)
}

// object of test
func NextMonth(ctx context.Context) time.Month {
	now := CurrentTime(ctx)
	return now.AddDate(0, 1, 0).Month()
}

func TestNextMonth(t *testing.T) {
	// include fix time in cotext
	ctx := SetFixTime(context.Background(), time.Date(1980, time.December, 1, 0, 0, 0, 0, time.Local))
	assert.Equal(t, time.January, NextMonth(ctx))
}
