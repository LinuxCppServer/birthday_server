package mysql

import (
	. "birthday_server/models/log"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"sync"
)

const (
	MYSQL_CMD_STUB = iota
	MYSQL_CMD_QUERY
	MYSQL_CMD_EXEC
)

var GMysqlConnPool *sql.DB = nil
var mysqlIniter, mysqlReleaser sync.Once

func GetDBConn() *sql.DB {
	dbUserName, dbPasswd, dbName := "root", "root_123", "birthday_gift"
	mysqlIniter.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp/%s?charset=utf8", dbUserName, dbPasswd, dbName)
		conn, err := sql.Open("mysql", dsn)
		if err != nil {
			//log.GLogger.Printf("connect mysql failed: err=%s", err.Error())
			Glog.Error("connect mysql failed: err=%s\n", err.Error())
			os.Exit(-1)
		}
		if err := conn.Ping(); err != nil {
			Glog.Error("connect mysql failed: err=%s\n", err.Error())
			os.Exit(-2)
		}
		conn.SetConnMaxLifetime(100)
		conn.SetMaxIdleConns(10)
		GMysqlConnPool = conn
	})
	return GMysqlConnPool
}

func Cmd(cmd int, sql string, args ...interface{}) (rows *sql.Rows, err error) {
	if cmd <= MYSQL_CMD_STUB {
		return nil, fmt.Errorf("unkown sql cmd: %d", cmd)
	}
	stmt, err := GMysqlConnPool.Prepare(sql)
	defer func() {
		if stmt != nil {
			stmt.Close()
		}
	}()
	if err != nil {
		Glog.Error("mysql prepare failed: %s\n", err.Error())
		return nil, err
	}
	if cmd == MYSQL_CMD_QUERY {
		if len(args) == 0 {
			rows, err = stmt.Query()
		} else {
			rows, err = stmt.Query(args...)
		}
	} else {
		_, err = stmt.Exec(args...)
	}
	return
}

//此接口的query_sql参数不是官方?占位符的用法
func FindRecordById(id int, query_sql string, args ...interface{}) *sql.Rows {
	return nil
}

func ReleaseDB() {
	mysqlReleaser.Do(func() {
		if GMysqlConnPool != nil {
			Glog.Debug("close mysql ok!\n")
			GMysqlConnPool.Close()
			GMysqlConnPool = nil
		}
	})
}

func TestCmd() {
	//query_sql := "update t_guest_money set money=?, name=?, attend_count=? where name=?;"
	//_, err := Cmd(MYSQL_CMD_EXEC, query_sql, 1000, "李婷", 4, "张峰")
	query_sql := "insert into t_guest_money(name, money) values(?,?)"
	_, err := Cmd(MYSQL_CMD_EXEC, query_sql, "王妙")
	if err != nil {
		Glog.Error("mysql query failed: err=%s", err.Error())
		os.Exit(-1)
	}
	/*
		for rows.Next() {
			var id, money, attend_cnt int
			var name string
			err := rows.Scan(&id, &name, &money, &attend_cnt)
			if err != nil {
				log.Printf("mysql scan failed: err=%s", err.Error())
				continue
			}
			log.Printf("id=%d name=%s money=%d attend_count=%d", id, name, money, attend_cnt)
		}
	*/
}
