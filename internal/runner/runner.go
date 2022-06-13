package runner

import (
	"fmt"
	"github.com/canc3s/cDomain/internal/fileutil"
	"github.com/canc3s/cDomain/internal/filters"
	"github.com/canc3s/cDomain/internal/gologger"
	"net/url"
	"os"
	"regexp"
	"strconv"
)

type Targets struct {
	ID			[]string
	Name		[]string
}

func GetDomainByKey(options *Options) []string {
	var domains []string
	gologger.Infof("正在查询关键字 %s\n",options.KeyWord)
	url := "https://beian.tianyancha.com/search/"+ url.QueryEscape(options.KeyWord)
	resp := GetPage(url, options)
	EnuDomainByKey(resp.Page, &domains)
	num := JudgePagesK(resp.Page)
	if num > 5 && options.Cookie == "" {
		gologger.Errorf("域名过多，需要cookie才能完整查询\n")
		num = 5
	}
	if num > 1 {
		for i := 2; i <= num; i++ {
			resp := GetPage(url+"/p"+strconv.Itoa(num), options)
			EnuDomainByKey(resp.Page, &domains)
		}
	}
	return domains
}

func GetDomainByID(options *Options) []string {
	domains := GetDomain(options)
	return domains
}

func RunEnumeration(options *Options) {
	var domains []string
	if options.InputFile != "" {
		fin, error := os.OpenFile(options.InputFile, os.O_RDONLY, 0)
		if error != nil {
			gologger.Fatalf("文件读取失败：%s",error)
		}
		defer fin.Close()
		imf := fileutil.ReadImf(fin)
		targets := TransImf(imf)
		for _,id := range targets.ID {
			options.CompanyID = id
			domains = append(domains,GetDomainByID(options)...)
		}
		for _,name := range targets.Name {
			options.KeyWord = name
			domains = append(domains,GetDomainByKey(options)...)
		}
	} else {
		if options.KeyWord != "" && options.CompanyID == "" {
			domains = GetDomainByKey(options)
		} else if options.CompanyID != "" {
			domains = GetDomainByID(options)
		}
	}

	results := filters.FilterIP(domains)
	for _,ip := range results.Ips {
		gologger.Warningf("find IP : %s\n",ip)
	}
	for _,domain := range results.Domains {
		fmt.Println(domain)
	}
	if options.Output != "" {
		file, err := os.OpenFile(options.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			gologger.Fatalf("结果无法写入文件：%s", err)
		}
		defer file.Close()

		for _,domain := range results.Domains {
			file.WriteString(domain+"\n")
		}
	}
}

func TransImf(imf []string) Targets {
	var targets Targets
	for _,i := range imf {
		re := regexp.MustCompile(`(\d{5,11})`)
		buf := re.FindStringSubmatch(i)
		if buf == nil {
			targets.Name = append(targets.Name, i)
		}else{
			targets.ID = append(targets.ID, buf[0])
		}
	}
	return targets
}