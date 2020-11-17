package checks

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/math2001/gocmt/cmt"
)

func URLs(c *cmt.CheckResult, args map[string]interface{}) {
	name := args["name"].(string)
	url := args["url"].(string)
	pattern := args["pattern"].(string)

	var allowRedirects bool
	if v, ok := args["allow_redirects"]; ok {
		allowRedirects = v.(bool)
	}

	var sslVerify bool
	if v, ok := args["ssl_verify"]; ok {
		sslVerify = v.(bool)
	}

	c.AddItem(&cmt.CheckItem{
		Name:  "url_name",
		Value: name,
	})

	c.AddItem(&cmt.CheckItem{
		Name:  "url",
		Value: url,
	})

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !sslVerify,
		},
	}

	var httpclient = &http.Client{
		Transport: transport,
		Timeout:   5 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if allowRedirects {
				return nil // follow redirect
			}
			return http.ErrUseLastResponse // stop sequence here
		},
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	start := time.Now()
	res, err := httpclient.Do(req)
	if err != nil {
		fmt.Fprintf(c.DebugBuffer(), "httpclient.Do: %s", err)
		c.AddItem(&cmt.CheckItem{
			Name:         "url_status",
			Value:        "nok",
			IsAlert:      true,
			AlertMessage: fmt.Sprintf("check_url - %s - error (%s)", name, err),
		})
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		c.AddItem(&cmt.CheckItem{
			Name:         "url_status",
			Value:        "nok",
			IsAlert:      true,
			AlertMessage: fmt.Sprintf("check_url - %s - bad http code response (%d)", name, res.StatusCode),
		})
		return
	}
	dt := time.Since(start)
	c.AddItem(&cmt.CheckItem{
		Name:  "url_msec",
		Value: dt.Milliseconds(),
		Unit:  "ms",
	})

	matched, err := regexp.MatchReader(pattern, bufio.NewReader(res.Body))
	if err != nil {
		panic(err)
	}
	if !matched {
		c.AddItem(&cmt.CheckItem{
			Name:         "url_status",
			Value:        "nok",
			Description:  "expected pattern not found",
			IsAlert:      true,
			AlertMessage: fmt.Sprintf("check_url - %s - expected pattern not found", name),
		})
		return
	}
	c.AddItem(&cmt.CheckItem{
		Name:  "url_status",
		Value: "ok",
	})
}
