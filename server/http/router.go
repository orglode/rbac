package http

import (
	"github.com/gin-gonic/gin"
)

func initRouter(rbac *gin.Engine) {

	//页面按钮接口
	rbac.POST("background/login", backgroundLogin) //CRM登录
	rbac.GET("background/list", backgroundList)    //菜单列表

	//用户管理
	rbac.GET("background/user/list", userList)
	rbac.POST("background/user/add", userAdd)
	rbac.POST("background/user/modify", userModify)
	rbac.GET("background/user/del", userDel)

	//角色管理
	rbac.GET("background/role/list", roleList)
	rbac.POST("background/role/add", roleAdd)
	rbac.POST("background/role/modify", roleModify)
	rbac.GET("background/role/del", roleDel)

	//角色类型管理
	rbac.GET("background/role_type/list", roleTypeList)
	rbac.GET("background/role_type/all", roleTypeAll)
	rbac.POST("background/role_type/add", roleTypeAdd)
	rbac.POST("background/role_type/modify", roleTypeModify)
	rbac.GET("background/role_type/del", roleTypeDel)

	//页面与按钮管理
	rbac.GET("background/page/list", pageList)
	rbac.POST("background/page/add", pageAdd)
	rbac.POST("background/page/modify", pageModify)
	rbac.GET("background/page/del", pageDel)

	//日志
	rbac.GET("background/log/list", logList)

}
