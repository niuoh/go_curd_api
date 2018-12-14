# go_curd_api
查一条：<br>
GET /tableName/get/:id<br>
查多条<br>
GET /tableName/list?page=1&size=10<br>
创建<br>
post /tableName/create   (data={"name":"lalaal","age":"24",...})<br>
更新<br>
post /tableName/update/:id  (data={"name":"lalaal","age":"24",...})<br>
删除<br>
get /tableName/delete/:id (支持 :id="1,2,3" 多个删除) <br>
