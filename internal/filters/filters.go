package filters

import "regexp"

type Result struct {
	Domains 	[]string
	Ips			[]string
}

func FilterIP(domains []string) Result {
	var results Result
	for _,domain := range domains {
		re := regexp.MustCompile(`((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)`)
		if re.MatchString(domain) {
			results.Ips  = append(results.Ips , domain)
		} else {
			results.Domains  = append(results.Domains , domain)
		}
	}
	return results
}