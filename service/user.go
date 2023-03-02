package service

import (
	"fmt"
	"github.com/orglode/navigator/api"
	"github.com/orglode/navigator/dao"
	"github.com/orglode/navigator/model"
	"time"
)

func (s *Service) CrmUserListInfo(req model.CrmUserListRequest) {
	result := model.CrmUserListRes{}
	res, count, err := s.dao.GetCrmUserList(dao.Paging{
		Size: req.Size,
		Page: req.Page,
	})
	if err != nil {
		s.Response.Code = api.SystemErr
		s.logger.Sugar().Errorf("err :%v", err)
	}
	result.List = res
	result.Total = count
	s.Response.Data = result
}

func (s *Service) UserAdd(req model.CrmUserRequest) {
	req.CreateTime = time.Now().Unix()
	row, err := s.dao.AddUserInfo(req.Users)
	if err != nil {
		s.logger.Sugar().Errorf("err :%v", err)
	}
	s.Response.Data = row
}

func (s *Service) UserModify(req model.CrmUserRequest) {
	req.UpdateTime = time.Now().Unix()
	row, err := s.dao.ModifyUserInfo(req.Id, req.Users)
	if err != nil {
		s.logger.Sugar().Errorf("err :%v", err)
	}
	s.Response.Data = row
}

func (s *Service) UserDel(id, operatorUid int64) {
	fmt.Println(operatorUid)
	row, err := s.dao.DelUserInfo(id)
	if err != nil {
		s.logger.Sugar().Errorf("err :%v", err)
	}
	s.Response.Data = row
}
