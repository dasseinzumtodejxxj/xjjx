package routes

import (
	"github.com/gin-contrib/multitemplate"

	"ginblog/api/v1"
	"ginblog/middleware"
	"ginblog/utils"
	"github.com/gin-gonic/gin"
)

func createMyRender() multitemplate.Renderer {
	p := multitemplate.NewRenderer()
	p.AddFromFiles("admin", "web/admin/dist/index.html")
	p.AddFromFiles("front", "web/front/dist/index.html")
	return p
}
func InitRouter() {
	gin.SetMode(utils.AppMode)
	r := gin.New()
	_ = r.SetTrustedProxies(nil)

	r.HTMLRender = createMyRender()
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())

	r.Static("/static", "./web/front/dist/static")
	r.Static("/admin", "./web/admin/dist")
	r.StaticFile("/favicon.ico", "/web/front/dist/favicon.ico")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "front", nil)
	})

	r.GET("/admin", func(c *gin.Context) {
		c.HTML(200, "admin", nil)
	})

	auth := r.Group("api/v1")
	auth.Use(middleware.JwtToken()) //需要鉴权
	{
		//用户模块的路由接口
		auth.GET("admin/users", v1.GetUsers)
		auth.PUT("user/:id", v1.EditUser)
		auth.DELETE("user/:id", v1.DeleteUser)
		//修改密码
		auth.PUT("admin/changepw/:id", v1.ChangeUserPassword)
		//分类模块的路由接口
		auth.POST("category/add", v1.AddCategory)
		auth.PUT("category/:id", v1.EditCate)
		auth.DELETE("category/:id", v1.DeleteCate)
		//文章模块的路由接口
		auth.GET("admin/article", v1.GetArt)
		auth.POST("article/add", v1.AddArticle)
		auth.PUT("article/:id", v1.EditArt)
		auth.DELETE("article/:id", v1.DeleteArt)
		//上传文件
		auth.POST("upload", v1.Upload)
		// 更新个人设置
		auth.GET("admin/profile/:id", v1.GetProfile)
		auth.PUT("profile/:id", v1.UpdateProfile)
		// 评论模块
		auth.GET("comment/list", v1.GetCommentList)
		auth.DELETE("delcomment/:id", v1.DeleteComment)
		auth.PUT("checkcomment/:id", v1.CheckComment)
		auth.PUT("uncheckcomment/:id", v1.UncheckComment)
	}
	routerv1 := r.Group("api/v1") //无需鉴权
	{
		routerv1.POST("user/add", v1.AddUser)
		routerv1.GET("user/:id", v1.GetUserInfo)
		routerv1.GET("users", v1.GetUsers)
		routerv1.GET("category/:id", v1.GetCateInfo)
		routerv1.GET("category", v1.GetCate)
		routerv1.GET("article", v1.GetArt)
		routerv1.GET("article/list/:id", v1.GetCateArt)
		routerv1.GET("article/info/:id", v1.GetArtInfo)
		// 登录控制模块
		routerv1.POST("login", v1.Login)
		routerv1.POST("loginfront", v1.LoginFront)
		// 获取个人设置信息
		routerv1.GET("profile/:id", v1.GetProfile)
		// 评论模块
		routerv1.POST("addcomment", v1.AddComment)
		routerv1.GET("comment/info/:id", v1.GetComment)
		routerv1.GET("commentfront/:id", v1.GetCommentListFront)
		routerv1.GET("commentcount/:id", v1.GetCommentCount)
	}

	_ = r.Run(utils.HttpPort)
}
