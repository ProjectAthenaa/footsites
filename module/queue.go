package module

import (
	b64 "encoding/base64"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ProjectAthenaa/sonic-core/protos/module"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var(
	userIdRe = regexp.MustCompile(`userid="([\w-]+)"`)
	queueItInputRe = regexp.MustCompile(`"input":"([\w-]+)"`)
	queueItZeroCountRe = regexp.MustCompile(`"zeroCount":(\d+)}`)
	metaIdRe = regexp.MustCompile(`"meta":"([^"]+)"`)
	sessionIdRe = regexp.MustCompile(`"sessionId":"([\w-]+)"`)
	eventIdRe = regexp.MustCompile(`e=(.*?)&`)
)

type QueueItEncode struct {
	UserID    string `json:"userId"`
	Meta      string `json:"meta"`
	SessionID string `json:"sessionId"`
	Solution  struct {
		Postfix int    `json:"postfix"`
		Hash    string `json:"hash"`
	} `json:"solution"`
	Tags  []string `json:"tags"`
	Stats struct {
		Duration       int    `json:"duration"`
		Tries          int    `json:"tries"`
		UserAgent      string `json:"userAgent"`
		Screen         string `json:"screen"`
		Browser        string `json:"browser"`
		BrowserVersion string `json:"browserVersion"`
		IsMobile       bool   `json:"isMobile"`
		Os             string `json:"os"`
		OsVersion      string `json:"osVersion"`
		CookiesEnabled bool   `json:"cookiesEnabled"`
	} `json:"stats"`
	Parameters struct {
		Input     string `json:"input"`
		ZeroCount int    `json:"zeroCount"`
	} `json:"parameters"`
}
type QueueItVerify struct {
	ChallengeType string `json:"challengeType"`
	SessionID     string `json:"sessionId"`
	CustomerID    string `json:"customerId"`
	EventID       string `json:"eventId"`
	Version       int    `json:"version"`
}

//302 is queue-it
//529 is fastly
func (tk *Task) Poll(){
	req, err := tk.NewRequest("GET", fmt.Sprintf("https://www.%s.come/en/product/%s.html", sites[tk.Site].name, tk.PID), nil)
	if err != nil{
		tk.SetStatus(module.STATUS_ERROR, "couldnt create product get req")
		tk.Stop()
		return
	}
	req.Headers = tk.GenerateDefaultHeaders(fmt.Sprintf("https://www.%s.com", sites[tk.Site].name))

	res, err := tk.Do(req)
	if err != nil{
		tk.SetStatus(module.STATUS_ERROR, "couldnt get product page")
		tk.Stop()
		return
	}

	switch res.StatusCode {
	case 302:
		tk.QueueItRedirect = res.Headers["location"][0]
		tk.QueueIt()
		return
	case 529:
		tk.Fastly()
		return
	}
	return
}

func (tk *Task) QueueIt(){
	req, err := tk.NewRequest("GET", tk.QueueItRedirect, nil)
	if err != nil{
		tk.SetStatus(module.STATUS_ERROR, "couldnt create queue-it init request")
		tk.Stop()
		return
	}
	req.Headers = tk.GenerateDefaultHeaders(fmt.Sprintf("https://www.%s.com", sites[tk.Site].name))

	res, err := tk.Do(req)
	if err != nil{
		tk.SetStatus(module.STATUS_ERROR, "couldnt get queue-it page")
		tk.Stop()
		return
	}

	tk.QueueItUserId = string(userIdRe.FindSubmatch(res.Body)[1])

	req, err = tk.NewRequest("POST", fmt.Sprintf("https://www.%s.com/challengeapi/pow/challenge/%s", sites[tk.Site].name, tk.QueueItUserId), nil)
	if err != nil{
		tk.SetStatus(module.STATUS_ERROR, "couldnt create queue-it challenge request")
		tk.Stop()
		return
	}
	req.Headers = tk.GenerateDefaultHeaders(tk.QueueItRedirect)

	res, err = tk.Do(req)
	if err != nil{
		tk.SetStatus(module.STATUS_ERROR, "couldnt get queue-it challenge")
		tk.Stop()
		return
	}

	metaId := string(metaIdRe.FindSubmatch(res.Body)[1])
	sessionId := string(sessionIdRe.FindSubmatch(res.Body)[1])
	eventId := string(eventIdRe.FindString(tk.QueueItRedirect)[1])
	input := string(queueItInputRe.FindSubmatch(res.Body)[1])
	count, err := strconv.Atoi(string(queueItZeroCountRe.FindSubmatch(res.Body)[1]))
	if err != nil{
		tk.SetStatus(module.STATUS_ERROR, "couldnt convert zero count to number")
		tk.Stop()
		return
	}

	postfix, hash := getHash(input, count)

	queueItPayload, err := json.Marshal(QueueItEncode{
		UserID:    tk.QueueItUserId,
		Meta:      metaId,
		SessionID: sessionId,
		Solution: struct {
			Postfix int    `json:"postfix"`
			Hash    string `json:"hash"`
		}{
			Postfix: postfix,
			Hash: hash,
		},
		Tags: []string{
			"powTag-CustomerId:"+sites[tk.Site].name,
			"powTag-EventId:"+eventId,
			"powTag-UserId:"+tk.QueueItUserId,
		},
		Stats: struct {
			Duration       int    `json:"duration"`
			Tries          int    `json:"tries"`
			UserAgent      string `json:"userAgent"`
			Screen         string `json:"screen"`
			Browser        string `json:"browser"`
			BrowserVersion string `json:"browserVersion"`
			IsMobile       bool   `json:"isMobile"`
			Os             string `json:"os"`
			OsVersion      string `json:"osVersion"`
			CookiesEnabled bool   `json:"cookiesEnabled"`
		}{
			Duration:3000+rand.Intn(2000),
			Tries:1,
			UserAgent:"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36",
			Screen:"1920 x 1080",
			Browser:"Chrome",
			BrowserVersion:"93.0.4577.63",
			IsMobile:false,
			Os:"Windows",
			OsVersion:"10",
			CookiesEnabled:true,
		},
		Parameters: struct {
			Input     string `json:"input"`
			ZeroCount int    `json:"zeroCount"`
		}{
			Input: input,
			ZeroCount: count,
		},
	})
	if err != nil{
		tk.SetStatus(module.STATUS_ERROR, "couldnt serialize queue-it payload")
		tk.Stop()
		return
	}
	sEnc := b64.StdEncoding.EncodeToString(queueItPayload)

	req, err = tk.NewRequest("POST", fmt.Sprintf(`https://www.%s.com/challengeapi/verify`, sites[tk.Site].name), []byte(fmt.Sprintf(`{"challengeType":"proofofwork","sessionId":"%s","customerId":"%s","eventId":"%s","version":5}`, sEnc, sites[tk.Site].name, eventId)))
	if err != nil{
		tk.SetStatus(module.STATUS_ERROR, "could not create queue-it verify post")
		tk.Stop()
		return
	}
	req.Headers = tk.GenerateDefaultHeaders(tk.QueueItRedirect)
	_, err = tk.Do(req)
	if err != nil{
		tk.SetStatus(module.STATUS_ERROR, "could not read queue-it post response")
		tk.Stop()
		return
	}
}

func getHash(input string, zeroCount int) (int, string) {
	zeros := strings.Repeat("0", zeroCount)
	for postfix := 0; ; postfix++ {
		str := input + strconv.Itoa(postfix)
		hash := sha256.New()
		hash.Write([]byte(str))
		encodedHash := hex.EncodeToString(hash.Sum(nil))
		if strings.HasPrefix(encodedHash, zeros) {
			return postfix, encodedHash
		}
	}
}

func (tk *Task) Fastly(){
	for {
		req, err := tk.NewRequest("GET", fmt.Sprintf("https://www.%s.come/en/product/%s.html", sites[tk.Site].name, tk.PID), nil)
		if err != nil{
			tk.SetStatus(module.STATUS_ERROR, "couldnt create product get req")
			tk.Stop()
			return
		}
		req.Headers = tk.GenerateDefaultHeaders(fmt.Sprintf("https://www.%s.com", sites[tk.Site].name))

		res, err := tk.Do(req)
		if err != nil{
			tk.SetStatus(module.STATUS_ERROR, "couldnt get product page")
			tk.Stop()
			return
		}

		if res.StatusCode != 529{
			return
		}

		time.Sleep(15*time.Second)
	}
}