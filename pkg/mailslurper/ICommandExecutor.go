package mailslurper

/*
ICommandExecutor describes an interface for a component that executes
a command issued by the SMTP worker. A command does some type of
processing, such as parse a piece of the mail stream, make replies,
etc...
*/
type ICommandExecutor interface {
	Process(streamInput string, mailItem *MailItem) error
}
