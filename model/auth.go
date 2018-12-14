package model

import(
  "github.com/gin-gonic/gin"
  . "../function"
  . "../config"
  "strconv"
)

func (this *Base)Auth_init(){
  this.forbbiden=map[string]bool{
    "/auth/list":true,
    "/auth/get":true,
    "/auth/delete":true,
  }

  if(this.action=="update"){

    this.args=strconv.Itoa(this.sid)

  }

}

func (this *Base)Auth_getToken(){

  var sid int = 0

  account:=this.c.PostForm("account")
  password:=this.c.PostForm("password")

  sql:="select id from c_auth where account='"+account+"' and password='"+password+"'"

  err:= this.db.QueryRow(sql).Scan(&sid)

  if(err!=nil){
    this.Error(-18,Lang["auth_error"])
  }


  if(sid!=0){
    token:=CreateToken(sid)

    this.R(0,"",gin.H{
      "token":token,
      "expires":Config["token_time"],
    })
  }



}
