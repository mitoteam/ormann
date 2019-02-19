package ormann

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mitoteam/mysqlann"
)

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
		mysqlann.SetDB(db)
	}
}

func (storage *OrmMysqlStorage) IsConnected() bool {
	return storage.db != nil
}

func (storage *OrmMysqlStorage) Disconnect() {
	storage.db.Close()
	storage.db = nil
}

func (storage *OrmMysqlStorage) PutObjectData(o *OrmObjectBase) OrmId {
	if o.Id() > 0 { //existing object
		var q = mysqlann.Update(o.TableName).Where(o.IdFieldName, o.Id())

		for _, field_name := range o.FieldNames {
			if o.HasFieldValue(field_name) {
				q.Set(field_name, o.GetFieldValue(field_name))
			}
		}

		q.Exec()
	} else { //new object
		var q = mysqlann.Insert(o.TableName)

		for _, field_name := range o.FieldNames {
			if o.HasFieldValue(field_name) {
				q.Set(field_name, o.GetFieldValue(field_name))
			}
		}

		id, err := q.Exec()
		if err == nil {
			o.id = OrmId(id)
		}
		//fmt.Println("save user:", q.Sql())
	}

	return o.id
}

func (storage *OrmMysqlStorage) GetObjectData(o *OrmObjectBase) {

}

func (storage *OrmMysqlStorage) DeleteObject(o *OrmObjectBase) {
	if o.Id() == 0 {
		panic("can not delete unsaved object")
	}

	var q = mysqlann.Delete(o.TableName).Where(o.IdFieldName, o.Id())
	q.Sql()
}
