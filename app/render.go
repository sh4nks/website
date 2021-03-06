package app

import (
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/unrolled/render"
	"someblocks/context"
)

// D stands for Data and is a shortcut for map[string]interface{}
type D map[string]interface{}

func (app *App) JSON(w http.ResponseWriter, status int, v interface{}) {
	app.Render.JSON(w, status, v)
}

func (app *App) Text(w http.ResponseWriter, status int, v string) {
	app.Render.Text(w, status, v)
}

func (app *App) HTML(w http.ResponseWriter, r *http.Request, tmpl string, data D) {
	htmlOpts := render.HTMLOptions{
		Funcs: template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
			"getFlashedMessages": func() *Flash {
				flash := app.Session.getFlash(r)
				if flash != nil {
					return flash
				}
				return nil
			},
			"isActive": func(endpoint string) template.HTML {
				if r.URL != nil && r.URL.String() == endpoint {
					return "active"
				}
				return ""
			},
		},
	}

	data["CurrentUser"] = context.CurrentUser(r.Context())
	data["RequestEndpoint"] = r.URL.String()
	app.Render.HTML(w, http.StatusOK, tmpl, data, htmlOpts)
}
