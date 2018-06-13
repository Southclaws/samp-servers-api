package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/alecthomas/template"
	"github.com/dyninc/qstring"

	"github.com/Southclaws/samp-servers-api/types"
)

type routeForTemplate struct {
	types.Route
	Version           string
	ParamsSerialised  string
	AcceptsSerialised string
	ReturnsSerialised string
}

const documentationHeader = `# %s

This is an automatically generated documentation page for the %s API endpoints.

`

const documentationRouteTemplate = `## {{ .Name }}

` + "`" + `{{ .Method }}` + "`" + `: ` + "`" + `/{{ .Version }}{{ .Path }}` + "`" + `

{{ .Description }}
{{ if .Params }}
### Query parameters

Example: ` + "`" + `{{ .ParamsSerialised }}` + "`" + `
{{ end }}{{ if .AcceptsSerialised}}
### Accepts

` + "```json" + `
{{ .AcceptsSerialised }}
` + "```" + `
{{ else }}{{ end }}{{ if .ReturnsSerialised}}
### Returns

` + "```json" + `
{{ .ReturnsSerialised }}
` + "```" + `
{{ else }}{{ end }}
`

func docsWrapper(handler types.RouteHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf(documentationHeader, handler.Version(), handler.Version())))
		for _, route := range handler.Routes() {
			docsForRoute(handler.Version(), route, w)
		}
		w.Header().Set("Content-Type", "text/markdown")
		return
	}
}

func docsForRoute(version string, route types.Route, w io.Writer) {
	var err error

	obj := routeForTemplate{
		Route:   route,
		Version: version,
	}

	if route.Params != nil {
		tmp := types.ServerListParams{}
		err = qstring.Unmarshal(route.Params, &tmp)
		if err != nil {
			fmt.Println(err)
			return
		}
		obj.ParamsSerialised, err = qstring.MarshalString(&tmp)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if route.Accepts != nil {
		acceptsSerialised, err2 := json.MarshalIndent(route.Accepts, "", "    ")
		if err2 != nil {
			fmt.Println(err2)
			return
		}
		obj.AcceptsSerialised = string(acceptsSerialised)
	}

	if route.Returns != nil {
		returnsSerialised, err2 := json.MarshalIndent(route.Returns, "", "    ")
		if err2 != nil {
			fmt.Println(err2)
			return
		}
		obj.ReturnsSerialised = string(returnsSerialised)
	}

	tpl, err := template.New("doc").Parse(documentationRouteTemplate)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = tpl.Execute(w, obj)
	if err != nil {
		fmt.Println(err)
		return
	}
}
