package db

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	_ "github.com/WlayRay/order-demo/common/config"
	_ "github.com/lib/pq" // 驱动导入
	"github.com/spf13/viper"
)

var (
	host       = viper.GetString("postgres.host")
	port       = viper.GetInt("postgres.port")
	user       = viper.GetString("postgres.user")
	password   = viper.GetString("postgres.password")
	dbName     = viper.GetString("postgres.dbname")
	searchPath = viper.GetString("postgres.search-path")
)

func GetPGSQLConn() (*sql.Driver, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s search_path=%s sslmode=disable connect_timeout=5",
		host, port, user, password, dbName, searchPath,
	)

	drv, err := sql.Open(dialect.Postgres, dsn)
	if err != nil {
		return nil, err
	}

	db := drv.DB()
	db.SetMaxOpenConns(80)
	db.SetMaxIdleConns(15)
	db.SetConnMaxLifetime(20 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	// 验证数据库连接
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return drv, nil
}
