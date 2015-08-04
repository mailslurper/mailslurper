package model

/*
ServiceSettings represents the necessary settings to connect to
and talk to the MailSlurper service tier.
*/
type ServiceSettings struct {
	ServiceAddress string `json:"serviceAddress"`
	ServicePort    int    `json:"servicePort"`
	Version        string `json:"version"`
}
