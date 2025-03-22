package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
)

type Renderer interface {
	Render(w http.ResponseWriter, r *http.Request) error
}

type JsonRenderer struct {
	Data any
}

func (jr *JsonRenderer) Render(w http.ResponseWriter, r *http.Request) error {
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(jr.Data)
}

type TemplRenderer struct {
	Component templ.Component
}

func (tr *TemplRenderer) Render(w http.ResponseWriter, r *http.Request) error {
	w.Header().Add("Content-Type", "text/html")
	return tr.Component.Render(r.Context(), w)
}

// Options
// - text/html, templ -> Render Templ Components
// - application/json, json -> Return Json
func RenderFactory(acceptHeader string, data any) (Renderer, error) {
	switch acceptHeader {
	case "text/html", "templ":
		{
			comp, ok := data.(templ.Component)

			if !ok {
				return nil, fmt.Errorf("invalid data type provided for text/html, %v is not a templ component", data)
			}

			return &TemplRenderer{
				Component: comp,
			}, nil
		}
	case "application/json", "json":
		{
			return &JsonRenderer{
				Data: data,
			}, nil
		}
	default:
		{
			return nil, fmt.Errorf("unsupported header type %s", acceptHeader)
		}
	}
}

// Default Render Method to Render as Either JSON or HTMX Component
func Render(w *http.ResponseWriter, r *http.Request, data any) error {
	acceptHeader := r.Header["Accept"]

	renderer, err := RenderFactory(acceptHeader[0], data)

	if err != nil {
		return err
	}

	return renderer.Render(*w, r)
}

func RenderTempl(w *http.ResponseWriter, r *http.Request, data templ.Component) error {
	renderer := &TemplRenderer{
		Component: data,
	}

	return renderer.Render(*w, r)
}

func RenderJson(w *http.ResponseWriter, r *http.Request, data any) error {
	renderer := &JsonRenderer{
		Data: data,
	}

	return renderer.Render(*w, r)
}

func TemplHandler(comp templ.Component) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := RenderTempl(&w, r, comp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

}
