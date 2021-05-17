package ui

const srcPublicEnvConfJson = `
{
    "copyright": {
        "company": "四川千行你我科技股份有限公司",
        "code": "蜀ICP备20003360号"
    },
    "system": {
        "systemName": "xxx系统",
        "themes": "bg-danger|bg-danger|bg-dark dark-danger",
        "logo": ""
    },
    "api": {
        "host": "http://localhost:8089",
        "enumURL": "/system/enums/query"
    },
    "menus": [
        {
            "name": "日常管理",
            "path": "-",
            "children": [
                {
                    "name": "交易管理",
                    "is_open": "1",
                    "icon": "fa fa-line-chart text-danger",
                    "path": "-",
                    "children": [
                        {
                            "name": "交易订单",
                            "icon": "fa fa-user-circle text-primary",
                            "path": "/order/info"
                        }
                    ]
                }
            ]
        }
    ],
    "sysList": [],
    "enums": [{"name":"全国","type":"province","value":"*","group":"fltr"},{"name":"全市","type":"city","value":"*","group":"fltr"}],
    "textColor": {},
    "bgColor": {}
}
`

const srcSSOPublicEnvConfJson = `
{
    "copyright": {
        "company": "四川千行你我科技股份有限公司",
        "code": "蜀ICP备20003360号"
    },
    "api": {
        "host": "http://localhost:8089",
        "enumURL": "/system/enums/query"
    }
}
`
