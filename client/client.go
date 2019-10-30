package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"strings"

	"gitlab.com/nl2go/hrobot-cli/client/models"
)

const baseURL string = "https://robot-ws.your-server.de"

const Version = "1.0.0"

const UserAgent = "hrobot-cli/" + Version

type Client struct {
	Username string
	Password string
}

func NewBasicAuthClient(username, password string) *Client {
	return &Client{
		Username: username,
		Password: password,
	}
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("User-Agent", UserAgent)
	req.SetBasicAuth(c.Username, c.Password)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if 200 != resp.StatusCode {
		return nil, fmt.Errorf("%s", body)
	}
	return body, nil
}

func (c *Client) ServerGetList() ([]models.Server, error) {
	url := baseURL + "/server"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	bytes, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var servers []models.ServerResponse
	err = json.Unmarshal(bytes, &servers)
	if err != nil {
		return nil, err
	}

	var data []models.Server
	for _, server := range servers {
		data = append(data, server.Server)
	}

	return data, nil
}

func (c *Client) ServerGet(ip string) (*models.Server, error) {
	url := fmt.Sprintf(baseURL+"/server/%s", ip)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	bytes, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var serverResp models.ServerResponse
	err = json.Unmarshal(bytes, &serverResp)
	if err != nil {
		return nil, err
	}

	return &serverResp.Server, nil
}

func (s *Client) ServerSetName(ip, name string) error {
	url := fmt.Sprintf(baseURL+"/server/%s", ip)

	formData := neturl.Values{}
	formData.Set("server_name", name)

	req, err := http.NewRequest("POST", url, strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err = s.doRequest(req)
	return err
}

func (c *Client) KeyGetList() ([]models.Key, error) {
	url := baseURL + "/key"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	bytes, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var keys []models.KeyResponse
	err = json.Unmarshal(bytes, &keys)
	if err != nil {
		return nil, err
	}

	var data []models.Key
	for _, key := range keys {
		data = append(data, key.Key)
	}

	return data, nil
}
