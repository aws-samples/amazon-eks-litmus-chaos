package handler

import (
	"counter-service/common"
	"counter-service/repository"
	"io"
	"log"
	"net/http"
)

type Counter struct {
	l        *log.Logger
	database *repository.Database
}

func NewCounter(l *log.Logger, d *repository.Database) *Counter {
	return &Counter{l: l,
		database: d}
}

func GetHealthz(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Healthy")
}

func (c *Counter) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		c.getCount(rw, r)
		return
	}

	if r.Method == http.MethodPost {
		c.addCount(rw, r)
		return
	}

	// catch all other http verb with 405
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (c *Counter) getCount(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle %s %s", r.Method, r.URL)

	count, err := c.database.GetCount(r.Context())
	if err != nil {
		c.l.Printf("[error] %v", err.Error())
		http.Error(rw, "Error retrieving count", http.StatusInternalServerError)
		return
	}

	countResp := CounterResponse{
		Count:    count,
		Hostname: common.Hostname,
	}
	err = countResp.ToJSON(rw)
	if err != nil {
		c.l.Printf("[error] %v", err.Error())
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
}

func (c *Counter) addCount(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle %s %s", r.Method, r.URL)

	count, err := c.database.IncrAndGetCount(r.Context())
	if err != nil {
		http.Error(rw, "Error incrementing and retrieving count", http.StatusInternalServerError)
	}

	countResp := CounterResponse{
		Count:    count,
		Hostname: common.Hostname,
	}
	err = countResp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}
