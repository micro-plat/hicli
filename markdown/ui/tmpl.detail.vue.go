package ui

const TmplDetail = `
{{- $len := 32 -}}
{{- $rows := .Rows|detail -}}
{{- $pks := .|pks -}}
{{- $tb :=. -}}
{{- $tabs := .TabTables -}}
<template>
  <div>
    <el-tabs v-model="tabName" type="border-card" @tab-click="handleClick">
      <el-tab-pane label="{{.Desc|shortName}}" name="{{.Name|rmhd|varName}}Detail">
        <div class="table-responsive">
          <table :date="info" class="table table-striped m-b-none">
            <tbody class="table-border">
            {{- range $i,$c:=$rows -}}
            {{- if eq 0 (mod $i 2)}}
              <tr>
                <td>
            {{- end}}                 
                  <el-col :span="6">
                    <div class="pull-right" style="margin-right: 10px">{{$c.Desc|shortName}}:</div>
                  </el-col>
            {{- if or ($c.Con|SL) ($c.Con|SLM) ($c.Con|RD) ($c.Con|CB) ($c.Con|reCon)}}
                  <el-col :span="6">
                    <div {{if ($c.Con|CC)}}:class="info.{{$c.Name}}|fltrTextColor"{{end}}>{{"{{ info."}}{{$c.Name}} | fltrEnum("{{(or (dicName $c.Con ($c.Con|reCon) $tb) $c.Name)|lower}}") }}</div>
                  </el-col>
            {{- else if and ($c.Type|isString) (or (gt $c.Len $len) (eq $c.Len 0) )}}
                  <el-col :span="6">
                    <el-tooltip class="item" v-if="info.{{$c.Name}} && info.{{$c.Name}}.length > {{or ($c.Con|rfCon) "50"}}" effect="dark" placement="top">
                      <div slot="content" style="width: 110px">{{"{{info."}}{{$c.Name}}}}</div>
                      <div >{{"{{ info."}}{{$c.Name}} | fltrSubstr({{or ($c.Con|rfCon) "50"}}) }}</div>
                    </el-tooltip>
                    <div v-else>{{"{{ info."}}{{$c.Name}} | fltrEmpty }}</div>
                  </el-col>
          	{{- else if and (or ($c.Type|isInt64) ($c.Type|isInt)) (ne $c.Name ($pks|firstStr)) }}
                  <el-col :span="6">
                    <div>{{"{{ info."}}{{$c.Name}} |  fltrNumberFormat({{or ($c.Con|rfCon) "0"}})}}</div>
                  </el-col>
            {{- else if $c.Type|isDecimal }}
                  <el-col :span="6">
                    <div>{{"{{ info."}}{{$c.Name}} |  fltrNumberFormat({{or ($c.Con|rfCon) "2"}})}}</div>
                  </el-col>
            {{- else if $c.Type|isTime }}
                  <el-col :span="6">
                    <div>{{"{{ info."}}{{$c.Name}} | fltrDate("{{ or (dateFormat $c.Con ($c.Con|rfCon)) "yyyy-MM-dd"}}") }}</div>
                  </el-col>
            {{- else}}
                  <el-col :span="6">
                    <div>{{"{{ info."}}{{$c.Name}} | fltrEmpty }}</div>
                  </el-col>
            {{- end}}
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
      {{- range $index,$tab:=$tabs }}
      <el-tab-pane label="{{$tab.Desc|shortName}}" name="{{$tab.Name|rmhd|varName}}Detail">
      <div class="table-responsive">
        <table :date="{{$tab.Name|rmhd|lowerName}}Info" class="table table-striped m-b-none">
          <tbody class="table-border">
          {{- range $i,$c:=$tab.Rows|detail -}}
          {{- if eq 0 (mod $i 2)}}
            <tr>
              <td>
          {{- end}}                 
                <el-col :span="6">
                  <div class="pull-right" style="margin-right: 10px">{{$c.Desc|shortName}}:</div>
                </el-col>
          {{- if or ($c.Con|SL) ($c.Con|SLM) ($c.Con|RD) ($c.Con|CB) ($c.Con|reCon)}}
                <el-col :span="6">
                  <div {{if ($c.Con|CC)}}:class="{{$tab.Name|rmhd|lowerName}}Info.{{$c.Name}}|fltrTextColor"{{end}}>{{"{{ "}}{{$tab.Name|rmhd|lowerName}}Info.{{$c.Name}} | fltrEnum("{{(or (dicName $c.Con ($c.Con|reCon) $tb) $c.Name)|lower}}") }}</div>
                </el-col>
          {{- else if and ($c.Type|isString) (or (gt $c.Len $len) (eq $c.Len 0) )}}
                <el-col :span="6">
                  <el-tooltip class="item" v-if="{{$tab.Name|rmhd|lowerName}}Info.{{$c.Name}} && {{$tab.Name|rmhd|lowerName}}Info.{{$c.Name}}.length > {{or ($c.Con|rfCon) "50"}}" effect="dark" placement="top">
                    <div slot="content" style="width: 110px">{{"{{ "}}{{$tab.Name|rmhd|lowerName}}Info.{{$c.Name}}}}</div>
                    <div >{{"{{ "}}{{$tab.Name|rmhd|lowerName}}Info.{{$c.Name}} | fltrSubstr({{or ($c.Con|rfCon) "50"}}) }}</div>
                  </el-tooltip>
                  <div v-else>{{"{{ "}}{{$tab.Name|rmhd|lowerName}}Info.{{$c.Name}} | fltrEmpty }}</div>
                </el-col>
          {{- else if and (or ($c.Type|isInt64) ($c.Type|isInt)) (ne $c.Name ($pks|firstStr)) }}
                <el-col :span="6">
                  <div>{{"{{ "}}{{$tab.Name|rmhd|lowerName}}Info.{{$c.Name}} |  fltrNumberFormat({{or ($c.Con|rfCon) "0"}})}}</div>
                </el-col>
          {{- else if $c.Type|isDecimal }}
                <el-col :span="6">
                  <div>{{"{{ "}}{{$tab.Name|rmhd|lowerName}}Info.{{$c.Name}} |  fltrNumberFormat({{or ($c.Con|rfCon) "2"}})}}</div>
                </el-col>
          {{- else if $c.Type|isTime }}
                <el-col :span="6">
                  <div>{{"{{ "}}{{$tab.Name|rmhd|lowerName}}Info.{{$c.Name}} | fltrDate("{{ or (dateFormat $c.Con ($c.Con|rfCon)) "yyyy-MM-dd"}}") }}</div>
                </el-col>
          {{- else}}
                <el-col :span="6">
                  <div>{{"{{ "}}{{$tab.Name|rmhd|lowerName}}Info.{{$c.Name}} | fltrEmpty }}</div>
                </el-col>
          {{- end}}
          {{- if and (eq (mod $i 2) 1) (ne ($tab.Rows|detail|maxIndex) $i) }}
              </td>
            </tr>
          {{- end}}
          {{- if eq ($tab.Rows|detail|maxIndex) $i }}
              </td>
            </tr>
          {{- end -}}
          {{- end}}            
          </tbody>
        </table>
      </div>
      </el-tab-pane>
      {{- end}}
	  </el-tabs>
	</div>
</template>

<script>
export default {
  data(){
    return {
      tabName: "{{.Name|rmhd|varName}}Detail",
      info: {},
      {{- range $index,$tab:=$tabs }}
      {{$tab.Name|rmhd|lowerName}}Info:{},
      {{- end}}
    }
  },
  mounted() {
    this.init();
  },
  created(){
  },
  methods: {
    init(){
      this.queryData()
    },
    queryData() {
      this.info = this.$http.xget("/{{.Name|rmhd|rpath}}",this.$route.query)
    },
    {{- range $index,$tab:=$tabs }}
    query{{$tab.Name|rmhd|lowerName}}Data() {
      this.{{$tab.Name|rmhd|lowerName}}Info = this.$http.xget("/{{$tab.Name|rmhd|rpath}}",{ {{$tab|pks|firstStr}}: this.info.{{$tab|pks|firstStr}} })
    },
    {{- end }}
    handleClick(tab) {
      switch (tab.name) {
      case "{{.Name|rmhd|varName}}Detail":
        this.queryData();
        break;
      {{- range $index,$tab:=$tabs }}
      case "{{$tab.Name|rmhd|varName}}Detail":
        this.query{{$tab.Name|rmhd|lowerName}}Data();
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
</script>`
