package api

import (
	"SearchEngine/core"
	"SearchEngine/router/mysql"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var db *sql.DB
var Engine *core.PictureEngine


func SetEngine(e *core.PictureEngine) {
	Engine = e
}

func Set(e *core.PictureEngine){
	//初始化MySql
	db = mysql.InitDB()

	SetEngine(e)
}


//注册方法
func Register(router *gin.Engine) {


	//加载静态网页
	router.LoadHTMLGlob("/home/wqk/GolandProjects/SearchEngine/router/template/*.html")
	//加载静态资源
	router.Static("statics","/home/wqk/GolandProjects/SearchEngine/router/statics")
	//首页
	router.GET("",welcome)

	//登录注册
	router.GET("/login",login).POST("/login",login2)
	router.GET("/signup",signup).POST("/signup",signup2)
	//查询
	router.POST("query",query)
}




func welcome(c *gin.Context){
	c.HTML(http.StatusOK,"welcome.html",nil)
}


func login2(c *gin.Context) {

	user_name := c.PostForm("username")
	user_password := c.PostForm("password")

	var verify_password string

	err := db.QueryRow("SELECT user_password FROM user WHERE user_name=?", user_name).Scan(&verify_password)

	switch {
	case err == sql.ErrNoRows || err != nil:
		c.HTML(http.StatusNotFound,"404.html",nil)
	default:
		if verify_password == user_password {
			c.HTML(http.StatusOK,"query.html",nil)
		} else if verify_password != user_password {
			c.HTML(http.StatusNotFound,"error.html",nil)
		}
	}
}
func login(c *gin.Context) {
	c.HTML(http.StatusOK,"login.html",nil)
}

func signup(c *gin.Context) {
	c.HTML(http.StatusOK,"signup.html",nil)
}

func signup2(c *gin.Context) {

	user_name := c.PostForm("username")
	user_password := c.PostForm("password")

	var name string

	err := db.QueryRow("SELECT user_name FROM user WHERE user_name=?", user_name).Scan(&name)

	switch {
	case err == sql.ErrNoRows:
		//不存在这样的用户名，实行注册
		_, err := db.Exec("INSERT INTO user (user_name,user_password)VALUES(?,?);", user_name, user_password)
		if err != nil {
			log.Fatal(err)
		}
		c.HTML(http.StatusOK,"query.html",nil)
	case err != nil:
		c.HTML(http.StatusNotFound,"404.html",nil)
	default:
		c.HTML(http.StatusNotFound,"error2.html",nil)
	}
}

func query(c *gin.Context){
	querytext := c.PostForm("text")
	keytext := c.PostForm("keyswords")
	hatetext := c.PostForm("hatewords")

	request := core.Searchrequest{
		QueryText: querytext,
		KeyWords: keytext,
		HateWords: hatetext,
	}

	result :=Engine.Search(request)

	c.JSON(http.StatusOK,result)

}

