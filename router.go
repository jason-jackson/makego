package main

import (
	"fmt"
	"strings"
)

type Router struct {
	Name    string
	Matches []string

	Main     string
	Object   string
	App      string
	HomeView string

	Package string
}

func findRouter(name string) (Router, error) {
	name = strings.ToLower(name)

	if r, ok := routers[name]; ok {
		return r, nil
	}

	for k, r := range routers {
		r.Name = k
		if strings.EqualFold(name, r.Name) {
			return r, nil
		}

		for _, m := range r.Matches {
			if strings.EqualFold(name, m) {
				return r, nil
			}
		}
	}

	return Router{}, fmt.Errorf("no router matching: %s", name)
}

var routers = map[string]Router{
	"echo": {
		Name:   "echo",
		Main:   `return e.Start(":8080")`,
		Object: "*echo.Echo",
		App: `func App() *echo.Echo {
	app := echo.New()

	app.GET("/", homeView)

	return app
}`,
		HomeView: `func homeView(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome!")
}`,
		Package: "github.com/labstack/echo/v4",
	},
	"gin": {
		Name:   "gin",
		Main:   `return router.Run(":8080")`,
		Object: "*gin.Engine",
		App: `func App() *gin.Engine {
	app := gin.Default()

	app.GET("/", homeView)

	return app
}`,
		HomeView: `func homeView(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome!",})
}`,
		Package: "github.com/gin-gonic/gin",
	},
	"mux": {
		Name:   "mux",
		Main:   `return http.ListenAndServe(":8080", router)`,
		Object: "*mux.Router",
		App: `func App() {{ *mux.Router }} {
	app := mux.NewRouter()

	app.HandleFunc("/", homeView).Methods(http.MethodGet)

	return app
}

func json(w http.ResponseWriter, status int, body any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	contents, err := json.Marshal(body)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, contents)
	return nil
}
`,
		HomeView: `func homeView(w http.ResponseWriter, r *http.Request) {
	json(w, http.StatusOK, "Welcome!")
}`,
		Package: "github.com/gorilla/mux",
	},
}
