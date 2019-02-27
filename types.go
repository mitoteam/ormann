package ormann

type OrmId int64

type ormCoreParameters map[string]string

type ormType struct {
	IdFieldName string
	TableName   string
	FieldNames  []string
}

type ormTypeRegistry map[string]*ormType
