package controllers

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/mailslurper/mailslurper/pkg/auth/auth"
	"github.com/mailslurper/mailslurper/pkg/auth/authfactory"
	"github.com/mailslurper/mailslurper/pkg/contexts"
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
	context := contexts.GetAdminContext(ctx)

	data := mailslurper.Page{
		Theme: c.config.GetTheme(),
		Title: "Admin",
		User:  context.User,
	}

	return context.Render(http.StatusOK, "mainLayout:admin", data)
}

/*
ApplyTheme updates the theme in the config file, and refreshes the renderer

	POST: /theme
*/
func (c *AdminController) ApplyTheme(ctx echo.Context) error {
	context := contexts.GetAdminContext(ctx)
	c.lock.Lock()
	defer c.lock.Unlock()

	var err error
	var applyThemeRequest *mailslurper.ApplyThemeRequest

	if applyThemeRequest, err = mailslurper.NewApplyThemeRequest(context); err != nil {
		return context.String(http.StatusBadRequest, "Invalid request")
	}

	c.config.Theme = applyThemeRequest.Theme

	if err = c.config.SaveConfiguration(c.configFileName); err != nil {
		c.logger.Errorf("Error saving configuration file in ApplyTheme: %s", err.Error())
		return context.String(http.StatusOK, fmt.Sprintf("Error saving configuration file: %s", err.Error()))
	}

	c.renderer.LoadTemplates(c.debugMode)
	return context.String(http.StatusOK, "OK")
}

/*
Index is the main view. This endpoint provides the email list and email detail
views.
*/
func (c *AdminController) Index(ctx echo.Context) error {
	context := contexts.GetAdminContext(ctx)

	data := mailslurper.Page{
		Theme: c.config.GetTheme(),
		Title: "Mail",
		User:  context.User,
	}

	return context.Render(http.StatusOK, "mainLayout:index", data)
}

/*
Login renders the login page
*/
func (c *AdminController) Login(ctx echo.Context) error {
	data := mailslurper.Page{
		Theme: c.config.GetTheme(),
	}

	if ctx.QueryParam("message") != "" {
		data.Message = ctx.QueryParam("message")
		data.Error = true
	}

	return ctx.Render(http.StatusOK, "loginLayout:login", data)
}

/*
ManageSavedSearches is the page for managing saved searches
*/
func (c *AdminController) ManageSavedSearches(ctx echo.Context) error {
	context := contexts.GetAdminContext(ctx)

	data := mailslurper.Page{
		Theme: c.config.GetTheme(),
		Title: "Manage Saved Searches",
		User:  context.User,
	}

	return context.Render(http.StatusOK, "mainLayout:manageSavedSearches", data)
}

/*
GetPruneOptions returns a set of valid pruning options.

	GET: /v1/pruneoptions
*/
func (c *AdminController) GetPruneOptions(ctx echo.Context) error {
	context := contexts.GetAdminContext(ctx)
	return context.JSON(http.StatusOK, mailslurper.PruneOptions)
}

/*
GetServiceSettings returns the settings necessary to talk to the MailSlurper
back-end service tier.
*/
func (c *AdminController) GetServiceSettings(ctx echo.Context) error {
	context := contexts.GetAdminContext(ctx)

	settings := mailslurper.ServiceSettings{
		IsSSL:          c.config.IsServiceSSL(),
		ServiceAddress: c.config.ServiceAddress,
		ServicePort:    c.config.ServicePort,
		Version:        c.serverVersion,
	}

	return context.JSON(http.StatusOK, settings)
}

/*
GetVersion outputs the current running version of this MailSlurper server instance
*/
func (c *AdminController) GetVersion(ctx echo.Context) error {
	context := contexts.GetAdminContext(ctx)

	result := mailslurper.Version{
		Version: c.serverVersion,
	}

	return context.JSON(http.StatusOK, result)
}

/*
GetVersionFromMaster returns the current MailSlurper version from GitHub
*/
func (c *AdminController) GetVersionFromMaster(ctx echo.Context) error {
	context := contexts.GetAdminContext(ctx)
	var err error
	var result *mailslurper.Version

	if result, err = mailslurper.GetServerVersionFromMaster(); err != nil {
		c.logger.Errorf("Error getting version file from Github: %s", err.Error())
		return context.String(http.StatusInternalServerError, "There was an error reading the version file from GitHub")
	}

	return context.JSON(http.StatusOK, result)
}

/*
PerformLogin validates credentials, and if correct, sets the session
and redirects to the home page
*/
func (c *AdminController) PerformLogin(ctx echo.Context) error {
	var err error
	var authProvider auth.IAuthProvider

	authFactory := &authfactory.AuthFactory{
		Config: c.config,
	}

	authProvider = authFactory.Get()
	credentials := &auth.AuthCredentials{
		UserName: ctx.FormValue("userName"),
		Password: ctx.FormValue("password"),
	}

	if err = authProvider.Login(credentials); err != nil {
		c.logger.WithError(err).Errorf("Admin authentication error")
		return ctx.Redirect(http.StatusFound, "/login?message=Invalid user name or password")
	}

	s, _ := session.Get("session", ctx)
	s.Options = &sessions.Options{
		Path:   "/",
		MaxAge: 0,
	}
	s.Values["user"] = credentials.UserName

	s.Save(ctx.Request(), ctx.Response())
	return ctx.Redirect(http.StatusFound, "/")
}

/*
Logout logs a user out, and deletes the cookie
*/
func (c *AdminController) Logout(ctx echo.Context) error {
	s, _ := session.Get("session", ctx)
	s.Options = &sessions.Options{
		Path:   "/",
		MaxAge: -1,
	}
	s.Save(ctx.Request(), ctx.Response())

	return ctx.Redirect(http.StatusFound, "/login")
}
