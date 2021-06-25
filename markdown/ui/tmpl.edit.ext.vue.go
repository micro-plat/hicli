package ui

//TmplEditVue 添加编辑弹框页面
const TmplEditExtVue = `
{{- $btnInfo := (index .BtnInfo .TempIndex) -}}
{{- $rows := $btnInfo.Rows -}}
{{- $empty := "" -}}
{{- $pks := .|pks -}}
{{- $choose:= false -}}
<template>
	<el-dialog title="编辑{{.Desc}}"{{if gt ($rows|len) 5}} width="720px" {{- else}} width="500px" {{- end}} @closed="closed" :visible.sync="dialogFormVisible">
		<el-form :model="editData" size="small" {{if gt ($rows|len) 5 -}}:inline="true"{{- end}} :rules="rules" ref="editForm" label-width="110px">
    	{{- range $i,$c:=$rows}}
      {{if $c.Con|TA -}}
			<el-form-item label="{{$c.Desc|shortName}}:" prop="{{$c.Name}}">
				<el-input {{if $c.Disable}}:disabled="true"{{end}} size="small" maxlength="{{or ($c.Con|cfCon) $c.Len}}" type="textarea" :rows="2" placeholder="请输入{{$c.Desc|shortName}}" v-model="editData.{{$c.Name}}">
        </el-input>
			</el-form-item>
			{{- else if $c.Con|RD }}
			<el-form-item  label="{{$c.Desc|shortName}}:" prop="{{$c.Name}}">
				<el-radio-group {{if $c.Disable}}:disabled="true"{{end}} size="small" v-model="editData.{{$c.Name}}" style="margin-left:5px">
        	<el-radio v-for="(item, index) in {{$c.Name|lowerName}}" :key="index" :label="item.value">{{"{{item.name}}"}}</el-radio>
				</el-radio-group>
			</el-form-item>
			{{- else if $c.Con|SL }}
			<el-form-item label="{{$c.Desc|shortName}}:" prop="{{$c.Name}}">
				<el-select {{if $c.Disable}}:disabled="true"{{end}} size="small" style="width: 100%;"	v-model="editData.{{$c.Name}}" clearable filterable class="input-cos" placeholder="---请选择---"
				 {{- if (uDicCName $c.Name $c.BelongTable) }} @change="set{{(uDicCName $c.Name $c.BelongTable)|upperName}}(editData.{{$c.Name}})"
				 {{- else if (uGroupCName $c.Name $c.BelongTable) }} @change="set{{$c.Name|upperName}}Group" 
				 {{- else if (uDicPName $c.Con $c.BelongTable) }} @change="handleChooseTool()"{{$choose = true}}
				 {{- else if (uGroupPName $c.Con $c.BelongTable) }} disabled @change="handleChooseTool()"{{$choose = true}}{{- end}}	>
					<el-option v-for="(item, index) in {{$c.Name|lowerName}}" :key="index" :value="item.value" :label="item.name"></el-option>
				</el-select>
			</el-form-item>
			{{- else if $c.Con|SLM }}
			<el-form-item label="{{$c.Desc|shortName}}:" prop="{{$c.Name}}">
				<el-select {{if $c.Disable}}:disabled="true"{{end}} size="small" placeholder="---请选择---" clearable filterable v-model="{{$c.Name|lowerName}}Array" multiple style="width: 100%;">
					<el-option v-for="(item, index) in {{$c.Name|lowerName}}" :key="index" :value="item.value" :label="item.name" ></el-option>
				</el-select>
			</el-form-item>
			{{- else if $c.Con|CB }}
			<el-form-item label="{{$c.Desc|shortName}}:" prop="{{$c.Name}}"> 
				<el-checkbox-group {{if $c.Disable}}:disabled="true"{{end}} size="small" v-model="{{$c.Name|lowerName}}Array">
					<el-checkbox v-for="(item, index) in {{$c.Name|lowerName}}" :key="index" :value="item.value" :label="item.value">{{"{{item.name}}"}}</el-checkbox>
				</el-checkbox-group>
			</el-form-item>
			{{- else if or ($c.Con|DTIME) ($c.Con|DATE) ($c.Type|isTime) }}
			<el-form-item prop="{{$c.Name}}" label="{{$c.Desc|shortName}}:">
					<el-date-picker {{if $c.Disable}}:disabled="true"{{end}} size="small" class="input-cos"  v-model="editData.{{$c.Name}}" type="{{dateType $c.Con ($c.Con|ueCon)}}" value-format="{{dateFormat $c.Con ($c.Con|ueCon)}}"  placeholder="选择日期"></el-date-picker>
			</el-form-item>
      {{- else -}}
      <el-form-item label="{{$c.Desc|shortName}}:" prop="{{$c.Name}}">
				<el-input {{if $c.Disable}}:disabled="true"{{end}} size="small" {{if gt $c.Len 0}}maxlength="{{$c.Len}}"{{end}} 
				{{- if gt $c.DecimalLen 0}} oninput="if(isNaN(value)) { value = null } if(value.indexOf('.')>0){value=value.slice(0,value.indexOf('.')+{{$c.DecimalLen|add1}})}"{{end}}
				clearable v-model="editData.{{$c.Name}}" placeholder="请输入{{$c.Desc|shortName}}">
				</el-input>
      </el-form-item>
      {{- end}}
      {{end}}
    </el-form>
		<div slot="footer" class="dialog-footer">
			<el-button size="small" @click="resetForm('editForm')">取 消</el-button>
			<el-button type="success" size="small" @click="edit('editForm')">确 定</el-button>
		</div>
	</el-dialog>
</template>

<script>
export default {
	data() {
		return {
			dialogFormVisible: false,    //编辑表单显示隐藏
			editData: {},                //编辑数据对象
      {{- range $i,$c:=$rows -}}
      {{if or ($c.Con|SL) ($c.Con|RD) }}
      {{$c.Name|lowerName}}:{{if (uDicPName $c.Con $c.BelongTable) }} []{{else}} this.$enum.get("{{or (dicName $c.Con ($c.Con|ueCon) $c.BelongTable) ($c.Name|lower)}}"){{end}},
			{{- else if $c.Con|SLM }}
			{{$c.Name|lowerName}}: this.$enum.get("{{or (dicName $c.Con ($c.Con|ueCon) $c.BelongTable) ($c.Name|lower)}}"),
			{{$c.Name|lowerName}}Array: [],
			{{- else if $c.Con|CB }}
			{{$c.Name|lowerName}}:{{if (uDicPName $c.Con $c.BelongTable) }} []{{else}}this.$enum.get("{{or (dicName $c.Con ($c.Con|ueCon) $c.BelongTable) ($c.Name|lower)}}"){{end}},
			{{$c.Name|lowerName}}Array: [],
			{{- end}}
      {{- end}}
			rules: {                    //数据验证规则
				{{- range $i,$c:=$rows -}}
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
		resetForm(formName) {
			this.dialogFormVisible = false;
			this.$refs[formName].resetFields();
		},
		{{- if $choose}}
		handleChooseTool() {
      this.$forceUpdate()
    },{{end}}
		show() {
			{{range $i,$c:=$pks}}var {{$c}} = this.editData.{{$c}}{{end}}
			this.editData = this.$http.xget("/{{.Name|rmhd|rpath}}/get{{$btnInfo.Name}}", { {{range $i,$c:=$pks}}{{$c}}: {{$c}}{{end}} })
			{{range $i,$c:=$pks}}this.editData.{{$c}} = {{$c}}{{end}}
			{{- range $i,$c:=$rows -}}
			{{if (uDicPName $c.Con $c.BelongTable)  }}
			this.{{$c.Name|lowerName}} = this.$enum.get("{{or (dicName $c.Con ($c.Con|ueCon) $c.BelongTable) ($c.Name|lower)}}",this.editData.{{uDicPName $c.Con $c.BelongTable}})
			{{- end}}
			{{- if or ($c.Con|SLM) ($c.Con|CB) }}
			this.{{$c.Name|lowerName}}Array = this.editData.{{$c.Name}}.split(",")
			{{- end -}}
      {{- end }}
			this.dialogFormVisible = true;
		},
		{{- range $i,$c:=$rows -}}
		{{if (uDicPName $c.Con $c.BelongTable)  }}
		set{{$c.Name|upperName}}(pid){
			this.editData.{{$c.Name}} = ""
			this.{{$c.Name|lowerName}}=this.$enum.get("{{or (dicName $c.Con ($c.Con|ueCon) $c.BelongTable) ($c.Name)|lower)}}",pid)
		},
		{{- end}}
		{{- if (uGroupCName $c.Name $c.BelongTable) }}
		set{{$c.Name|upperName}}Group(value){
			var obj = this.{{$c.Name|lowerName}}.find((item) => {
        return item.value === value
      })
			if (obj){
				{{- range $i,$c1:=(ugroup $c.Name $c.BelongTable)}}
				{{- if (uDicCName $c1.Name $c.BelongTable)  }}
				this.set{{(uDicCName $c1.Name $c.BelongTable)|upperName}}(obj.{{$c1.Name}})
				{{- end}}
				this.editData.{{$c1.Name}} = obj.{{$c1.Name}}
				{{- end}}
			}
		},
		{{- end}}
		{{- end }}
		edit(formName) {
			{{- range $i,$c:=$rows -}}
			{{- if or ($c.Con|DTIME) ($c.Con|DATE) ($c.Type|isTime)}}
			this.editData.{{$c.Name}} = this.$utility.dateFormat(this.editData.{{$c.Name}},"{{dateFormat $c.Con ($c.Con|ueCon)}}")
			{{- else if or ($c.Con|SLM) ($c.Con|CB) }}
			this.editData.{{$c.Name}} = this.{{$c.Name|lowerName}}Array.toString()
			{{- end -}}
			{{- end}}
			this.$refs[formName].validate((valid) => {
				if (valid) {
					{{- if $btnInfo.Confirm}}
					this.$confirm("{{$btnInfo.Confirm}}?", "提示", { confirmButtonText: "确定", cancelButtonText: "取消", type: "warning" })
						.then(() => {
					{{- end}}
					this.$http.post("/{{.Name|rmhd|rpath}}/{{$btnInfo.Name}}", this.editData, {}, true, true)
					.then(res => {			
						this.dialogFormVisible = false;
						this.refresh()
					})
					{{- if $btnInfo.Confirm}}
					});
					{{- end}}
				} else {
						console.log("error submit!!");
						return false;
				}
			});
		},
	}
}
</script>

<style scoped>
</style>
`
