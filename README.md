# odi
对象依赖注入(Object Dependency Injection), 通过配置数据创建对象, 并递归地对依赖字段进行注入

## 安装与使用
```go get -u github.com/rcpqc/odi@latest```

```
package main

import (
	"encoding/json"
	"log"

	"github.com/rcpqc/odi/odi"
)

type Animal interface {
	Feed()
}

type Cat struct {
	Name  string `odi:"name"`
	Breed string `odi:"breed"`
}

func (c *Cat) Feed() {
	log.Printf("Feed Cat (%s %s)", c.Breed, c.Name)
}

type Dog struct {
	Name  string `odi:"name"`
	Breed string `odi:"breed"`
}

func (c *Dog) Feed() {
	log.Printf("Feed Dog (%s %s)", c.Breed, c.Name)
}

type Zoo struct {
	Name    string   `odi:"name"`
	Animals []Animal `odi:"animals"`
}

func (z *Zoo) Resolve(_ any) error {
	log.Printf("Zoo(%s) Open", z.Name)
	return nil
}

func (z *Zoo) FeedAll() {
	for _, a := range z.Animals {
		a.Feed()
	}
}

var cfg = `{
	"object": "zoo",
	"name": "Maple Zoo", 
	"animals": [{
		"object": "cat",
		"breed": "Garfield",
		"name": "erika"
	},{
		"object": "cat",
		"breed": "Ragdoll",
		"name": "bebe"
	},{
		"object": "dog",
		"breed": "Alaskan Malamute",
		"name": "ulrica"
	}]
}`

func main() {
	// 注册组建到IoC容器
	odi.Provide("cat", func() any { return &Cat{} })
	odi.Provide("dog", func() any { return &Dog{} })
	odi.Provide("zoo", func() any { return &Zoo{} })

	// 解析配置为any类型的中间结果
	var data any
	if err := json.Unmarshal([]byte(cfg), &data); err != nil {
		log.Fatal(err)
	}

	// 从any类型的中间结果，构建对象
	object, err := odi.Resolve(data)
	if err != nil {
		log.Fatal(err)
	}

	// 业务代码
	object.(*Zoo).FeedAll()
	
	// 释放对象
	odi.Dispose(object)
}
```

## 特性
- 支持属性注入  
  包括int,bool,string,float,map,slice,struct...，类似json.Unmarshal的能力
- 支持接口注入（interface）  
  可以将已注册到ioc容器的组件，实例化后注入到interface字段
- 支持Resolve回调接口（类似构造函数）  
  实现了该接口的类型，在完成注入后ODI会回调该接口用以一些自定义的处理，比如初始化
- 支持Dispose回调接口（类似析构函数）  
  实现了该接口的类型，在释放对象时会进行回调

## 例子
可参考[单测](github.com/rcpqc/odi/test/)