package handlers

import (
	"net/http"

	"github.com/SujithSubhash/bookings/pkg/config"
	"github.com/SujithSubhash/bookings/pkg/models"
	"github.com/SujithSubhash/bookings/pkg/render"
)

// Repo is the repository used by the handlers
var Repo *Repository

// its the repository type
type Repository struct {
	App *config.AppConfig
}

// its create a newRepository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// its sets the repository for the new handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "home.page.html", &models.TemplateData{})
}

// About page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	//perform some logic here
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello again!"
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	//send the data to template
	render.RenderTemplate(w, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}
