package main

import (
	"fmt"
	"strings"
)

type Database struct {
	Name      string
	Matches   []string
	Version   string
	Admin     string
	Port      string
	DockerEnv string
}

func findDatabase(name string) (Database, error) {
	name = strings.ToLower(name)

	if db, ok := databases[name]; ok {
		return db, nil
	}

	for k, db := range databases {
		db.Name = k

		if strings.EqualFold(name, db.Name) {
			return db, nil
		}

		for _, m := range db.Matches {
			if strings.EqualFold(name, m) {
				return db, nil
			}
		}
	}

	return Database{}, fmt.Errorf("no database matching: %s", name)
}

var databases = map[string]Database{
	"mariadb": {
		Name:    "mariadb",
		Version: "latest",
		Admin:   "phpmyadmin",
		Port:    "3306",
		DockerEnv: `MYSQL_ROOT_PASSWORD: r00tp@ss!
      MYSQL_DATABASE: app
      MYSQL_USER: user
      MYSQL_PASSWORD: pass`,
	},
	"mysql": {
		Name:    "mysql",
		Version: "latest",
		Admin:   "phpmyadmin",
		Port:    "3306",
		DockerEnv: `MYSQL_ROOT_PASSWORD: r00tp@ss!
      MYSQL_DATABASE: app
      MYSQL_USER: user
      MYSQL_PASSWORD: pass`,
	},
	"postgres": {
		Name:    "postgres",
		Matches: []string{"pg", "postgres", "postgresql"},
		Version: "latest",
		Admin:   "adminer",
		Port:    "5432",
		DockerEnv: `POSTGRES_USER: user
      POSTGRES_PASSWORD: pass
      POSTGRES_DB: app`,
	},
}
