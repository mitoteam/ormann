package ormann

import "reflect"

type OrmObject interface {
	Id() OrmId

	HasFieldValue(field_name string) bool
	GetFieldValue(field_name string) interface{}
	SetFieldValue(field_name string, value interface{})

	Save() OrmId
	Load(id OrmId) bool
	MustLoad(id OrmId)
	Delete()
	//Select() []OrmObject
}

//region OrmObjectBase
type OrmObjectBase struct {
	id          OrmId
	data        map[string]interface{}
	ormTypeName string
}

func (o *OrmObjectBase) Init(nilPointer interface{}, idFieldName, tableName string, fieldNames []string) {
	o.ormTypeName = Core().registerOrmType(nilPointer, idFieldName, tableName, fieldNames)

	if o.data == nil {
		o.data = make(map[string]interface{}, len(fieldNames))
	}
}

func (o *OrmObjectBase) Id() OrmId {
	return o.id
}

func (o *OrmObjectBase) HasFieldValue(field_name string) bool {
	var _, ok = o.data[field_name]

	return ok
}

func (o *OrmObjectBase) GetFieldValue(field_name string) interface{} {
	var v, ok = o.data[field_name]

	if ok {
		return v
	} else {
		return ""
	}
}

func (o *OrmObjectBase) SetFieldValue(field_name string, value interface{}) {
	o.data[field_name] = value
}

func (o *OrmObjectBase) Save() OrmId {
	o.id = Core().s().PutObjectData(o)

	return o.id
}

func (o *OrmObjectBase) Load(id OrmId) bool {
	o.id = id
	return Core().s().GetObjectData(o)
}

func (o *OrmObjectBase) MustLoad(id OrmId) {
	if !Core().s().GetObjectData(o) {
		panic("can not load object")
	}
}

func (o *OrmObjectBase) Delete() {
	Core().s().DeleteObject(o)
}

//endregion

//region Fetching lists
func Select(emptyO interface{}, list interface{}) {
	var emptyObjectPointerT reflect.Type = reflect.TypeOf(emptyO)
	var emptyObjectPointerV reflect.Value = reflect.ValueOf(emptyO)
	var emptyObjectV reflect.Value = emptyObjectPointerV.Elem()

	var oobValue = emptyObjectV.FieldByName("OrmObjectBase")

	var o *OrmObjectBase = oobValue.Addr().Interface().(*OrmObjectBase)

	idList := Core().s().SelectIdList(o)

	slice := reflect.MakeSlice(reflect.SliceOf(emptyObjectPointerT), len(idList), len(idList))

	for i := 0; i < len(idList); i++ {
		newObjectPointerV := reflect.New(emptyObjectPointerT.Elem())
		newObjectPointerV.Elem().Set(emptyObjectV)

		_, ok := emptyObjectPointerT.MethodByName("Load")
		if ok {
			newObjectPointerV.MethodByName("Load").Call([]reflect.Value{reflect.ValueOf(idList[i])})
		}

		slice.Index(i).Set(newObjectPointerV)
	}

	listValue := reflect.ValueOf(list)
	listValue.Elem().Set(slice)
}

//endregion
