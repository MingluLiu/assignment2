package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []R

var routes = Routes{
	R{
		"Create",
		"POST",
		"/locations",
		Create,
	},
	R{
		"Query",
		"GET",
		"/locations/{location_id}",
		Query,
	},
	R{
		"Update",
		"PUT",
		"/locations/{location_id}",
		Update,
	},
	R{
		"Delete",
		"DELETE",
		"/locations/{location_id}",
		Delete,
	},
}
