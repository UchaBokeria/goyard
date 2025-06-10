package htmx

import "github.com/a-h/templ"

func Link(path string) templ.Attributes {
	return templ.Attributes{
		"hx-get":       path,
		"hx-swap":      "innerHTML show:window:top",
		"hx-push-url":  "true",
		"hx-target":    "#Content",
		"hx-encoding":  "text/html",
		"hx-indicator": ".Loading",
		"hx-trigger":   "click",
	}
}

func PostLink(path string, params string) templ.Attributes {
	x := templ.Attributes{
		"hx-post":      path,
		"hx-swap":      "innerHTML show:window:top",
		"hx-push-url":  "true",
		"hx-target":    "#Content",
		"hx-encoding":  "text/html",
		"hx-indicator": ".Loading",
	}
	if params != "" {
		x["hx-params"] = params
	}
	return x
}
