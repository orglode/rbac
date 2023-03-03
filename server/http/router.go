package http

import (
	"github.com/gin-gonic/gin"
)

func initRouter(rbac *gin.Engine) {

	//页面按钮接口
	rbac.GET("background/login", backgroundLogin)                               //CRM登录
	rbac.GET("background/login_out", checkUserLoginSession, backgroundLoginOut) //CRM登录

	rbac.GET("background/list", checkUserLoginSession, backgroundList) //菜单列表

	//用户管理
	rbac.GET("background/user/list", checkUserLoginSession, userList)
	rbac.POST("background/user/add", checkUserLoginSession, userAdd)
	rbac.POST("background/user/modify", checkUserLoginSession, userModify)
	rbac.GET("background/user/del", checkUserLoginSession, userDel)

	//角色管理-ing
	rbac.GET("background/role/list", checkUserLoginSession, roleList)
	rbac.POST("background/role/add", checkUserLoginSession, roleAdd)
	rbac.GET("background/role/all", checkUserLoginSession, roleAll)
	rbac.POST("background/role/modify", checkUserLoginSession, roleModify)
	rbac.GET("background/role/del", checkUserLoginSession, roleDel)

	//角色类型管理-ok
	rbac.GET("background/role_type/list", checkUserLoginSession, roleTypeList)
	rbac.GET("background/role_type/all", checkUserLoginSession, roleTypeAll)
	rbac.POST("background/role_type/add", checkUserLoginSession, roleTypeAdd)
	rbac.POST("background/role_type/modify", checkUserLoginSession, roleTypeModify)
	rbac.GET("background/role_type/del", checkUserLoginSession, roleTypeDel)

	//页面与按钮管理
	rbac.GET("background/page/list", checkUserLoginSession, pageList)
	rbac.POST("background/page/add", checkUserLoginSession, pageAdd)
	rbac.POST("background/page/modify", checkUserLoginSession, pageModify)
	rbac.GET("background/page/del", checkUserLoginSession, pageDel)

	//日志
	rbac.GET("background/log/list", checkUserLoginSession, logList)

}
