package service

import (
	"github.com/orglode/hades/logging"
	"github.com/orglode/rbac/api"
	"github.com/orglode/rbac/dao"
	"github.com/orglode/rbac/model"
	"time"
)

func (s *Service) CrmRoleListInfo(req model.RoleListRequest) (interface{}, error) {
	result := model.RoleListRes{}
	res, count, err := s.dao.GetCrmRoleList(req, dao.Paging{
		Size: req.Size,
		Page: req.Page,
	})
	if err != nil {
		s.Response.Code = api.SystemErr
		logging.Errorf("err :%v", err)
	}
	result.List = res
	result.Total = count
	s.Response.Data = result
	return nil, nil
}

func (s *Service) GetRoleAll() (interface{}, error) {
	res, err := s.dao.GetRoleAll()
	if err != nil {
		s.Response.Code = api.SystemErr
		logging.Errorf("err :%v", err)
	}
	s.Response.Data = res
	return s.Response, err
}

func (s *Service) RoleAdd(req model.RoleRequest) (interface{}, error) {
	row, err := s.dao.AddRoleInfo(req.Role)
	if err != nil {
		logging.Errorf("err :%v", err)
	}
	if row <= 0 {
		s.Response.Code = api.SystemErr
		return s.Response, nil
	}
	if row <= 0 {
		s.Response.Code = api.SystemErr
		return s.Response, nil
	}
	s.Response.Data = row
	if len(req.RoleModuleId) <= 0 {
		return s.Response, nil
	}
	for _, v := range req.RoleModuleId {
		s.dao.AddRoleModuleInfo(model.RoleModule{
			RoleId:   row,
			ModuleId: v,
		})
	}
	return row, err
}

func (s *Service) RoleModify(req model.RoleRequest) (interface{}, error) {
	row, err := s.dao.ModifyRoleInfo(req.Id, req.Role)
	if err != nil {
		logging.Errorf("err :%v", err)
	}
	s.Response.Data = row
	if len(req.RoleModuleId) <= 0 {
		return s.Response, nil
	}
	//删除原有配置
	s.dao.DelRoleModuleByRoleIdInfo(req.Id)
	for _, v := range req.RoleModuleId {
		s.dao.AddRoleModuleInfo(model.RoleModule{
			RoleId:   req.Id,
			ModuleId: v,
		})
	}
	return row, err
}

func (s *Service) RoleDel(id, operatorUid int64) (interface{}, error) {
	row, err := s.dao.DelRoleInfo(id)
	if err != nil {
		logging.Errorf("err :%v", err)
	}
	s.Response.Data = row
	return row, err
}

func (s *Service) CrmRoleTypeListInfo(req model.RoleTypeListRequest) (interface{}, error) {
	result := model.RoleTypeListRes{}
	res, count, err := s.dao.GetCrmRoleTypeList(dao.Paging{
		Size: req.Size,
		Page: req.Page,
	})
	if err != nil {
		s.Response.Code = api.SystemErr
		logging.Errorf("err :%v", err)
	}
	result.List = res
	result.Total = count
	s.Response.Data = result
	return nil, nil
}

func (s *Service) GetRoleTypeAll() (interface{}, error) {
	res, err := s.dao.GetRoleTypeAll()
	if err != nil {
		s.Response.Code = api.SystemErr
		logging.Errorf("err :%v", err)
	}
	s.Response.Data = res
	return s.Response, err
}

func (s *Service) RoleTypeAdd(req model.RoleTypeRequest) (interface{}, error) {
	req.CreateTime = time.Now().Unix()
	row, err := s.dao.AddRoleTypeInfo(req.RoleType)
	s.Response.Code = api.Success
	if err != nil {
		s.Response.Code = api.SystemErr
		logging.Errorf("err :%v", err)
	}
	s.Response.Data = row
	return s.Response, err
}

func (s *Service) RoleTypeModify(req model.RoleTypeRequest) (interface{}, error) {
	req.UpdateTime = time.Now().Unix()
	s.Response.Code = api.Success
	row, err := s.dao.ModifyRoleTypeInfo(req.Id, req.RoleType)
	if err != nil {
		logging.Errorf("err :%v", err)
	}
	s.Response.Data = row
	return s.Response, err
}

func (s *Service) DelRoleTypeInfo(id, operatorUid int64) (interface{}, error) {
	row, err := s.dao.DelRoleInfo(id)
	if err != nil {
		logging.Errorf("err :%v", err)
	}
	s.Response.Data = row
	return row, err
}
