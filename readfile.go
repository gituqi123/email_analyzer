package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func get_page(link string) (int, string) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	link = strings.Trim(link, "\r")
	//fmt.Println(link)
	resp, err := http.Get(link)
	if err != nil {
		//fmt.Printf("[ERROR] can't connect to website %s\n", link)
		//fmt.Println(err)
		return -1, "error"
	}
	resp.Header.Set("User-Agent", "fta15")
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		//fmt.Printf("[ERROR] can't read from %s \n", link)
		return -1, "error"
	}
	return resp.StatusCode, buf.String()
}

func read_file(path string) string {
	dat, err := os.ReadFile(path)
	check(err)
	return string(dat)
}

func get_version_roundcube(webpage string) string {
	version_index := strings.Index(webpage, "rcversion")
	next_word_index := strings.Index(webpage, "cookie_domain")
	version := next_word_index - version_index
	if version_index != -1 {
		fmt.Println("webpage_len", len(webpage))
		if len(webpage) == 5570 {
			fmt.Println(webpage)
		}
		rc_version := webpage[version_index : version_index+version]
		rc_version_array := strings.Split(rc_version, ":")
		fmt.Println(rc_version)
		rcversion := rc_version_array[1]
		rcversion = strings.Replace(rcversion, "\"", "", -1)
		rcversion = strings.Replace(rcversion, ",", "", -1)
		rcversion = strings.Replace(rcversion, "0", ".", -1)
		return rcversion
	} else {
		return "-1"
	}
}

func check_webmail(domain string) {
	prefix := "http://"
	status_code, web_page := get_page(prefix + domain)
	if status_code == 200 {
		if strings.Contains(web_page, "/owa/auth") {
			//s := []string{"https://", domain, "/ecp/Current/exporttool/microsoft.exchange.ediscovery.exporttool.application"}
			//version_endpoint := strings.Join(s, "")
			//fmt.Println(version_endpoint)
			//status_code, web_page = get_page(version_endpoint)
			//fmt.Printf("https://%s/ecp/Current/exporttool/microsoft.exchange.ediscovery.exporttool.application", domain)
			//fmt.Println("exchange", domain)
			//if status_code == 200 {
			//	fmt.Println(web_page)
			//	regex_match := regexp.MustCompile("<assemblyIdentity.*version=\"(\\d+.\\d+.\\d+.\\d+)\"")
			//	exchange_version := regex_match.FindAllString(web_page, -1)
			//	fmt.Println(domain)
			//	fmt.Println(exchange_version)
			//} else {
			//	status_code, web_page = get_page("https://" + domain + "/ecp/" + domain + "/exporttool/microsoft.exchange.ediscovery.exporttool.application'")
			//	fmt.Println(status_code)
			//	fmt.Println(web_page)
			//}
			fmt.Println("exchange mail detected ", prefix+domain)
		}
		if strings.Contains(web_page, "Roundcube") {
			version := get_version_roundcube(web_page)
			fmt.Println("roundcube  mail detected ", prefix+domain)
			fmt.Println(version)
		}
		if strings.Contains(web_page, "mdaemon") {
			fmt.Println("mdaemon  mail detected ", prefix+domain)
		}
		if strings.Contains(web_page, "SmarterMail") {
			fmt.Println("SmarterMail  mail detected ", prefix+domain)
		}
		if strings.Contains(web_page, "Zimbra") {
			fmt.Println("Zimbra  mail detected ", prefix+domain)
		}
		if strings.Contains(web_page, "MailEnable") {
			fmt.Println("MailEnable  mail detected ", prefix+domain)
		}
		if strings.Contains(web_page, "Zoho") {
			fmt.Println("Zoho  mail detected ", prefix+domain)
		}
		if strings.Contains(web_page, "IceWarp") {
			fmt.Println("IceWarp  mail detected ", prefix+domain)
		}
		if strings.Contains(web_page, "Index of") {
			fmt.Println("nomail  detected ", prefix+domain)
		}
	}
}

func main() {
	domains := read_file("get_csv.txt")
	splitted_string := strings.Split(domains, "\n")
	for i := 0; i < len(splitted_string); i++ {
		go check_webmail(splitted_string[i])
	}
	select {}
}
