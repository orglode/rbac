package service

import (
	"github.com/orglode/hades/logging"
	"github.com/orglode/rbac/api"
	"github.com/orglode/rbac/model"
)

func (s *Service) GetUserPageShow(req model.UserPageList) (interface{}, error) {
	//获取当前角色
	roleInfo, err := s.dao.GetUserRole(req.OperatorUid)
	if err != nil {
		s.Response.Code = api.SystemErr
		logging.Errorf("err :%v", err)
		return nil, nil
	}
	if roleInfo.RoleId <= 0 {
		s.Response.Code = api.UserNotFindRole
		return nil, nil
	}
	//获取这个角色下的所有菜单
	roleModule, err := s.dao.GetRoleModuleByRoleIdAll(roleInfo.RoleId)
	if err != nil {
		logging.Errorf("err :%v", err)
		s.Response.Code = api.SystemErr
		return nil, nil
	}
	if len(roleModule) <= 0 {
		s.Response.Code = api.RoleNotFindModule
		return nil, nil
	}
	roleIdArr := make([]int64, 0)
	for _, v := range roleModule {
		roleIdArr = append(roleIdArr, v.ModuleId)
	}
	res, err := s.dao.GetModuleByIdAll(roleIdArr)
	if err != nil {
		s.Response.Code = api.SystemErr
		logging.Errorf("err :%v", err)
	}
	data := make([]SystemMenu, 0)
	for _, v := range res {
		temp := SystemMenu{
			Id:     int(v.Id),
			Parent: int(v.Parent),
			Name:   v.Name,
			Sort:   v.Sort,
			Code:   v.Code,
			Type:   v.Type,
			Root:   v.Root,
		}
		data = append(data, temp)
	}
	nodes := SystemMenus.ConvertToINodeArray(data)
	newNodes := GenerateTree(nodes, nil)
	s.Response.Data = newNodes
	return s.Response, err
}
