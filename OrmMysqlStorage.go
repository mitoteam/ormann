package ormann

import _ "github.com/go-sql-driver/mysql"
import "database/sql"

type OrmMysqlStorage struct {
	connection_string string
	db                *sql.DB
}

func (storage *OrmMysqlStorage) Connect(parameters *ormCoreParameters) {
	storage.connection_string = (*parameters)["user"] + ":" + (*parameters)["password"] +
		"@" + (*parameters)["host"] +
		"/" + (*parameters)["database"] +
		"?charset=utf8"

	db, err := sql.Open("mysql", storage.connection_string)
	if err == nil {
		storage.db = db
	}
}

func (storage *OrmMysqlStorage) IsConnected() bool {
	return storage.db != nil
}

func (storage *OrmMysqlStorage) Disconnect() {
	storage.db.Close()
	storage.db = nil
}

func (storage *OrmMysqlStorage) PutObjectData(object *OrmObject) {

}

func (storage *OrmMysqlStorage) GetObjectData(object *OrmObject) {

}
