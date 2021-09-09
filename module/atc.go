package module

import (
	"fmt"
	"github.com/ProjectAthenaa/sonic-core/protos/module"
	"time"
)

func (tk *Task) ATC(){
	for {
		req, err := tk.NewRequest("POST", fmt.Sprintf(`https://%s.hosts.fastly.net/api/users/carts/current/entries?timestamp=%d`, GetCacheNode(), time.Now().Unix()), []byte(fmt.Sprintf(`{"productQuantity":1,"productId":"%s"}`, tk.VariantId)))
		if err != nil{
			tk.SetStatus(module.STATUS_ERROR, "couldnt create atc request")
			tk.Stop()
			return
		}
		req.Headers = tk.GenerateDefaultHeaders(fmt.Sprintf("https://www.%s.com", sites[tk.Site].name))

		res, err := tk.Do(req)
		if err != nil{
			tk.SetStatus(module.STATUS_ERROR, "couldnt make atc request")
			tk.Stop()
			return
		}

		if res.StatusCode == 200{
			return
		}

		time.Sleep(3*time.Second)
	}
}

func (tk *Task) CartAuthenticate() {
	req, err := tk.NewRequest("PUT", fmt.Sprintf(`https://www.%s.com/api/users/carts/current/email/%s?timestamp=%d`, sites[tk.Site].name, tk.Data.Profile.Email, time.Now().Unix()), nil)
	if err != nil{
		tk.SetStatus(module.STATUS_ERROR, "couldnt create auth req")
		tk.Stop()
		return
	}
	req.Headers = tk.GenerateDefaultHeaders(fmt.Sprintf("https://www.%s.com", sites[tk.Site].name))

	_, err = tk.Do(req)
	if err != nil{
		tk.SetStatus(module.STATUS_ERROR, "couldnt put auth req")
		tk.Stop()
		return
	}
}