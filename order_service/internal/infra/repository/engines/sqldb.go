package engines

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/app/config"
)

type SqlDbEngine struct {
	Client *sql.DB
}

var dbEngine *SqlDbEngine

func GetSqlDbEngine() *SqlDbEngine {
	return dbEngine
}

func SetupSqlDBEngine(cfg *config.AppConfig) (*SqlDbEngine, error) {
	if dbEngine == nil {
		dbEngine = new(SqlDbEngine)
		ConnectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", cfg.SqlUser, cfg.SqlPassword, cfg.SqlHost, cfg.SqlPort, cfg.SqlDatabaseName)

		db, err := sql.Open("mysql", ConnectionString)
		if err != nil {
			return nil, err
		}
		engine := &SqlDbEngine{}
		engine.Client = db
		dbEngine = engine
		return dbEngine, nil
	}
	return dbEngine, nil
}
