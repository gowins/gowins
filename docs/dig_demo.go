package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"go.uber.org/dig"
)

// DBOutput --- 1. 定义命名依赖的输出 ---
type DBOutput struct {
	dig.Out
	ReadDB  *sql.DB `name:"readDB"`
	WriteDB *sql.DB `name:"writeDB"`
}

// DBs --- 2. 初始化数据库连接 ---
func NewDBs() DBOutput {
	readDB := &sql.DB{}
	writeDB := &sql.DB{}
	return DBOutput{
		ReadDB:  readDB,
		WriteDB: writeDB,
	}
}

// ServiceDependencies --- 3. 定义依赖消费者 ---
type ServiceDependencies struct {
	dig.In            // 标记为依赖消费者
	DeclareDB *sql.DB `name:"DeclareDB"` //dig.Name //显式注入
	ReadDB    *sql.DB `name:"readDB"`    //隐式provide
	WriteDB   *sql.DB `name:"writeDB"`   //隐式provide
	Logger    *log.Logger
	Option    *Option `optional:"true"` // 可选依赖
}

// ServiceContainer --- 4. 定义服务容器 ---
type ServiceContainer struct {
	dig.Out // 标记为依赖提供者

	UserService  *UserService
	OrderService *OrderService
}

// NewServiceContainer --- 5. 构造函数（需要依赖注入） ---
func NewServiceContainer(deps ServiceDependencies) ServiceContainer {
	// 使用注入的 DB 和 Logger 初始化服务
	if deps.Option == nil {
		log.Print("Option is nil")
	}
	return ServiceContainer{
		UserService: &UserService{
			DB:      deps.DeclareDB,
			ReadDB:  deps.ReadDB,
			WriteDB: deps.WriteDB,
			Logger:  deps.Logger,
		},
		OrderService: &OrderService{
			ReadDB:  deps.ReadDB,
			WriteDB: deps.WriteDB,
			Logger:  deps.Logger,
		},
	}
}

// UserService --- 6. 业务服务定义 ---
type UserService struct {
	DB      *sql.DB
	ReadDB  *sql.DB
	WriteDB *sql.DB
	Logger  *log.Logger
}

type OrderService struct {
	ReadDB  *sql.DB
	WriteDB *sql.DB
	Logger  *log.Logger
}

type Option struct {
	Name string
}

// BuildContainer --- 7. 容器配置 ---
func BuildContainer() *dig.Container {
	container := dig.New()

	// 提供基础依赖（DB 和 Logger）
	errdb := container.Provide(NewDBs)
	if errdb != nil {
		fmt.Println("Error:", errdb)
	}

	errlog := container.Provide(log.Default)
	if errlog != nil {
		fmt.Println("Error:", errlog)
	}

	//
	container.Provide(func() (*sql.DB, error) {
		return &sql.DB{}, nil
	}, dig.Name("DeclareDB"))

	// 提供 ServiceContainer（自动注入其依赖）
	errsvc := container.Provide(NewServiceContainer)
	if errsvc != nil {
		fmt.Println("Error:", errsvc)
	}

	return container
}

// --- 8. 使用 ---
func main() {
	container := BuildContainer()

	err := container.Invoke(func(userSvc *UserService) {
		userSvc.Logger.Println("UserService 初始化成功")
		userSvc.Logger.Printf("ReadDB: %v, WriteDB: %v, DB:%v\n", &userSvc.ReadDB, &userSvc.WriteDB, &userSvc.DB)
	})
	if err != nil {
		log.Fatal(err)
	}

	// 可视化依赖图并输出到文件
	// dot -Tpng  dependency_graph.dot -o dependency_graph.png
	file, err := os.Create("dependency_graph.dot")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	err = dig.Visualize(container, file, dig.VisualizeError(err))

}
