package main

import (
	"log"
	"net/url"
	"os"

	badger_server "badger-api/internal/server"
	"badger-api/pkg/server"
)

func GetEnvWithDefault(env_var string, default_val string) string {
	env_val := os.Getenv(env_var)

	if len(env_val) == 0 {
		return default_val
	}

	return env_val
}

func GetMongoConnUri(
	scheme string,
	user string,
	pass string,
	host string,
	path string,
) string {
	uri := url.URL{
		Scheme: scheme,
		Host:   host,
		User:   url.UserPassword(user, pass),
		Path:   path,
	}

	return uri.String()
}

func main() {
	mongo_scheme := os.Getenv("MONGO_SCHEME")
	mongo_user := os.Getenv("MONGO_USER")
	mongo_pass := os.Getenv("MONGO_PASS")
	mongo_host := os.Getenv("MONGO_HOST")
	mongo_path := os.Getenv("MONGO_PATH")

	mongo_conn_uri := GetMongoConnUri(mongo_scheme, mongo_user, mongo_pass, mongo_host, mongo_path)

	log.Println("Creating new server context.")
	ctx := server.NewServerContext(mongo_conn_uri)
	log.Println("Creating new server context done.")
	defer ctx.Cleanup()

	server_host := GetEnvWithDefault("HOST", "0.0.0.0")
	server_port := GetEnvWithDefault("PORT", "8080")
	server_addr := server_host + ":" + server_port

	server := badger_server.NewServer(ctx)

	server.Listen(server_addr)
}
