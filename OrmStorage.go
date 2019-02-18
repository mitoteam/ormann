package ormann

type OrmStorage interface {
	Connect(parameters * ormCoreParameters)
	IsConnected() bool
	Disconnect()

	PutObjectData(object *OrmObject)
	GetObjectData(object *OrmObject)
}
