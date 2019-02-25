package ormann

import (
	"log"
)

type ormCoreParameters map[string]string

type ormCore struct {
	storage    OrmStorage
	parameters ormCoreParameters

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
returns storage or dies
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
