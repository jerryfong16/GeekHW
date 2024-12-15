package config

var Config = config{
	DB:    DBConfig{DSN: "root:admin111@tcp(localhost:3306)/webook"},
	Redis: RedisConfig{Addr: "localhost:6379"},
}