package ui

const TmplDetail = `
{{- $len := 32 -}}
{{- $rows := .Rows|detail -}}
{{- $pks := .|pks -}}
{{- $tb :=. -}}
{{- $tabs := .TabTables -}}
{{- $choose :=false -}}
{{- $panelbody :=false -}}
{{- $name:=.Name -}}
{{- $components:=false -}}
<template>
  <div>
    <div>
      <el-tabs v-model="tabName" type="border-card" @tab-click="handleClick">
        <el-tab-pane label="{{.Desc|shortName}}" name="{{.Name|rmhd|varName}}Detail">
          <div class="table-responsive">
            <table :date="info" class="table table-striped m-b-none">
              <tbody class="table-border">
              {{- range $i,$c:=$rows -}}
              {{- if eq 0 (mod $i 2) }}
                <tr>
                  <td>
              {{- end}}
                    <el-col :span="6">
                      <div class="pull-right" style="margin-right: 10px">{{$c.Desc|shortName}}:</div>
                    </el-col>
                    <el-col :span="6">
                    {{- if gt ($c.DetailBtnInfo|len) 0}}
                      <el-row type="flex">
                    {{- end}}
              {{- if or ($c.Con|SL) ($c.Con|SLM) ($c.Con|RD) ($c.Con|CB) ($c.Con|reCon) }}
                      <div {{- if ($c.Con|CC) }} :class="info.{{$c.Name}} | fltrTextColor"{{end}}>{{"{{ info."}}{{$c.Name}} | fltrEnum("{{or (dicName $c.Con ($c.Con|reCon) $tb) ($c.Name|lower) }}") }}</div>
              {{- else if and ($c.Type|isString) (or (gt $c.Len $len) (eq $c.Len 0) ) }}
                      <el-tooltip class="item" v-if="info.{{$c.Name}} && info.{{$c.Name}}.length > {{or ($c.Con|rfCon) "50"}}" effect="dark" placement="top">
                        <div slot="content" style="width: 110px">{{"{{ info."}}{{$c.Name}} }}</div>
                        <div>{{"{{ info."}}{{$c.Name}} | fltrSubstr({{or ($c.Con|rfCon) "50"}}) }}</div>
                      </el-tooltip>
                      <div v-else>{{"{{ info."}}{{$c.Name}} | fltrEmpty }}</div>
              {{- else if ($c.Con|fIsNofltr) }}
                      <div>{{"{{ info."}}{{$c.Name}} | fltrEmpty }}</div>
              {{- else if and (or ($c.Type|isInt64) ($c.Type|isInt)) (ne $c.Name ($pks|firstStr)) }}
                      <div>{{"{{ info."}}{{$c.Name}} | fltrNumberFormat({{or ($c.Con|rfCon) "0"}}) }}</div>
              {{- else if $c.Type|isDecimal }}
                      <div>{{"{{ info."}}{{$c.Name}} | fltrNumberFormat({{or ($c.Con|rfCon) "2"}}) }}</div>
              {{- else if $c.Type|isTime }}
                      <div>{{"{{ info."}}{{$c.Name}} | fltrDate("{{ or (dateFormat $c.Con ($c.Con|rfCon)) "yyyy-MM-dd HH:mm:ss"}}") }}</div>
              {{- else}}
                      <div>{{"{{ info."}}{{$c.Name}} | fltrEmpty }}</div>
              {{- end}}
              {{- range $i,$c:= $c.DetailBtnInfo }}
                        <el-button style="font-size: 14px; margin-left: 10px; padding: 0" type="text" size="mini" {{- if $c.Condition }} v-if="{{$c.Condition}}"{{end}} @click="{{$c.Name}}(info)">{{$c.DESC}}</el-button>
              {{- end}}
              {{- if gt ($c.DetailBtnInfo|len) 0}}
                      </el-row>
              {{- end}}
                    </el-col>
              {{- if and (eq (mod $i 2) 1) (ne ($rows|maxIndex) $i) }}
                  </td>
                </tr>
              {{- end}}
              {{- if eq ($rows|maxIndex) $i }}
                  </td>
                </tr>
              {{- end -}}
              {{- end}}
              </tbody>
            </table>
          </div>
        </el-tab-pane>
        {{- range $index,$tab:=$tabs}}
        <el-tab-pane label="{{$tab.Desc|shortName}}" name="{{$tab.Name|rmhd|varName}}Detail">
        {{- if (index $tab.TabInfo.TabTable $name) }}
          <div class="table-responsive">
            <table :date="{{$tab.Name|rmhd|lowerName}}Info" class="table table-striped m-b-none">
              <tbody class="table-border">
              {{- range $i,$c:=$tab.Rows|tabdetail -}}
              {{- if eq 0 (mod $i 2) }}
                <tr>
                  <td>
              {{- end}}
                    <el-col :span="6">
                      <div class="pull-right" style="margin-right: 10px">{{$c.Desc|shortName}}:</div>
                    </el-col>
                    <el-col :span="6">
                    {{- if gt ($c.DetailBtnInfo|len) 0}}
                      <el-row type="flex">
                    {{- end}}
              {{- if or ($c.Con|SL) ($c.Con|SLM) ($c.Con|RD) ($c.Con|CB) ($c.Con|reCon) }}
                      <div {{if ($c.Con|CC) }}:class="{{$tab.Name|rmhd|lowerName}}Info.{{$c.Name}} | fltrTextColor"{{end}}>{{"{{ "}}{{$tab.Name|rmhd|lowerName}}Info.{{$c.Name}} | fltrEnum("{{or (dicName $c.Con ($c.Con|reCon) $tb) ($c.Name|lower) }}") }}</div>
              {{- else if and ($c.Type|isString) (or (gt $c.Len $len) (eq $c.Len 0) ) }}
                      <el-tooltip class="item" v-if="{{$tab.Name|rmhd|lowerName}}Info.{{$c.Name}} && {{$tab.Name|rmhd|lowerName}}Info.{{$c.Name}}.length > {{or ($c.Con|rfCon) "50"}}" effect="dark" placement="top">
                        <div slot="content" style="width: 110px">{{"{{ "}}{{$tab.Name|rmhd|lowerName}}Info.{{$c.Name}} }}</div>
                        <div>{{"{{ "}}{{$tab.Name|rmhd|lowerName}}Info.{{$c.Name}} | fltrSubstr({{or ($c.Con|rfCon) "50"}}) }}</div>
                      </el-tooltip>
                      <div v-else>{{"{{ "}}{{$tab.Name|rmhd|lowerName}}Info.{{$c.Name}} | fltrEmpty }}</div>
              {{- else if ($c.Con|fIsNofltr) }}
                      <div>{{"{{ "}}{{$tab.Name|rmhd|lowerName}}Info.{{$c.Name}} | fltrEmpty }}</div>
              {{- else if and (or ($c.Type|isInt64) ($c.Type|isInt)) (ne $c.Name ($tab|pks|firstStr)) }}
                      <div>{{"{{ "}}{{$tab.Name|rmhd|lowerName}}Info.{{$c.Name}} | fltrNumberFormat({{or ($c.Con|rfCon) "0"}}) }}</div>
              {{- else if $c.Type|isDecimal }}
                      <div>{{"{{ "}}{{$tab.Name|rmhd|lowerName}}Info.{{$c.Name}} | fltrNumberFormat({{or ($c.Con|rfCon) "2"}}) }}</div>
              {{- else if $c.Type|isTime }}
                      <div>{{"{{ "}}{{$tab.Name|rmhd|lowerName}}Info.{{$c.Name}} | fltrDate("{{ or (dateFormat $c.Con ($c.Con|rfCon)) "yyyy-MM-dd HH:mm:ss"}}") }}</div>
              {{- else}}
                      <div>{{"{{ "}}{{$tab.Name|rmhd|lowerName}}Info.{{$c.Name}} | fltrEmpty }}</div>
              {{- end}}
              {{- range $i,$c:= $c.DetailBtnInfo }}
                        <el-button style="font-size: 14px; margin-left: 10px; padding: 0" type="text" size="mini" {{- if $c.Condition }} v-if="{{$c.Condition}}"{{end}} @click="{{$c.Name}}(info)">{{$c.DESC}}</el-button>
              {{- end}}
              {{- if gt ($c.DetailBtnInfo|len) 0}}
                      </el-row>
              {{- end}}              
                    </el-col>
              {{- if and (eq (mod $i 2) 1) (ne ($tab.Rows|tabdetail|maxIndex) $i) }}
                  </td>
                </tr>
              {{- end}}
              {{- if eq ($tab.Rows|tabdetail|maxIndex) $i }}
                  </td>
                </tr>
              {{- end -}}
              {{- end}}    
              </tbody>
            </table>
          </div>
        {{- else if (index $tab.TabInfo.TabTableList $name) }}{{$panelbody = true}}
          <el-scrollbar style="height: 100%" id="panel-body">
            <el-table :data="{{$tab.Name|rmhd|varName}}List.items" size="small" stripe style="width: 100%" :height="maxHeight">
              {{- if gt $tab.ELTableIndex 0}}
              <el-table-column type="index" fixed	:index="indexMethod"></el-table-column>
              {{- end}}
              {{- range $i,$c:=$tab.Rows|tablist}}
              <el-table-column {{- if $c.Con|FIXED}} fixed{{end}} {{- if $c.Con|SORT}} sortable{{end}} prop="{{$c.Name}}" label="{{$c.Desc|shortName}}" align="center">
                {{- if or ($c.Con|SL) ($c.Con|SLM)  ($c.Con|CB) ($c.Con|RD) ($c.Con|leCon) }}
                <template slot-scope="scope">
                  <span {{- if ($c.Con|CC) }} :class="scope.row.{{$c.Name}} | fltrTextColor"{{end}}>{{"{{ scope.row."}}{{$c.Name}} | fltrEnum("{{or (dicName $c.Con ($c.Con|leCon) $tab) ($c.Name|lower) }}") }}</span>
                </template>
                {{- else if and ($c.Type|isString) (or (gt $c.Len $len) (eq $c.Len 0) ) }}
                <template slot-scope="scope">
                  <el-tooltip class="item" v-if="scope.row.{{$c.Name}} && scope.row.{{$c.Name}}.length > {{or ($c.Con|lfCon) "20"}}" effect="dark" placement="top">
                    <div slot="content" style="width: 110px">{{"{{ scope.row."}}{{$c.Name}} }}</div>
                    <span>{{"{{ scope.row."}}{{$c.Name}} | fltrSubstr({{or ($c.Con|lfCon) "20"}}) }}</span>
                  </el-tooltip>
                  <span v-else>{{"{{ scope.row."}}{{$c.Name}} | fltrEmpty }}</span>
                </template>
                {{- else if ($c.Con|fIsNofltr) }}
                <template slot-scope="scope">
                  <span>{{"{{ scope.row."}}{{$c.Name}} | fltrEmpty }}</span>
                </template>
                {{- else if and (or ($c.Type|isInt64) ($c.Type|isInt) ) (ne $c.Name ($tab|pks|firstStr)) }}
                <template slot-scope="scope">
                  <span>{{"{{ scope.row."}}{{$c.Name}} | fltrNumberFormat({{or ($c.Con|lfCon) "0"}}) }}</span>
                </template>
                {{- else if $c.Type|isDecimal }}
                <template slot-scope="scope">
                  <span>{{"{{ scope.row."}}{{$c.Name}} | fltrNumberFormat({{or ($c.Con|lfCon) "2"}}) }}</span>
                </template>
                {{- else if $c.Type|isTime }}
                <template slot-scope="scope">
                  <div>{{"{{ scope.row."}}{{$c.Name}} | fltrDate("{{ or (dateFormat $c.Con ($c.Con|lfCon)) "yyyy-MM-dd HH:mm:ss"}}") }}</div>
                </template>
                {{- else}}
                <template slot-scope="scope">
                  <span>{{"{{ scope.row."}}{{$c.Name}} | fltrEmpty }}</span>
                </template>
                {{- end}}
              </el-table-column>
              {{- end}}
              {{- if gt ($tab.BtnInfo|len) 0}}
              <el-table-column label="操作" align="center">
                <template slot-scope="scope">
                  {{- range $i,$c:= $tab.BtnInfo }}
                  <el-button type="text" size="mini" {{- if $c.Condition }} v-if="{{$c.Condition}}"{{end}} @click="{{$c.Name}}(scope.row)">{{$c.DESC}}</el-button>
                  {{- end}}
                </template>
              </el-table-column>
              {{- end}}
            </el-table>
          </el-scrollbar>
        {{- end}}
        </el-tab-pane>
        {{- end }}
      </el-tabs>
    </div>
    {{- range $index,$tab:=$tabs -}}  
    {{- if (index $tab.TabInfo.TabTableList $name) }}
    <div class="page-pagination" v-show="tabName == '{{$tab.Name|rmhd|varName}}Detail'">
      <el-pagination
        @size-change="page{{$tab.Name|rmhd|varName}}SizeChange"
        @current-change="page{{$tab.Name|rmhd|varName}}IndexChange"
        :current-page="paging{{$tab.Name|rmhd|varName}}.pi"
        :page-size="paging{{$tab.Name|rmhd|varName}}.ps"
        :page-sizes="paging{{$tab.Name|rmhd|varName}}.sizes"
        layout="total, sizes, prev, pager, next, jumper"
        :total="{{$tab.Name|rmhd|varName}}List.count"
      >
      </el-pagination>
    </div>
    {{- end}}
    {{- end}}
  </div>
</template>

<script>
export default {
  {{- if $components}}
  components: {
  {{- end}}
  {{- if $components}}
  },
  {{- end}}
  data() {
    return {
      tabName: "{{.Name|rmhd|varName}}Detail",
      info: {},
      {{- range $index,$tab:=$tabs }}
      {{- if (index $tab.TabInfo.TabTable $name) }}
      {{$tab.Name|rmhd|lowerName}}Info:{},
      {{- else if (index $tab.TabInfo.TabTableList $name) }}
      paging{{$tab.Name|rmhd|varName}}: { ps: 10, pi: 1, total: 0, sizes: [5, 10, 20, 50] },
      {{$tab.Name|rmhd|varName}}List: { count: 0, items: [] }, //表单数据对象,
      query{{$tab.Name|rmhd|varName}}Params: {},  //查询数据对象
      {{- end}}
      {{- end}}
      maxHeight: 0
    }
  },
  {{- if not $panelbody}}
  watch: {
    "$route.fullPath"(val) {
      this.init();
    }
  },
  {{- end}}
  mounted() {
    {{- if $panelbody}}
    this.$nextTick(() => {
      this.maxHeight = this.$utility.getTableHeight("panel-body")
    })
    {{- end}}
    this.init();
  },
  created() {
  },
  methods: {
    init() {
      this.queryDetailData()
    },
    {{- if $choose}}
		handleChooseTool() {
      this.$forceUpdate()
    },{{end}}
    queryDetailData() {
      this.info = this.$http.xget("/{{.Name|rmhd|rpath}}", this.$route.query)
    },
    {{- range $index,$tab:=$tabs }}
    {{- if (index $tab.TabInfo.TabTable $name) }}
    query{{$tab.Name|rmhd|varName}}Data() {
      this.{{$tab.Name|rmhd|lowerName}}Info = this.$http.xget("/{{$tab.Name|rmhd|rpath}}/detail", { {{or (index $tab.TabInfo.TabTableProField $name) ($pks|firstStr) }}: this.info.{{or (index $tab.TabInfo.TabTablePreField $name) ($pks|firstStr) }} })
    },
    {{- else if (index $tab.TabInfo.TabTableList $name) }}
    /**查询数据并赋值*/
    query{{$tab.Name|rmhd|varName}}Datas() {
      this.paging{{$tab.Name|rmhd|varName}}.pi = 1
      this.query{{$tab.Name|rmhd|varName}}Data()
    },
    query{{$tab.Name|rmhd|varName}}Data() {
      this.query{{$tab.Name|rmhd|varName}}Params.pi = this.paging{{$tab.Name|rmhd|varName}}.pi
      this.query{{$tab.Name|rmhd|varName}}Params.ps = this.paging{{$tab.Name|rmhd|varName}}.ps
      var data = this.$utility.delEmptyProperty(this.query{{$tab.Name|rmhd|varName}}Params)
      data.{{or (index $tab.TabInfo.TabTableProField $name) ($pks|firstStr) }} = this.info.{{or (index $tab.TabInfo.TabTablePreField $name) ($pks|firstStr) }} || ""
      let res = this.$http.xget("/{{.Name|rmhd|rpath}}/querydetail", data)
      this.{{$tab.Name|rmhd|varName}}List.items = res.items || []
      this.{{$tab.Name|rmhd|varName}}List.count = res.count
    },
    /**改变页容量*/
    page{{$tab.Name|rmhd|varName}}SizeChange(val) {
      this.paging{{$tab.Name|rmhd|varName}}.ps = val
      this.query{{$tab.Name|rmhd|varName}}Data()
    },
    /**改变当前页码*/
    page{{$tab.Name|rmhd|varName}}IndexChange(val) {
      this.paging{{$tab.Name|rmhd|varName}}.pi = val
      this.query{{$tab.Name|rmhd|varName}}Data()
    },
    {{- end}}
    {{- end }}
    
    {{- range $index,$r:=$rows -}}  
    {{- range $i,$c:= $r.DetailBtnInfo }}
    {{$c.Name}}(val) {
      var data = {
        {{- range $i,$c:=$c.Rows}}
        {{$c.Name}}: val.{{$c.Name}},
				{{- end}}
      }
			{{- if $c.Confirm}}
      this.$confirm("{{$c.Confirm}}?", "提示", { confirmButtonText: "确定", cancelButtonText: "取消", type: "warning" })
        .then(() => {
			{{- end}}
          this.$http.post("/{{$tb.Name|rmhd|rpath}}/{{or $c.URL ($c.Name|lowerName) }}", data, {}, true, true)
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

    {{- range $index,$tab:=$tabs -}}  
    {{- if (index $tab.TabInfo.TabTableList $name) }}  
    {{- range $i,$c:= $tab.BtnInfo }}
    {{$c.Name}}(val) {
      var data = {
        {{- range $i,$c:=$c.Rows}}
        {{$c.Name}}: val.{{$c.Name}},
				{{- end}}
      }
			{{- if $c.Confirm}}
      this.$confirm("{{$c.Confirm}}?", "提示", { confirmButtonText: "确定", cancelButtonText: "取消", type: "warning" })
        .then(() => {
			{{- end}}
          this.$http.post("/{{$tb.Name|rmhd|rpath}}/{{or $c.URL ($c.Name|lowerName) }}", data, {}, true, true)
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
		{{- end}}
    handleClick(tab) {
      switch (tab.name) {
        case "{{.Name|rmhd|varName}}Detail":
          this.queryDetailData();
          break;
        {{- range $index,$tab:=$tabs }}
        case "{{$tab.Name|rmhd|varName}}Detail":
          this.query{{$tab.Name|rmhd|varName}}Data();
          break;
        {{- end}}
        default:
          this.$notify({
            title: "警告",
            message: "选项卡错误！"
          });
          return;
      }
    }
  },
}
</script>
<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.page-pagination {
  padding: 10px 15px;
  text-align: right;
}
</style>
`
