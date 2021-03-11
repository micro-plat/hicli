package ui

//TmplEditVue 添加编辑弹框页面
const TmplEditVue = `
{{- $rows := .Rows -}}
{{- $empty := "" -}}
{{- $tb :=. -}}
{{- $pks := .|pks -}}
{{- $choose:= false -}}
<template>
	<el-dialog title="编辑{{.Desc}}"{{if gt ($rows|len) 5}} width="65%" {{- else}} width="25%" {{- end}} @closed="closed" :visible.sync="dialogFormVisible">
		<el-form :model="editData" {{if gt ($rows|update|len) 5 -}}:inline="true"{{- end}} :rules="rules" ref="editForm" label-width="110px">
    	{{- range $i,$c:=$rows|update}}
      {{if $c.Con|TA -}}
			<el-form-item label="{{$c.Desc|shortName}}:" prop="{{$c.Name}}">
				<el-input size="medium" type="textarea" :rows="2" placeholder="请输入{{$c.Desc|shortName}}" v-model="editData.{{$c.Name}}">
        </el-input>
			</el-form-item>
			{{- else if $c.Con|RD }}
			<el-form-item  label="{{$c.Desc|shortName}}:" prop="{{$c.Name}}">
				<el-radio-group  size="medium" v-model="editData.{{$c.Name}}" style="margin-left:5px">
        	<el-radio v-for="(item, index) in {{$c.Name|lowerName}}" :key="index" :label="item.value">{{"{{item.name}}"}}</el-radio>
				</el-radio-group>
			</el-form-item>
			{{- else if $c.Con|SL }}
			<el-form-item label="{{$c.Desc|shortName}}:" prop="{{$c.Name}}">
				<el-select size="medium" style="width: 100%;"	v-model="editData.{{$c.Name}}" clearable filterable class="input-cos" placeholder="---请选择---"
				 {{- if (qDicPName $c.Con $tb) }} @change="handleChooseTool()"{{$choose = true}}{{end}}
				 {{- if (uDicCName $c.Name $tb) }} @change="set{{(uDicCName $c.Name $tb)|upperName}}(editData.{{$c.Name}})"	{{- end}}	>
					<el-option v-for="(item, index) in {{$c.Name|lowerName}}" :key="index" :value="item.value" :label="item.name"></el-option>
				</el-select>
			</el-form-item>
			{{- else if $c.Con|SLM }}
			<el-form-item label="{{$c.Desc|shortName}}:" prop="{{$c.Name}}">
				<el-select size="medium" placeholder="---请选择---" clearable filterable v-model="{{$c.Name|lowerName}}Array" multiple style="width: 100%;">
					<el-option v-for="(item, index) in {{$c.Name|lowerName}}" :key="index" :value="item.value" :label="item.name" ></el-option>
				</el-select>
			</el-form-item>
			{{- else if $c.Con|CB }}
			<el-form-item label="{{$c.Desc|shortName}}:" prop="{{$c.Name}}"> 
				<el-checkbox-group size="medium" v-model="{{$c.Name|lowerName}}Array">
					<el-checkbox v-for="(item, index) in channelNo" :key="index" :value="item.value" :label="item.value">{{"{{item.name}}"}}</el-checkbox>
				</el-checkbox-group>
			</el-form-item>
			{{- else if or ($c.Con|DTIME) ($c.Con|DATE) }}
			<el-form-item prop="{{$c.Name}}" label="{{$c.Desc|shortName}}:">
					<el-date-picker size="medium" class="input-cos"  v-model="editData.{{$c.Name}}" type="{{dateType $c.Con ($c.Con|ueCon)}}" value-format="{{dateFormat $c.Con ($c.Con|ueCon)}}"  placeholder="选择日期"></el-date-picker>
			</el-form-item>
      {{- else -}}
      <el-form-item label="{{$c.Desc|shortName}}:" prop="{{$c.Name}}">
				<el-input size="medium" {{if gt $c.Len 0}}maxlength="{{$c.Len}}"{{end}} 
				{{- if gt $c.DecimalLen 0}} oninput="if(isNaN(value)) { value = null } if(value.indexOf('.')>0){value=value.slice(0,value.indexOf('.')+{{$c.DecimalLen|add1}})}"{{end}}
				clearable v-model="editData.{{$c.Name}}" placeholder="请输入{{$c.Desc|shortName}}">
				</el-input>
      </el-form-item>
      {{- end}}
      {{end}}
    </el-form>
		<div slot="footer" class="dialog-footer">
			<el-button size="medium" @click="dialogFormVisible = false">取 消</el-button>
			<el-button type="success" size="medium" @click="edit">确 定</el-button>
		</div>
	</el-dialog>
</template>

<script>
export default {
	data() {
		return {
			dialogFormVisible: false,    //编辑表单显示隐藏
			editData: {},                //编辑数据对象
      {{- range $i,$c:=$rows|update -}}
      {{if or ($c.Con|SL) ($c.Con|RD) }}
      {{$c.Name|lowerName}}:{{if (uDicPName $c.Con $tb) }} []{{else}} this.$enum.get("{{(or (dicName $c.Con ($c.Con|ueCon) $tb) $c.Name)|lower}}"){{end}},
			{{- else if $c.Con|SLM }}
			{{$c.Name|lowerName}}: this.$enum.get("{{(or (dicName $c.Con ($c.Con|ueCon) $tb) $c.Name)|lower}}"),
			{{$c.Name|lowerName}}Array: [],
			{{- else if $c.Con|CB }}
			{{$c.Name|lowerName}}:{{if (cDicPName $c.Con $tb) }} []{{else}}this.$enum.get("{{(or (dicName $c.Con ($c.Con|ceCon) $tb) $c.Name)|lower}}"){{end}},
			{{$c.Name|lowerName}}Array: [],
			{{- end}}
      {{- end}}
			rules: {                    //数据验证规则
				{{- range $i,$c:=$rows|update -}}
				{{if ne ($c|isNull) $empty}}
				{{$c.Name}}: [
					{ required: true, message: "请输入{{$c.Desc|shortName}}", trigger: "blur" }
				],
				{{- end}}
				{{- end}}
			},
		}
	},
	props: {
		refresh: {
			type: Function,
				default: () => {
				},
		}
	},
	created(){
	},
	methods: {
		closed() {
			this.refresh()
		},
		{{- if $choose}}
		handleChooseTool() {
      this.$forceUpdate()
    },{{end}}
		show() {
			{{range $i,$c:=$pks}}var {{$c}} = this.editData.{{$c}}{{end}}
			this.editData = this.$http.xget("/{{.Name|rmhd|rpath}}", { {{range $i,$c:=$pks}}{{$c}}: {{$c}}{{end}} })
			{{range $i,$c:=$pks}}this.editData.{{$c}} = {{$c}}{{end}}
			{{- range $i,$c:=$rows|update -}}
			{{- if or ($c.Con|SLM) ($c.Con|CB) }}
			this.{{$c.Name|lowerName}}Array = this.editData.{{$c.Name}}.split(",")
			{{- end -}}
      {{- end }}
			this.dialogFormVisible = true;
		},
		{{- range $i,$c:=$rows|update -}}
		{{if (uDicPName $c.Con $tb)  }}
		set{{$c.Name|upperName}}(pid){
			this.editData.{{$c.Name}} = ""
			this.{{$c.Name|lowerName}}=this.$enum.get("{{(or (dicName $c.Con ($c.Con|ueCon) $tb) $c.Name)|lower}}",pid)
		},
		{{- end}}
		{{- end }}
		edit() {
			{{- range $i,$c:=$rows|update -}}
			{{- if or ($c.Con|DTIME) ($c.Con|DATE) }}
			this.editData.{{$c.Name}} = this.$utility.dateFormat(this.editData.{{$c.Name}},"{{dateFormat $c.Con ($c.Con|ueCon)}}")
			{{- else if or ($c.Con|SLM) ($c.Con|CB) }}
			this.editData.{{$c.Name}} = this.{{$c.Name|lowerName}}Array.toString()
			{{- end -}}
			{{- end}}
			this.$http.put("/{{.Name|rmhd|rpath}}", this.editData, {}, true, true)
			.then(res => {			
				this.dialogFormVisible = false;
				this.refresh()
			})
		},
	}
}
</script>

<style scoped>
</style>
`
