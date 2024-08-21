package apsara

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

const (
	tianjiApiTimeout = 3
)

// tianji api获取云平台全量（集群列表,[集群-集群物理机列表]map,error
type Machine struct {
	Cluster string `json:"m.cluster"`
	ID      string `json:"m.id"`
}

func GetAllClustrMachineList() ([]string, map[string][]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*tianjiApiTimeout)
	defer cancel()
	url := "http://127.0.0.1:7070/api/v3/column/m.id,m.cluster?m.sm_name!=VM&m.state=GOOD&m.cpu_arch=x86_64"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var machines []Machine
	if err := json.Unmarshal(ret, &machines); err != nil {
		return nil, nil, err
	}

	// cluster slice
	clusterSlice := make([]string, 0)
	clusterSet := make(map[string]struct{}) // 辅助生成clusterSlice用：map key唯一的原理
	// [ cluster-machine_list ]map
	clusterMap := make(map[string][]string)

	for _, v := range machines {
		if _, exist := clusterSet[v.Cluster]; !exist {
			clusterSlice = append(clusterSlice, v.Cluster)
			// 在map中标记已记录！
			clusterSet[v.Cluster] = struct{}{}
		}
		clusterMap[v.Cluster] = append(clusterMap[v.Cluster], v.ID)
	}

	return clusterSlice, clusterMap, nil
}

// tianji api获取云平台全量project
type Project struct {
	Proejct string `json:"c.project"`
}

func GetProjectList() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*tianjiApiTimeout)
	defer cancel()
	url := "localhost:7070/api/v3/column/c.project"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var projectList []Project
	if err := json.Unmarshal(ret, &projectList); err != nil {
		return nil, err
	}

	// 去重project
	projectMap := make(map[string]struct{})
	var ret_project_list []string

	for _, v := range projectList {
		if _, exist := projectMap[v.Proejct]; !exist {
			ret_project_list = append(ret_project_list, v.Proejct)
			// 在map中标记已记录！
			projectMap[v.Proejct] = struct{}{}
		}
	}

	return ret_project_list, nil
}

type MachineID struct {
	MID string `json:"m.id"`
}

// tianji api获取云平台全量物理机主机名列表
func GetAllMachineList() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*tianjiApiTimeout)
	defer cancel()
	url := "http://127.0.0.1:7070/api/v3/column/m.id?m.sm_name!=VM&m.state=GOOD&m.cpu_arch=x86_64"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var machines []MachineID

	if err := json.Unmarshal(ret, &machines); err != nil {
		return nil, err
	}

	var machineList []string
	for _, m := range machines {
		machineList = append(machineList, m.MID)
	}
	return machineList, nil

}

func GetAllMachineListNonGood() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*tianjiApiTimeout)
	defer cancel()
	url := "http://127.0.0.1:7070/api/v3/column/m.id?m.sm_name!=VM"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var machines []MachineID

	if err := json.Unmarshal(ret, &machines); err != nil {
		return nil, err
	}

	var machineList []string
	for _, m := range machines {
		machineList = append(machineList, m.MID)
	}
	return machineList, nil

}

// 外层映射
type Wrapper struct {
	ServiceRegistration string `json:"c.sr.service_registration"`
}

// 内层映射
type RegistrationData struct {
	// domain
	PopAnyTunnelVIP  string `json:"pop.anytunnel.vip"`
	PopDomainACS     string `json:"pop.domain.acs"`
	PopDomainACSBiz  string `json:"pop.domain.acs-biz"`
	PopDomainACSMgmt string `json:"pop.domain.acs-mgmt"`
	PopDomainACSOps  string `json:"pop.domain.acs-ops"`
	PopDomainFt      string `json:"pop.domain.ft"`
	// vip
	PopInternetVIP string `json:"pop.internet.vip"`
	PopIplist      string `json:"pop.iplist"`
	PopVIP         string `json:"pop.vip"`
	PopVIPBiz      string `json:"pop.vip.biz"`
	PopVIPMgmt     string `json:"pop.vip.mgmt"`
	PopVIPOps      string `json:"pop.vip.ops"`
	// ak/sk
	SLSPopAccessID  string `json:"sls_pop_accessId"`
	SLSPopAccessKey string `json:"sls_pop_accessKey"`
	SLSPopAliUUID   string `json:"sls_pop_aliUid"`
	SLSPopLogstore  string `json:"sls_pop_logstore"`
	SLSPopProject   string `json:"sls_pop_project"`
}

// 获取tianji kv.json获取sls ep/ak/sk
func GetAkSk(product string) (string, string, string, error) {
	var (
		ep string
		ak string
		sk string
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*tianjiApiTimeout)
	defer cancel()
	url := "http://127.0.0.1:7070/api/v3/column/c.sr.service_registration?c.sr.id=webapp-pop"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", "", "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()
	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", err
	}

	// 解析外层
	var wrapperData []Wrapper
	if err := json.Unmarshal(ret, &wrapperData); err != nil {
		return "", "", "", err
	}

	// 解析内层
	var registration RegistrationData
	if err := json.Unmarshal([]byte(wrapperData[0].ServiceRegistration), &registration); err != nil {
		return "", "", "", err
	}

	ak = registration.SLSPopAccessID
	sk = registration.SLSPopAccessKey

	// 这里默认返回sls-inr的endpoint
	ep, err = getSlsInnerEp()
	if err != nil {
		return "", "", "", err
	}

	return ep, ak, sk, nil
}

// 随机获取一台ops的ip
func GetRandomOpsIp() (string, error) {
	// 调用tianjiapi获取ops ip
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*tianjiApiTimeout)
	defer cancel()
	url := "http://127.0.0.1:7070/api/v3/column/m.ip?m.sr.id=ops.OpsNtp%23"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var (
		opsIpSlice []string
		// 这里可以确认每个m.ip的值都是string类型，所以不需要使用interface+断言的方式
		data []map[string]string
	)
	if err = json.Unmarshal(ret, &data); err != nil {
		return "", err
	}

	for _, v := range data {
		if ip, ok := v["m.ip"]; ok {
			opsIpSlice = append(opsIpSlice, ip)
		} else {
			return "", fmt.Errorf("不存在m.ip字段")
		}
	}

	if len(opsIpSlice) == 0 {
		return "", fmt.Errorf("GetRandomOpsIp(): len(opsIpSlice) == 0")
	}

	// 随机返回一个ops ip
	rand.Seed(time.Now().UnixNano())
	randomOpsIp := opsIpSlice[rand.Intn(len(opsIpSlice))]

	return randomOpsIp, nil

}

// 外层映射
type ServiceResult struct {
	Result string `json:"service.res.result"`
}

// 内层映射
type InnerResult struct {
	Ip     string `json:"ip"`
	Domain string `json:"domain"`
	Dns    string `json:"dns"`
	Alias  string `json:"alias"`
}

// 获取服务的endpoint(默认是sls-inr)
func getSlsInnerEp() (ep string, err error) {
	product := "sls-common"

	// 调用tianjiapi获取sls_common的服务注册变量
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*tianjiApiTimeout)
	defer cancel()

	url := fmt.Sprintf("http://127.0.0.1:7070/api/v3/column/c.sr.service_registration?c.sr.id=%s", product)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 解析服务注册变量
	// 外层切片
	var data []map[string]interface{}

	err = json.Unmarshal([]byte(ret), &data)
	if err != nil {
		msg := fmt.Sprintf("parse sls-common service registation failed [outer]! err:%v\n", err)
		err = errors.New(msg)
		return "", err
	}

	// 内层c.sr.service_registration值（字符串）
	serviceRegistrationString := data[0]["c.sr.service_registration"].(string)

	// c.sr.service_registration对象
	var serviceRegistration map[string]interface{}
	err = json.Unmarshal([]byte(serviceRegistrationString), &serviceRegistration)
	if err != nil {
		msg := fmt.Sprintf("parse sls-common service registation failed [inner]! err:%v\n", err)
		err = errors.New(msg)
		return "", err
	}

	// sls_data.endpoint
	slsDataEndpoint, ok := serviceRegistration["sls_data.endpoint"].(string)
	if !ok {
		msg := fmt.Sprintln("get sls_data.endpoint failed!")
		err = errors.New(msg)
		return "", err
	}

	return slsDataEndpoint, nil
}
