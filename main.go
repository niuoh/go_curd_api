package main

import(
  "github.com/gin-gonic/gin"
  "./model"
  . "./config"

  "database/sql"
  _ "github.com/go-sql-driver/mysql"

  "fmt"
)



func main(){

  db, err := sql.Open("mysql", Db["mysql_user"]+":"+Db["mysql_passwd"]+"@tcp("+Db["mysql_host"]+")/"+Db["mysql_dbname"]+"?charset="+Db["mysql_charset"])

  cache := make(map[string]interface{})

  if(err!=nil){

    fmt.Println(err)

  }else{

    r := gin.Default()

  	r.GET("/*path", func(c *gin.Context) {
      var m model.Base
      m.Start(c,db,cache)
  	})

    r.POST("/*path", func(c *gin.Context) {
      var m model.Base
      m.Start(c,db,cache)
  	})


  	r.Run("0.0.0.0:8080")

  }



}
