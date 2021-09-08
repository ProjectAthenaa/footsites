package module

import (
	"github.com/ProjectAthenaa/sonic-core/protos/module"
	"github.com/ProjectAthenaa/sonic-core/sonic/frame"
	"strings"
)

func (tk *Task) AwaitMonitor(){
	pubsub, err := frame.SubscribeToChannel(tk.Data.Channels.MonitorChannel)
	if err != nil{
		tk.Stop()
		return
	}
	defer pubsub.Close()

	var variantid, size, color, price string
	tk.SetStatus(module.STATUS_MONITORING)
	for monitorData := range pubsub.Chan(tk.Ctx){
		variantid = monitorData["variantid"].(string)
		size = monitorData["size"].(string)
		color = monitorData["color"].(string)
		price = monitorData["price"].(string)
		if inString(size, tk.Data.TaskData.Size) && inString(color, tk.Data.TaskData.Color){
			tk.PID = variantid
			tk.ReturningFields.Price = price
		}
	}
}

func inString(keyword string, stringslice []string) bool{
	kwlower := strings.ToLower(keyword)
	for _, kw := range stringslice{
		if strings.Contains(kwlower,strings.ToLower(kw)){
			return true
		}
	}
	return false
}