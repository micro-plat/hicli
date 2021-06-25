package ui

//TmplEditVue 添加编辑弹框页面
const TmplEditVue = `
{{- $rows := .Rows -}}
{{- $empty := "" -}}
{{- $tb :=. -}}
{{- $pks := .|pks -}}
{{- $choose:= false -}}
<template>
	<el-dialog title="编辑{{.Desc}}"{{if gt ($rows|update|len) 5}} width="720px" {{- else}} width="500px" {{- end}} @closed="closed" :visible.sync="dialogFormVisible">
		<el-form :model="editData" size="small" {{if gt ($rows|update|len) 5 -}}:inline="true"{{- end}} :rules="rules" ref="editForm" label-width="110px">
    	{{- range $i,$c:=$rows|update}}
      {{if $c.Con|TA -}}
			<el-form-item label="{{$c.Desc|shortName}}:" prop="{{$c.Name}}">
				<el-input size="small" maxlength="{{or ($c.Con|cfCon) $c.Len}}" type="textarea" :rows="2" placeholder="请输入{{$c.Desc|shortName}}" v-model="editData.{{$c.Name}}">
        </el-input>
			</el-form-item>
			{{- else if $c.Con|RD }}
			<el-form-item  label="{{$c.Desc|shortName}}:" prop="{{$c.Name}}">
				<el-radio-group  size="small" v-model="editData.{{$c.Name}}" style="margin-left:5px">
        	<el-radio v-for="(item, index) in {{$c.Name|lowerName}}" :key="index" :label="item.value">{{"{{item.name}}"}}</el-radio>
				</el-radio-group>
			</el-form-item>
			{{- else if $c.Con|SL }}
			<el-form-item label="{{$c.Desc|shortName}}:" prop="{{$c.Name}}">
				<el-select size="small" style="width: 100%;"	v-model="editData.{{$c.Name}}" clearable filterable class="input-cos" placeholder="---请选择---"
					{{- if (uDicCName $c.Name $tb) }} @change="set{{(uDicCName $c.Name $tb)|upperName}}(editData.{{$c.Name}})"
					{{- else if (uGroupCName $c.Name $tb) }} @change="set{{$c.Name|upperName}}Group" 
					{{- else if or  (uDicPName $c.Con $tb) (uGroupPName $c.Con $tb) }} @change="handleChooseTool()"{{$choose = true}}{{- end}}
					{{- if (uGroupPName $c.Con $tb)}} disabled{{end}}	>
					<el-option v-for="(item, index) in {{$c.Name|lowerName}}" :key="index" :value="item.value" :label="item.name"></el-option>
				</el-select>
			</el-form-item>
			{{- else if $c.Con|SLM }}
			<el-form-item label="{{$c.Desc|shortName}}:" prop="{{$c.Name}}">
				<el-select size="small" placeholder="---请选择---" clearable filterable v-model="{{$c.Name|lowerName}}Array" multiple style="width: 100%;">
					<el-option v-for="(item, index) in {{$c.Name|lowerName}}" :key="index" :value="item.value" :label="item.name" ></el-option>
				</el-select>
			</el-form-item>
			{{- else if $c.Con|CB }}
			<el-form-item label="{{$c.Desc|shortName}}:" prop="{{$c.Name}}"> 
				<el-checkbox-group size="small" v-model="{{$c.Name|lowerName}}Array">
					<el-checkbox v-for="(item, index) in {{$c.Name|lowerName}}" :key="index" :value="item.value" :label="item.value">{{"{{item.name}}"}}</el-checkbox>
				</el-checkbox-group>
			</el-form-item>
			{{- else if or ($c.Con|DTIME) ($c.Con|DATE) ($c.Type|isTime) }}
			<el-form-item prop="{{$c.Name}}" label="{{$c.Desc|shortName}}:">
					<el-date-picker size="small" class="input-cos"  v-model="editData.{{$c.Name}}" type="{{dateType $c.Con ($c.Con|ueCon)}}" value-format="{{dateFormat $c.Con ($c.Con|ueCon)}}"  placeholder="选择日期"></el-date-picker>
			</el-form-item>
			{{- else if $c.Con|UP }}
			<el-form-item label="{{$c.Desc|shortName}}:" prop="{{$c.Name}}">
				<el-upload id="sbtn{{$i}}"
					class="upload-demo"
					ref="upload"
					:with-credentials="options{{$i}}.withCredentials"
					:accept="options{{$i}}.accept"
					:headers="options{{$i}}.headers"
					:action="options{{$i}}.target"
					:limit="1"
					:on-exceed="handleExceed{{$i}}"
					:file-list="fileList{{$i}}"
					:on-success="uploadSuccess{{$i}}"
					:on-error="onError{{$i}}"
					:before-upload="beforeUpload{{$i}}"
					:on-remove="onRemove{{$i}}"
				>
					<el-button size="small" type="primary">点击上传</el-button>
					<div slot="tip" class="el-upload__tip" style="margin-top: 0px">建议尺寸格式，大小在2M以内</div>
				</el-upload>
			</el-form-item>
      {{- else -}}
      <el-form-item label="{{$c.Desc|shortName}}:" prop="{{$c.Name}}">
				<el-input size="small" {{if gt $c.Len 0}}maxlength="{{$c.Len}}"{{end}} 
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
      {{- range $i,$c:=$rows|update -}}
      {{if or ($c.Con|SL) ($c.Con|RD) }}
      {{$c.Name|lowerName}}:{{if (uDicPName $c.Con $tb) }} []{{else}} this.$enum.get("{{or (dicName $c.Con ($c.Con|ueCon) $tb) ($c.Name|lower)}}"){{end}},
			{{- else if $c.Con|SLM }}
			{{$c.Name|lowerName}}: this.$enum.get("{{or (dicName $c.Con ($c.Con|ueCon) $tb) ($c.Name|lower)}}"),
			{{$c.Name|lowerName}}Array: [],
			{{- else if $c.Con|CB }}
			{{$c.Name|lowerName}}:{{if (uDicPName $c.Con $tb) }} []{{else}}this.$enum.get("{{or (dicName $c.Con ($c.Con|ueCon) $tb) ($c.Name|lower)}}"){{end}},
			{{$c.Name|lowerName}}Array: [],
			{{- else if $c.Con|UP }}
			fileList{{$i}}: [],
      nameList{{$i}}: [],
      options{{$i}}: {
        accept: "image/jpg,image/jpeg,image/png,image/gif,image/bmp,.doc,.docx,application/msword,application/vnd.openxmlformats-officedocument.wordprocessingml.document", //上传文件类型
        target: this.$env.conf.api.host + '/{{$.Name|rmhd|rpath}}/upload', //上传地址
        withCredentials: true, //携带cookies
        headers: { "authorization": window.localStorage.getItem("authorization") },//根据后端配置在cookie还是header
      },
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
			this.editData = this.$http.xget("/{{.Name|rmhd|rpath}}/getupdate", { {{range $i,$c:=$pks}}{{$c}}: {{$c}}{{end}} })
			{{range $i,$c:=$pks}}this.editData.{{$c}} = {{$c}}{{end}}
			{{- range $i,$c:=$rows|update -}}
			{{if (uDicPName $c.Con $tb)  }}
			this.{{$c.Name|lowerName}} = this.$enum.get("{{or (dicName $c.Con ($c.Con|ueCon) $tb) ($c.Name|lower)}}",this.editData.{{uDicPName $c.Con $tb}})
			{{- end}}
			{{- if or ($c.Con|SLM) ($c.Con|CB) }}
			this.{{$c.Name|lowerName}}Array = this.editData.{{$c.Name}}.split(",")
			{{- end -}}
      {{- end }}
			this.dialogFormVisible = true;
		},
		{{- range $i,$c:=$rows|update -}}
		{{- if (uGroupCName $c.Name $tb) }}
		set{{$c.Name|upperName}}Group(value){
			var obj = this.{{$c.Name|lowerName}}.find((item) => {
        return item.value === value
      })
			if (obj){
				{{- range $i,$c1:=(ugroup $c.Name $tb)}}
				{{- if (uDicCName $c1.Name $tb)  }}
				this.set{{(uDicCName $c1.Name $tb)|upperName}}(obj.{{$c1.Name}})
				{{- end}}
				{{- end}}
				{{- range $i,$c1:=(ugroup $c.Name $tb)}}
				this.editData.{{$c1.Name}} = obj.{{$c1.Name}}
				{{- end}}
			}
		},
		{{- end}}
		{{if (uDicPName $c.Con $tb)  }}
		set{{$c.Name|upperName}}(pid){
			this.editData.{{$c.Name}} = ""
			this.{{$c.Name|lowerName}}=this.$enum.get("{{or (dicName $c.Con ($c.Con|ueCon) $tb) ($c.Name|lower)}}",pid)
		},
		{{- else if $c.Con|UP }}
		handleExceed{{$i}}(files, fileList) {
      this.$message.warning("文件上传数量超出设置范围");
    },
    uploadSuccess{{$i}}(response, file, fileList) {
			let info = {
        original_name: file.name,
        store_name: response.file_name,
        uid: file.uid,
      }
      this.nameList{{$i}}.push(info)
    },
    onError{{$i}}(err, file, fileList) {
      this.$notify({
        title: "错误",
        message: "上传失败，请稍后再试",
        type: "error",
        offset: 50,
        duration: 2000
      });
    },
    beforeUpload{{$i}}(file) {
      const isLt2M = file.size / 1024 / 1024 < 2 //这里做文件大小限制
      if (!isLt2M) {
        this.$message({
          message: '上传文件大小不能超过 2MB!',
          type: 'warning'
        });
        return false
      }
      return isLt2M
    },
    onRemove{{$i}}(file, fileList) {
      this.nameList{{$i}}.forEach((item, idx, array) => {
        if (array[idx] != undefined && item.uid == file.uid) {
          this.nameList{{$i}}.splice(idx, 1)
          return
        }
      })
    },
		{{- end}}
		{{- end }}
		edit(formName) {
			{{- range $i,$c:=$rows|update -}}
			{{- if or ($c.Con|DTIME) ($c.Con|DATE) ($c.Type|isTime) }}
			this.editData.{{$c.Name}} = this.$utility.dateFormat(this.editData.{{$c.Name}},"{{dateFormat $c.Con ($c.Con|ueCon)}}")
			{{- else if or ($c.Con|SLM) ($c.Con|CB) }}
			this.editData.{{$c.Name}} = this.{{$c.Name|lowerName}}Array.toString()
			{{- else if $c.Con|UP }}
			var list{{$i}} =[]
			this.nameList{{$i}}.forEach((v, i) => {list{{$i}}.push(v.store_name)})
			this.editData.{{$c.Name}} = list{{$i}}.join(",")
			{{- end -}}
			{{- end}}
			this.$refs[formName].validate((valid) => {
				if (valid) {
					this.$http.put("/{{.Name|rmhd|rpath}}", this.editData, {}, true, true)
					.then(res => {			
						this.dialogFormVisible = false;
						this.refresh()
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

{{- range $i,$c:=$rows|update -}}
{{- if $c.Con|UP }}
<style>
#sbtn{{$i}} > div > input {
  display: none !important;
}
</style>
{{- end}}
{{- end}}
`
