package main;

import
(
	"log"
	"os"
	"database/sql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

const retry_timeout time.Duration = 1000000000 * 5;

type Environ struct {
	// In the standard postgres:// url format.
	// Required.
	PostgresUrl  string;
	// In the format without the colon.
	Port         string;
};

func main() {
	e := Environ{};
	e.PostgresUrl = os.Getenv("POSTGRES_URL");
	e.Port = os.Getenv("QXM_PORT");
	
	if e.Port == "" {
		e.Port = "8080"
	}
	// repeat until database is live.
	s, err := sql.Open("pgx", e.PostgresUrl);
	for {
		pingerr := s.Ping();
		if pingerr == nil {
			break;
		}

		s.Close();
		log.Printf("%v... Retrying in 5s...", pingerr);
		s, _ = sql.Open("pgx", e.PostgresUrl);
		time.Sleep(retry_timeout);	
	}

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: s}), &gorm.Config{});	
	if err != nil {
		log.Fatal(err);
	}

	err = db.AutoMigrate(&User{}, &Event{})

	h := HttpInitServer(db, e.Port);

	err = h.HttpListen()
	if err != nil {
		log.Print(err)
	}
}
