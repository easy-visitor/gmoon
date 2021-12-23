package gmoon

import (
	"fmt"
	"reflect"
	"strings"
)

//注解处理
type Annotation interface {
	SetTag(tag reflect.StructTag)
	String() string
}

var AnnotationList []Annotation

//判断当前的注入对象是否为注解
func IsAnnotation(t reflect.Type) bool {
	for _, item := range AnnotationList {
		if reflect.TypeOf(item) == t {
			return true
		}
	}
	return false
}

func init() {
	AnnotationList = make([]Annotation, 0)
	AnnotationList = append(AnnotationList, new(Value))
}

type Value struct {
	tag         reflect.StructTag
	Beanfactory *BeanFactory
}

func (this *Value) SetTag(tag reflect.StructTag) {
	this.tag = tag
}

func (this *Value) String() string {
	getPrefix := this.tag.Get("prefix")
	if getPrefix == "" {
		return ""
	}
	prefix := strings.Split(getPrefix, ".")
	if config := this.Beanfactory.GetBean(new(SysConfig)); config != nil {
		getValue := GetConfigValue(config.(*SysConfig).Config, prefix, 0)
		if getValue != nil {
			return fmt.Sprintf("%v", getValue)
		} else {
			return ""
		}
	}
	return ""
}
