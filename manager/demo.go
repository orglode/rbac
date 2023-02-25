package manager

import "fmt"

func (m *Manager) DemoHttp() (string, error) {
	res, err := m.httpClient.
		SetTimeout(5).
		Get("")
	if err != nil {
		fmt.Println(err)
	}
	body, err := res.Body()
	return string(body), err
}
