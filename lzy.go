package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"regexp"
    "os"
    "fmt"
)

func main()  {
    var R1, R2, R3, R4 Rr
    u := "https://lanzous.com"
    q := ""
    p := ""
    R4.u = ""
    if len(os.Args) == 2 {
        q = os.Args[1]
        if match(q, `http.*?//.*?\.lanzous\.com/.*?`, "") != nil {
            if strings.Index(q, "密码")!= -1 {
                u = match(q, `^(http.{4}.*?)/`, "ERROR!0x00")[0][1]
                p = match(q, `密码.(.*?)$`, "ERROR!0x01")[0][1]
                q = match(q, `^http.*?//.*?/(.*?)密码`, "ERROR!0x02")[0][1]
            }else if match(q, `http.*?//`, "") != nil{
                u = match(q, `^(http.{4}.*?)/`, "ERROR!0x03")[0][1]
                q = match(q, `^http.*?//.*?/(.*?)$`, "ERROR!0x04")[0][1]
            }            
        }
    }else if len(os.Args) == 3 {
        q = os.Args[1]
        p = os.Args[2]
    }else if len(os.Args) == 1 {
        //print("缺少参数\n")
        fmt.Print("请输入id：")
        fmt.Scan(&q)
        fmt.Print("请输入密码：")
        fmt.Scanln(&p)
    }else{
        print("参数异常\n")
        os.Exit(0)
    }
    R1.u = u + "/" + q
    r1 := Http(R1)
    m1 := match(r1.b, `<iframe.class="ifr2".*?src="(.*?)"`, "")
    m11 := match(r1.b, `data.:.'(.*?)'`, "")
    sign := ""
    if m1 != nil {
        for i := 0; i < len(m1); i++ {
            if len(m1[i][1])>30 {
                R2.u = u + m1[i][1]
                r2 := Http(R2)
                m2 := match(r2.b, `.data.:.*?sign\':(.*?),`, "")
                for j := 0; j < len(m2); j++ {
                    if strings.Index(m2[j][0], "/") == -1 {
                        if len(m2[j][1]) > 33 {
                            sign = m2[j][1]
                        }else if len(m2[j][1]) < 20 {
                            m21 := match(r2.b, m2[j][1] + `.=.(.*?);`, "ERROR!0x05")
                            sign = m21[0][1]
                        }
                        R3.u = u + "/ajaxm.php"
                        R3.h = [][]string{
                            {"Referer", R2.u},
                        }
                        R3.m = "POST"
                        R3.d = "action=downprocess&ves=1&sign=" + strings.Replace(sign, "'", "", -1)
                        r3 := Http(R3)
                        _ = R4
                        m31 := match(r3.b, `dom":"(.*?)"`, "ERROR!0x06")
                        dom := strings.Replace(m31[0][1], "\\", "", -1)
                        m32 := match(r3.b, `url":"(.*?)"`, "ERROR!0x07")
                        url := m32[0][1]
                        R4.u = dom + "/file/" + url
                    }
                }
                
            }
        }
    }else if m11 != nil {
        if p != "" {
            R2.d = m11[0][1] + p
            R2.h = [][]string{
                {"Referer", R1.u},
            }
            R2.u = u + "/ajaxm.php"
            R2.m = "POST"
            r2 := Http(R2)
            m2 := match(r2.b, `zt":(.*?),`, "ERROR!0x08")
            if m2[0][1] != "0" {
                m21 := match(r2.b, `dom":"(.*?)"`, "ERROR!0x09")
                dom := strings.Replace(m21[0][1], "\\", "", -1)
                m22 := match(r2.b, `url":"(.*?)"`, "ERROR!0x10")
                url := m22[0][1]
                R4.u = dom + "/file/" + url
            }else{
                print("密码错误\n")
                os.Exit(0)
            }
        }else{
            print("缺少密码\n")
            os.Exit(0)
        }
    }
    if R4.u != "" {
        R4.h = [][]string{
            {"accept-language", "zh-CN,zh;q=0.9"},
        }
        r4 := Http(R4)
        m4 := match(r4.rh, `Location:(.*?)\r\n`, "ERROR!0x11")
        print(m4[0][1] + "\n")
        os.Exit(0)
    }
    print("来晚啦...文件取消分享了\n")

}

type Rr struct {
	Url,u string
	Header,h [][]string
	Status,s string
	Data,d string
	Body,b string
	Method,m string
	ReturnHeader,rh string
}

func Http(Req Rr) Rr{
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		 },
	}
	D := strings.NewReader(Req.d)
	request, err := http.NewRequest(Req.m, Req.u, D)
	if(len(Req.h) != 0){
        request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		for i := 0; i < len(Req.h); i++ {
			request.Header.Set(Req.h[i][0], Req.h[i][1])
        }
	}
	if err != nil {
		print("ERROR!!0x22\n")
		os.Exit(0)
	}
    response, err := client.Do(request)
    if err != nil {
        print("ERROR!!0x23\n")
        os.Exit(0)
    }
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		print("ERROR!!0x33\n")
		os.Exit(0)
	}
	h := response.Header
	hs := ""
	for k := range h {
        hs = hs + k + ":" + h[k][0] + "\r\n"
	}
	var ret Rr
	ret.rh = hs + response.Status + "\r\n"
	ret.b = string(body)
	ret.s = response.Status
	return ret
}

func match(nr, reg, err string) [][]string{
	p := regexp.MustCompile(reg)
	result := p.FindAllStringSubmatch(nr,-1)
	if len(result) != 0 {
		return result
	}else{
        if err != "" {
            print(err + "\n")
            os.Exit(0)
        }
        return nil
	}
}