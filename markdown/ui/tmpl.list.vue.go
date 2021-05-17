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
{{- $btn:=.BtnInfo -}}
<template>
	<div class="panel panel-default">
    <!-- query start -->
		<div class="panel-body" id="panel-body">
			<el-form ref="form" size="small" :inline="true" class="form-inline pull-left">
			{{- range $i,$c:=$rows|query}}
				{{- if $c.Con|TA}}
				<el-form-item>
					<el-input size="small" maxlength="{{$c.Len}}" type="textarea" :rows="2" placeholder="请输入{{$c.Desc|shortName}}" v-model="queryData.{{$c.Name}}">
					</el-input>
				</el-form-item>
				{{- else if or ($c.Con|SL) ($c.Con|SLM) }}
				<el-form-item>
					<el-select size="small" v-model="queryData.{{$c.Name}}"  clearable filterable class="input-cos" placeholder="请选择{{$c.Desc|shortName}}"
					{{- if (qDicPName $c.Con $tb) }} @change="handleChooseTool()"{{$choose = true}}{{end}} 
					{{- if (qDicCName $c.Name $tb) }} @change="set{{(qDicCName $c.Name $tb)|upperName}}(queryData.{{$c.Name}})" {{- end}}>
						<el-option value="" label="全部"></el-option>
						<el-option v-for="(item, index) in {{$c.Name|lowerName}}" :key="index" :value="item.value" :label="item.name"></el-option>
					</el-select>
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
				{{- else}}
				<el-form-item>
					<el-input clearable size="small" prefix-icon="el-icon-search" v-model="queryData.{{$c.Name}}" placeholder="请输入{{$c.Desc|shortName}}">
					</el-input>
				</el-form-item>
				{{- end}}
			{{end}}
				{{- if gt ($rows|query|len) 0}}
				<el-form-item>
					<el-button  type="primary" @click="queryDatas" size="small">查询</el-button>
				</el-form-item>
				{{end}}
				{{- if gt ($rows|create|len) 0}}
				<el-form-item>
					<el-button type="success" size="small" @click="showAdd">添加</el-button>
				</el-form-item>
				{{end}}
				{{- if gt ($rows|export|len) 0}}
				<el-form-item>
					<el-button type="success" @click="exportExcl" size="small">导出excel</el-button>
				</el-form-item>
				{{- end}}
				{{- if gt ($tb.DownloadInfo.Title|len) 0}}
				<el-form-item>
					<el-button type="text" @click="download" size="small" style="font-size:14px">下载模版</el-button>
				</el-form-item>
				{{- end}}
				{{- if not $tb.SelectInfo.IsEmpty }}
				<el-form-item>
					<el-button size="small" @click="toggleSelection()">批量操作</el-button>
				</el-form-item>
				{{- end}}
			</el-form>
		</div>
    <!-- query end -->

    <!-- list start-->
		<el-scrollbar style="height:100%">
			<el-table :data="dataList.items" stripe style="width: 100%" size="small" :height="maxHeight" {{if gt ($sort|len) 0}}@sort-change="sort"{{end}}
			{{- if not $tb.SelectInfo.IsEmpty }}  @selection-change="handleSelectionChange" {{end}}>
			  {{if not $tb.SelectInfo.IsEmpty }}<el-table-column type="selection" width="55"> </el-table-column>{{end}}
				{{if gt $tb.ELTableIndex 0}}<el-table-column type="index" fixed	:index="indexMethod" label="序号"></el-table-column>{{end}}
				{{- range $i,$c:=$rows|list}}
				<el-table-column {{if $c.Con|FIXED}}fixed{{end}} {{if $c.Con|SORT}}sortable="custom"{{end}} prop="{{$c.Name}}" label="{{$c.Desc|shortName}}" align="center">
				{{- if or ($c.Con|SL) ($c.Con|SLM)  ($c.Con|CB) ($c.Con|RD) ($c.Con|leCon)}}
					<template slot-scope="scope">
						<span {{if ($c.Con|CC)}}:class="scope.row.{{$c.Name}}|fltrTextColor"{{end}}>{{"{{scope.row."}}{{$c.Name}} | fltrEnum("{{or (dicName $c.Con ($c.Con|leCon) $tb) ($c.Name|lower)}}")}}</span>
					</template>
				{{- else if and ($c.Type|isString) (or (gt $c.Len $len) (eq $c.Len 0) )}}
					<template slot-scope="scope">
						<el-tooltip class="item" v-if="scope.row.{{$c.Name}} && scope.row.{{$c.Name}}.length > {{or ($c.Con|lfCon) "20"}}" effect="dark" placement="top">
							<div slot="content" style="width: 110px">{{"{{scope.row."}}{{$c.Name}}}}</div>
							<span>{{"{{scope.row."}}{{$c.Name}} | fltrSubstr({{or ($c.Con|lfCon) "20"}}) }}</span>
						</el-tooltip>
						<span v-else>{{"{{scope.row."}}{{$c.Name}} | fltrEmpty }}</span>
					</template>
				{{- else if and (or ($c.Type|isInt64) ($c.Type|isInt) ) (ne $c.Name ($pks|firstStr))}}
				<template slot-scope="scope">
					<span>{{"{{scope.row."}}{{$c.Name}} | fltrNumberFormat({{or ($c.Con|lfCon) "0"}})}}</span>
				</template>
				{{- else if $c.Type|isDecimal }}
				<template slot-scope="scope">
					<span>{{"{{scope.row."}}{{$c.Name}} | fltrNumberFormat({{or ($c.Con|lfCon) "2"}})}}</span>
				</template>
				{{- else if $c.Type|isTime }}
				<template slot-scope="scope">
					<div>{{"{{scope.row."}}{{$c.Name}} | fltrDate("{{ or (dateFormat $c.Con ($c.Con|lfCon)) "yyyy-MM-dd"}}") }}</div>
				</template>
				{{- else}}
				<template slot-scope="scope">
					<span>{{"{{scope.row."}}{{$c.Name}} | fltrEmpty }}</span>
				</template>
				{{end}}
				</el-table-column>
				{{- end}}
				<el-table-column  label="操作" align="center">
					<template slot-scope="scope">
						{{- if gt ($rows|update|len) 0}}
						<el-button type="text" size="mini" @click="showEdit(scope.row)">编辑</el-button>
						{{- end}}
						{{- if gt ($rows|detail|len) 0}}
						<el-button type="text" size="mini" @click="showDetail(scope.row)">详情</el-button>
						{{- end}}
						{{- if gt ($rows|delete|len) 0}}
						<el-button type="text" size="mini" @click="del(scope.row)">删除</el-button>
						{{- end}}

						{{- range $i,$c:=$tb.ComponentsInfo}}
						<el-button type="text" size="mini" @click="show{{$c.Name|upperName}}(scope.row)">form{{$i}}</el-button>
						{{- end}}

						{{- range $i,$c:= $btn }}
							{{- if gt ($c.VIF|len) 0}}
								{{- range $k,$v:= $c.VIF}}
									{{- if eq $k 0}}
						<el-button v-if="scope.row.{{(index $c.Rows 0).Name}} == {{$v.IfName}}" type="text" size="mini" @click="{{$c.Name}}(scope.row)">{{$v.IfDESC}}</el-button>
									{{- else if lt $k ($c.VIF|maxIndex) }}		
						<el-button v-else-if="scope.row.{{(index $c.Rows 0).Name}} == {{$v.IfName}}" type="text" size="mini" @click="{{$c.Name}}(scope.row)">{{$v.IfDESC}}</el-button>
									{{- else}}
						<el-button v-else type="text" size="mini" @click="{{$c.Name}}(scope.row)">{{$v.IfDESC}}</el-button>
									{{- end}}
								{{- end}}
							{{- else}}	
						<el-button type="text" size="mini" @click="show{{$c.Name}}(scope.row)">{{$c.DESC}}</el-button>
							{{- end }}
						{{- end}}
					</template>
				</el-table-column>
			</el-table>
		</el-scrollbar>
		<!-- list end-->

		{{- range $i,$c:= $btn }}
		{{- if eq ($c.VIF|len) 0}}
		<!-- {{$c.Name|upperName}} Form -->
		<{{$c.Name|upperName}} ref="{{$c.Name|upperName}}" :refresh="query"></{{$c.Name|upperName}}>
		<!--{{$c.Name|upperName}} Form -->
		{{- end}}
		{{- end}}

		{{if gt ($rows|create|len) 0 -}}
		<!-- Add Form -->
		<Add ref="Add" :refresh="query"></Add>
		<!--Add Form -->
		{{- end}}

		{{if gt ($rows|update|len) 0 -}}
		<!-- edit Form start-->
		<Edit ref="Edit" :refresh="query"></Edit>
		<!-- edit Form end-->
		{{- end}}

		{{- range $i,$c:=$tb.ComponentsInfo}}

		<!-- {{$c.Name|upperName}} Form -->
		<{{$c.Name|upperName}} ref="{{$c.Name|upperName}}" :refresh="query"></{{$c.Name|upperName}}>
		<!--{{$c.Name|upperName}} Form -->
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
{{- if gt ($rows|create|len) 0}}
import Add from "./{{.Name|rmhd|l2d}}.add"
{{- end}}
{{- if gt ($rows|update|len) 0}}
import Edit from "./{{.Name|rmhd|l2d}}.edit"
{{- end}}
{{- range $i,$c:=$tb.ComponentsInfo}}
import {{$c.Name|upperName}} from "{{$c.Path}}"
{{- end}}
{{- range $i,$c:= $btn }}
{{- if eq ($c.VIF|len) 0}}
import {{$c.Name|upperName}} from "./{{$name|rmhd|l2d}}.{{$c.Name}}"
{{- end}}
{{- end}}
export default {
	name: "{{$name|rmhd|varName}}",
  components: {
		{{- if gt ($rows|create|len) 0}}
		Add,
		{{- end}}
		{{- if gt ($rows|update|len) 0}}
		Edit,
		{{- end}}
		{{- range $i,$c:=$tb.ComponentsInfo}}
		{{$c.Name|upperName}},
		{{- end}}
		{{- range $i,$c:= $btn }}
		{{- if eq ($c.VIF|len) 0 -}},
		{{$c.Name|upperName}}
		{{- end}}
		{{- end}}
  },
  data () {
		return {
			paging: {ps: 10, pi: 1,total:0,sizes:[5, 10, 20, 50]},
			editData:{},                //编辑数据对象
			addData:{},                 //添加数据对象 
      queryData:{},               //查询数据对象 
			{{- range $i,$c:=$rows|query -}}
			{{if or ($c.Con|SL) ($c.Con|SLM) ($c.Con|RD) }}
			{{$c.Name|lowerName}}: {{if (qDicPName $c.Con $tb) }}[]{{else}}this.$enum.get("{{or (dicName $c.Con ($c.Con|qeCon) $tb) ($c.Name|lower)}}"){{end}},
			{{- else if $c.Con|CB }}
			{{$c.Name|lowerName}}: {{if (qDicPName $c.Con $tb) }}[]{{else}}this.$enum.get("{{or (dicName $c.Con ($c.Con|qeCon) $tb) ($c.Name|lower)}}"){{end}},
			{{$c.Name|lowerName}}Array: [],
			{{- end}}
			{{- if or ($c.Con|DTIME) ($c.Con|DATE) ($c.Type|isTime) }}
			{{$c.Name|lowerName}}: this.$utility.dateFormat(new Date(),"{{dateFormatDef $c.Con ($c.Con|qfCon)}}"),{{end}}
      {{- end}}
			{{- if gt ($sort|len) 0}}
			order: "{{range $i,$c:=$sort|sortSort}}t.{{$c.Name}} {{or ($c.Con|sortCon) "desc"}}{{if lt $i ($sort|maxIndex)}}, {{end}}{{end}}",
			{{- end}}
			{{- if not $tb.SelectInfo.IsEmpty }}
			multipleSelection: [],
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
      this.query()
		},
		{{- if gt ($sort|len) 0}}
		sort(column) {
      if (column.order == "ascending") {
        this.order ="t." +  column.prop + " " + "asc"
      } else if (column.order == "descending") {
        this.order ="t." +  column.prop + " " + "desc"
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
		{{- range $i,$c:=$rows|query -}}
		{{if (qDicPName $c.Con $tb)  }}
		set{{$c.Name|upperName}}(pid){
			this.queryData.{{$c.Name}} = ""
			this.{{$c.Name|lowerName}}=this.$enum.get("{{or (dicName $c.Con ($c.Con|qeCon) $tb) ($c.Name|lower)}}",pid)
		},
		{{- end}}
		{{- end }}
    /**查询数据并赋值*/
		queryDatas() {
      this.paging.pi = 1
      this.query()
    },
    query(){
      this.queryData.pi = this.paging.pi
			this.queryData.ps = this.paging.ps
			{{- range $i,$c:=$rows|query -}}
			{{- if or ($c.Con|DTIME) ($c.Con|DATE) ($c.Type|isTime) }}
			this.queryData.{{$c.Name}} = this.$utility.dateFormat(this.{{$c.Name|lowerName}},"{{dateFormat $c.Con ($c.Con|qfCon)}}")
			{{- else if ($c.Con|CB) }}
			this.queryData.{{$c.Name}} = this.{{$c.Name|lowerName}}Array.toString()
			{{- end -}}
      {{- end}}
			{{- if gt ($sort|len) 0}}
			this.queryData.order_by = this.order
			{{- end}}
      let res = this.$http.xpost("/{{.Name|rmhd|rpath}}/query",this.$utility.delEmptyProperty(this.queryData))
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
		{{- if gt ($rows|create|len) 0}}
    showAdd(){
      this.$refs.Add.show();
		},
		{{- end}}
		{{- if gt ($rows|update|len) 0}}
    showEdit(val) {
      this.$refs.Edit.editData = val
      this.$refs.Edit.show();
		},
		{{- end}}
		{{- range $i,$c:=$tb.ComponentsInfo}}
		show{{$c.Name|upperName}}(val) {
      this.$refs.{{$c.Name|upperName}}.show();
		},
		{{- end}}
		{{- range $i,$c:= $btn }}
		{{- if eq ($c.VIF|len) 0}}
	  show{{$c.Name}}(val) {
      this.$refs.{{$c.Name|upperName}}.editData = val
      this.$refs.{{$c.Name|upperName}}.show();
		},
		{{- end}}
		{{- end}}

		{{- range $i,$c:= $btn }}
		{{- if gt ($c.VIF|len) 0}}
		{{$c.Name}}(val){
			var data = {
				{{- range $i,$c:=$c.Rows}}
				{{$c.Name}} : ""
				{{- end}}
			}
			{{- if $c.Confirm}}
      this.$confirm("{{$c.Confirm}}?", "提示", { confirmButtonText: "确定", cancelButtonText: "取消", type: "warning" })
        .then(() => {
			{{- end}}
					this.$http.post("/{{$tb.Name|rmhd|rpath}}/{{or $c.URL ($c.Name|lowerName)}}", data, {}, true, true)
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

		{{- if gt ($rows|export|len) 0}}
		exportExcl() {
			this.queryData.pi = this.paging.pi
			this.queryData.ps = this.paging.ps
			{{- range $i,$c:=$rows|query -}}
			{{- if or ($c.Con|DTIME) ($c.Con|DATE) ($c.Type|isTime) }}
			this.queryData.{{$c.Name}} = this.$utility.dateFormat(this.{{$c.Name|lowerName}},"{{dateFormat $c.Con ($c.Con|qfCon)}}")
			{{- else if ($c.Con|CB) }}
			this.queryData.{{$c.Name}} = this.{{$c.Name|lowerName}}Array.toString()
			{{- end -}}
      {{- end}}
			{{- if gt ($sort|len) 0}}
			this.queryData.order_by = this.order
			{{- end}}
      this.$http.post("/{{.Name|rmhd|rpath}}/query",this.$utility.delEmptyProperty(this.queryData))
        .then(res => {
          var header = [
					{{- range $i,$c:=$rows|export}}
					{ Key: "{{$c.Name}}",	Txt: "{{$c.Desc}}" },
					{{- end}}
          ];
          this.BuildExcel("{{$desc}}.xlsx", [header], res.items || [], {
						{{- range $i,$c:=$rows|list}}
						{{- if or ($c.Con|SL) ($c.Con|SLM)  ($c.Con|CB) ($c.Con|RD) ($c.Con|leCon)}}
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

		{{- if not $tb.SelectInfo.IsEmpty }}
		toggleSelection() {
      var data = []
      this.multipleSelection.forEach(row => {
        data.push(row.{{range $i,$c:=$pks}}{{$c}}{{end}})
      });
			this.$http.put("{{$tb.SelectInfo.URL}}", {data:data}, {}, true, true)
			.then(res => {			
				this.query()
			})
    },
    handleSelectionChange(val) {
      this.multipleSelection = val;
    },
		{{- end}}

		{{- if gt ($rows|delete|len) 0}}
    del(val){
			this.$confirm("此操作将永久删除该数据, 是否继续?", "提示", {confirmButtonText: "确定",  cancelButtonText: "取消", type: "warning"})
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
