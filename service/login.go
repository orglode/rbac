package service

import (
	"github.com/gin-contrib/sessions"
	"github.com/orglode/hades/logging"
	"github.com/orglode/rbac/api"
	"github.com/orglode/rbac/model"
)

func (s *Service) AccountLogin(req model.AccountLoginRequest, session sessions.Session) {
	userInfo, err := s.dao.GetUserInfoByAccount(req.Account)
	if err != nil {
		s.Response.Code = api.SystemErr
		logging.Errorf("err :%v", err)
		return
	}
	if userInfo.Id <= 0 {
		s.Response.Code = api.UserNoTRegister
		return
	}
	if userInfo.PassWord != req.PassWord {
		s.Response.Code = api.PassWordErr
		return
	}
	userInfo.PassWord = ""
	session.Set("uid", userInfo.Id)
	session.Save()
	s.Response.Data = userInfo
	return
}
