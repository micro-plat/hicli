package tmpl

const EntityTmpl = `
//{{.Name|rmhd|varName}} {{.Desc}} 
type {{.Name|rmhd|varName}} struct {
			
	{{range $i,$c:=.Rows -}}
	//{{$c.Name|varName}} {{$c.Desc}}
	{{$c.Name|varName}} {{$c.Type|codeType}} {###}json:"{{$c.Name|lower}}" form:"{{$c.Name|lower}}"{###}

	{{end -}}	
}
`
