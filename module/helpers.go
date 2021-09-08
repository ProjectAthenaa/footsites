package module

import (
	"github.com/ProjectAthenaa/sonic-core/fasttls"
	"math/rand"
	"regexp"
)

var (
	csrfMatchRe = regexp.MustCompile(`"csrfToken":"(\w+)"`)
)

func GetCacheNode() string{
	return nodeList[rand.Intn(3299)]
}

func (tk *Task) GenerateDefaultHeaders(referrer string) fasttls.Headers{
	return nil
}