package gmoon

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
)

var Logger *zap.Logger

type GMoon struct {
	*gin.Engine
	g           *gin.RouterGroup
	beanFactory *BeanFactory
	props       []interface{}

	exprData map[string]interface{}
}

//全局配置文件
var globalConfig = &SysConfig{}

func init() {
	globalConfig = InitConfig()
	logger := NewZapLogger(globalConfig.Zap).Logger()
	Logger = logger
}

//初始化
func Ignite() *GMoon {
	g := &GMoon{
		Engine:      gin.New(),
		beanFactory: NewBeanFactory(),
		exprData:    map[string]interface{}{},
	}

	//gin.SetMode(gin.ReleaseMode)

	g.Use(ErrorHandler())               //强迫加载的异常处理中间件
	g.beanFactory.setBean(globalConfig) //整个配置加载进bean中
	if globalConfig.Server.Html != "" {
		// g.FuncMap = funs.FuncMap
		g.LoadHTMLGlob(globalConfig.Server.Html)
	}
	return g
}

//实现路由
func (this *GMoon) Handle(httpMethod, relativePath string, handler interface{}) *GMoon {
	if h := Convert(handler); h != nil {
		this.g.Handle(httpMethod, relativePath, h)
	}
	return this
}

//启动
func (this *GMoon) Launch() { //启动函数
	var port = 8080
	if config := this.beanFactory.GetBean(new(SysConfig)); config != nil {
		port = config.(*SysConfig).Server.Port
	}
	getCronTask().Start()
	this.Run(fmt.Sprintf(":%d", port))

}

//设置数据库链接对象
func (this *GMoon) Beans(beans ...Bean) *GMoon {

	for _, bean := range beans {
		fmt.Println(bean)
		this.exprData[bean.Name()] = bean
	}
	this.beanFactory.setBean(beans...)
	return this
}

//注册中间件
func (this *GMoon) Attach(f Fairing) *GMoon {
	this.Use(func(context *gin.Context) {
		err := f.OnRequest(context)
		if err != nil {
			context.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		} else {
			context.Next()
		}
	})
	return this
}

//定时任务
func (this *GMoon) Task(cron string, expr interface{}) *GMoon {
	var err error
	if f, ok := expr.(func()); ok {
		_, err = getCronTask().AddFunc(cron, f)
	}
	if exp, ok := expr.(Expr); ok {
		_, err = getCronTask().AddFunc(cron, func() {
			_, expErr := ExecExpr(exp, this.exprData)
			if expErr != nil {
				log.Println(expErr)
			}
		})
	}

	if err != nil {
		log.Println(err)
	}
	return this
}

//挂载控制器与路由
func (this *GMoon) Mount(group string, classes ...IClass) *GMoon {
	this.g = this.Group(group)
	for _, class := range classes {
		class.Build(this) //这一步是关键 。 这样在main里面 就不需要 调用了
		this.beanFactory.inject(class)
		this.Beans(class) //控制器 也作为bean加入到bean容器
	}
	return this
}
