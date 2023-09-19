package ding2

import (
	"fmt"

	// "github.com/beego/beego/v2/client/httplib"
	dingtalk "github.com/icepy/go-dingtalk/src"
)

type PackData struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type Token struct {
	PackData
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// 打卡记录
type WorkRecord struct {
	PackData
	HasMore      bool `json:"hasMore"`
	Recordresult []struct {
		CheckType      string `json:"checkType"`
		CorpId         string `json:"corpId"`
		LocationResult string `json:"locationResult"`
		BaseCheckTime  int64  `json:"baseCheckTime"`
		GroupId        int    `json:"groupId"`
		TimeResult     string `json:"timeResult"`
		UserId         string `json:"userId"`
		RecordId       int64  `json:"recordId"`
		WorkDate       int64  `json:"workDate"`
		SourceType     string `json:"sourceType"`
		UserCheckTime  int64  `json:"userCheckTime"`
		PlanId         int64  `json:"planId"`
		Id             int64  `json:"id"`
	} `json:"recordresult"`
}

type Department struct {
	PackData
	SubDeptIDList []int `json:"sub_dept_id_list"`
}

type Ding struct {
	Client *dingtalk.DingTalkClient
}

func NewClient(corp_id string, corp_secret string) *Ding {
	config := &dingtalk.DTConfig{
		CorpID:     corp_id,
		CorpSecret: corp_secret,
	}
	return &Ding{Client: dingtalk.NewDingTalkCompanyClient(config)}
}

func (d *Ding) GetAccessToken() string {
	d.Client.RefreshCompanyAccessToken()
	return d.Client.AccessToken
}

// func (d *Ding) GetDepartment() (deptIds []int, err error) {
// 	var root Department
// 	req := httplib.Get("https://oapi.dingtalk.com/department/list_ids")
// 	req.Param("access_token", d.GetAccessToken())
// 	req.Param("fetch_child", "true")
// 	req.Param("id", "1")
// 	err = req.ToJSON(&root)
// 	if err != nil {
// 		return
// 	}
// 	for i := 0; i < len(root.SubDeptIDList); i++ {
// 		var dept Department
// 		req := httplib.Get("https://oapi.dingtalk.com/department/list_ids")
// 		req.Param("access_token", d.GetAccessToken())
// 		req.Param("id", fmt.Sprintf("%d", root.SubDeptIDList[i]))
// 		err = req.ToJSON(&dept)
// 		if err != nil {
// 			return
// 		}
// 		deptIds = append(deptIds, dept.SubDeptIDList...)
// 	}
// 	deptIds = append(deptIds, root.SubDeptIDList...) //父级
// 	deptIds = append(deptIds, 1)                     //根级
// 	//父级
// 	fmt.Println("=======deptIds=========")
// 	fmt.Println(deptIds)
// 	return
// }
