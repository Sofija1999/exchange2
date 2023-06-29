package log

import (
	"strings"
	"time"

	"github.com/Bloxico/exchange-gateway/sofija/server/auth"
	"github.com/Bloxico/exchange-gateway/sofija/server/params"
	"github.com/emicklei/go-restful/v3"
)

// Returns a filter function that outputs request logs according to the NCSA standard
// The passed logger will be used to output the message (this is likely the app logger itself - zap_logger)
func NCSACommonLogFormatLogger(logger Logger) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {

		// sanity check: nil during tests
		if logger == nil {
			return
		}

		var username = "-"

		// extract user making the request from the context, if possible
		reqId, err := params.StringFrom(req.Request, auth.USER_EMAIL_CTX_KEY)
		if reqId != "" && err == nil {
			username = reqId
		}

		chain.ProcessFilter(req, resp)

		logger.Infof("%s - %s [%s] \"%s %s %s\" %d %d",
			strings.Split(req.Request.RemoteAddr, ":")[0],
			username,
			time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			"req.Request.Method",
			req.Request.URL.RequestURI(),
			req.Request.Proto,
			resp.StatusCode(),
			resp.ContentLength(),
		)
	}
}
