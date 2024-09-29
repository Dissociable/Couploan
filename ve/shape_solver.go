package ve

import "github.com/pkg/errors"

func (ve *VE) SolveShape() (err error) {
	if ve.proxy == nil || ve.proxy.IsDirect() {
		err = errors.New("no proxy is set")
		return
	}
	return
}
