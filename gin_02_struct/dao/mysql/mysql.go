package mysql

import (
	"fmt"
	"go_gin_advanced/gin_02_struct/settings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var db *sqlx.DB

func Init(cfg *settings.MysqlConfig) (err error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Dbname,
	)

	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("found err in sqlx.Connect", zap.Error(err))
		return
	}

	db.SetMaxIdleConns(cfg.Max_idle_conns)
	db.SetMaxOpenConns(cfg.Max_open_conns)
	return
}

func Close() {
	_ = db.Close()
}
