package dao

import (
	"context"
	"x-server/core/apollo"
	"xy3-proto/pkg/conf/paladin"
	"xy3-proto/pkg/log"

	_ "github.com/go-sql-driver/mysql"

	"x-server/logger/internal/model"

	"gorm.io/gorm"
)

type MysqlConfig struct {
	Addr       string `toml:"addr"`
	User       string `toml:"user"`
	Password   string `toml:"password"`
	DBName     string `toml:"dbname"`
	Parameters string `toml:"parameters"`
	MaxIdle    int    `toml:"maxidle"`
	MaxOpen    int    `toml:"maxopen"`
}

type MssqlConfig struct {
	Addr     string `toml:"addr"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	DBName   string `toml:"dbname"`
	MaxIdle  int    `toml:"maxidle"`
	MaxOpen  int    `toml:"maxopen"`
}

var (
	mssqlDB  *gorm.DB
	mssqlCfg *MssqlConfig
)

func NewDB() (db *gorm.DB, err error) {
	var (
		dbct paladin.TOML
	)
	v := apollo.Get(apollo.DbNS)
	if v == nil || v.Unmarshal(&dbct) != nil {
		err = paladin.Get(apollo.DbNS).Unmarshal(&dbct)
		if err != nil {
			log.Error("db ns unmarshal err %v", err)
			return nil, err
		}
	}
	err = dbct.Get("Mssql").UnmarshalTOML(&mssqlCfg)
	if err != nil {
		log.Error("db ns mssql toml err %v", err)
		return nil, err
	}
	db, err = openMssql(mssqlCfg)
	if err != nil {
		panic(err)
	}
	return db, nil
}

//func openMysql(cfg *MysqlConfig) (db *gorm.DB, err error) {
//
//	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4", cfg.User, cfg.Password, cfg.Addr, cfg.DBName)
//	if cfg.Parameters != "" {
//		dsn = dsn + "&" + cfg.Parameters
//	}
//
//	log.Info("mysql dataSource: %v", dsn)
//	db, err = gorm.Open(mysql.New(mysql.Config{
//		DSN:                       dsn,   // data source name
//		DefaultStringSize:         256,   // default size for string fields
//		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
//		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
//		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
//		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
//
//	}), &gorm.Config{})
//	if err != nil {
//		panic(err)
//	}
//
//	var sqlDB *sql.DB
//	sqlDB, err = db.DB()
//	if err != nil {
//		panic(err)
//	}
//	sqlDB.SetMaxIdleConns(cfg.MaxIdle)
//	sqlDB.SetMaxOpenConns(cfg.MaxOpen)
//	return
//}

func openMssql(cfg *MssqlConfig) (db *gorm.DB, err error) {
	//log.Info("mssql dataSource: %v", cfg)
	//query := url.Values{}
	//query.Add("database", cfg.DBName)
	//query.Add("encrypt", "disable")
	//dns := &url.URL{
	//	Scheme:   "sqlserver",
	//	User:     url.UserPassword(cfg.User, cfg.Password),
	//	Host:     cfg.Addr,
	//	RawQuery: query.Encode(),
	//}
	//db, err = gorm.Open(sqlserver.Open(dns.String()), &gorm.Config{})
	//if err != nil {
	//	panic(err)
	//}
	//
	sqlDB := &gorm.DB{}
	//sqlDB, err = db.DB()
	//if err != nil {
	//	panic(err)
	//}
	//sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	//sqlDB.SetMaxOpenConns(cfg.MaxOpen)
	return sqlDB, nil
}

//func openMssql(cfg *MssqlConfig) (db *gorm.DB, err error) {
//	log.Info("mssql dataSource: %v", cfg)
//	query := url.Values{}
//	query.Add("database", cfg.DBName)
//	query.Add("encrypt", "disable")
//	dns := &url.URL{
//		Scheme:   "sqlserver",
//		User:     url.UserPassword(cfg.User, cfg.Password),
//		Host:     cfg.Addr,
//		RawQuery: query.Encode(),
//	}
//	db, err = gorm.Open(sqlserver.Open(dns.String()), &gorm.Config{})
//	if err != nil {
//		panic(err)
//	}
//
//	var sqlDB *sql.DB
//	sqlDB, err = db.DB()
//	if err != nil {
//		panic(err)
//	}
//	sqlDB.SetMaxIdleConns(cfg.MaxIdle)
//	sqlDB.SetMaxOpenConns(cfg.MaxOpen)
//	return
//}

func GetMysqlDB() *gorm.DB {
	if mssqlDB != nil {
		return mssqlDB
	}
	//
	if mssqlCfg != nil {
		var err error
		mssqlDB, err = openMssql(mssqlCfg)
		if err != nil {
			panic(err)
		}
		return mssqlDB
	}
	panic("not find mssql db")
}

func (d *dao) RawArticle(ctx context.Context, id int64) (art *model.Article, err error) {
	// get data from db
	return
}
