package client

import "gitlab.com/nl2go/hrobot-cli/client/models"

type RobotClient interface {
	SetBaseURL(baseURL string)
	SetUserAgent(userAgent string)
	ServerGetList() ([]models.Server, error)
	ServerGet(ip string) (*models.Server, error)
	ServerSetName(ip, name string) error
	ServerReverse(ip string) error
	KeyGetList() ([]models.Key, error)
	IPGetList() ([]models.IP, error)
	RDnsGetList() ([]models.Rdns, error)
	RDnsGet(ip string) (*models.Rdns, error)
}
