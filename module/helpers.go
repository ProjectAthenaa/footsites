package module

import (
	"fmt"
	http "github.com/ProjectAthenaa/sonic-core/fasttls"
	"math/rand"
	"regexp"
)

var (
	csrfMatchRe = regexp.MustCompile(`"csrfToken":"(\w+)"`)
)

func GetCacheNode() string{
	return nodeList[rand.Intn(3299)]
}

func (tk *Task) GenerateDefaultHeaders(referrer string) http.Headers {
	return http.Headers{
		`user-agent`:         {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36"},
		`accept`:             {`application/json`},
		`accept-encoding`:    {`gzip, deflate, br`},
		`accept-language`:    {`en-us`},
		`content-type`:       {`application/x-www-form-urlencoded; charset=UTF-8`},
		`sec-ch-ua`:          {`"Chromium";v="91", " Not A;Brand";v="99", "Google Chrome";v="91"`},
		`sec-ch-ua-mobile`:   {`?0`},
		`Sec-Fetch-Site`:     {`same-site`},
		`Sec-Fetch-Dest`:     {`empty`},
		`Sec-Fetch-Mode`:     {`cors`},
		`referer`:            {referrer},
		`X-Requested-With`:   {`XMLHttpRequest`},
		`origin`:             {fmt.Sprintf(`https://www.%s.com`, sites[tk.Site].name)},
		`Pragma`:             {`no-cache`},
		`Cache-Control`:      {`no-cache`},
		`Connection`:         {`keep-alive`},
	}
}