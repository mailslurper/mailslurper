package controllers

import (
	"bytes"
	"encoding/base64"
	"math"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/mailslurper/mailslurper/pkg/auth/auth"
	"github.com/mailslurper/mailslurper/pkg/auth/authfactory"
	"github.com/mailslurper/mailslurper/pkg/auth/jwt"
	"github.com/mailslurper/mailslurper/pkg/contexts"
	"github.com/mailslurper/mailslurper/pkg/mailslurper"
	"github.com/sirupsen/logrus"
)

/*
ServiceController provides methods for handling service endpoints.
This is to primarily support the API
*/
type ServiceController struct {
	config        *mailslurper.Configuration
	database      mailslurper.IStorage
	logger        *logrus.Entry
	serverVersion string
}

/*
NewServiceController creates a new admin controller
*/
func NewServiceController(logger *logrus.Entry, serverVersion string, config *mailslurper.Configuration, database mailslurper.IStorage) *ServiceController {
	return &ServiceController{
		config:        config,
		database:      database,
		logger:        logger,
		serverVersion: serverVersion,
	}
}

/*
DeleteMail is a request to delete mail items. This expects a body containing
a DeleteMailRequest object.

	DELETE: /mail/{pruneCode}
*/
func (c *ServiceController) DeleteMail(ctx echo.Context) error {
	var err error
	var deleteMailRequest *mailslurper.DeleteMailRequest
	var rowsDeleted int64

	context := contexts.GetAdminContext(ctx)

	if err = ctx.Bind(&deleteMailRequest); err != nil {
		c.logger.Errorf("Invalid delete request in DeleteMail: %s", err.Error())
		return context.String(http.StatusBadRequest, "Invalid delete request")
	}

	if !deleteMailRequest.PruneCode.IsValid() {
		c.logger.Errorf("Attempt to use invalid prune code - %s", deleteMailRequest.PruneCode)
		return context.String(http.StatusBadRequest, "Invalid prune type")
	}

	startDate := deleteMailRequest.PruneCode.ConvertToDate()

	if rowsDeleted, err = c.database.DeleteMailsAfterDate(startDate); err != nil {
		c.logger.Errorf("Problem deleting mails with code %s - %s", deleteMailRequest.PruneCode.String(), err.Error())
		return context.String(http.StatusInternalServerError, "There was a problem deleting mails")
	}

	c.logger.Infof("Deleting %d mails, code %s before %s", rowsDeleted, deleteMailRequest.PruneCode.String(), startDate)
	return context.String(http.StatusOK, strconv.Itoa(int(rowsDeleted)))
}

/*
GetMail returns a single mail item by ID.

	GET: /mail/{id}
*/
func (c *ServiceController) GetMail(ctx echo.Context) error {
	var mailID string
	var result *mailslurper.MailItem
	var err error

	context := contexts.GetAdminContext(ctx)

	mailID = context.Param("id")

	/*
	 * Retrieve the mail item
	 */
	if result, err = c.database.GetMailByID(mailID); err != nil {
		c.logger.Errorf("Problem getting mail item %s - %s", mailID, err.Error())
		return context.String(http.StatusInternalServerError, "Problem getting mail item")
	}

	c.logger.Infof("Mail item %s retrieved", mailID)
	return context.JSON(http.StatusOK, result)
}

/*
GetMailCollection returns a collection of mail items. This is constrianed
by a page number. A page of data contains 50 items.

	GET: /mails?pageNumber={pageNumber}
*/
func (c *ServiceController) GetMailCollection(ctx echo.Context) error {
	var err error
	var pageNumberString string
	var pageNumber int
	var mailCollection []*mailslurper.MailItem
	var totalRecordCount int

	context := contexts.GetAdminContext(ctx)

	/*
	 * Validate incoming arguments. A page is currently 50 items, hard coded
	 */
	pageNumberString = context.QueryParam("pageNumber")
	if pageNumberString == "" {
		pageNumber = 1
	} else {
		if pageNumber, err = strconv.Atoi(pageNumberString); err != nil {
			c.logger.Errorf("Invalid page number passed to GetMailCollection - %s", pageNumberString)
			return context.String(http.StatusBadRequest, "A valid page number is required")
		}
	}

	length := 50
	offset := (pageNumber - 1) * length

	/*
	 * Retrieve mail items
	 */
	mailSearch := &mailslurper.MailSearch{
		Message: context.QueryParam("message"),
		Start:   context.QueryParam("start"),
		End:     context.QueryParam("end"),
		From:    context.QueryParam("from"),
		To:      context.QueryParam("to"),

		OrderByField:     context.QueryParam("orderby"),
		OrderByDirection: context.QueryParam("dir"),
	}

	if mailCollection, err = c.database.GetMailCollection(offset, length, mailSearch); err != nil {
		c.logger.Errorf("Problem getting mail collection - %s", err.Error())
		return context.String(http.StatusInternalServerError, "Problem getting mail collection")
	}

	if totalRecordCount, err = c.database.GetMailCount(mailSearch); err != nil {
		c.logger.Errorf("Problem getting record count in GetMailCollection - %s", err.Error())
		return context.String(http.StatusInternalServerError, "Error getting record count")
	}

	totalPages := int(math.Ceil(float64(totalRecordCount / length)))
	if totalPages*length < totalRecordCount {
		totalPages++
	}

	c.logger.Infof("Mail collection page %d retrieved", pageNumber)

	result := &mailslurper.MailCollectionResponse{
		MailItems:    mailCollection,
		TotalPages:   totalPages,
		TotalRecords: totalRecordCount,
	}

	return context.JSON(http.StatusOK, result)
}

/*
GetMailCount returns the number of mail items in storage.

	GET: /mailcount
*/
func (c *ServiceController) GetMailCount(ctx echo.Context) error {
	var err error
	var mailItemCount int

	context := contexts.GetAdminContext(ctx)

	/*
	 * Get the count
	 */
	if mailItemCount, err = c.database.GetMailCount(&mailslurper.MailSearch{}); err != nil {
		c.logger.Errorf("Problem getting mail item count in GetMailCount - %s", err.Error())
		return context.String(http.StatusInternalServerError, "Problem getting mail count")
	}

	c.logger.Infof("Mail item count - %d", mailItemCount)

	result := &mailslurper.MailCountResponse{
		MailCount: mailItemCount,
	}

	return context.JSON(http.StatusOK, result)
}

/*
GetMailMessage returns the message contents of a single mail item

	GET: /mail/{id}/message
*/
func (c *ServiceController) GetMailMessage(ctx echo.Context) error {
	var mailID string
	var mailItem *mailslurper.MailItem
	var err error

	context := contexts.GetAdminContext(ctx)

	mailID = context.Param("id")

	/*
	 * Retrieve the mail item
	 */
	if mailItem, err = c.database.GetMailByID(mailID); err != nil {
		c.logger.Errorf("Problem getting mail item %s in GetMailMessage - %s", mailID, err.Error())
		return context.String(http.StatusInternalServerError, "Problem getting mail item")
	}

	c.logger.Infof("Mail item %s retrieved", mailID)
	return context.HTML(http.StatusOK, mailItem.Body)
}

/*
GetPruneOptions retrieves the set of options available to users for pruning

	GET: /pruneoptions
*/
func (c *ServiceController) GetPruneOptions(ctx echo.Context) error {
	context := contexts.GetAdminContext(ctx)
	return context.JSON(http.StatusOK, mailslurper.PruneOptions)
}

/*
DownloadAttachment retrieves binary database from storage and streams
it back to the caller

	GET: /mail/{mailID}/attachment/{attachmentID}
*/
func (c *ServiceController) DownloadAttachment(ctx echo.Context) error {
	var err error
	var attachmentID string
	var mailID string

	var attachment *mailslurper.Attachment
	var data []byte

	context := contexts.GetAdminContext(ctx)
	mailID = context.Param("mailID")
	attachmentID = context.Param("attachmentID")

	/*
	 * Retrieve the attachment
	 */
	if attachment, err = c.database.GetAttachment(mailID, attachmentID); err != nil {
		c.logger.Errorf("Problem getting attachment %s - %s", attachmentID, err.Error())
		return context.String(http.StatusInternalServerError, "Error getting attachment")
	}

	/*
	 * Decode the base64 data and stream it back
	 */
	if attachment.IsContentBase64() {
		if data, err = base64.StdEncoding.DecodeString(attachment.Contents); err != nil {
			c.logger.Errorf("Problem decoding attachment %s - %s", attachmentID, err.Error())
			return context.String(http.StatusInternalServerError, "Cannot decode attachment")
		}
	} else {
		data = []byte(attachment.Contents)
	}

	c.logger.Infof("Attachment %s retrieved", attachmentID)

	reader := bytes.NewReader(data)
	return context.Stream(http.StatusOK, attachment.Headers.ContentType, reader)
}

func (c *ServiceController) Login(ctx echo.Context) error {
	var err error
	var token string

	authFactory := &authfactory.AuthFactory{
		Config: c.config,
	}

	jwtService := &jwt.JWTService{
		Config: c.config,
	}

	authService := authFactory.Get()
	credentials := &auth.AuthCredentials{
		UserName: ctx.FormValue("userName"),
		Password: ctx.FormValue("password"),
	}

	if err = authService.Login(credentials); err != nil {
		c.logger.WithError(err).Errorf("Invalid service login attempt")
		return ctx.String(http.StatusForbidden, "Invalid credentials")
	}

	if token, err = jwtService.CreateToken(c.config.AuthSecret, credentials.UserName); err != nil {
		c.logger.WithError(err).Errorf("Problem creating token in service login")
		return ctx.String(http.StatusInternalServerError, "Problem creating JWT token")
	}

	return ctx.String(http.StatusOK, token)
}

func (c *ServiceController) Version(ctx echo.Context) error {
	context := contexts.GetAdminContext(ctx)
	return context.String(http.StatusOK, c.serverVersion)
}
