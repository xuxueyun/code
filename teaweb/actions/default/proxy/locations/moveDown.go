package locations

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/global"
)

type MoveDownAction actions.Action

func (this *MoveDownAction) Run(params struct {
	Filename string
	Index    int
}) {
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	if params.Index >= 0 && params.Index < len(proxy.Locations)-1 {
		next := proxy.Locations[params.Index+1]
		current := proxy.Locations[params.Index]
		proxy.Locations[params.Index+1] = current
		proxy.Locations[params.Index] = next
	}

	proxy.Save()

	global.NotifyChange()

	this.Refresh().Success()
}
