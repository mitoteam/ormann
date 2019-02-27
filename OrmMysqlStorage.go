package ormann

import (
	"database/sql"
	"strconv"
	"reflect"
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
	if o.Id() > 0 { //existing object
		var q = mysqlann.Update(o.TableName).Where(o.IdFieldName, o.Id())

		for _, field_name := range o.FieldNames {
			if o.HasFieldValue(field_name) {
				q.Set(field_name, o.GetFieldValue(field_name))
			}
		}

		q.Exec() //no ID changed by this query, so nothing to assign
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

func (storage *OrmMysqlStorage) GetObjectData(o *OrmObjectBase) bool {
	if o == nil {
		panic("o is nil")
	}

	if o.id == 0 {
		panic("can not load object with ID=0")
	}

	//prepare query
	args := make([]string, len(o.FieldNames) + 1)
	args[0] = "t" //table alias
	copy(args[1:], o.FieldNames)
	var q = mysqlann.Select(o.TableName, args...)

	q.Where(o.IdFieldName, o.Id())

	var err error
	o.data, err = q.QueryRowMap()

	return err != nil
}

func (storage *OrmMysqlStorage) DeleteObject(o *OrmObjectBase) {
	if o.Id() == 0 {
		panic("can not delete unsaved object")
	}

	var q = mysqlann.Delete(o.TableName).Where(o.IdFieldName, o.Id())
	q.Exec()
}

func (storage *OrmMysqlStorage) SelectIdList(empty_o *OrmObjectBase) []OrmId {
	var q = mysqlann.Select(empty_o.TableName, "t", empty_o.IdFieldName)

	list, _ := q.QueryColumn()

	id_list := make([]OrmId, len(list))

	for i := 0; i < len(list); i++ {
		id_list[i] = mysqlAnythingToOrmId(list[i])
	}

	return id_list
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