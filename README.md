# go_curd_api
查一条：
GET /tableName/get/:id
查多条
GET /tableName/list?page=1&size=10
创建
post /tableName/create   (data={"name":"lalaal","age":"24",...})
更新
post /tableName/update/:id  (data={"name":"lalaal","age":"24",...})
删除
get /tableName/delete/:id (支持 :id="1,2,3" 多个删除)
