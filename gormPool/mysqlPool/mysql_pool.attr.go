package mysqlPool

import "github.com/aid297/aid/v2/gormPool"

type MySQLPoolAttr func(mysqlPool MySQLPool)

func Username(username string) MySQLPoolAttr {
	return func(mysqlPool MySQLPool) { mysqlPool.SetUsername(username) }
}
func Password(password string) MySQLPoolAttr {
	return func(mysqlPool MySQLPool) { mysqlPool.SetPassword(password) }
}
func Host(host string) MySQLPoolAttr { return func(mysqlPool MySQLPool) { mysqlPool.SetHost(host) } }
func Port(port uint16) MySQLPoolAttr { return func(mysqlPool MySQLPool) { mysqlPool.SetPort(port) } }
func Sources(sources map[string]gormPool.MySQLEndpoint) MySQLPoolAttr {
	return func(mysqlPool MySQLPool) { mysqlPool.SetSources(sources) }
}
func Replicas(replicas map[string]gormPool.MySQLEndpoint) MySQLPoolAttr {
	return func(mysqlPool MySQLPool) { mysqlPool.SetReplicas(replicas) }
}
func MaxIdleTime(maxIdleTime int) MySQLPoolAttr {
	return func(mysqlPool MySQLPool) { mysqlPool.SetMaxIdleTime(maxIdleTime) }
}
func MaxLifetime(maxLifetime int) MySQLPoolAttr {
	return func(mysqlPool MySQLPool) { mysqlPool.SetMaxLifetime(maxLifetime) }
}
func MaxIdleConnections(maxIdleConnections int) MySQLPoolAttr {
	return func(mysqlPool MySQLPool) { mysqlPool.SetMaxIdleConnections(maxIdleConnections) }
}
func MaxOpenConnections(maxOpenConnections int) MySQLPoolAttr {
	return func(mysqlPool MySQLPool) { mysqlPool.SetMaxOpenConnections(maxOpenConnections) }
}
