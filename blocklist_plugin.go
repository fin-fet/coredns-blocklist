package blocklist

import (
	"context"
	"time"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type BlocklistPlugin struct {
	options       *BlocklistOptions
	blocklist     Blocklist
	blockResponse BlockResponse

	nextPlugin plugin.Handler

	lastFileTime time.Time
	lastFileSize int64
}

func (bl BlocklistPlugin) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	if bl.shouldBlock(state.Name()) {
		msg, err := bl.blockResponse.MakeResponse(w, r)
		if err != nil {
			log.Errorf("failed to build block response for '%s': %v", state.Name(), err)
		}

		err = w.WriteMsg(msg)
		if err != nil {
			log.Errorf("failed to write DNS message for '%s': %v", state.Name(), err)
		}

		log.Warningf("%s blocked %s %s %s %s (%s)",
			state.RemoteAddr(),
			state.Type(),
			state.Class(),
			state.Name(),
			state.Proto(),
			bl.options.ResponseType)

		return dns.RcodeNameError, nil
	}

	return plugin.NextOrFailure(bl.Name(), bl.nextPlugin, ctx, w, r)
}

func (blp *BlocklistPlugin) Name() string {
	return "blocklist"
}

func (blp BlocklistPlugin) shouldBlock(name string) bool {
	// Don't block localhost
	if name == "localhost." {
		return false
	}

	return blp.blocklist.Contains(name)
}
