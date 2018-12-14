package model

import(
  "github.com/gin-gonic/gin"

  "database/sql"
  _ "github.com/go-sql-driver/mysql"

  . "../function"
  . "../config"

  "strings"
  "fmt"
  "reflect"
  "math"

  "strconv"
)

type Base struct{
  path string
  controller string
  action string
  c *gin.Context
  cCp *gin.Context
  db *sql.DB
  args string
  sid int
  cache map[string]interface{}
  forbbiden map[string]bool
}

var not_need_token_path = []string{
  "/auth/getToken",
  "/auth/create",
  "/org/list",
}

func (this *Base)InvokeObjectMethod(methodName string, args ...interface{}){

    inputs := make([]reflect.Value, len(args))
    for i, _ := range args {
        inputs[i] = reflect.ValueOf(args[i])
    }

    func_kind:=reflect.ValueOf(func(){}).Kind()

    init:=reflect.ValueOf(this).MethodByName(this.controller+"_init")

    if( init.Kind() == func_kind ){
      init.Call(inputs)
    }

    if ( this.checkAccess()==true ){

      method:=reflect.ValueOf(this).MethodByName(methodName)

      if( method.Kind() == func_kind ){

        method.Call(inputs)

      }else{
        this.R(-404,"not found",gin.H{})
      }

    }



}

func (this *Base) checkAccess() bool{

  if(this.forbbiden[this.path]==true){
    this.R(-403,"forbbiden",gin.H{})
    return false
  }

  return true
}

func (this *Base) Start(c *gin.Context,db *sql.DB,cache map[string]interface{}) {

  this.c=c
  this.cCp = this.c.Copy()
  this.cache=cache
  method:=this.RouterToMethod(this.c.Param("path"))



  if( !InArray(this.path,not_need_token_path) ){
    if(!this.haveToken() ){
      this.R(-403,Lang["token_unvalid"],gin.H{})
      return
    }else{
        sid:=CheckToken(this.c.Request.Header.Get(Config["token_key"].(string)))
        if(sid==0){
          this.R(-403,Lang["token_unvalid"],gin.H{})
          return
        }else{
          this.sid=sid
        }
    }
  }

  this.db=db

  this.InvokeObjectMethod(method)
}

func(this *Base) haveToken() bool{

  for k,_ := range this.c.Request.Header{
    if(k==Config["token_key"]){
      return true
    }
  }

  return false
}

func (this *Base) R(code int,message string,r gin.H){

  var data interface{}

  data=r
  if(this.action=="get"){
    data=r["get"]
  }

  this.c.JSON(200,gin.H{
    "code":code,
    "message":message,
    "data":data,
  })


}

func (this *Base) Error(code int,err interface{}){

  if(err!=nil){
    err_msg := fmt.Sprintf("%s", err)
    this.R(code,err_msg,gin.H{})
    this.c.Abort()
  }

}

func (this *Base) RouterToMethod(url string) string{

  url= strings.Split(url,"?")[0]
  url= strings.Split(url,"#")[0]


  this.path=url
  path := strings.Split( url , "/" )

  var controller string = "index"
  var action string = "index"

  if(len(path)>1 && path[1]!=""){
    controller=path[1]
  }

  if(len(path)>2 && path[2]!=""){
    action=path[2]
  }

  var args string
  for i,arg := range path{
    if(i<3){
      continue
    }
    if(i>3){
      args+="/"
    }
    args+=arg
  }
  this.args=args


  controller=Ucfirst(controller)

  this.controller=controller
  this.action=action


  method:=controller+"_"+action


  base_method := []string{"get","list","update","create","delete"}

  if(InArray(action,base_method)){
    return Ucfirst(action)
  }


  return method

}


func (this *Base) getTableName() string{
  name := Db["mysql_prefix"]+this.controller
  return name
}



func (this *Base) getFrom(fields []string) []string{

  var forms []string

  for _,field :=range fields{
    forms=append(forms,this.c.PostForm(field))
  }

  return forms
}

func (this *Base) prepare(sql string) *sql.Stmt{

  stmt,err := this.db.Prepare(sql)
  this.Error(-1,err)
  return stmt
}

func (this *Base) exec(stmt *sql.Stmt,args ...interface{}) sql.Result{

  res,err := stmt.Exec()
  this.Error(-2,err)

  return res
}

func (this *Base) query(sql string) *sql.Rows{
  rows,err := this.db.Query(sql)
  this.Error(-14,err)
  return rows
}

func (this *Base) Get(){
  dbname:=this.getTableName()
  sql:="select  * from "+dbname+" where id='"+this.args+"'"

  cache_key:="get_"+this.controller+"_"+this.args+"_"+sql

  record := make(map[string]string)

  if cache,ok := this.cache[cache_key] ; ok{


    record=cache.(map[string]string)

  }else{

    res,err:= this.db.Query(sql)

    this.Error(-17,err)

    columns, _ := res.Columns()
    scanArgs := make([]interface{}, len(columns))
    values := make([]interface{}, len(columns))

    for i := range values {
    	scanArgs[i] = &values[i]
    }

    res.Next()
    err=res.Scan(scanArgs...)
    this.Error(-18,err)

    if(err==nil){

      for i, col := range values {
        if col != nil {
          record[columns[i]] = string(col.([]byte))
        }
      }

      this.cache[cache_key]=record

    }



  }


  this.R(0,Lang["get_success"],gin.H{
    "get":record,
  })

}


func (this *Base) List(){


  page:=this.c.DefaultQuery("page","1")
  size:=this.c.DefaultQuery("size","10")
  orderBy:=this.c.DefaultQuery("orderBy","id")
  orderMethod:=this.c.DefaultQuery("orderMethod","desc")

  pageInt,_:=strconv.Atoi(page)
  sizeInt,_:=strconv.Atoi(size)

  sql:="select * from "+this.getTableName()+" "
  sql+="order by "+orderBy+" "+orderMethod+" "
  sql+="limit "+strconv.Itoa((pageInt-1)*sizeInt)+","+size+" "






  cache_key:="list_"+this.controller+"_"+sql
  cache_page_key:="list_page_"+this.controller+"_"+sql

  var records []interface{}
  listParam:=make(map[string]int)

  if cache,ok := this.cache[cache_key] ; ok{

    records=cache.([]interface{})
    listParam=this.cache[cache_page_key].(map[string]int)

  }else{


    var totalSize int
    this.db.QueryRow(`SELECT COUNT(*) AS count FROM `+this.getTableName()).Scan(&totalSize)
    totalPage := int(math.Ceil(float64(totalSize)/float64(sizeInt)))

    rows:=this.query(sql)

    columns, _ := rows.Columns()
    scanArgs := make([]interface{}, len(columns))
    values := make([]interface{}, len(columns))
    for i := range values {
    	scanArgs[i] = &values[i]
    }



    for rows.Next(){
      err := rows.Scan(scanArgs...)
      this.Error(-15,err)


        record := make(map[string]string)
        for i, col := range values {
      		if col != nil {
      			record[columns[i]] = string(col.([]byte))
      		}
      	}


        records=append(records,record)
    }

    listParam=map[string]int{
      "page":pageInt,
      "size":sizeInt,
      "totalPage":totalPage,
      "totalSize":totalSize,
    }

    this.cache[cache_key]=records
    this.cache[cache_page_key]=listParam

  }




  this.R(0,Lang["get_success"],gin.H{
    "list":records,
    "param":listParam,
  })

}



func (this *Base) Create() {

  dbname:=this.getTableName()
  var fields_key string
  var fields_value string


  this.c.Request.ParseMultipartForm(32 << 20)

  fields_key+="create_time"
  fields_value+="'"+GetTime()+"'"

  for k, v := range this.c.Request.PostForm {
      fields_key+=","
      fields_value+=","
      fields_key+=k
      fields_value+="'"+strings.Join(v,"")+"'"
  }

  sql := "insert into "+dbname+" ("+fields_key+") value ("+fields_value+")"
  fmt.Println(sql)

  stmt:=this.prepare(sql)
  res:=this.exec(stmt)

  lastInsertId,err := res.LastInsertId()

  this.Error(-3,err)

  go func(){
    this.clearCache()
  }()

  this.R(0,Lang["create_success"],gin.H{
    "insertId":lastInsertId,
  })
}

func (this *Base)Update() {
  dbname:=this.getTableName()

  sql:="update "+dbname+" set update_time='"+GetTime()+"'"

  this.c.Request.ParseMultipartForm(32 << 20)
  for k, v := range this.c.Request.PostForm {
      sql+=","
      sql+=k+"='"+strings.Join(v,"")+"'"
  }

  sql+=" where id='"+this.args+"'"


  stmt:=this.prepare(sql)
  res:=this.exec(stmt)
  update_count,_ := res.RowsAffected()

  go func(){
    this.clearCache()
  }()

  this.R(0,Lang["update_success"],gin.H{
    "update_count":update_count,
  })
}


func (this *Base)Delete(){
  dbname:=this.getTableName()
  sql := "delete from "+dbname+" where id in ("+this.args+")"
  stmt := this.prepare(sql)
  this.exec(stmt)

  go func(){
    this.clearCache()
  }()

  this.R(0,Lang["delete_success"],gin.H{})
}

func (this *Base) clearCache(){
  for k,_ :=range this.cache{

    if( strings.HasPrefix(k,"get_"+this.controller) || strings.HasPrefix(k,"list_"+this.controller) || strings.HasPrefix(k,"list_page_"+this.controller) ){
      delete(this.cache,k)
    }

  }

}
