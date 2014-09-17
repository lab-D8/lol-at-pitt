package main

import (
	"github.com/codegangsta/martini"
	"github.com/lab-d8/lol-at-pitt/ols"
	"github.com/martini-contrib/render"
	"labix.org/v2/mgo"
	"net/http"
	"net/url"
)

var DatabaseName = "lolpitt"                  // Name that we are storing this all under
var MongoLocation = "mongodb://localhost"     // Location of the db
var TemplatesLocation = "resources/templates" // Location of templates to render..

func main() {
	m := martini.Classic()
	// Setup middleware to be attached to the controllers on every call.
	m.Use(DB())
	m.Use(render.Renderer(render.Options{Directory: TemplatesLocation}))
	m.Use(PARAMS)
	m.Use(martini.Static("public", martini.StaticOptions{Prefix: "/public"}))

	// TODO: Individual variables not sustainable. Need a better system.
	handler := func(mongo *mgo.Database, urls url.Values, renderer render.Render) {
		teams := ols.QueryAllTeams(mongo)
		renderer.HTML(200, "teams", teams)
	}

	handler2 := func(db *mgo.Database, params martini.Params, renderer render.Render) {
		team := ols.QueryTeam(db, params["name"])
		renderer.HTML(200, "team", team)
	}
	m.Get("/teams", handler)
	m.Get("/team/:name", handler2)

	http.ListenAndServe(":80", m)

}

// PARAMS is a middleware binder for injecting the params into each handler
func PARAMS(req *http.Request, c martini.Context) {
	req.ParseForm()
	response := req.Form
	c.Map(response)
	c.Next()
}

// DB is a middleware binder that injects the mongo db into each handler
func DB() martini.Handler {
	session, err := mgo.Dial(MongoLocation)
	if err != nil {
		panic(err)
	}

	return func(c martini.Context) {
		s := session.Clone()
		c.Map(s.DB(DatabaseName))
		defer s.Close()
		c.Next()
	}
}