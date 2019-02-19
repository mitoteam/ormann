package ormann

type OrmId uint64

type OrmObject interface {
	Id() OrmId

	HasFieldValue(field_name string) bool
	GetFieldValue(field_name string) string
	SetFieldValue(field_name string, value string)

	Save() OrmId
}

//region OrmObjectBase
type OrmObjectBase struct {
  id OrmId
  data map[string]string

  TableName   string
  FieldNames  []string
  IdFieldName string
}

func (o *OrmObjectBase) Init(){
	if o.data == nil {
		o.data = make(map[string]string, len(o.FieldNames))
	}
}

func (o *OrmObjectBase) Id() OrmId{
	return o.id
}

func (o *OrmObjectBase) HasFieldValue(field_name string) bool {
	var _, ok = o.data[field_name]

	return ok
}

func (o *OrmObjectBase) GetFieldValue(field_name string) string {
	var v, ok = o.data[field_name]

	if(ok) {
		return v
	}	else {
		return ""
	}
}

func (o *OrmObjectBase) SetFieldValue(field_name string, value string) {
	o.data[field_name] = value
}

func (o *OrmObjectBase) Save() OrmId {
	return Core().s().PutObjectData(o)
}
//endregion