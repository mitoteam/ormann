package ormann

import (
	"reflect"
	"fmt"
)

type OrmCore struct
{
	cache map[string]OrmObject

	types map[string]reflect.Type
}

var core *OrmCore

/*
ORM Core
*/
func Core() *OrmCore{
	if(core == nil) {
		core = new(OrmCore)
		core.init()
	}

	return core
}

func typed_nil_to_name(typedNil OrmObject) string{
	var t = reflect.TypeOf(typedNil)

	if(t.Kind() == reflect.Ptr){
		t = t.Elem()
	}

 return t.PkgPath() + "." + t.Name()
}

/*
Initialization
*/
func (core *OrmCore) init(){
	core.cache = make(map[string]OrmObject)
	core.types = make(map[string]reflect.Type)
}

func (core *OrmCore) RegisterType(typedNil OrmObject){
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
func (core *OrmCore) Create(type_name string) *OrmObject{
 var o = reflect.New(core.types[type_name])

 fmt.Println(o)

 return nil;
}

/*
Load
*/
func (core *OrmCore) Load(type_name string, id OrmId) *OrmObject{
  return nil;
}

/*
Save
*/
func (core *OrmCore) Save(o *OrmObject) OrmId{
	return 1;
}