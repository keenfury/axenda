package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/keenfury/axenda/config"
)

func TruncateTimeToMinute(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, GetLocation(t))
}

func GetLocation(t time.Time) *time.Location {
	if config.UseUTC == "true" {
		return time.UTC
	}
	return t.Location()
}

func GetNow() time.Time {
	if config.UseUTC == "true" {
		return time.Now().UTC()
	}
	return time.Now()
}

func SimpleRequest(mode, url string, bodyIn, bodyOut interface{}, expectedCode int, hdrArgs map[string]string) (err error) {
	var readerIn io.Reader
	if bodyIn != nil {
		bBodyIn, errMarshal := json.Marshal(bodyIn)
		if errMarshal != nil {
			err = errMarshal
			return
		}
		readerIn = bytes.NewReader(bBodyIn)
	}
	req, errReq := http.NewRequest(mode, url, readerIn)
	if errReq != nil {
		err = errReq
		return
	}
	if hdrArgs != nil {
		for k, v := range hdrArgs {
			req.Header.Add(k, v)
		}
	}
	resp, errResp := http.DefaultClient.Do(req)
	if errResp != nil {
		err = errResp
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != expectedCode {
		err = fmt.Errorf("Unexpected code: %d, wanted: %d, reason: %s", resp.StatusCode, expectedCode, resp.Status)
		return
	}
	body, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		err = errRead
		return
	}
	if bodyOut != nil {
		if err = json.Unmarshal(body, bodyOut); err != nil {
			return
		}
	}
	return
}
