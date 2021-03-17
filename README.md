# hicli
hydra项目协助开发工具，用于管理hydra后端、前端项目的创建，监控，自动启动等

##### 数据字典约束配置
 ```
PK:主键
SEQ:SEQ，mysql自增，oracle序列
C: 创建数据时的字段
R: 单条数据读取时的字段 
U: 修改数据时需要的字段
D: 删除，默认为更新字段状态值为1，D[(更新状态值)]
Q: 查询条件的字段
L：(前端页面)列表里列出的字段
order：查询时的order by字段；默认降序； OB[(顺序)]，越小越先排序
DI: 字典编号，数据表作为字典数据时的id字段
DN: 字典名称，数据表作为字典数据时的name字段
SL: "select"      //表单下拉框,默认使用dds字典表枚举,指定表名的SL[(字典表名)]
CB: "checkbox"    //表单复选框,默认使用dds字典表枚举,指定表名的CB[(字典表名)]
RD: "radio"       //表单单选框,默认使用dds字典表枚举,指定表名的RB[(字典表名)]
TA: "textarea"    //表单文本域
CC: "color-class"  //状态颜色过滤器
DATE: "date-picker" //表单日期选择器
DTIME: "datetime-picker" //表单日期时间选择器,
FIXED: 列表表单固定列
SORT: 列表表单排序列 ,sort(asc,顺序)
AFTER:在某个字段后面
列表自定义索引， 约定给当前表添加一行，字段名为_el_table_index,约束为索引大小

//C,R,U,Q,L子约束
f:前端过滤器，L(f:过滤器参数)
e:枚举参数

//枚举级联  #字段名，组件约束和子约束都可以使用

//排除表 ^表名

//{el_tab(表名，字段名/字段名，list),el_index(2)}

```