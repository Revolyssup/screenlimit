package policy

import (
	"fmt"
)

type ActionType string

const Default = "failed attempt"
const (
	RESTART  ActionType = "ActionRestart"
	SHUTDOWN ActionType = "ActionShutdown"
	Log      ActionType = "ActionLog"
)

type PolicyRequest struct {
	Action string     `json:"event" gorm:"event"`
	Type   ActionType `json:"action" gorm:"action"`
}

func (p *PolicyRequest) Run(ch <-chan PolicyRequest) {
	for {
		select {
		case req := <-ch:
			req.Type.Exec()
		}

	}
}

//All the below actions will implement the Action interface
type Restart struct{}

func (r ActionType) Exec() {
	fmt.Println(r.String())
	switch r {
	case RESTART:
		// err := exec.Command("reboot").Run()
		// if err != nil {
		// 	fmt.Println("Could not reboot ", err.Error())
		// 	return
		// }
		panic("rebooted")
	case Log:

	}

}

func (r ActionType) String() string {
	return string(r)
}
