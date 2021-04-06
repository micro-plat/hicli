package ui

const srcUtilityExportJS = `import XLSX from "xlsx";

export function ExportTemplate(value, name){
    var sheet = XLSX.utils.aoa_to_sheet(value);
    openDownloadDialog(sheet2blob(sheet), name);
}

function openDownloadDialog(url, saveName) {
    if (typeof url == "object" && url instanceof Blob) {
        url = URL.createObjectURL(url);
    }
    var aLink = document.createElement("a");
    aLink.href = url;
    aLink.download = saveName || "";
    var event;
    if (window.MouseEvent) event = new MouseEvent("click");
    else {
        event = document.createEvent("MouseEvents");
        event.initMouseEvent(
            "click",
            true,
            false,
            window,
            0,
            0,
            0,
            0,
            0,
            false,
            false,
            false,
            false,
            0,
            null
        );
    }
    aLink.dispatchEvent(event);
}

function sheet2blob(sheet, sheetName) {
    sheetName = sheetName || "sheet1";
    var workbook = {
        SheetNames: [sheetName],
        Sheets: {}
    };
    workbook.Sheets[sheetName] = sheet;
    // 生成excel的配置项
    var wopts = {
        bookType: "xlsx", // 要生成的文件类型
        bookSST: false, // 是否生成Shared String Table，官方解释是，如果开启生成速度会下降，但在低版本IOS设备上有更好的兼容性
        type: "binary"
    };
    var wbout = XLSX.write(workbook, wopts);
    var blob = new Blob([s2ab(wbout)], {
        type: "application/octet-stream"
    });
    // 字符串转ArrayBuffer
    function s2ab(s) {
        var buf = new ArrayBuffer(s.length);
        var view = new Uint8Array(buf);
        for (var i = 0; i != s.length; ++i) view[i] = s.charCodeAt(i) & 0xff;
        return buf;
    }
    return blob;
}


export function BuildExcel(fileName,headers,rows,filter={}){
    
    var newfilter = {}
    for(var k in filter){
        var list = filter[k]
        var newVal = list 
        if(Array.isArray(list)){
            newVal = {}
            for(var i in list){
                newVal[list[i].value] = list[i].name
            }
        }
        newfilter[k] = newVal;
    }
    filter = newfilter
    var resultData = [];
 
    for(var h in  headers){
        var headerTxt  = [];
        var r = headers[h]
        for(var i in r){
            headerTxt.push(r[i].Txt)
        }
        resultData.push(headerTxt);    
    }
    var cols = headers.slice(headers.length-1,headers.length)[0];
    var key = "",val = "";
    for (var i in rows) {
        var row = []
        var item = rows[i];
        for (var j in cols) {  
            key = cols[j].Key 
            val = item[key]
            if(filter[key]){
                if(typeof filter[key] == 'function'){
                    val = filter[key](item[key])
                }else{
                    val = filter[key][item[key]]
                }                
            }
            row.push(val)
        }
        resultData.push(row)
    }
    
    ExportTemplate(resultData, fileName)
}
`
