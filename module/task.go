package module

import (
	"github.com/ProjectAthenaa/sonic-core/protos/module"
	"github.com/ProjectAthenaa/sonic-core/sonic/base"
	"github.com/ProjectAthenaa/sonic-core/sonic/database/ent/product"
	"github.com/ProjectAthenaa/sonic-core/sonic/face"
)

var _ face.ICallback = (*Task)(nil)

type Task struct {
	*base.BTask

	PID string
	VariantId string
	Site product.Site

	QueueItRedirect string
	QueueItUserId string

	CsrfToken string
	CartGuid string
}

func NewTask(data *module.Data) *Task {
	task := &Task{
		BTask: &base.BTask{
			Data: data,
		},
		Site: product.Site(data.Metadata["site"]),
	}
	task.Callback = task
	task.Init()
	return task
}

func (tk *Task) OnInit() {
	return
}
func (tk *Task) OnPreStart() error {
	return nil
}
func (tk *Task) OnStarting() {
	tk.FastClient.CreateCookieJar()
	tk.InitializeSession()
	tk.AwaitMonitor()
	tk.Flow()
}
func (tk *Task) OnPause() error {
	return nil
}
func (tk *Task) OnStopping() {
	tk.FastClient.Destroy()
	//panic("")
	return
}

func (tk *Task) Flow() {
	funcarr := []func(){
		//wait queue
		tk.ATC,
		tk.CartAuthenticate,
		tk.SubmitShipping,
		tk.SubmitBilling,
		tk.AdyenConfirm,
	}

	for _, f := range funcarr {
		select {
		case <-tk.Ctx.Done():
			return
		default:
			f()
		}
	}
}
