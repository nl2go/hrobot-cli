package cmd_test

import (
	"bytes"
	"testing"

	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	. "gopkg.in/check.v1"

	"github.com/nl2go/hrobot-cli/cmd"
	"github.com/nl2go/hrobot-cli/test/mock"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type AppSuite struct{}

var _ = Suite(&AppSuite{})

func (s *AppSuite) TestDefaultRunNoCommand(c *C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	mockRobotClient := mock.NewMockRobotClient(ctrl)
	mockRobotClient.EXPECT().SetUserAgent(gomock.Any()).Times(1)

	app := cmd.NewRobotApp(mockRobotClient, log.StandardLogger())
	app.Run()
}

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = executeCommandC(root, args...)
	return output, err
}

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOutput(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}
