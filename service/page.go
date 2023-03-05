package service

import (
	"github.com/orglode/hades/logging"
	"github.com/orglode/rbac/api"
	"github.com/orglode/rbac/model"
)

func (s *Service) GetPageAll() (interface{}, error) {
	res, err := s.dao.GetModuleAll()
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

func (s *Service) PageAdd(req model.ModuleRequest) (interface{}, error) {
	row, err := s.dao.AddModuleInfo(req.Module)
	if err != nil {
		logging.Errorf("err :%v", err)
	}
	s.Response.Data = row
	return row, err
}

func (s *Service) PageModify(req model.ModuleRequest) (interface{}, error) {
	row, err := s.dao.ModifyModuleInfo(req.Id, req.Module)
	if err != nil {
		logging.Errorf("err :%v", err)
	}
	s.Response.Data = row
	return row, err
}

func (s *Service) PageDel(id, operatorUid int64) (interface{}, error) {
	row, err := s.dao.DelModuleInfo(id)
	if err != nil {
		logging.Errorf("err :%v", err)
	}
	s.Response.Data = row
	return row, err
}
