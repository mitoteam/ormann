package ormann

type OrmId uint64

type OrmObject interface {
	Id() OrmId
	TableName() string
	GetFieldValue(field_name string) string
	SetFieldValue(field_name string, value string)
	GetFieldsList() []string
	Save() OrmId
}

//region OrmObjectBase
type OrmObjectBase struct {
  id OrmId
  data map[string]string
}

func (*OrmObjectBase) TableName() string{
	return "dummy"
}

func (o *OrmObjectBase) Id() OrmId{
	return o.id
}

func (o *OrmObjectBase) GetFieldValue(field_name string) string{
	if(o.data == nil) {
		o.data = make(map[string]string)
	}

	return o.data[field_name]
}

func (o *OrmObjectBase) SetFieldValue(field_name string, value string) {
	if(o.data == nil) {
		o.data = make(map[string]string)
	}

	o.data[field_name] = value
}

func (o *OrmObjectBase) Save() OrmId {
  return 1;
}
//endregion