package main

import (
	"fmt"
	"strings"
)

type ORM struct {
	Name     string
	Matches  []string
	Object   string
	Package  string
	DBDriver map[string]string
	Driver   string
	Init     string
}

func findORM(name string) (ORM, error) {
	name = strings.ToLower(name)

	if o, ok := orms[name]; ok {
		return o, nil
	}

	for k, o := range orms {
		o.Name = k

		if strings.EqualFold(name, o.Name) {
			return o, nil
		}

		for _, m := range o.Matches {
			if strings.EqualFold(name, m) {
				return o, nil
			}
		}
	}

	return ORM{}, fmt.Errorf("no ORM matching: %s", name)
}

var orms = map[string]ORM{
	"gorm": {
		Name:    "gorm",
		Object:  "*gorm.DB",
		Package: "gorm.io/gorm",
		Driver:  "gorm.io/driver/",
		DBDriver: map[string]string{
			"mysql":     "mysql",
			"postgres":  "postgres",
			"sqlite":    "sqlite",
			"sqlserver": "sqlserver",
			"tidb":      "mysql",
		},
		Init: `DB, err = gorm.Open({{ .DBDriver }}.Open(app.DatabaseDsn), &gorm.Config{})`,
	},
}
