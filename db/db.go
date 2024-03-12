package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

const (
	dataBaseName = "root:1qazzaq1@tcp(127.0.0.1:3306)/DB_VIRTMANAGER"
	CreateVm     = `
		CREATE TABLE IF NOT EXISTS vm_info (
		id INT AUTO_INCREMENT PRIMARY KEY,
		create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		device_ip VARCHAR(64),
		vm_name VARCHAR(64),
		vm_ip VARCHAR(64),
		vm_cpu VARCHAR(64),
	    vm_mem VARCHAR(64),
		vm_disk VARCHAR(64),
		image_type VARCHAR(64),
		image_name VARCHAR(2048),
		vm_remarks VARCHAR(2048),
		expand1 VARCHAR(64),
		expand2 VARCHAR(64),
		expand3 VARCHAR(64)
	);`
)

type VirtSql struct {
	Db *sql.DB
}

var initSqlList = []string{
	CreateVm,
}

func DBInit() (*VirtSql, error) {
	db, err := sql.Open("mysql", dataBaseName)
	if err != nil {
		return nil, err
	}
	if err = softTableInit(db); err != nil {
		return nil, err
	}

	return &VirtSql{
		Db: db,
	}, nil
}

func softTableInit(db *sql.DB) error {
	for i := 0; i < len(initSqlList); i++ {
		_, err := db.Exec(initSqlList[i])
		if err != nil {
			logrus.Errorf("create table error, sql:[%s]", initSqlList[i])
			return err
		}
	}
	logrus.Debugf("create tool tables success !!!")
	return nil
}
