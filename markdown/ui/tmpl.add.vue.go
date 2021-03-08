package ui

//TmplCreateVue 添加创建弹框页面
const TmplCreateVue = `
{{- $empty := "" -}}
{{- $rows := .Rows -}}
{{- $tb :=. -}}
{{- $choose:= false -}}
<template>
  <!-- Add Form -->
  <el-dialog title="添加{{.Desc}}" {{- if gt ($rows|create|len) 5}} width="65%" {{else}} width="25%" {{- end}} :visible.sync="dialogAddVisible">
    <el-form :model="addData" {{if gt ($rows|create|len) 5 -}}:inline="true"{{- end}} :rules="rules" ref="addForm" label-width="110px">
    	{{- range $i,$c:=$rows|create }}
      {{if $c.Con|TA -}}
			<el-form-item label="{{$c.Desc|shortName}}:" prop="{{$c.Name}}">
				<el-input size="medium" type="textarea" :rows="2" placeholder="请输入{{$c.Desc|shortName}}" v-model="addData.{{$c.Name}}">
        </el-input>
			</el-form-item>
			{{- else if $c.Con|RD }}
			<el-form-item  label="{{$c.Desc|shortName}}:" prop="{{$c.Name}}">
				<el-radio-group size="medium" v-model="addData.{{$c.Name}}" style="margin-left:5px">
        	<el-radio v-for="(item, index) in {{$c.Name|lowerName}}" :key="index" :label="item.value">{{"{{item.name}}"}}</el-radio>
				</el-radio-group>
			</el-form-item>
			{{- else if $c.Con|SL }}
			<el-form-item label="{{$c.Desc|shortName}}:" prop="{{$c.Name}}">
				<el-select size="medium" style="width: 100%;"	v-model="addData.{{$c.Name}}"	clearable filterable class="input-cos" placeholder="---请选择---"
				{{- if (qDicPName $c.Con $tb) }} @change="handleChooseTool()"{{$choose = true}}{{end}}
				{{- if (cDicCName $c.Name $tb) }} @change="set{{(cDicCName $c.Name $tb)|upperName}}(addData.{{$c.Name}})" {{- end}}>
					<el-option v-for="(item, index) in {{$c.Name|lowerName}}" :key="index" :value="item.value" :label="item.name"></el-option>
				</el-select>
			</el-form-item>
			{{- else if $c.Con|SLM }}
			<el-form-item label="{{$c.Desc|shortName}}:" prop="{{$c.Name}}">
				<el-select size="medium"  placeholder="---请选择---" clearable filterable v-model="{{$c.Name|lowerName}}Array" multiple style="width: 100%;">
					<el-option v-for="(item, index) in {{$c.Name|lowerName}}" :key="index" :value="item.value" :label="item.name" ></el-option>
				</el-select>
			</el-form-item>
			{{- else if $c.Con|CB }}
			<el-form-item label="{{$c.Desc|shortName}}:" prop="{{$c.Name}}"> 
				<el-checkbox-group size="medium" v-model="addData.{{$c.Name}}">
					<el-checkbox v-for="(item, index) in {{$c.Name|lowerName}}" :key="index" :value="item.value" :label="item.name"></el-checkbox>
				</el-checkbox-group>
			</el-form-item>
			{{- else if or ($c.Con|DTIME) ($c.Con|DATE) }}
			<el-form-item prop="{{$c.Name}}" label="{{$c.Desc|shortName}}:">
					<el-date-picker size="medium" class="input-cos"  v-model="addData.{{$c.Name}}" type="{{dateType $c.Con ($c.Con|ceCon)}}" value-format="{{dateFormat $c.Con ($c.Con|ceCon)}}"  placeholder="选择日期"></el-date-picker>
			</el-form-item>
      {{- else -}}
      <el-form-item label="{{$c.Desc|shortName}}:" prop="{{$c.Name}}">
				<el-input size="medium" {{if gt $c.Len 0}}maxlength="{{$c.Len}}"{{end}} 
				{{- if gt $c.DecimalLen 0}} oninput="if(isNaN(value)) { value = null } if(value.indexOf('.')>0){value=value.slice(0,value.indexOf('.')+{{$c.DecimalLen|add1}})}"{{end}}
				 clearable v-model="addData.{{$c.Name}}" placeholder="请输入{{$c.Desc|shortName}}">
				</el-input>
      </el-form-item>
      {{- end}}
      {{end}}
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button size="medium" @click="resetForm('addForm')">取 消</el-button>
      <el-button size="medium" type="success" @click="add('addForm')">确 定</el-button>
    </div>
  </el-dialog>
  <!--Add Form -->
</template>

<script>
export default {
	data() {
		return {
			addData: {},
			dialogAddVisible: false,
			{{- range $i,$c:=$rows|create -}}
			{{if or ($c.Con|SL) ($c.Con|CB) ($c.Con|RD) }}
			{{$c.Name|lowerName}}:{{if (cDicPName $c.Con $tb) }} []{{else}}this.$enum.get("{{(or (dicName $c.Con ($c.Con|ceCon) $tb) $c.Name)|lower}}"){{end}},
			{{- else if $c.Con|SLM }}
			{{$c.Name|lowerName}}: this.$enum.get("{{(or (dicName $c.Con ($c.Con|ceCon) $tb) $c.Name)|lower}}"),
			{{$c.Name|lowerName}}Array: [],
      {{- end}}
			{{- end}}
			rules: {                    //数据验证规则
				{{- range $i,$c:=$rows|create -}}
				{{if ne ($c|isNull) $empty}}
				{{$c.Name}}: [{ required: true, message: "请输入{{$c.Desc|shortName}}", trigger: "blur" }],
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
			this.dialogAddVisible = false;
			this.$refs[formName].resetFields();
		},
		show(){
			this.dialogAddVisible = true;
		},
		{{- if $choose}}
		handleChooseTool() {
      this.$forceUpdate()
    },{{end}}
		{{- range $i,$c:=$rows|create -}}
		{{if (cDicPName $c.Con $tb)  }}
		set{{$c.Name|upperName}}(pid){
			this.addData.{{$c.Name}} = ""
			this.{{$c.Name|lowerName}}=this.$enum.get("{{(or (dicName $c.Con ($c.Con|ueCon) $tb) $c.Name)|lower}}",pid)
		},
		{{- end}}
		{{- end }}
		add(formName) {
			{{- range $i,$c:=$rows|create -}}
			{{- if or ($c.Con|DTIME) ($c.Con|DATE) }}
			this.addData.{{$c.Name}} = this.$utility.dateFormat(this.addData.{{$c.Name}},"{{dateFormat $c.Con ($c.Con|ueCon)}}")
			{{- else if $c.Con|SLM }}
			this.addData.{{$c.Name}} = this.{{$c.Name|lowerName}}Array.toString()
			{{- end -}}
			{{- end}}
			this.$refs[formName].validate((valid) => {
				if (valid) {
					this.$http.post("/{{.Name|rmhd|rpath}}", this.addData, {}, true, true)
						.then(res => {
							this.$refs[formName].resetFields()
							this.dialogAddVisible = false
							this.refresh()
						})
						.catch(err => {
							this.$message({
								type: "error",
								message: err.response.data
							});
						})
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
