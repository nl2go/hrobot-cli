package cmd_test

import (
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	. "gopkg.in/check.v1"

	"gitlab.com/newsletter2go/hrobot-cli/cmd"
	"gitlab.com/newsletter2go/hrobot-cli/test/mock"
	"gitlab.com/newsletter2go/hrobot-go/models"
)

func (s *AppSuite) TestServerListCommandEmptyList(c *C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	result := []models.Server{}

	mockRobotClient := mock.NewMockRobotClient(ctrl)
	mockRobotClient.EXPECT().SetUserAgent(gomock.Any()).Times(1)
	mockRobotClient.EXPECT().ServerGetList().Times(1).Return(result, nil)

	app := cmd.NewRobotApp(mockRobotClient, log.StandardLogger())

	rootCmd := app.NewRootCommand(log.StandardLogger())
	rootCmd.SetErr(log.StandardLogger().Out)

	_, err := executeCommand(rootCmd, "server:list")
	c.Assert(err, IsNil)
}

func (s *AppSuite) TestServerListCommandSuccess(c *C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	result := []models.Server{
		{
			ServerIP:   "123,123,123,123",
			ServerName: "app-prod-42",
		},
		{
			ServerIP:   "124,124,124,124",
			ServerName: "app-prod-84",
		},
	}

	mockRobotClient := mock.NewMockRobotClient(ctrl)
	mockRobotClient.EXPECT().SetUserAgent(gomock.Any()).Times(1)
	mockRobotClient.EXPECT().ServerGetList().Times(1).Return(result, nil)

	app := cmd.NewRobotApp(mockRobotClient, log.StandardLogger())

	rootCmd := app.NewRootCommand(log.StandardLogger())
	rootCmd.SetErr(log.StandardLogger().Out)

	_, err := executeCommand(rootCmd, "server:list")
	c.Assert(err, IsNil)
}

func (s *AppSuite) TestServerGenerateAnsibleInventoryCommandSuccess(c *C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	result := []models.Server{
		{
			ServerIP:   "123,123,123,123",
			ServerName: "app-prod-42",
		},
		{
			ServerIP:   "124,124,124,124",
			ServerName: "app-prod-84",
		},
	}

	mockRobotClient := mock.NewMockRobotClient(ctrl)
	mockRobotClient.EXPECT().SetUserAgent(gomock.Any()).Times(1)
	mockRobotClient.EXPECT().ServerGetList().Times(1).Return(result, nil)

	app := cmd.NewRobotApp(mockRobotClient, log.StandardLogger())

	rootCmd := app.NewRootCommand(log.StandardLogger())
	rootCmd.SetErr(log.StandardLogger().Out)

	_, err := executeCommand(rootCmd, "server:gen-ansible-inv")
	c.Assert(err, IsNil)
}
