package main;

import (
	"net/http"
	"time"
	"gorm.io/gorm"
	"github.com/gorilla/mux"
	"log"
	"encoding/json"
	"fmt"
)

type HttpServer struct {
	s *http.Server;
	mux *mux.Router
	db *gorm.DB;
};

type okS struct {
	ok bool `json:"ok"`
}

func httpOkHandler(w http.ResponseWriter, r *http.Request) {
	ok := &okS{ok : true};
	okb, _ := json.Marshal(ok)

	w.Write(okb)
}

func (s *HttpServer) authHandler(w http.ResponseWriter, r *http.Request) {
	u := &User{}
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024);
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&u)

	if(err != nil) {
		log.Printf("JSON Decoding error: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	var count int64

	err = s.db.Session(&gorm.Session{}).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&User{}).Unscoped().Where("name = ?", u.Name).Count(&count)
		if res.Error != nil {
			return res.Error
		}

		if count == 0 {
			res = tx.Create(&u)
			if res != nil {
				return res.Error
			}
		} else {
			res = tx.Where("name = ?", u.Name).First(&u)
			if res != nil {
				return res.Error;
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("GORM TRANSACTION FAILURE: %v", err)
		http.Error(w, "Database op failed", http.StatusInternalServerError);
		return
	}
	rets, err := json.Marshal(u);
	if(err != nil) {
		log.Printf("JSON Encoding error: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	w.Write(rets)
}


func HttpInitServer(db *gorm.DB, port string) (*HttpServer) {
	h := &HttpServer{};
	h.db = db;
	h.mux = mux.NewRouter()
	h.mux.HandleFunc("/ok", httpOkHandler)
	h.mux.HandleFunc("/auth", h.authHandler).Methods("POST")

	//h.mux.HandleFunc("/sql/ok", httpSqlOkHandler);

	h.s = &http.Server{
		Addr:	fmt.Sprintf(":%s", port),
		Handler: h.mux,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	
	log.Printf("Listening on port :%s", port);

	return h;
}

func (self *HttpServer) HttpListen() error {
	return self.s.ListenAndServe()
}
