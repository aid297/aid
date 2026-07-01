package pgPool

import "github.com/aid297/aid/v2/gormPools"

type PGPoolAttr func(*PGPoolImpl)

func Host(host string) PGPoolAttr     { return func(pgPool *PGPoolImpl) { pgPool.SetHost(host) } }
func Port(port uint16) PGPoolAttr     { return func(pgPool *PGPoolImpl) { pgPool.SetPort(port) } }
func SSLMode(sslMode bool) PGPoolAttr { return func(pgPool *PGPoolImpl) { pgPool.SetSSLMode(sslMode) } }
func TimeZone(timezone string) PGPoolAttr {
	return func(pgPool *PGPoolImpl) { pgPool.SetTimeZone(timezone) }
}
func MaxIdleTime(maxIdleTime int) PGPoolAttr {
	return func(pgPool *PGPoolImpl) { pgPool.SetMaxIdleTime(maxIdleTime) }
}
func MaxLifetime(maxLifetime int) PGPoolAttr {
	return func(pgPool *PGPoolImpl) { pgPool.SetMaxLifetime(maxLifetime) }
}
func MaxIdleConns(maxIdleConns int) PGPoolAttr {
	return func(pgPool *PGPoolImpl) { pgPool.SetMaxIdleConns(maxIdleConns) }
}
func MaxOpenConns(maxOpenConns int) PGPoolAttr {
	return func(pgPool *PGPoolImpl) { pgPool.SetMaxOpenConns(maxOpenConns) }
}
func Sources(sources map[string]*gormPools.PGEndpoint) PGPoolAttr {
	return func(pgPool *PGPoolImpl) { pgPool.SetSources(sources) }
}
func Replicas(replicas map[string]*gormPools.PGEndpoint) PGPoolAttr {
	return func(pgPool *PGPoolImpl) { pgPool.SetReplicas(replicas) }
}
