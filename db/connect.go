package db

import (
	"database/sql"
	"fmt"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// Connection 数据库
type Connection struct {
	db *sql.DB
}

// NewConnection 新建连接
func NewConnection(name, driver, dsn string) (conn *Connection, err error) {

	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Create %s connection successful.\n", name)

	c := new(Connection)
	c.db = db

	return c, nil
}

// GetDBNowTime 获取数据库系统时间
func (conn *Connection) GetDBNowTime() string {

	rows, err := conn.db.Query("select current_time() from dual")

	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	nowTime := ""

	rows.Next()
	rows.Scan(&nowTime)

	return nowTime
}

// Select 根据传入的SQL文, 返回rows对象指针
func (conn *Connection) Select(sql string) (result *sql.Rows, err error) {

	result, err = conn.db.Query(sql)

	if err != nil {
		return nil, err
	}

	return
}

// Insert 执行传入的SQL, 向数据库写入数据.
func (conn *Connection) Insert(sql string) (lastInsertID, rowsAffected int64, err error) {

	ins, err := conn.db.Prepare(sql)
	defer ins.Close()

	if err != nil {
		return 0, 0, err
	}

	result, err := ins.Exec()

	if err != nil {
		return 0, 0, err
	}

	lastInsertID, _ = result.LastInsertId()
	rowsAffected, _ = result.RowsAffected()

	return
}

// Update 执行传入的SQL, 更新数据库数据.
func (conn *Connection) Update(sql string) (rowsAffected int64, err error) {
	ins, err := conn.db.Prepare(sql)
	defer ins.Close()

	if err != nil {
		return 0, err
	}

	result, err := ins.Exec()

	if err != nil {
		return 0, err
	}

	rowsAffected, _ = result.RowsAffected()

	return
}
