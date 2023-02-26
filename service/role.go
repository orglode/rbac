package service

import "github.com/orglode/navigator/model"

func (s *Service) RoleAdd(req model.Role) (interface{}, error) {
	row, err := s.dao.AddRoleInfo(req)
	if err != nil {
		s.logger.Sugar().Errorf("err :%v", err)
	}
	return row, err
}

func (s *Service) RoleModify(req model.Role) (interface{}, error) {
	row, err := s.dao.ModifyRoleInfo(req.Id, req)
	if err != nil {
		s.logger.Sugar().Errorf("err :%v", err)
	}
	return row, err
}

func (s *Service) RoleDel(id int64) (interface{}, error) {
	row, err := s.dao.DelRoleInfo(id)
	if err != nil {
		s.logger.Sugar().Errorf("err :%v", err)
	}
	return row, err
}

func (s *Service) RoleTypeAdd(req model.Role) (interface{}, error) {
	row, err := s.dao.AddRoleInfo(req)
	if err != nil {
		s.logger.Sugar().Errorf("err :%v", err)
	}
	return row, err
}
