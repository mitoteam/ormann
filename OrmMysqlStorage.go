package ormann

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mitoteam/mysqlann"
	"reflect"
	"strconv"
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

	db, _ := sql.Open("mysql", storage.connection_string)

	//check connection
	err := db.Ping()

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
	ot := Core().getOrmType(o.ormTypeName)

	if o.Id() > 0 { //existing object
		var q = mysqlann.Update(ot.TableName).Where(ot.IdFieldName, o.Id())

		for _, field_name := range ot.FieldNames {
			if o.HasFieldValue(field_name) {
				q.Set(field_name, o.GetFieldValue(field_name))
			}
		}

		_, _ = q.Exec() //no ID changed by this query, so nothing to assign
	} else { //new object
		var q = mysqlann.Insert(ot.TableName)

		for _, field_name := range ot.FieldNames {
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

func (storage *OrmMysqlStorage) GetObjectData(o *OrmObjectBase) bool {
	if o == nil {
		panic("o is nil")
	}

	if o.id == 0 {
		panic("can not load object with ID=0")
	}

	ot := Core().getOrmType(o.ormTypeName)

	//prepare query
	args := make([]string, len(ot.FieldNames)+1)
	args[0] = "t" //table alias
	copy(args[1:], ot.FieldNames)
	var q = mysqlann.Select(ot.TableName, args...)

	q.Where(ot.IdFieldName, o.Id())

	var err error
	o.data, err = q.QueryRowMap()

	return err != nil
}

func (storage *OrmMysqlStorage) DeleteObject(o *OrmObjectBase) {
	if o.Id() == 0 {
		panic("can not delete unsaved object")
	}

	ot := Core().getOrmType(o.ormTypeName)

	var q = mysqlann.Delete(ot.TableName).Where(ot.IdFieldName, o.Id())
	_, _ = q.Exec()
}

func (storage *OrmMysqlStorage) SelectIdList(empty_o *OrmObjectBase) []OrmId {
	ot := Core().getOrmType(empty_o.ormTypeName)

	var q = mysqlann.Select(ot.TableName, "t", ot.IdFieldName)

	list, _ := q.QueryColumn()

	idList := make([]OrmId, len(list))

	for i := 0; i < len(list); i++ {
		idList[i] = mysqlAnythingToOrmId(list[i])
	}

	return idList
}

func mysqlAnythingToOrmId(value mysqlann.Anything) (r OrmId) {
	switch v := value.(type) {
	case uint:
		r = OrmId(v)
	case uint8:
		r = OrmId(v)
	case uint16:
		r = OrmId(v)
	case uint32:
		r = OrmId(v)
	case uint64:
		r = OrmId(v)
	case int:
		r = OrmId(v)
	case int8:
		r = OrmId(v)
	case int16:
		r = OrmId(v)
	case int32:
		r = OrmId(v)
	case int64:
		r = OrmId(v)
	case string:
		int_v, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			r = OrmId(int_v)
		} else {
			panic("can't convert string to OrmId")
		}
	default:
		panic("unknown type to convert to OrmId: " + reflect.TypeOf(value).Name())
	}

	return r
}
