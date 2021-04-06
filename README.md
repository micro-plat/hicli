# hicli

对gitlab仓库的分组目录提供clone,pull,reset等操作指令

### 一、克隆项目

  hicli clone gitlab.100bm.cn/micro-plat/fas/apiserver

```sh
> hicli clone gitlab.100bm.cn/micro-plat/fas/apiserver
get clone https://gitlab.100bm.cn/micro-plat/fas/apiserver 
    /home/yanglei/work/src/gitlab.100bm.cn/micro-plat/fas/apiserver
```


### 二、克隆分组下所有项目

 hicli clone gitlab.100bm.cn/micro-plat/oms

```sh
> hicli clone gitlab.100bm.cn/micro-plat/oms
get clone https://gitlab.100bm.cn/micro-plat/oms/release/lxl/mgrsys/oms-web /home/yanglei/work/src/gitlab.100bm.cn/micro-plat/oms/release/lxl/mgrsys/oms-web
get clone https://gitlab.100bm.cn/micro-plat/oms/release/lxl/mgrsys/oms-api /home/yanglei/work/src/gitlab.100bm.cn/micro-plat/oms/release/lxl/mgrsys/oms-api

```


### 三、拉取分组下所有项目的指定分支

 hicli clone gitlab.100bm.cn/micro-plat/fas -branch dev

 ```sh
> hicli pull gitlab.100bm.cn/micro-plat/fas -branch dev
get clone https://gitlab.100bm.cn/micro-plat/fas/docs /home/yanglei/work/src/gitlab.100bm.cn/micro-plat/fas/docs
/home/yanglei/work/src/gitlab.100bm.cn/micro-plat/fas/docs > git pull origin dev:dev
....
 ```

 ### 四、撤销分组下所有项目的修改

 hicli reset gitlab.100bm.cn/micro-plat/fas -branch dev

 ```sh
> hicli reset gitlab.100bm.cn/micro-plat/fas -branch dev
/home/yanglei/work/src/gitlab.100bm.cn/micro-plat/fas/fds > git reset --hard
HEAD 现在位于 1575b23 RPC SDK
/home/yanglei/work/src/gitlab.100bm.cn/micro-plat/fas/fds > git checkout dev
切换到分支 'dev'
/home/yanglei/work/src/gitlab.100bm.cn/micro-plat/fas/fds > git reset --hard
HEAD 现在位于 1575b23 RPC SDK
/home/yanglei/work/src/gitlab.100bm.cn/micro-plat/fas/apiserver > git reset --hard
HEAD 现在位于 420caa3 
/home/yanglei/work/src/gitlab.100bm.cn/micro-plat/fas/apiserver > git checkout dev
切换到分支 'dev'
/home/yanglei/work/src/gitlab.100bm.cn/micro-plat/fas/apiserver > git reset --hard
HEAD 现在位于 420caa3
....
 ```
### 五、监控并自动重启hydra服务
对hydra项目文件进行监控，并自动重启
```
示例：
hicli server [项目路径] -run="run_tags" -install="-tags install_tags"
```


### 六、根据markdown创建项目
**Notice:生成文件会根据安装的tags生成对应类型的文件，默认为mysql**
```
生成mysql: go install
生成oracle: go install -tags "oracle"
```
#### 根据数据库创建对应markdown数据字典
```
hicli dic create  -db connect_str -cover -f
connect_str可参考hydra配置db连接时的连接串
oracle数据库：oracle:scheme/pwd@orcl136
mysql数据库：mysql:root:pwd@tcp(192.168.0.36:3306)/scheme
```
#### 根据markdown数据字典创建对应的数据库建表文件
```
hicli db create docs/dic.md [outpath] -t "" -d -s -v -g
```

##### 创建前端vue项目
1.创建vue项目文件
```
hicli ui create project_name
```
2.创建vue项目页面文件
```
hicli ui page md文档路径 [输出文件路径] -table [指定表名] -f -cover
```
##### 创建后端hydra项目
1.创建web服务文件
```
hicli app create server_name
```
2.创建服务层和sql文件
```
hicli app service md文档路径 [输出文件路径] -table [指定表名] --exclude [排除表名] -f -cover
```

##### 数据字典约束配置（约束关键字不区分大小写，内容区分大小写）
 ```
字段约束关键字
    PK:主键
    SEQ:序列，mysql:SEQ，自增大小为字典的默认值，oracle:SEQ(序列名,起始值,增长值)
    UNQ:唯一索引，UNQ[(索引名,字段在联合索引中的位置)]
    IDX:索引，IDX[(索引名,字段在联合索引中的位置)]
    C: 创建数据时的字段
    R: 单条数据读取时的字段 
    U: 修改数据时需要的字段
    D: 删除，默认为更新字段状态值为1，D[(更新状态值)]
    Q: 查询条件的字段
    L：(前端页面)列表里列出的字段
    ORDER：后端sql查询时的order by字段；默认降序； ORDER[(asc|desc,字段排序顺序)]，越小越先排序
    SORT：前后端联合排序；默认降序； SORT[(asc|desc,字段排序顺序)]，越小越先排序；不可与ORDER同时使用
    DI: 字典编号，数据表作为字典数据时的id字段
    DN: 字典名称，数据表作为字典数据时的name字段
    SL: 前端页面下拉框,默认使用dds字典表枚举,指定表名的SL[(表名)]
    SLM: 前端页面可多选表单下拉框,默认使用dds字典表枚举,指定表名的SL[(表名)]
    CB: 前端页面复选框,默认使用dds字典表枚举,指定表名的CB[(表名)]
    RD: 前端页面单选框,默认使用dds字典表枚举,指定表名的RB[(表名)]
    TA: 前端页面文本域
    CC: 前端页面状态颜色过滤器
    DATE: 前端页面日期选择器
    DTIME: 前端页面日期时间选择器
    FIXED: 前端页面列表表单固定列
    AFTER: 前端列表字段在指定字段后面，AFTER(字段名)

关键字C,R,U,Q,L约定子约束
    f:前端字段展示过滤器参数，示例:L(f:过滤器参数)
    e:枚举参数，示例:L(e:#字段名),含义为与指定字段为级联枚举；L(e:表名)，含义为指定表枚举；L(e:类型名)，含义为dds指定枚举

字典表名的约束关键字
	^:排除表，示例：^表名
  复合功能:{el_tab(),el_index(),el_btn(),el_btn1(),el_btn2}
    el_tab:生成关联的详情页，{el_tab(表名，字段名/字段名，list)}
    el_index:生成列表页面的索引，示例，索引大小为2，{el_index(2)}
    el_btn:生成页面的按钮操作,多个按钮,以el_btn,el_btn1,el_btn2...顺序配置，最多为10个
        状态扭转:el_btn(name:方法名,desc:1-禁用|2-启用,confirm:确定要进行此操作吗,key:)
        关联表数据及数据更新:el_btn(name:方法名,desc:按钮名称,confirm:确定要进行此操作吗,table:关联表:字段1/字段2|关联表2,key:btn_key)
	
```