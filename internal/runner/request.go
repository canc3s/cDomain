package runner

import (
	"github.com/antchfx/htmlquery"
	"github.com/canc3s/cDomain/internal/gologger"
	"github.com/canc3s/cDomain/internal/requests"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func GetPage(url string, options *Options) requests.Response {

	var transport = requests.DefaultTransport()
	var client = &http.Client{
		Transport: transport,
		//Timeout:       time.Duration(options.Timeout),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse /* 不进入重定向 */
		},
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 14_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/7.0.18(0x17001231) NetType/4G Language/zh_CN")
	req.Header.Set("Authorization", "0###oo34J0fXrrGHr0A8plQcS2-uWwn4###1614567391527###a851050f2a3276bd71b91b39789cda9a")
	req.Header.Set("version", "TYC-XCX-WX")
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

func EnuDomainByKey(page *html.Node, domains *[]string) {
	list := htmlquery.Find(page, "/html/body/div[2]/div/div[2]/div[1]/div[2]/div[2]/table/tbody/tr/td[5]/span")
	for _,node  := range list {
		domain := htmlquery.InnerText(node)
		*domains = append(*domains, domain)
	}
}


func GetDomain(options *Options) [][][]byte {
	resp := GetPage("https://api9.tianyancha.com/services/v3/ar/icp/"+options.CompanyID+".json", options)
	re := regexp.MustCompile(`"ym":"(.*?)"`)
	return re.FindAllSubmatch(resp.Body,-1)
}