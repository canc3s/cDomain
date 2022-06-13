package runner

import (
	"github.com/antchfx/htmlquery"
	"github.com/canc3s/cDomain/internal/gologger"
	"github.com/canc3s/cDomain/internal/requests"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetPage(url string, options *Options) requests.Response {

	time.Sleep(time.Duration(options.Delay) * time.Second)
	var transport = requests.DefaultTransport()
	var client = &http.Client{
		Transport: transport,
		//Timeout:       time.Duration(options.Timeout),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse /* 不进入重定向 */
		},
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5026.0 Safari/537.36 Edg/103.0.1254.0")
	if options.Cookie != "" {
		req.Header.Set("Cookie", options.Cookie)
	}
	resp, err := client.Do(req)
	if err != nil {
		gologger.Fatalf("请求发生错误，请检查网络连接\n%s\n", err)
	}

	if resp.StatusCode == 403 {
		gologger.Fatalf("海外用户或者云服务器ip被禁止访问网站，请更换ip\n")
	} else if resp.StatusCode == 401 {
		gologger.Fatalf("天眼查Cookie有问题或过期，请重新获取\n")
	} else if resp.StatusCode == 302 {
		gologger.Fatalf("天眼查免费查询次数已用光，需要加Cookie\n")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	page,_ := htmlquery.Parse(strings.NewReader(string(body)))

	return requests.Response{
		Body: body,
		Page:page,
	}
}

func JudgePagesK(page *html.Node) int {
	list := htmlquery.Find(page, "/html/body/div[2]/div/div[2]/div[1]/div[2]/div[3]/ul/li/a")
	num := 1
	if len(list) > 2 {
		var err error
		pages := htmlquery.InnerText(list[len(list)-2])
		num,err = strconv.Atoi(strings.Trim(pages, "."))
		if err != nil {
			num = 1
		}
	}
	return num
}

func JudgePagesI(page *html.Node) int {
	list := htmlquery.Find(page, "/html/body/div/ul/li/a")
	return len(list) - 1
}

func EnuDomainByKey(page *html.Node, domains *[]string) {
	list := htmlquery.Find(page, "/html/body/div[2]/div/div[2]/div[1]/div[2]/div[2]/table/tbody/tr/td[5]/span")
	for _,node  := range list {
		domain := htmlquery.InnerText(node)
		*domains = append(*domains, domain)
	}
}

func GetInformation(page *html.Node) []string {
	list := htmlquery.Find(page, "/html/body/table/tbody/tr/td[5]")
	var domains []string
	for _,node := range list {
		domain := htmlquery.InnerText(node)
		domains = append(domains, domain)
	}
	return domains
}


func GetDomain(options *Options) []string {
	resp := GetPage("https://www.tianyancha.com/pagination/icp.xhtml?ps=30&isAjaxLoad=true&pn=1&id="+options.CompanyID, options)
	page := JudgePagesI(resp.Page)
	domains := GetInformation(resp.Page)
	for i := 2; i <= page; i++ {
		resp := GetPage("https://www.tianyancha.com/pagination/icp.xhtml?ps=30&isAjaxLoad=true&pn="+strconv.Itoa(i)+"&id="+options.CompanyID, options)
		domains = append(domains,GetInformation(resp.Page)...)
	}
	return domains
}
