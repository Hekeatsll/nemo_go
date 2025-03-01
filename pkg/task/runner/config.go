package runner

type PortscanRequestParam struct {
	Target             string `form:"target"`
	IsPortScan         bool   `form:"portscan"`
	IsIPLocation       bool   `form:"iplocation"`
	IsFofa             bool   `form:"fofasearch"`
	IsQuake            bool   `form:"quakesearch"`
	IsHunter           bool   `form:"huntersearch"`
	Port               string `form:"port"`
	Rate               int    `form:"rate"`
	NmapTech           string `form:"nmap_tech"`
	CmdBin             string `form:"bin"`
	OrgId              int    `form:"org_id"`
	IsHttpx            bool   `form:"httpx"`
	IsPing             bool   `form:"ping"`
	ExcludeIP          string `form:"exclude"`
	IsScreenshot       bool   `form:"screenshot"`
	IsFingerprintHub   bool   `form:"fingerprinthub"`
	IsIconHash         bool   `form:"iconhash"`
	TaskMode           int    `form:"taskmode"`
	IsTaskCron         bool   `form:"taskcron" json:"-"`
	TaskCronRule       string `form:"cronrule" json:"-"`
	TaskCronComment    string `form:"croncomment" json:"-"`
	IsLoadOpenedPort   bool   `form:"load_opened_port"`
	IsIgnoreOutofChina bool   `form:"ignoreoutofchina"`
	IsIgnoreCDN        bool   `form:"ignorecdn"`
}

type DomainscanRequestParam struct {
	Target             string `form:"target"`
	OrgId              int    `form:"org_id"`
	IsSubfinder        bool   `form:"subfinder"`
	IsSubdomainBrute   bool   `form:"subdomainbrute"`
	IsFldDomain        bool   `form:"fld_domain"`
	IsHttpx            bool   `form:"httpx"`
	IsIPPortscan       bool   `form:"portscan"`
	IsSubnetPortscan   bool   `form:"networkscan"`
	IsCrawler          bool   `form:"crawler"`
	IsFofa             bool   `form:"fofasearch"`
	IsQuake            bool   `form:"quakesearch"`
	IsHunter           bool   `form:"huntersearch"`
	IsScreenshot       bool   `form:"screenshot"`
	IsICPQuery         bool   `form:"icpquery"`
	IsWhoisQuery       bool   `form:"whoisquery"`
	IsFingerprintHub   bool   `form:"fingerprinthub"`
	IsIconHash         bool   `form:"iconhash"`
	TaskMode           int    `form:"taskmode"`
	PortTaskMode       int    `form:"porttaskmode"`
	IsTaskCron         bool   `form:"taskcron" json:"-"`
	TaskCronRule       string `form:"cronrule" json:"-"`
	TaskCronComment    string `form:"croncomment" json:"-"`
	IsIgnoreOutofChina bool   `form:"ignoreoutofchina"`
	IsIgnoreCDN        bool   `form:"ignorecdn"`
}

type PocscanRequestParam struct {
	Target           string `form:"target"`
	IsXrayVerify     bool   `form:"xrayverify"`
	XrayPocFile      string `form:"xray_poc_file"`
	IsNucleiVerify   bool   `form:"nucleiverify"`
	NucleiPocFile    string `form:"nuclei_poc_file"`
	IsDirsearch      bool   `form:"dirsearch"`
	DirsearchExtName string `form:"ext"`
	IsLoadOpenedPort bool   `form:"load_opened_port"`
	IsTaskCron       bool   `form:"taskcron" json:"-"`
	TaskCronRule     string `form:"cronrule" json:"-"`
	TaskCronComment  string `form:"croncomment" json:"-"`
}

type XScanRequestParam struct {
	XScanType       string `form:"xscan_type"`
	Target          string `form:"target"`
	Port            string `form:"port"`
	OrgId           int    `form:"org_id"`
	IsOrgIP         bool
	IsOrgDomain     bool
	IsFofaSearch    bool   `form:"fofa"`
	IsFingerprint   bool   `form:"fingerprint"`
	IsXrayPocscan   bool   `form:"xraypoc"`
	XrayPocFile     string `form:"xraypocfile"`
	IsTaskCron      bool   `form:"taskcron" json:"-"`
	TaskCronRule    string `form:"cronrule" json:"-"`
	TaskCronComment string `form:"croncomment" json:"-"`
	IsCn            bool   `form:"is_CN"`
}

type taskKeySearchParam struct {
	KeyWord      string `json:"key_word"`
	SearchTime   string `json:"search_time"`
	ExcludeWords string `json:"exclude_words"`
	CheckMod     string `json:"check_mod"`
	Count        int    `json:"count"`
	IsCN         bool   `json:"is_CN"`
}
