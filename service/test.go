package service

import (
	"github.com/orglode/navigator/model"
)

func (s *Service) Test(req model.TestAReq) (interface{}, error) {

	result := struct {
		Row  interface{} `json:"row"`
		Http interface{} `json:"http"`
	}{}
	row, err := s.dao.GetMachineAll()
	if err != nil {
		s.logger.Info("2333")
		s.logger.Error("asdjkl")
	}
	result.Row = row
	result.Http, _ = s.mgr.DemoHttp()
	return result, err
}
