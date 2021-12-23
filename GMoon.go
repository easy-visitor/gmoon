package gmoon

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gmoon/views/funs"
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

var config = &SysConfig{}

func init() {
	config = InitConfig()
	logger := NewZapLogger(config.Zap).Logger()
	Logger = logger
}

func Ignite() *GMoon { //这就是所谓的构造函数，ignite有 发射、燃烧， 很有激情。符合我们骚动的心情
	g := &GMoon{
		Engine:      gin.New(),
		beanFactory: NewBeanFactory(),
		exprData:    map[string]interface{}{},
	}
	g.Use(ErrorHandler())         //强迫加载的异常处理中间件
	g.beanFactory.setBean(config) //整个配置加载进bean中
	if config.Server.Html != "" {
		g.FuncMap = funs.FuncMap
		//g.LoadHTMLGlob(config.Server.Html)
	}
	return g
}

//func (this *GMoon) Handle(httpMethod, relativePath string, handlers ...gin.HandlerFunc) *GMoon {
func (this *GMoon) Handle(httpMethod, relativePath string, handler interface{}) *GMoon {

	if h := Convert(handler); h != nil {
		this.g.Handle(httpMethod, relativePath, h)
	}
	/*if h, ok := handlers.(func(ctx *gin.Context) string); ok {
		this.g.Handle(httpMethod, relativePath, func(context *gin.Context) {
			context.String(200, h(context))
		})
	}*/
	return this
}
func (this *GMoon) Launch() { //启动函数
	var port = 8080
	if config := this.beanFactory.GetBean(new(SysConfig)); config != nil {
		port = config.(*SysConfig).Server.Port
	}
	getCronTask().Start()
	Logger.Info("start success")
	this.Run(fmt.Sprintf(":%d", port))

}

//设置数据库链接对象
func (this *GMoon) Beans(beans ...Bean) *GMoon {

	for _, bean := range beans {
		this.exprData[bean.Name()] = bean
	}
	this.beanFactory.setBean(beans...)
	return this
}

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

func (this *GMoon) Mount(group string, classes ...IClass) *GMoon {
	this.g = this.Group(group)
	for _, class := range classes {
		class.Build(this) //这一步是关键 。 这样在main里面 就不需要 调用了
		this.beanFactory.inject(class)
		this.Beans(class) //控制器 也作为bean加入到bean容器
	}
	return this
}
