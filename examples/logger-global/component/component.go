package component

import (
	"github.com/pkg/errors"
	logger "github.com/toolboxlogger"
)

type CoolComponent struct {
	log *logger.Logger
}

func NewCoolComponent(log *logger.Logger) *CoolComponent {
	return &CoolComponent{
		log: log,
	}
}

func (cc *CoolComponent) DoSomething() {
	cc.log.Entry().Info("Component is doing something...")
}

func (cc *CoolComponent) CallSomethingElse() error {
	cc.log.Entry().Info("Now the component is calling something else")
	if err := somethingElse(); err != nil {
		return errors.Wrap(err, "CallSomething failed")
	}
	return nil
}

func somethingElse() error {
	return errors.New("Something went wrong")
}
