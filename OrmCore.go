package ormann

import (
	"log"
	"reflect"
)

type ormCore struct {
	storage    OrmStorage
	parameters ormCoreParameters
	registry   ormTypeRegistry

	cache map[string]OrmObject //TBD
}

var core *ormCore

/*
ORM Core
*/
func Core() *ormCore {
	if core == nil {
		core = new(ormCore)
		core.init()
	}

	return core
}

/*
Initialization
*/
func (core *ormCore) init() {
	core.parameters = make(ormCoreParameters)
	core.registry = make(ormTypeRegistry)
	core.cache = make(map[string]OrmObject)
}

func (core *ormCore) SetParam(name, value string) *ormCore {
	core.parameters[name] = value

	return core //method chaining
}

/*
Startup
*/
func (core *ormCore) Go() {
	storage_type, ok := core.parameters["storage"]

	if !ok {
		panic("\"storage\" parameter not set")
	}

	if storage_type == "mysql" {
		//core.storage = &OrmMysqlStorage{}
		core.storage = new(OrmMysqlStorage)
		core.storage.Connect(&core.parameters)
		log.Println("Connected to ORM storage")
	} else {
		panic("unknown storage type")
	}
}

/*
Shutdown
*/
func (core *ormCore) Shutdown() {
	if core.storage != nil {
		if core.storage.IsConnected() {
			core.storage.Disconnect()
			core.storage = nil
		}
	}
}

/*
OrmStorage returns connected storage or dies
*/
func (core *ormCore) s() OrmStorage {
	if core.storage == nil {
		core.Go()
		return core.storage
	}

	if !core.storage.IsConnected() {
		core.storage.Connect(&core.parameters)

		return core.storage
	}

	return core.storage
}

func (core *ormCore) registerOrmType(nilPointer interface{}, idFieldName, tableName string, fieldNames []string) (name string) {
	name = reflect.TypeOf(nilPointer).Elem().String()

	t := core.searchOrmType(name)

	if t == nil {
		if len(tableName) == 0 {
			panic("no table name")
		}

		if len(fieldNames) == 0 {
			panic("no fields")
		}

		if len(idFieldName) == 0 {
			idFieldName = "ID"
		}

		newOrmType := ormType{
			IdFieldName: idFieldName,
			TableName:   tableName,
			FieldNames:  fieldNames,
		}

		core.registry[name] = &newOrmType
	}

	return name
}

func (core *ormCore) searchOrmType(name string) (t *ormType) {
	t, ok := core.registry[name]
	if !ok {
		return nil
	}

	return t
}

func (core *ormCore) getOrmType(name string) (t *ormType) {
	t = core.searchOrmType(name)

	if t == nil {
		panic(name + "ORM type not registered")
	}

	return t
}
