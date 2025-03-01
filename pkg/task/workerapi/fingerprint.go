package workerapi

import (
	"encoding/json"
	"github.com/hanc00l/nemo_go/pkg/comm"
	"github.com/hanc00l/nemo_go/pkg/logging"
	"github.com/hanc00l/nemo_go/pkg/task/domainscan"
	"github.com/hanc00l/nemo_go/pkg/task/fingerprint"
	"github.com/hanc00l/nemo_go/pkg/task/portscan"
	"strings"
)

const (
	IPNumberPerSubTask     = 10
	DomainNumberPerSubTask = 20
)

type FingerprintTaskConfig struct {
	IsHttpx          bool
	IsFingerprintHub bool
	IsIconHash       bool
	IsScreenshot     bool
	IPTargetMap      map[string][]int
	DomainTargetMap  map[string]struct{}
	WorkspaceId      int
}

// Fingerprint 指纹识别任务，将所有指纹识别任务整合后，可由worker将端口扫描与域名收集后的结果进行二次拆分为多个任务，提高指纹识别效率
func Fingerprint(taskId, mainTaskId, configJSON string) (result string, err error) {
	var ok bool
	if ok, result, err = CheckTaskStatus(taskId); !ok {
		return result, err
	}
	config := FingerprintTaskConfig{}
	if err = ParseConfig(configJSON, &config); err != nil {
		return FailedTask(err.Error()), err
	}
	//
	_, _, result, err = doFingerPrintAndSave(taskId, mainTaskId, config)
	if err != nil {
		logging.RuntimeLog.Error(err)
		return FailedTask(err.Error()), err
	}

	return
}

// doFingerPrintAndSave 指纹的综合任务，包括IP与domain
func doFingerPrintAndSave(taskId string, mainTaskId string, config FingerprintTaskConfig) (resultPortScan portscan.Result, resultDomainScan domainscan.Result, result string, err error) {
	var domainPort map[string]map[int]struct{}
	// 返回的结果
	resultPortScan.IPResult = make(map[string]*portscan.IPResult)
	resultDomainScan.DomainResult = make(map[string]*domainscan.DomainResult)
	// 通过RPC调用保存结果
	// 保存结果
	resultArgs := comm.ScanResultArgs{
		TaskID:     taskId,
		MainTaskId: mainTaskId,
	}
	// IP的指纹任务
	if len(config.IPTargetMap) > 0 {
		for ip, ports := range config.IPTargetMap {
			resultPortScan.SetIP(ip)
			for _, port := range ports {
				resultPortScan.SetPort(ip, port)
			}
		}
		portscanConfig := portscan.Config{
			IsHttpx:          config.IsHttpx,
			IsFingerprintHub: config.IsFingerprintHub,
			IsIconHash:       config.IsIconHash,
			WorkspaceId:      config.WorkspaceId,
		}
		doIPFingerPrint(portscanConfig, &resultPortScan)
		resultArgs.IPConfig = &portscanConfig
		resultArgs.IPResult = resultPortScan.IPResult
	}
	// 域名的指纹任务
	if len(config.DomainTargetMap) > 0 {
		// 读取目标的数据库中已保存的开放端口
		args := comm.LoadDomainOpenedPortArgs{
			WorkspaceId: config.WorkspaceId,
			Target:      config.DomainTargetMap,
		}
		err = comm.CallXClient("LoadDomainOpenedPort", &args, &domainPort)
		if err != nil {
			logging.RuntimeLog.Error(err)
		}
		for domain := range config.DomainTargetMap {
			resultDomainScan.SetDomain(domain)
		}
		domainscanConfig := domainscan.Config{
			IsHttpx:          config.IsHttpx,
			IsFingerprintHub: config.IsFingerprintHub,
			IsIconHash:       config.IsIconHash,
			WorkspaceId:      config.WorkspaceId,
		}
		doDomainFingerPrint(domainscanConfig, &resultDomainScan, domainPort)
		resultArgs.DomainConfig = &domainscanConfig
		resultArgs.DomainResult = resultDomainScan.DomainResult
	}
	// 保存结果
	err = comm.CallXClient("SaveScanResult", &resultArgs, &result)
	if err != nil {
		logging.RuntimeLog.Error(err)
		return
	}
	// screenshot任务
	if config.IsScreenshot {
		resultScreenshot := doScreenshotAndSave(config.WorkspaceId, mainTaskId, &resultPortScan, &resultDomainScan, domainPort)
		result = strings.Join([]string{result, resultScreenshot}, ",")
	}

	return
}

// doIPFingerPrint 对 IP结果进行指纹识别
func doIPFingerPrint(config portscan.Config, resultPortScan *portscan.Result) {
	if config.IsHttpx {
		//httpx := fingerprint.NewHttpx()
		httpx := fingerprint.NewHttpxFinger()
		httpx.ResultPortScan = *resultPortScan
		httpx.DoHttpxAndFingerPrint()
	}
	if config.IsFingerprintHub {
		fp := fingerprint.NewFingerprintHub()
		fp.ResultPortScan = *resultPortScan
		fp.Do()
	}
	if config.IsIconHash {
		doIconHashAndSave(config.WorkspaceId, resultPortScan, nil, nil)
	}
}

// doDomainFingerPrint 对域名结果进行指纹识别
func doDomainFingerPrint(config domainscan.Config, resultDomainScan *domainscan.Result, domainPort map[string]map[int]struct{}) {
	// 指纹识别
	if config.IsHttpx {
		//httpx := fingerprint.NewHttpx()
		httpx := fingerprint.NewHttpxFinger()
		httpx.ResultDomainScan = *resultDomainScan
		httpx.DomainTargetPort = domainPort
		httpx.DoHttpxAndFingerPrint()
	}
	if config.IsFingerprintHub {
		fp := fingerprint.NewFingerprintHub()
		fp.ResultDomainScan = *resultDomainScan
		fp.DomainTargetPort = domainPort
		fp.Do()
	}
	if config.IsIconHash {
		doIconHashAndSave(config.WorkspaceId, nil, resultDomainScan, domainPort)
	}
}

// doScreenshotAndSave 执行Screenshot并保存
func doScreenshotAndSave(workspaceId int, mainTaskId string, resultIPScan *portscan.Result, resultDomainScan *domainscan.Result, domainPort map[string]map[int]struct{}) (result string) {
	ss := fingerprint.NewScreenShot()
	if resultIPScan != nil {
		ss.ResultPortScan = *resultIPScan
	}
	if resultDomainScan != nil {
		ss.ResultDomainScan = *resultDomainScan
		ss.DomainTargetPort = domainPort
	}
	ss.Do()
	args := comm.ScreenshotResultArgs{
		MainTaskId:  mainTaskId,
		FileInfo:    ss.LoadResult(),
		WorkspaceId: workspaceId,
	}
	err := comm.CallXClient("SaveScreenshotResult", &args, &result)
	if err != nil {
		logging.RuntimeLog.Error(err)
		return err.Error()
	}
	return
}

// doIconHashAndSave 获取icon，并将icon image保存到服务端
func doIconHashAndSave(workspaceId int, resultIPScan *portscan.Result, resultDomainScan *domainscan.Result, domainPort map[string]map[int]struct{}) (result string) {
	hash := fingerprint.NewIconHash()
	if resultIPScan != nil {
		hash.ResultPortScan = *resultIPScan
	}
	if resultDomainScan != nil {
		hash.ResultDomainScan = *resultDomainScan
		hash.DomainTargetPort = domainPort
	}
	hash.Do()

	if len(hash.IconHashInfoResult.Result) <= 0 {
		return ""
	}
	args := comm.IconHashResultArgs{
		WorkspaceId:  workspaceId,
		IconHashInfo: hash.IconHashInfoResult.Result,
	}
	err := comm.CallXClient("SaveIconImageResult", &args, &result)
	if err != nil {
		logging.RuntimeLog.Error(err)
		return err.Error()
	}
	return
}

// NewFingerprintTask 根据端口及域名扫描结果，根据设置的拆分规模，生成指纹识别子任务
func NewFingerprintTask(taskId string, mainTaskId string, portScanResult *portscan.Result, domainScanResult *domainscan.Result, config FingerprintTaskConfig) (result string, err error) {
	result, err = newFingerprintTask(taskId, mainTaskId, portScanResult, domainScanResult, config, "fingerprint")
	return
}

// newFingerprintTask 根据端口及域名扫描结果，根据设置的拆分规模，生成指纹识别子任务
func newFingerprintTask(taskId, mainTaskId string, portScanResult *portscan.Result, domainScanResult *domainscan.Result, config FingerprintTaskConfig, taskName string) (result string, err error) {
	if config.IsHttpx == false && config.IsFingerprintHub == false && config.IsIconHash == false && config.IsScreenshot == false {
		return
	}
	//拆分子任务
	ipTarget, domainTarget := MakeSubTaskTarget(portScanResult, domainScanResult)
	//生成任务
	for _, t := range ipTarget {
		newConfig := config
		newConfig.IPTargetMap = t
		result, err = sendTask(taskId, mainTaskId, newConfig, taskName)
		if err != nil {
			return
		}
	}
	for _, t := range domainTarget {
		newConfig := config
		newConfig.DomainTargetMap = t
		result, err = sendTask(taskId, mainTaskId, newConfig, taskName)
		if err != nil {
			return
		}
	}
	return
}

// MakeSubTaskTarget 根据端口及域名扫描结果，根据设置的拆分规模，生成指纹识别等子任务
func MakeSubTaskTarget(portScanResult *portscan.Result, domainScanResult *domainscan.Result) (ipTarget []map[string][]int, domainTarget []map[string]struct{}) {
	if portScanResult != nil && len(portScanResult.IPResult) > 0 {
		index := 0
		mapIpPort := make(map[string][]int)
		for ip, ipr := range portScanResult.IPResult {
			mapIpPort[ip] = make([]int, 0)
			for port := range ipr.Ports {
				mapIpPort[ip] = append(mapIpPort[ip], port)
			}
			index++
			if index%IPNumberPerSubTask == 0 || index == len(portScanResult.IPResult) {
				ipTarget = append(ipTarget, mapIpPort)
				mapIpPort = make(map[string][]int)
			}
		}
	}
	if domainScanResult != nil && len(domainScanResult.DomainResult) > 0 {
		index := 0
		mapDomain := make(map[string]struct{})
		for domain := range domainScanResult.DomainResult {
			mapDomain[domain] = struct{}{}
			index++
			if index%DomainNumberPerSubTask == 0 || index == len(domainScanResult.DomainResult) {
				domainTarget = append(domainTarget, mapDomain)
				mapDomain = make(map[string]struct{})
			}
		}
	}
	return
}

// sendTask 调用api发送任务
func sendTask(taskId string, mainTaskId string, config interface{}, taskName string) (result string, err error) {
	configMarshal, err := json.Marshal(config)
	if err != nil {
		return
	}
	newTaskArgs := comm.NewTaskArgs{
		MainTaskID:    mainTaskId,
		LastRunTaskId: taskId,
		TaskName:      taskName,
		ConfigJSON:    string(configMarshal),
	}
	err = comm.CallXClient("NewTask", &newTaskArgs, &result)
	if err != nil {
		logging.RuntimeLog.Errorf("Start %s task fail:%v", taskName, err)
		logging.CLILog.Errorf("Start %s task fail:%v", taskName, err)
	}
	return
}
