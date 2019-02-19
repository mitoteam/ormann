package ormann

type OrmStorage interface {
	Connect(parameters * ormCoreParameters)
	IsConnected() bool
	Disconnect()

	PutObjectData(*OrmObjectBase) OrmId
	GetObjectData(o *OrmObjectBase) bool
	DeleteObject(o *OrmObjectBase)
}
