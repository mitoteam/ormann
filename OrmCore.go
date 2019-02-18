package ormann

import (
	"reflect"
	"fmt"
)

type ormCoreParameters map[string]string

type ormCore struct
{
	storage OrmStorage
	parameters ormCoreParameters

	cache map[string]OrmObject
	types map[string]reflect.Type
}

var core *ormCore

/*
ORM Core
*/
func Core() *ormCore {
	if(core == nil) {
		core = new(ormCore)
		core.init()
	}

	return core
}

/*
Initialization
*/
func (core *ormCore) init(){
	core.parameters = make(ormCoreParameters)
	core.cache = make(map[string]OrmObject)
	core.types = make(map[string]reflect.Type)
}

func (core *ormCore) SetParam(name, value string) *ormCore {
	core.parameters[name] = value

	return core; //method chaining
}

/*
Startup
*/
func (core *ormCore) Go(){
	storage_type, ok := core.parameters["storage"]

	if(!ok){
		panic("\"storage\" parameter not set")
	}

	if storage_type == "mysql" {
		//core.storage = &OrmMysqlStorage{}
		core.storage = new(OrmMysqlStorage)
		core.storage.Connect(&core.parameters)
	}	else{
		panic("unknown storage type")
	}
}

/*
Shutdown
*/
func (core *ormCore) Shutdown(){
	if core.storage != nil {
		if(core.storage.IsConnected()){
			core.storage.Disconnect()
		}
	}
}

func (core *ormCore) RegisterType(typedNil OrmObject){
	var t = reflect.TypeOf(typedNil)

	if(t.Kind() == reflect.Ptr){
		t = t.Elem()
	}

	type_name := typed_nil_to_name(typedNil)

	core.types[type_name] = t
}

/*
Create
*/
func (core *ormCore) Create(type_name string) *OrmObject{
 var o = reflect.New(core.types[type_name])

 fmt.Println(o)

 return nil;
}

/*
Load
*/
func (core *ormCore) Load(type_name string, id OrmId) *OrmObject{
  return nil;
}

/*
Save
*/
func (core *ormCore) Save(o *OrmObject) OrmId{
	return 1;
}

func typed_nil_to_name(typedNil OrmObject) string{
	var t = reflect.TypeOf(typedNil)

	if(t.Kind() == reflect.Ptr){
		t = t.Elem()
	}

	return t.PkgPath() + "." + t.Name()
}

//region sdfsdfg

//endregion