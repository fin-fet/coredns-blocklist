package blocklist

import "github.com/miekg/dns"

type BlockResponse interface {
	MakeResponse(w dns.ResponseWriter, r *dns.Msg) (*dns.Msg, error)
}

// StandardResponse is a response with a standard response code
type StandardResponse struct {
	RCode int
}

// ExtendedResponse is a response that uses an extended response code
// via EDNS0 EDE (RFC 8914)
type ExtendedResponse struct {
	InfoCode  uint16
	ExtraText string
}

func NewStandardResponse(rcode int) *StandardResponse {
	sr := new(StandardResponse)
	sr.RCode = rcode
	return sr
}

func (sr *StandardResponse) MakeResponse(w dns.ResponseWriter, r *dns.Msg) (*dns.Msg, error) {
	resp := new(dns.Msg)
	resp.SetRcode(r, sr.RCode)
	return resp, nil
}

func NewExtendedResponse(infoCode uint16, extraText string) *ExtendedResponse {
	er := new(ExtendedResponse)
	er.InfoCode = infoCode
	er.ExtraText = extraText
	return er
}

func (er *ExtendedResponse) MakeResponse(w dns.ResponseWriter, r *dns.Msg) (*dns.Msg, error) {
	// TODO: Force INET?
	// OPT must have root name (RFC 6891)
	opt := new(dns.OPT)
	opt.Hdr = dns.RR_Header{Name: ".", Rrtype: dns.TypeOPT, Class: dns.ClassINET}

	// Build EDE
	ede := new(dns.EDNS0_EDE)
	ede.InfoCode = er.InfoCode
	if er.ExtraText != "" {
		ede.ExtraText = er.ExtraText
	}
	opt.Option = append(opt.Option, ede)

	// OPT goes into Extra section
	resp := new(dns.Msg)
	resp.Extra = append(resp.Extra, opt)
	resp.SetRcode(r, dns.RcodeRefused)

	return resp, nil
}
