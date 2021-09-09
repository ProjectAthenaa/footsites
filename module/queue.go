package module

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/ProjectAthenaa/sonic-core/protos/module"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var(
	userIdRe = regexp.MustCompile(`userid="([\w-]+)"`)
	queueItInputRe = regexp.MustCompile(`"input":"([\w-]+)"`)
	queueItZeroCountRe = regexp.MustCompile(`"zeroCount":(\d+)}`)
)

//302 is queueit
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

	input := string(queueItInputRe.FindSubmatch(res.Body)[1])
	count, err := strconv.Atoi(string(queueItZeroCountRe.FindSubmatch(res.Body)[1]))
	if err != nil{
		tk.SetStatus(module.STATUS_ERROR, "couldnt convert zero count to number")
		tk.Stop()
		return
	}

	postfix, hash := getHash(input, count)
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