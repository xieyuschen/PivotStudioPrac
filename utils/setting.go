package utils

import "fmt"
import "gopkg.in/ini.v1"
var (
	AppMode string
	HttpPort string

	Db string
	DbHost string
	DbPort string
	DbUser string
	DbPassWord string
	DbName string

	JwtKey string
)

func init(){
	file, err := ini.Load("config/config.ini")
	if err != nil {
		fmt.Println("配置文件读取错误，请检查文件路径:", err)
	}
	LoadServer(file)
	LoadData(file)
}

func LoadServer(file *ini.File){
	AppMode = file.Section("server").Key("AppMode").MustString("release")
	HttpPort = file.Section("server").Key("HttpPort").MustString(":3000")
	JwtKey = file.Section("server").Key("JwtKey").MustString("ac114514def")
}

func LoadData(file *ini.File){
	Db = file.Section("database").Key("Db").MustString("mysql")
	DbHost = file.Section("database").Key("DbHost").MustString("localhost")
	DbName = file.Section("database").Key("DbName").MustString("3306")
	DbPassWord = file.Section("database").Key("DbPassWord").MustString("ginblog")
	DbPort = file.Section("database").Key("DbPort").MustString("admin123")
	DbUser = file.Section("database").Key("DbUser").MustString("ginblog")

}