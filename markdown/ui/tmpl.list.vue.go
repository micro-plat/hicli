package ui

const TmplList = `
{{- $len := 32 -}}
{{- $rows := .Rows -}}
{{- $desc := .Desc -}}
{{- $name := .Name -}}
{{- $pks := .|pks -}}
{{- $tb :=. -}}
{{- $sort:=.Rows|sort -}}
{{- $choose:= false -}}
{{- $drange:= mkSlice -}}
{{- $btn:=.ListBtnInfo -}}
<template>
	<div class="panel panel-default">
    <!-- query start -->
		<div class="panel-body" id="panel-body">
			<el-form ref="form" size="small" :inline="true" class="form-inline pull-left">
			{{- range $i,$c:=$rows|query}}
				{{- if $c.Con|CSCR}}
				<el-form-item>
					<el-cascader
						placeholder="请选择{{$c.Desc|shortName}}"
						v-model="{{$c.Con|cscrCon|lowerName}}Value"
						@change="{{$c.Con|cscrCon|lowerName}}Change"
						:show-all-levels="false"
						collapse-tags
						:options="{{$c.Con|cscrCon|lowerName|lowerName}}"
						:props="{ multiple: true, label: 'name' }"
						clearable
					></el-cascader>
				</el-form-item>
				{{- else if $c.Con|TA}}
				<el-form-item>
					<el-input size="small" maxlength="{{$c.Len}}" type="textarea" :rows="2" placeholder="请输入{{$c.Desc|shortName}}" v-model="queryData.{{$c.Name}}">
					</el-input>
				</el-form-item>
				{{- else if or ($c.Con|SL) ($c.Con|SLM) }}
				<el-form-item>
					<el-select size="small" v-model="queryData.{{$c.Name}}" clearable filterable class="input-cos" placeholder="请选择{{$c.Desc|shortName}}"
					{{- if (qDicCName $c.Name $tb) }} @change="set{{(qDicCName $c.Name $tb)|upperName}}(queryData.{{$c.Name}})"
					{{- else if (qGroupCName $c.Name $tb) }} @change="set{{$c.Name|upperName}}Group" 
					{{- else if or (qGroupPName $c.Con $tb) (qDicPName $c.Con $tb) }} @change="handleChooseTool()"{{$choose = true}}{{- end}} >
						<el-option value="" label="全部"></el-option>
						<el-option v-for="(item, index) in {{$c.Name|lowerName}}" :key="index" :value="item.value" :label="item.name"></el-option>
					</el-select>
				</el-form-item>
				{{- else if ($c.Con|DRANGE) }}{{$drange = $c.Con|drangeCon|drangeValue }}
				<el-form-item>
					<el-date-picker
						v-model="times"
						type="{{dateType $c.Con ($c.Con|qfCon)}}range"
						:clearable="false"
						:picker-options="pickerOptions"
						range-separator="至"
						start-placeholder="开始日期"
						end-placeholder="结束日期"
						align="right"
						size="small"
					>
					</el-date-picker>
        </el-form-item>
				{{- else if or ($c.Con|DTIME) ($c.Con|DATE) ($c.Type|isTime) }}
				<el-form-item label="{{$c.Desc|shortName}}:">
					<el-date-picker size="small" class="input-cos" v-model="{{$c.Name|lowerName}}" type="{{dateType $c.Con ($c.Con|qfCon)}}" value-format="{{dateFormat $c.Con ($c.Con|qfCon)}}"  placeholder="选择日期"></el-date-picker>
				</el-form-item>
				{{- else if $c.Con|CB }}
				<el-form-item label="{{$c.Desc|shortName}}:">
          <el-checkbox-group size="small" v-model="{{$c.Name|lowerName}}Array">
						<el-checkbox v-for="(item, index) in channelNo" :key="index" :value="item.value" :label="item.value">{{"{{item.name}}"}}</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
				{{- else if $c.Con|DMI }}
				{{- else}}{{$c|setIsInput}}
				{{- end}}
			{{- end}}
			{{- if gt ($rows|query|dropmenurow|len) 0}}
				<el-form-item>
					<el-input clearable size="small" v-model="queryContent" placeholder="请输入查询内容">
						<el-dropdown size="small" slot="prepend">
							<el-button size="small" style="width: 100px; font-size: 0.677083vw">{{"{{ currentDropItem.name }}"}}<i class="el-icon-arrow-down el-icon--right"></i> </el-button>
							<el-dropdown-menu slot="dropdown">
								<el-dropdown-item v-for="(item, index) in dropMenu" @click.native="dropmenu(item)" :key="index" :value="item.value" :label="item.name">{{"{{item.name}}"}}</el-dropdown-item>
							</el-dropdown-menu>
						</el-dropdown>
					</el-input>
				</el-form-item>
			{{- end}}
			{{- range $i,$c:=$rows|query}}
				{{- if $c.IsInput}}
				<el-form-item>
					<el-input clearable size="small" prefix-icon="el-icon-search" v-model="queryData.{{$c.Name}}" placeholder="请输入{{$c.Desc|shortName}}">
					</el-input>
				</el-form-item>
				{{- end}}
			{{- end}}
				{{- if and (not .BtnShowQuery) (gt ($rows|query|len) 0)}}
				<el-form-item>
					<el-button type="primary" @click="queryDatas" size="small">查询</el-button>
				</el-form-item>
				{{- end}}
				{{- if and (not .BtnShowAdd) (gt ($rows|create|len) 0)}}
				<el-form-item>
					<el-button type="success" size="small" @click="showAdd">添加</el-button>
				</el-form-item>
				{{- end}}
				{{- if gt ($rows|export|len) 0}}
				<el-form-item>
					<el-button type="success" @click="exportExcl" size="small">导出</el-button>
				</el-form-item>
				{{- end}}
				{{- if gt ($tb.DownloadInfo.Title|len) 0}}
				<el-form-item>
					<el-button type="text" @click="download" size="small" style="font-size:14px">下载模版</el-button>
				</el-form-item>
				{{- end}}
				{{- range $i,$c := $tb.BatchInfo }}
				<el-form-item>
					<el-button type="primary" v-if="multipleSelection.length != 0" size="small" @click="{{$c.Method}}()">{{$c.Name}}</el-button>
					<el-button type="primary" v-else size="small" disabled>{{$c.Name}}</el-button>
				</el-form-item>
				{{- end}}
				{{- range $i,$c:=$tb.QueryDialogs}}
				<el-form-item>
					<el-button type="primary" size="small" @click="show{{$c.Method|upperName}}">{{$c.Name}}</el-button>
				</el-form-item>
				{{- end}}
				{{- range $i,$c:=$tb.QueryBtnInfo}}
				<el-form-item>
					<el-button type="primary" size="small" {{- if $c.Condition }} v-if="{{$c.Condition}}"{{end}} @click="{{$c.Method}}">{{$c.Name}}</el-button>
				</el-form-item>
				{{- end}}
			</el-form>
		</div>
    <!-- query end -->

    <!-- list start-->
		<el-scrollbar style="height:100%">
			<el-table :data="dataList.items" stripe style="width: 100%" size="small" :height="maxHeight" {{if gt ($sort|len) 0}}@sort-change="sort"{{end}}
			{{- if gt ($tb.BatchInfo|len) 0 }} @selection-change="handleSelectionChange" {{end}}>
			  {{- if gt ($tb.BatchInfo|len) 0 }}
				<el-table-column type="selection" :selectable="selectableCheckbox" width="24"></el-table-column>
				{{- end}}
				{{- if gt $tb.ELTableIndex 0}}
				<el-table-column type="index" fixed	:index="indexMethod" label="序号"></el-table-column>
				{{- end}}
				{{- range $i,$c:=$rows|list}}
				<el-table-column {{- if $c.Con|FIXED}} fixed{{end}} {{- if $c.Con|SORT}} sortable="custom"{{end}} prop="{{$c.Name}}" label="{{$c.Desc|shortName}}" align="center">
					<template slot-scope="scope">
						{{- if $c.Con|LINK}}
						<el-tooltip class="item" effect="dark" :content="scope.row.{{$c.Name}}" placement="top">
						<el-button type="text" size="small" @click="{{if ($c.Con|linkCon)}}link{{$c.Name|upperName}}{{else}}showDetail{{end}}(scope.row)">
						{{- end}}
				{{- if or ($c.Con|SL) ($c.Con|SLM) ($c.Con|CB) ($c.Con|RD) ($c.Con|leCon) ($c.Con|CSCR)}}
						<span {{if ($c.Con|CC)}}:class="scope.row.{{$c.Name}}|fltrTextColor"{{end}}>{{"{{scope.row."}}{{$c.Name}} | fltrEnum("{{or (dicName $c.Con ($c.Con|leCon) $tb) ($c.Name|lower)}}")}}</span>
				{{- else if and ($c.Type|isString) (or (gt $c.Len $len) (eq $c.Len 0) )}}
						<el-tooltip class="item" v-if="scope.row.{{$c.Name}} && scope.row.{{$c.Name}}.length > {{or ($c.Con|lfCon) "20"}}" effect="dark" placement="top">
							<div slot="content" style="width: 110px">{{"{{scope.row."}}{{$c.Name}}}}</div>
							<span>{{"{{scope.row."}}{{$c.Name}} | fltrSubstr({{or ($c.Con|lfCon) "20"}}) }}</span>
						</el-tooltip>
						<span v-else>{{"{{scope.row."}}{{$c.Name}} | fltrEmpty }}</span>
				{{- else if ($c.Con|fIsNofltr)}}
						<span>{{"{{scope.row."}}{{$c.Name}} | fltrEmpty }}</span>
				{{- else if and (or ($c.Type|isInt64) ($c.Type|isInt) ) (ne $c.Name ($pks|firstStr))}}
						<span>{{"{{scope.row."}}{{$c.Name}} | fltrNumberFormat({{or ($c.Con|lfCon) "0"}})}}</span>
				{{- else if $c.Type|isDecimal }}
						<span>{{"{{scope.row."}}{{$c.Name}} | fltrNumberFormat({{or ($c.Con|lfCon) "2"}})}}</span>
				{{- else if $c.Type|isTime }}
						<div>{{"{{scope.row."}}{{$c.Name}} | fltrDate("{{ or (dateFormat $c.Con ($c.Con|lfCon)) "yyyy-MM-dd HH:mm:ss"}}") }}</div>
				{{- else}}
						<span>{{"{{scope.row."}}{{$c.Name}} | fltrEmpty }}</span>
				{{- end}}
					{{- if $c.Con|LINK}}
						</el-button>
						</el-tooltip>
					{{- end}}
					</template>
				</el-table-column>
				{{- end}}
				<el-table-column label="操作" align="center">
					<template slot-scope="scope">
						{{- if and  (not .BtnShowEdit) (gt ($rows|update|len) 0)}}
						<el-button type="text" size="mini" @click="showEdit(scope.row)">编辑</el-button>
						{{- end}}
						{{- if and (not .BtnDel) (gt ($rows|delete|len) 0)}}
						<el-button type="text" size="mini" @click="del(scope.row)">删除</el-button>
						{{- end}}

						{{- range $i,$c:=$tb.ListDialogs}}
						<el-button type="text" {{if $c.Condition }}v-if="{{$c.Condition}}"{{end}} size="mini" @click="show{{$c.Method|upperName}}(scope.row)">{{$c.Name}}</el-button>
						{{- end}}

						{{- if and (not .BtnShowDetail) (gt ($rows|detail|len) 0)}}
						<el-button type="text" size="mini" @click="showDetail(scope.row)">详情</el-button>
						{{- end}}

						{{- range $i,$c:= $btn }}
						<el-button type="text" size="mini" {{- if $c.Condition }} v-if="{{$c.Condition}}"{{end}} @click="{{$c.Method}}(scope.row)">{{$c.Name}}</el-button>
						{{- end}}
					</template>
				</el-table-column>
			</el-table>
		</el-scrollbar>
		<!-- list end-->

		{{- if and (not .BtnShowAdd) (gt ($rows|create|len) 0)}}

		<!-- Add Form -->
		<Add ref="Add" :refresh="query"></Add>
		<!--Add Form -->
		{{- end}}

		{{- if and (not .BtnShowEdit) (gt ($rows|update|len) 0)}}

		<!-- edit Form start-->
		<Edit ref="Edit" :refresh="query"></Edit>
		<!-- edit Form end-->
		{{- end}}

		{{- range $i,$c:=$tb.ListDialogs}}

		<!-- {{$c.Method|upperName}} Form -->
		<{{$c.Method|upperName}} ref="{{$c.Method|upperName}}" :refresh="query"></{{$c.Method|upperName}}>
		<!--{{$c.Method|upperName}} Form -->
		{{- end}}

		{{- range $i,$c:=$tb.QueryDialogs}}

		<!-- {{$c.Method|upperName}} Form -->
		<{{$c.Method|upperName}} ref="{{$c.Method|upperName}}" :refresh="query"></{{$c.Method|upperName}}>
		<!--{{$c.Method|upperName}} Form -->
		{{- end}}

		<!-- pagination start -->
		<div class="page-pagination">
		<el-pagination
			@size-change="pageSizeChange"
			@current-change="pageIndexChange"
			:current-page="paging.pi"
			:page-size="paging.ps"
			:page-sizes="paging.sizes"
			layout="total, sizes, prev, pager, next, jumper"
			:total="dataList.count">
		</el-pagination>
		</div>
		<!-- pagination end -->

	</div>
</template>


<script>
{{- if and (not .BtnShowAdd) (gt ($rows|create|len) 0)}}
import Add from "./{{.Name|rmhd|l2d}}.add"
{{- end}}
{{- if and (not .BtnShowEdit) (gt ($rows|update|len) 0)}}
import Edit from "./{{.Name|rmhd|l2d}}.edit"
{{- end}}
{{- range $i,$c:=$tb.ListDialogs}}
import {{$c.Method|upperName}} from "{{$c.Path}}"
{{- end}}
{{- range $i,$c:=$tb.QueryDialogs}}
import {{$c.Method|upperName}} from "{{$c.Path}}"
{{- end}}
export default {
	name: "{{$name|rmhd|varName}}",
  components: {
		{{- if and (not .BtnShowAdd) (gt ($rows|create|len) 0)}}
		Add,
		{{- end}}
		{{- if and (not .BtnShowEdit) (gt ($rows|update|len) 0)}}
		Edit,
		{{- end}}
		{{- range $i,$c:=$tb.ListDialogs}}
		{{$c.Method|upperName}},
		{{- end}}
		{{- range $i,$c:=$tb.QueryDialogs}}
		{{$c.Method|upperName}},
		{{- end}}
  },
  data () {
		return {
			paging: {ps: 10, pi: 1,total:0,sizes:[5, 10, 20, 50]},
			queryData:{},               //查询数据对象
			{{- range $i,$c:=$rows|query -}}
			{{- if $c.Con|CSCR}}
			{{$c.Name|lowerName}}: this.$enum.get("{{or (dicName $c.Con ($c.Con|qeCon) $tb) ($c.Name|lower)}}"),
			{{$c.Con|cscrCon|lowerName}}Value: [],
			{{$c.Con|cscrCon|lowerName}}: [],
			{{$c.Con|cscrCon|lowerName}}Selected: "{{$c.Con|cscrDefault}}",
			{{- else if or ($c.Con|SL) ($c.Con|SLM) ($c.Con|RD) }}
			{{$c.Name|lowerName}}: this.$enum.get("{{or (dicName $c.Con ($c.Con|qeCon) $tb) ($c.Name|lower)}}"),
			{{- else if $c.Con|CB }}
			{{$c.Name|lowerName}}: this.$enum.get("{{or (dicName $c.Con ($c.Con|qeCon) $tb) ($c.Name|lower)}}"),
			{{$c.Name|lowerName}}Array: [],
			{{- end}}
			{{- if ($c.Con|DRANGE)}}
      times: '',
      pickerOptions: {
        shortcuts: [{
						text: '今天',
						onClick(picker) {
							const end = new Date();
							const start = new Date();
							picker.$emit('pick', [start, end]);
						}
					},
					{
						text: '昨天',
						onClick(picker) {
							const end = new Date();
							const start = new Date();
							start.setTime(start.getTime() - 3600 * 1000 * 24 * 1);
							end.setTime(end.getTime() - 3600 * 1000 * 24 * 1);
							picker.$emit('pick', [start, end]);
						}
					},
          {
            text: '最近三天',
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 3);
              picker.$emit('pick', [start, end]);
            }
          }, {
            text: '最近一周',
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 7);
              picker.$emit('pick', [start, end]);
            }
          }, {
            text: '最近一个月',
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 30);
              picker.$emit('pick', [start, end]);
            }
          }]
      },
			{{- else if or ($c.Con|DTIME) ($c.Con|DATE) ($c.Type|isTime) }}
			{{$c.Name|lowerName}}: this.$utility.dateFormat(new Date(),"{{dateFormatDef $c.Con ($c.Con|qfCon)}}"),{{end}}
      {{- end}}
			{{- if gt ($sort|len) 0}}
			order: "",
			{{- end}}
			{{- if gt ($tb.BatchInfo|len) 0 }}
			multipleSelection: [],
			{{- end}}
			{{- if gt ($rows|query|dropmenurow|len) 0}}
			queryContent:"",
      preDropItem: {},
      currentDropItem: {},
      dropMenu: [
				{{- range $i,$c:=$rows|query|dropmenurow}}
				{ name: "{{$c.Desc|shortName}}", value: "{{$c.Name}}" },
				{{- end}}
      ],
			{{- end}}
			dataList: {count: 0,items: []}, //表单数据对象,
			maxHeight: 0
		}
  },
  created(){
  },
  mounted(){
		this.$nextTick(()=>{
			this.maxHeight = this.$utility.getTableHeight("panel-body")
		})
    this.init()
  },
	methods:{
    /**初始化操作**/
    init(){
			{{- range $i,$c:=$rows|query}}
			{{- if ($c.Con|CSCR) }}
			this.set{{$c.Con|cscrCon|upperName}}()
			{{- end}}
			{{- end}}
			{{- if gt ($rows|query|dropmenurow|len) 0}}
			this.preDropItem = this.dropMenu[0]
      this.currentDropItem = this.dropMenu[0]
			{{- end}}
			{{- if gt ($drange|len) 0}}
			var now = new Date()
			now.setTime(new Date().getTime() - 3600 * 1000 * 24 * {{index $drange 0}})
			{{- if eq ($drange|len) 1}}
			this.times = [now, new Date()]
			{{- else }}
			var end = new Date()
			end.setTime(new Date().getTime() - 3600 * 1000 * 24 * {{index $drange 1}})
			this.times = [now, end]
			{{- end}}
			{{- end}}
      this.query()
		},
		{{- if gt ($sort|len) 0}}
		sort(column) {
      if (column.order == "ascending") {
        this.order ="t." + column.prop + " " + "asc"
      } else if (column.order == "descending") {
        this.order ="t." + column.prop + " " + "desc"
      } else {
        this.order = ""
      }
      this.query()
    },
		{{- end}}
		{{- if $choose}}
		handleChooseTool() {
      this.$forceUpdate()
    },{{end}}
		{{- if gt $tb.ELTableIndex 0}}
		indexMethod(index) {
			return index * {{$tb.ELTableIndex}};
		},
		{{- end}}
		{{- range $i,$c:=$rows|query}}
		{{- if ($c.Con|CSCR) }}
		set{{$c.Con|cscrCon|upperName}}() {
      var that = this
			var selected = this.{{$c.Con|cscrCon|lowerName}}Selected
      this.{{$c.Con|cscrCon|lowerName}} = this.$enum.get("{{$c.Con|cscrCon}}")
      this.{{$c.Con|cscrCon|lowerName}}.forEach(function (item) {
        for (const el of that.{{$c.Name|lowerName}}) {
          if (el.group_code == item.value) {
						if (!item.children) {
              item.children = []
            }
            item.children.push(el)
            if (selected == "" || selected.split(",").includes(item.value)){
              that.{{$c.Con|cscrCon|lowerName}}Value.push([item.value, el.value])
            }
          }
        }
      })
			this.{{$c.Con|cscrCon|lowerName}}Change(this.{{$c.Con|cscrCon|lowerName}}Value)
    },
    {{$c.Con|cscrCon|lowerName}}Change(val) {
      let vals = [];
      val.forEach((item) => {
        vals.push(item[item.length - 1]);
      })
      this.queryData.{{$c.Name}} = vals.join(',');
    },
		{{- end}}
		{{- if (qDicPName $c.Con $tb) }}
		set{{$c.Name|upperName}}(pid){
			this.queryData.{{$c.Name}} = ""
			this.{{$c.Name|lowerName}}=this.$enum.get("{{or (dicName $c.Con ($c.Con|qeCon) $tb) ($c.Name|lower)}}",pid)
		},
		{{- end}}
		{{- if (qGroupCName $c.Name $tb) }}
		set{{$c.Name|upperName}}Group(value){
			var obj = this.{{$c.Name|lowerName}}.find((item) => {
        return item.value === value
      })
			if (obj){
				{{- range $i,$c1:=(qgroup $c.Name $tb)}}
				{{- if (qDicCName $c1.Name $tb)}}
				this.set{{(qDicCName $c1.Name $tb)|upperName}}(obj.{{$c1.Name}})
				{{- end}}
				{{- end}}
				{{- range $i,$c1:=(qgroup $c.Name $tb)}}
				this.queryData.{{$c1.Name}} = obj.{{$c1.Name}}
				{{- if (qGroupCName $c1.Name $tb)}}
				this.set{{$c1.Name|upperName}}Group(this.queryData.{{$c1.Name}})
				{{- end}}
				{{- end}}
			}
		},
		{{- end}}
		{{- end }}
		{{- if gt ($rows|query|dropmenurow|len) 0}}
		dropmenu(item) {
      this.currentDropItem = item
      if (this.currentDropItem != this.preDropItem) {
        this.queryData[this.preDropItem.value] = ""
				this.preDropItem=item
      }
    },
		{{- end}}
    /**查询数据并赋值*/
		queryDatas() {
      this.paging.pi = 1
      this.query()
    },
    query(){
			{{- if gt ($rows|query|dropmenurow|len) 0}}
			this.queryData[this.currentDropItem.value] = this.queryContent
			{{- end}}
      this.queryData.pi = this.paging.pi
			this.queryData.ps = this.paging.ps
			{{- range $i,$c:=$rows|query -}}
			{{- if ($c.Con|DRANGE)}}
			this.queryData.start_time = this.$utility.dateFormat(this.times[0],"{{dateFormat $c.Con ($c.Con|qfCon)}}");
      this.queryData.end_time = this.$utility.dateFormat(this.times[1],"{{dateFormat $c.Con ($c.Con|qfCon)}}");
			{{- else if or ($c.Con|DTIME) ($c.Con|DATE) ($c.Type|isTime) }}
			this.queryData.{{$c.Name}} = this.$utility.dateFormat(this.{{$c.Name|lowerName}},"{{dateFormat $c.Con ($c.Con|qfCon)}}")
			{{- else if ($c.Con|CB) }}
			this.queryData.{{$c.Name}} = this.{{$c.Name|lowerName}}Array.toString()
			{{- end -}}
      {{- end}}
			{{- if gt ($sort|len) 0}}
			this.queryData.order_by = this.order
			{{- end}}
      let res = this.$http.xget("/{{.Name|rmhd|rpath}}/{{or .QueryHandler "query"}}",this.$utility.delEmptyProperty(this.queryData))
			this.dataList.items = res.items || []
			this.dataList.count = res.count
    },
    /**改变页容量*/
		pageSizeChange(val) {
      this.paging.ps = val
      this.query()
    },
    /**改变当前页码*/
    pageIndexChange(val) {
      this.paging.pi = val
      this.query()
    },
    /**重置添加表单*/
    resetForm(formName) {
      this.$refs[formName].resetFields();
		},
		{{- if gt ($rows|detail|len) 0}}
		showDetail(val){
			var data = {
        {{range $i,$c:=$pks}}{{$c}}: val.{{$c}},{{end}}
      }
      this.$emit("addTab","详情"+val.{{range $i,$c:=$pks}}{{$c}}{{end}},"/{{.Name|rmhd|rpath}}/detail",data);
		},
		{{- end}}
		{{- if and (not .BtnShowAdd) (gt ($rows|create|len) 0)}}
    showAdd(){
      this.$refs.Add.show();
		},
		{{- end}}
		{{- if and (not .BtnShowEdit) (gt ($rows|update|len) 0)}}
    showEdit(val) {
      this.$refs.Edit.{{range $i,$c:=$pks}}{{$c}} = val.{{$c}};{{end}}
      this.$refs.Edit.show();
		},
		{{- end}}
		{{- range $i,$c:=$tb.ListDialogs}}
		show{{$c.Method|upperName}}(val) {
			this.$refs.{{$c.Method|upperName}}.{{range $i,$c:=$pks}}{{$c}} = val.{{$c}};{{end}}
      this.$refs.{{$c.Method|upperName}}.show();
		},
		{{- end}}
		{{- range $i,$c:=$tb.QueryDialogs}}
		show{{$c.Method|upperName}}() {
      this.$refs.{{$c.Method|upperName}}.show();
		},
		{{- end}}

		{{- range $i,$c:=$tb.QueryBtnInfo}}
		{{- if not $c.IsQuery}}
		{{$c.Method}}(){
			var data = this.queryData
			{{- if $c.Confirm}}
      this.$confirm("{{$c.Confirm}}?", "提示", { confirmButtonText: "确定", cancelButtonText: "取消", type: "warning" })
        .then(() => {
			{{- end}}
					this.$http.post("/{{$tb.Name|rmhd|rpath}}/{{or $c.Handler ($c.Name|lowerName)}}", data, {}, true, true)
						.then(res => {
							this.dialogFormVisible = false;
							this.query()
						})
			{{- if $c.Confirm}}
				});
			{{- end}}
		},
		{{- end}}
		{{- end}}

		{{- range $i,$c:= $btn }}
		{{$c.Method}}(val){
			var data = {
				{{- range $i,$c:=$c.Rows}}
				{{$c.Name}} :val.{{$c.Name}},
				{{- end}}
			}
			{{- if $c.Confirm}}
      this.$confirm("{{$c.Confirm}}?", "提示", { confirmButtonText: "确定", cancelButtonText: "取消", type: "warning" })
        .then(() => {
			{{- end}}
					this.$http.post("/{{$tb.Name|rmhd|rpath}}/{{or $c.Handler ($c.Name|lowerName)}}", data, {}, true, true)
						.then(res => {
							this.dialogFormVisible = false;
							this.query()
						})
			{{- if $c.Confirm}}
				});
			{{- end}}
		},
		{{- end}}

		{{- if gt ($rows|export|len) 0}}
		exportExcl() {
			this.queryData.pi = this.paging.pi
			this.queryData.ps = this.dataList.count
			{{- range $i,$c:=$rows|query -}}
			{{- if ($c.Con|DRANGE)}}
			this.queryData.start_time = this.$utility.dateFormat(this.times[0],"{{dateFormat $c.Con ($c.Con|qfCon)}}");
      this.queryData.end_time = this.$utility.dateFormat(this.times[1],"{{dateFormat $c.Con ($c.Con|qfCon)}}");
			{{- else if or ($c.Con|DTIME) ($c.Con|DATE) ($c.Type|isTime) }}
			this.queryData.{{$c.Name}} = this.$utility.dateFormat(this.{{$c.Name|lowerName}},"{{dateFormat $c.Con ($c.Con|qfCon)}}")
			{{- else if ($c.Con|CB) }}
			this.queryData.{{$c.Name}} = this.{{$c.Name|lowerName}}Array.toString()
			{{- end -}}
      {{- end}}
			{{- if gt ($sort|len) 0}}
			this.queryData.order_by = this.order
			{{- end}}
      this.$http.post("/{{.Name|rmhd|rpath}}/export",this.$utility.delEmptyProperty(this.queryData))
        .then(res => {
          var header = [
					{{- range $i,$c:=$rows|export}}
					{ Key: "{{$c.Name}}",	Txt: "{{$c.Desc}}" },
					{{- end}}
          ];
          this.BuildExcel("{{$desc}}.xlsx", [header], res.items || [], {
						{{- range $i,$c:=$rows|export}}
						{{- if and (or ($c.Con|SL) ($c.Con|SLM) ($c.Con|CB) ($c.Con|RD) ($c.Con|leCon)) (not ($c.Con|eptCon|isTrue)) }}
						{{$c.Name}}: this.$enum.get("{{or (dicName $c.Con ($c.Con|leCon) $tb) ($c.Name|lower)}}"),
						{{- end}}
						{{- end}}
          });
        });
    },
		{{- end}}

		{{- if gt ($tb.DownloadInfo.Title|len) 0}}
		download() {
      var data = [
				[
				{{- range $i,$c:=$tb.DownloadInfo.Title}}
					"{{$c}}",
				{{- end}}
        ]
      ];
      this.ExportTemplate(data, "模板.xlsx")
    },
		{{- end}}

		{{- if gt ($tb.BatchInfo|len) 0 }}
		selectableCheckbox(row, rowIndex) {
			{{- if $tb.BatchInfo.IsCondition}}
      if ({{$tb.BatchInfo.Condition}}) {
        return true;
      }
      return false;
			{{- else}}
			return true;
			{{- end}}
    },
		handleSelectionChange(val) {
      this.multipleSelection = val;
    },
		{{- end}}
		{{- range $i,$c := $tb.BatchInfo }}
		{{$c.Method}}() {
      var data = []
      this.multipleSelection.forEach(row => {
        data.push(row.{{range $i,$c:=$pks}}{{$c}}{{end}})
      });
			{{- if $c.Confirm}}
      this.$confirm("{{$c.Confirm}}?", "提示", { confirmButtonText: "确定", cancelButtonText: "取消", type: "warning" })
        .then(() => {
			{{- end}}
					this.$http.post("/{{$tb.Name|rmhd|rpath}}/{{$c.Handler}}", { {{range $i,$c:=$pks}}{{$c}}s{{end}}: data.join(",") }, {}, true, true)
						.then(res => {			
							this.query()
						})
			{{- if $c.Confirm}}
			});
		{{- end}}
    },
		{{- end}}

		{{- range $i,$c:=$rows|list}}
		{{- if ($c.Con|linkCon)}}
		link{{$c.Name|upperName}}(val){
			var data = {
        {{$c.Name}}: val.{{$c.Name}},
      }
      this.$emit("addTab","详情"+val.{{$c.Name}},"/{{$c.Con|linkCon|rmhd|rpath}}",data);
		},
		{{- end}}
		{{- end }}

		{{- if gt ($rows|delete|len) 0}}
    del(val){
			this.$confirm("此操作将永久删除该数据, 是否继续?", "提示", {confirmButtonText: "确定", cancelButtonText: "取消", type: "warning"})
			.then(() => {
				this.$http.del("/{{.Name|rmhd|rpath}}",val, {}, true, true)
					.then(res => {
						this.query()
					})
      });
		}
		{{- end}}
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
  .page-pagination{padding: 10px 15px;text-align: right;}
</style>
`
