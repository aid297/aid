package redisPool

type RedisPoolAttr func(redisPool *RedisPool)

func Addr(addr string) RedisPoolAttr { return func(redisPool *RedisPool) { redisPool.SetAddr(addr) } }
func Username(username string) RedisPoolAttr {
	return func(redisPool *RedisPool) { redisPool.SetUsername(username) }
}
func Password(password string) RedisPoolAttr {
	return func(redisPool *RedisPool) { redisPool.SetPassword(password) }
}
func Prefix(prefix string) RedisPoolAttr {
	return func(redisPool *RedisPool) { redisPool.SetPrefix(prefix) }
}
func Pool(clientName string, dbNum int) RedisPoolAttr {
	return func(redisPool *RedisPool) {
		redisPool.SetPool(clientName, dbNum)
	}
}
