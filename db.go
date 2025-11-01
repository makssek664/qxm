package main;
import ( 
	"context"
	"github.com/jackc/pgx/v5"
)

type Db struct {
	db *pgx.Conn;
};

type DbInitArgs struct {
	postgres_uri string
		
};

func DbInit(a DbInitArgs) (*Db, error) {
	c := &Db{}
	db, err := pgx.Connect(context.Background(), a.postgres_uri)
	if err != nil { return nil, err }
	c.db = db;
	return c, nil
}
