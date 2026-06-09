package msPool

import "github.com/aid297/aid/v2/gormPool"

type MSSQLPoolAttr func(mssqlPool MSSQLPool)

func Host(host string) MSSQLPoolAttr {
	return func(mssqlPool MSSQLPool) { mssqlPool.SetHost(host) }
}
func Port(port uint16) MSSQLPoolAttr {
	return func(mssqlPool MSSQLPool) { mssqlPool.SetPort(port) }
}
func MaxIdleTime(maxIdleTime int) MSSQLPoolAttr {
	return func(mssqlPool MSSQLPool) { mssqlPool.SetMaxIdleTime(maxIdleTime) }
}
func MaxLifetime(maxLifetime int) MSSQLPoolAttr {
	return func(mssqlPool MSSQLPool) { mssqlPool.SetMaxLifetime(maxLifetime) }
}
func MaxIdleConns(maxIdleConns int) MSSQLPoolAttr {
	return func(mssqlPool MSSQLPool) { mssqlPool.SetMaxIdleConns(maxIdleConns) }
}
func MaxOpenConns(maxOpenConns int) MSSQLPoolAttr {
	return func(mssqlPool MSSQLPool) { mssqlPool.SetMaxOpenConns(maxOpenConns) }
}
func Sources(sources map[string]*gormPool.MSSQLEndpoint) MSSQLPoolAttr {
	return func(mssqlPool MSSQLPool) { mssqlPool.SetSources(sources) }
}
func Replicas(replicas map[string]*gormPool.MSSQLEndpoint) MSSQLPoolAttr {
	return func(mssqlPool MSSQLPool) { mssqlPool.SetReplicas(replicas) }
}
