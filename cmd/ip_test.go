package cmd_test

import (
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	. "gopkg.in/check.v1"

	"gitlab.com/newsletter2go/hrobot-cli/cmd"
	"gitlab.com/newsletter2go/hrobot-cli/test/mock"
	"gitlab.com/newsletter2go/hrobot-go/models"
)

func (s *AppSuite) TestIPListCommandEmptyList(c *C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	result := []models.IP{}

	mockRobotClient := mock.NewMockRobotClient(ctrl)
	mockRobotClient.EXPECT().SetUserAgent(gomock.Any()).Times(1)
	mockRobotClient.EXPECT().IPGetList().Times(1).Return(result, nil)

	app := cmd.NewRobotApp(mockRobotClient, log.StandardLogger())

	rootCmd := app.NewRootCommand(log.StandardLogger())
	rootCmd.SetErr(log.StandardLogger().Out)

	_, err := executeCommand(rootCmd, "ip:list")
	c.Assert(err, IsNil)
}

func (s *AppSuite) TestIPListCommandSuccess(c *C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	result := []models.IP{
		{
			ServerIP:     "123,123,123,123",
			ServerNumber: 321,
		},
		{
			ServerIP:     "124,124,124,124",
			ServerNumber: 123,
		},
	}

	mockRobotClient := mock.NewMockRobotClient(ctrl)
	mockRobotClient.EXPECT().SetUserAgent(gomock.Any()).Times(1)
	mockRobotClient.EXPECT().IPGetList().Times(1).Return(result, nil)

	app := cmd.NewRobotApp(mockRobotClient, log.StandardLogger())

	rootCmd := app.NewRootCommand(log.StandardLogger())
	rootCmd.SetErr(log.StandardLogger().Out)

	_, err := executeCommand(rootCmd, "ip:list")
	c.Assert(err, IsNil)
}
