package module

import (
	"fmt"
	"github.com/ProjectAthenaa/sonic-core/protos/module"
)

func (tk *Task) InitializeSession(){
	req, err := tk.NewRequest("GET", fmt.Sprintf(`https://www.%s.com/api/v5/session`, tk.Site), nil)
	if err != nil{
		tk.SetStatus(module.STATUS_ERROR, "couldnt create init request")
		tk.Stop()
		return
	}
	req.Headers = tk.GenerateDefaultHeaders(fmt.Sprintf("https://www.%s.com", sites[tk.Site].name))

	res, err := tk.Do(req)
	if err != nil{
		tk.SetStatus(module.STATUS_ERROR, "couldnt make init request")
		tk.Stop()
		return
	}

	tk.CsrfToken = string(csrfMatchRe.FindSubmatch(res.Body)[1])
}