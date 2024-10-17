//go:build k8s

package config

var Config = config{
	DB: DBConfig{
		DSN: "root:admin111@tcp(webook-mysql:3308)/webook",
	},
	Redis: RedisConfig{
		Addr: "webook-redis:6379",
	},
}
