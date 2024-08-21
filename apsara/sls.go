package apsara

import (
	"os"
	"reflect"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
)

// new sls client： default v1 client
func NewSlsCLient() sls.ClientInterface {
	var client sls.ClientInterface

	ep, ak, sk, err := GetAkSk("sls")
	if err != nil {
		os.Exit(1)
	}

	// fmt.Printf("ep: %v\nak: %v\nsk: %v\n", ep, ak, sk)

	client = sls.CreateNormalInterface(ep, ak, sk, "")

	return client
}

// 删除project
func DeleteProject(client sls.ClientInterface, project string) error {
	projectExist, _ := client.CheckProjectExist(project)
	if projectExist {
		err := client.DeleteProject(project)
		if err != nil {
			return err
		}
	}
	return nil
}

// 创建project
func CreateProject(client sls.ClientInterface, project string) error {
	_, err := client.CreateProject(project, "Support by kaitianpidi team!")
	if err != nil {
		e, ok := err.(*sls.Error)
		if ok && e.Code == "ProjectAlreadyExist" {
		} else {
			return err
		}
	} else {
		time.Sleep(time.Second * 1)
	}

	return nil
}

// Exist? project
func ExistProject(client sls.ClientInterface, project string) (bool, error) {
	return client.CheckProjectExist(project)
}

// 创建logstore
func CreateLogstore(client sls.ClientInterface, project, logstore string) error {
	err := client.CreateLogStore(project, logstore, 1, 2, true, 6)
	if err != nil {
		e, ok := err.(*sls.Error)
		if ok && e.Code == "LogStoreAlreadyExist" {
		} else {
			return err
		}
	} else {
		time.Sleep(time.Second * 1)
	}

	return nil
}

// 删除logstore
func DeleteLogstore(client sls.ClientInterface, project, logstore string) error {
	err := client.DeleteLogStore(project, logstore)
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 1)

	return nil
}

// Exist? logstore
func ExistLogstore(client sls.ClientInterface, project, logstore string) (bool, error) {
	return client.CheckLogstoreExist(project, logstore)
}

// 字段索引（适用于有明确字段结构的日志，可以对每个字段进行精确的查询和分析，适合结构化日志数据。
// 全文索引（对整个日志行进行搜索，适用于那些不易于拆分成独立字段的日志，或者对整个日志内容进行模糊匹配的场景。将整行日志视为一个长字符串!
// 创建CollectorIndex对象 (变长参数)，适用于所有字段的IndexKey配置都相同的情况
func NewVariableIndex(colNames ...string) sls.Index {
	index := sls.Index{
		Keys: make(map[string]sls.IndexKey),
	}

	for _, v := range colNames {
		index.Keys[v] = sls.IndexKey{
			Token:         []string{";"},
			CaseSensitive: false,
			DocValue:      true, // 开启统计（sql语法查询）
			Type:          "text",
		}
	}

	return index
}

// 字段索引（适用于有明确字段结构的日志，可以对每个字段进行精确的查询和分析，适合结构化日志数据。
// 全文索引（对整个日志行进行搜索，适用于那些不易于拆分成独立字段的日志，或者对整个日志内容进行模糊匹配的场景。将整行日志视为一个长字符串!
// 创建DetectorIndex对象（定长参数），适用于所有字段的IndexKey配置不相同的情况
func NewFixedIndex() sls.Index {
	return sls.Index{
		Keys: map[string]sls.IndexKey{
			"col_1": {
				Token:         []string{";"},
				CaseSensitive: false,
				DocValue:      true, // 开启统计（sql语法查询）
				Type:          "text",
			},
			"col_2": {
				Token:         []string{";"},
				CaseSensitive: false,
				DocValue:      true,
				Type:          "text",
			},
			"col_3": {
				Token:         []string{";"},
				CaseSensitive: false,
				DocValue:      true,
				Type:          "text",
			},
			"col_4": {
				Token:         []string{";"},
				CaseSensitive: false,
				DocValue:      true,
				Type:          "text",
			},
			"col_5": {
				Token:         []string{";"},
				CaseSensitive: false,
				DocValue:      true,
				Type:          "text",
			},
		},
	}
}

// 获取index配置
func GetIndex(client sls.ClientInterface, project, logstore string) (*sls.Index, error) {
	idx, err := client.GetIndex(project, logstore)
	if err != nil {
		return nil, err
	}

	return idx, nil
}

// 创建index实例
func CreateIndex(client sls.ClientInterface, project, logstore string, idx sls.Index) error {
	err := client.CreateIndex(project, logstore, idx)
	if err != nil {
		e, ok := err.(*sls.Error)
		if ok && e.Code == "IndexAlreadyExist" {
			return nil
		} else {
			return err
		}
	}

	return nil
}

// 删除index实例
func DeleteIndex(client sls.ClientInterface, project, logstore string) error {
	err := client.DeleteIndex(project, logstore)
	if err != nil {
		return err
	}

	return nil
}

// 更新index
func UpdateIndex(client sls.ClientInterface, project, logstore string, idx sls.Index) error {
	// get
	old_idx, err := client.GetIndex(project, logstore)
	if err != nil {
		return err
	}

	// 对比
	if reflect.DeepEqual(idx, old_idx) {
		return nil
	}

	// (idx != old_idx)才更新
	err = client.UpdateIndex(project, logstore, idx)
	if err != nil {
		return err
	}

	return nil
}
