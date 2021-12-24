package gmoon

import (
	"github.com/gin-gonic/gin"
	"reflect"
)

type Responder interface {
	RespondTo() gin.HandlerFunc
}

var ResponderList []Responder

func init() {

	ResponderList = []Responder{new(StringResponder), new(ModelResponder), new(ModelsResponder), new(ViewResponder),new(JsonResponder)}
}

func Convert(handler interface{}) gin.HandlerFunc {
	hRef := reflect.ValueOf(handler)
	for _, r := range ResponderList {
		rRef := reflect.ValueOf(r).Elem()
		if hRef.Type().ConvertibleTo(rRef.Type()) {
			rRef.Set(hRef)
			return rRef.Interface().(Responder).RespondTo()
		}
	}
	return nil
}

type StringResponder func(ctx *gin.Context) string

func (this StringResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.String(200, this(context))
	}
}

type Json interface{}
type JsonResponder func(*gin.Context) Json

func (this JsonResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(200, this(context))
	}
}

type ModelResponder func(ctx *gin.Context) Model

func (this ModelResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(200, this(context))
	}
}

type ModelsResponder func(ctx *gin.Context) Models

func (this ModelsResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Writer.Header().Set("Content-type", "application/json")
		context.Writer.WriteString(string(this(context)))
		//context.JSON(200, this(context))
	}
}

type View string
type ViewResponder func(ctx *gin.Context) View

func (this ViewResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.HTML(200, string(this(context))+".html", context.Keys)

	}
}
