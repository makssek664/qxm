package main;

import (
	"net/http"
	"time"
	"encoding/json"
)

type HttpServer struct {
	s *http.Server;
	mux *http.ServeMux;
	d *Db;
};

type OkS struct {
	ok bool `json:"ok"`
}

func httpOkHandler(w http.ResponseWriter, r *http.Request) {
	ok := &OkS{ok : true};
	okb, _ := json.Marshal(ok)

	w.Write(okb)
}

func httpInitServer(d *Db) (*HttpServer) {
	h := &HttpServer{};
	h.mux = http.NewServeMux()
	h.mux.HandleFunc("/ok", httpOkHandler)
	//h.mux.HandleFunc("/sql/ok", httpSqlOkHandler);

	h.s = &http.Server{
		Addr:	":8080",
		Handler: h.mux,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	
	return h;
}

func (self *HttpServer) httpListen() error {
	return self.s.ListenAndServe()
}
