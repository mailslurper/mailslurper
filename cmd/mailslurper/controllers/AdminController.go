package controllers

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo"
	"github.com/mailslurper/mailslurper/pkg/mailslurper"
	"github.com/mailslurper/mailslurper/pkg/ui"
	"github.com/sirupsen/logrus"
)

/*
AdminController provides methods for handling admin endpoints.
This is to primarily support the front-end
*/
type AdminController struct {
	config         *mailslurper.Configuration
	configFileName string
	debugMode      bool
	renderer       *ui.TemplateRenderer
	lock           *sync.Mutex
	logger         *logrus.Entry
	serverVersion  string
}

/*
NewAdminController creates a new admin controller
*/
func NewAdminController(logger *logrus.Entry, renderer *ui.TemplateRenderer, serverVersion string, config *mailslurper.Configuration, configFileName string, debugMode bool) *AdminController {
	return &AdminController{
		config:         config,
		configFileName: configFileName,
		debugMode:      debugMode,
		lock:           &sync.Mutex{},
		logger:         logger,
		renderer:       renderer,
		serverVersion:  serverVersion,
	}
}

/*
Admin is the page for performing administrative tasks in MailSlurper
*/
func (c *AdminController) Admin(ctx echo.Context) error {
	data := mailslurper.Page{
		Theme: c.config.GetTheme(),
		Title: "Admin",
	}

	return ctx.Render(http.StatusOK, "mainLayout:admin", data)
}

/*
ApplyTheme updates the theme in the config file, and refreshes the renderer

	POST: /theme
*/
func (c *AdminController) ApplyTheme(ctx echo.Context) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	var err error
	var applyThemeRequest *mailslurper.ApplyThemeRequest

	if applyThemeRequest, err = mailslurper.NewApplyThemeRequest(ctx); err != nil {
		return ctx.String(http.StatusBadRequest, "Invalid request")
	}

	c.config.Theme = applyThemeRequest.Theme

	if err = c.config.SaveConfiguration(c.configFileName); err != nil {
		c.logger.Errorf("Error saving configuration file in ApplyTheme: %s", err.Error())
		return ctx.String(http.StatusOK, fmt.Sprintf("Error saving configuration file: %s", err.Error()))
	}

	c.renderer.LoadTemplates(c.debugMode)
	return ctx.String(http.StatusOK, "OK")
}

/*
Index is the main view. This endpoint provides the email list and email detail
views.
*/
func (c *AdminController) Index(ctx echo.Context) error {
	data := mailslurper.Page{
		Theme: c.config.GetTheme(),
		Title: "Mail",
	}

	return ctx.Render(http.StatusOK, "mainLayout:index", data)
}

/*
ManageSavedSearches is the page for managing saved searches
*/
func (c *AdminController) ManageSavedSearches(ctx echo.Context) error {
	data := mailslurper.Page{
		Theme: c.config.GetTheme(),
		Title: "Manage Saved Searches",
	}

	return ctx.Render(http.StatusOK, "mainLayout:manageSavedSearches", data)
}

/*
GetPruneOptions returns a set of valid pruning options.

	GET: /v1/pruneoptions
*/
func (c *AdminController) GetPruneOptions(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, mailslurper.PruneOptions)
}

/*
GetServiceSettings returns the settings necessary to talk to the MailSlurper
back-end service tier.
*/
func (c *AdminController) GetServiceSettings(ctx echo.Context) error {
	settings := mailslurper.ServiceSettings{
		IsSSL:          c.config.IsServiceSSL(),
		ServiceAddress: c.config.ServiceAddress,
		ServicePort:    c.config.ServicePort,
		Version:        c.serverVersion,
	}

	return ctx.JSON(http.StatusOK, settings)
}

/*
GetVersion outputs the current running version of this MailSlurper server instance
*/
func (c *AdminController) GetVersion(ctx echo.Context) error {
	result := mailslurper.Version{
		Version: c.serverVersion,
	}

	return ctx.JSON(http.StatusOK, result)
}

/*
GetVersionFromMaster returns the current MailSlurper version from GitHub
*/
func (c *AdminController) GetVersionFromMaster(ctx echo.Context) error {
	var err error
	var result *mailslurper.Version

	if result, err = mailslurper.GetServerVersionFromMaster(); err != nil {
		c.logger.Errorf("Error getting version file from Github: %s", err.Error())
		return ctx.String(http.StatusInternalServerError, "There was an error reading the version file from GitHub")
	}

	return ctx.JSON(http.StatusOK, result)
}
