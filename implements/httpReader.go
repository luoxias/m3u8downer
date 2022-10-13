package implements

import (
	"io"
	"net/http"
)

type ReadFromUrl struct {
	res    *http.Response
	reader io.Reader
}

func (h *ReadFromUrl) SetUrl(url string) error {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header = map[string][]string{
		"accept":          {"*/*"},
		"accept-encoding": {"gzip, deflate, br"},
		"accept-language": {"zh-CN,zh;q=0.9"},
		"dnt":             {"1"},
		"origin":          {"https://www.yhdmp.cc"},
		"user-agent":      {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36"},
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	h.res = res
	h.reader = res.Body
	return nil
}

func (h *ReadFromUrl) GetReader() io.Reader {
	return h.reader
}

func (h *ReadFromUrl) Close() error {
	return h.res.Body.Close()
}
