package router

import (
	"context"

	"github.com/yuzuki616/xray-core/app/observatory"
	"github.com/yuzuki616/xray-core/common"
	"github.com/yuzuki616/xray-core/core"
	"github.com/yuzuki616/xray-core/features/extension"
)

type LeastPingStrategy struct {
	ctx         context.Context
	observatory extension.Observatory
}

func (l *LeastPingStrategy) InjectContext(ctx context.Context) {
	l.ctx = ctx
}

func (l *LeastPingStrategy) PickOutbound(strings []string) string {
	if l.observatory == nil {
		common.Must(core.RequireFeatures(l.ctx, func(observatory extension.Observatory) error {
			l.observatory = observatory
			return nil
		}))
	}

	observeReport, err := l.observatory.GetObservation(l.ctx)
	if err != nil {
		newError("cannot get observe report").Base(err).WriteToLog()
		return ""
	}
	outboundsList := outboundList(strings)
	if result, ok := observeReport.(*observatory.ObservationResult); ok {
		status := result.Status
		leastPing := int64(99999999)
		selectedOutboundName := ""
		for _, v := range status {
			if outboundsList.contains(v.OutboundTag) && v.Alive && v.Delay < leastPing {
				selectedOutboundName = v.OutboundTag
				leastPing = v.Delay
			}
		}
		return selectedOutboundName
	}

	//No way to understand observeReport
	return ""
}

type outboundList []string

func (o outboundList) contains(name string) bool {
	for _, v := range o {
		if v == name {
			return true
		}
	}
	return false
}
