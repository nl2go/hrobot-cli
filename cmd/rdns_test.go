package cmd_test

import (
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	. "gopkg.in/check.v1"

	"github.com/nl2go/hrobot-cli/cmd"
	"github.com/nl2go/hrobot-cli/test/mock"
	"github.com/nl2go/hrobot-go/models"
)

func (s *AppSuite) TestRdnsListCommandEmptyList(c *C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	result := []models.Rdns{
		{
			IP:  "123,123,123,123",
			Ptr: "info42.com",
		},
		{
			IP:  "124,124,124,124",
			Ptr: "example.com",
		},
	}

	mockRobotClient := mock.NewMockRobotClient(ctrl)
	mockRobotClient.EXPECT().SetUserAgent(gomock.Any()).Times(1)
	mockRobotClient.EXPECT().RDnsGetList().Times(1).Return(result, nil)

	app := cmd.NewRobotApp(mockRobotClient, log.StandardLogger())

	rootCmd := app.NewRootCommand(log.StandardLogger())
	rootCmd.SetErr(log.StandardLogger().Out)

	_, err := executeCommand(rootCmd, "rdns:list")
	c.Assert(err, IsNil)
}
