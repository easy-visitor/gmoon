# gmoon

gmoon 是基于gin构建的脚手架，主要实现了简单依赖注入，控制器，中间件，任务组件，zap日志

## 安装

go get -u github.com/easy-visitor/gmoon

## 使用
gmoon.Ignite().Mount("v1",app.NewTestClass()).Launch()
