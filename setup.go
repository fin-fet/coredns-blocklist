package blocklist

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/miekg/dns"
)

type BlocklistOptions struct {
	Url             string
	ReloadPeriod    time.Duration
	SourceType      int
	ResponseType    string
	ResponseExtra   string
	MatchSubdomains bool
}

const (
	SourceTypeFile = 0
	SourceTypeHttp = 1
)

var log = clog.NewWithPlugin("blocklist")

func init() {
	plugin.Register("blocklist", setup)
}

func setup(c *caddy.Controller) error {
	options, err := parseArguments(c)
	if err != nil {
		return err
	}

	config := dnsserver.GetConfig(c)

	if strings.HasPrefix(options.Url, "http") {
		options.SourceType = SourceTypeHttp
	} else {
		options.SourceType = SourceTypeFile

		// Convert path
		if !filepath.IsAbs(options.Url) && config.Root != "" {
			options.Url = filepath.Join(config.Root, options.Url)
		}

		// Check to see if file exists
		stat, err := os.Stat(options.Url)
		if err != nil {
			if os.IsNotExist(err) {
				return plugin.Error("blocklist", fmt.Errorf("blocklist file '%s' does not exist", options.Url))
			} else {
				return plugin.Error("blocklist", fmt.Errorf("error opening blocklist file '%s': '%v", options.Url, err))
			}
		}

		// Doublecheck its not a directory
		if stat != nil && stat.IsDir() {
			return plugin.Error("blocklist", fmt.Errorf("blocklist file '%s' is a directory", options.Url))
		}
	}

	blp := BlocklistPlugin{options: &options}

	blp.blockResponse = buildResponse(options)
	blp.loadDomains()

	reloadLoop := reloadLoop(&blp)

	c.OnFinalShutdown(func() error {
		close(reloadLoop)
		return nil
	})

	config.AddPlugin(func(next plugin.Handler) plugin.Handler {
		blp.nextPlugin = next
		return &blp
	})

	return nil
}

func reloadLoop(blp *BlocklistPlugin) chan bool {
	runLoop := make(chan bool)

	if blp.options.ReloadPeriod == 0 {
		return runLoop
	}

	go func() {
		ticker := time.NewTicker(blp.options.ReloadPeriod)
		for {
			select {
			case <-runLoop:
				log.Info("Tick die")
				ticker.Stop()
				return
			case <-ticker.C:
				blp.loadDomains()
			}
		}
	}()

	return runLoop
}

func parseArguments(c *caddy.Controller) (BlocklistOptions, error) {
	options := BlocklistOptions{ResponseType: "nxdomain", ResponseExtra: "", MatchSubdomains: true}

	for c.Next() {
		c.Args(&options.Url)

		if options.Url == "" {
			return options, plugin.Error("blocklist", errors.New("missing file path or URL"))
		}

		for c.NextBlock() {
			option := c.Val()

			switch option {
			case "reload":
				if !c.NextArg() {
					return options, plugin.Error("blocklist", errors.New("reload requires a value"))
				}

				duration, err := time.ParseDuration(c.Val())
				if err != nil {
					return options, plugin.Error("blocklist", fmt.Errorf("unable to parse duration '%s'", c.Val()))
				}

				options.ReloadPeriod = duration
			case "response":
				if !c.NextArg() {
					return options, plugin.Error("blocklist", errors.New("response requires a value"))
				}

				options.ResponseType = c.Val()

				if c.NextArg() {
					options.ResponseExtra = c.Val()
				}
			case "match_subdomains":
				if !c.NextArg() {
					return options, plugin.Error("blocklist", errors.New("match_subdomains requires a value"))
				}

				parsed, err := strconv.ParseBool(c.Val())
				if err != nil {
					return options, plugin.Error("blocklist", fmt.Errorf("invalid option '%s' for match_subdomains, must be true or false", c.Val()))
				}

				options.MatchSubdomains = parsed
			}
		}
	}

	return options, nil
}

func buildResponse(options BlocklistOptions) BlockResponse {
	switch options.ResponseType {
	case "nxdomain":
		return NewStandardResponse(dns.RcodeNameError)
	case "refused":
		return NewStandardResponse(dns.RcodeRefused)
	case "other":
		return NewExtendedResponse(dns.ExtendedErrorCodeOther, options.ResponseExtra)
	case "blocked":
		return NewExtendedResponse(dns.ExtendedErrorCodeBlocked, options.ResponseExtra)
	case "censored":
		return NewExtendedResponse(dns.ExtendedErrorCodeCensored, options.ResponseExtra)
	case "filtered":
		return NewExtendedResponse(dns.ExtendedErrorCodeFiltered, options.ResponseExtra)
	case "prohibited":
		return NewExtendedResponse(dns.ExtendedErrorCodeProhibited, options.ResponseExtra)
	}

	return nil
}
