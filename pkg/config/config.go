package config

import (
	"github.com/alexedwards/scs/v2"
	"html/template"
	"log"
	"myapp2/pkg/models"
)

//AppConfig holds the application config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData
}
