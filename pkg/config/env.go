package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"sync"
)

var lock = &sync.Mutex{}

type EnvConfig struct {
	Port string `env:"port"`
}

func MustGetEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalln("Missing env variable:", key)
		return ""
	}
	return value
}

func getEnv() *EnvConfig {
	env := EnvConfig{}
	st := reflect.TypeOf(env)
	sv := reflect.ValueOf(&env).Elem()
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		key := field.Tag.Get("env")
		value := MustGetEnv(key)
		sv.Field(i).SetString(value)
	}
	return &env
}

var env *EnvConfig

func initConfig() *EnvConfig {
	if env == nil {
		lock.Lock()
		defer lock.Unlock()
		if env == nil {
			fmt.Println("init config")
			env = getEnv()
		} else {
			fmt.Println("config already initialized")
		}
	} else {
		fmt.Println("config already initialized")
	}
	return env
}

func GetConfig() *EnvConfig {
	return initConfig()
}

var Env = GetConfig()
