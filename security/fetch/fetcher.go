package fetch

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func mapToString(m map[string]string) string {
	s := ""
	for k, v := range m {
		s += k + "=" + v + "&"
	}
	return s
}

type BaseRequest struct {
	method       string
	url          string
	headers      map[string]string
	params       map[string]string
	data         interface{}
	json_payload interface{}

	request  *http.Request
	app_code string
}

func NewBaseRequest(method string, url string, headers map[string]string, params map[string]string, data interface{}, json_payload interface{},
	app_code string) *BaseRequest {
	r := &BaseRequest{
		method:       method,
		url:          url,
		headers:      headers,
		params:       params,
		data:         data,
		json_payload: json_payload,
		app_code:     app_code,
	}
	var err error
	r.request, err = http.NewRequest(method, url, nil)
	if err != nil {
		log.Println(err)
	}
	return r
}

func (r *BaseRequest) Execute(retries_left int, is_retry bool) (http.Response, error) {
	starttime_ns := time.Now().UnixNano()
	log_request := map[string]string{
		"method":  r.method,
		"url":     r.url,
		"params":  mapToString(r.params),
		"headers": mapToString(r.headers),
		"time_ns": fmt.Sprintf("%d", starttime_ns),
	}
	log.Printf("%s", log_request)

	response, err := http.DefaultClient.Do(r.request)
	if err != nil {
		return http.Response{}, err
	}
	endtime_ns := time.Now().UnixNano()
	log_response := map[string]string{
		"status":  response.Status,
		"time_ns": fmt.Sprintf("%d", endtime_ns),
	}
	log.Printf("%s", log_response)
	return *response, err
}
