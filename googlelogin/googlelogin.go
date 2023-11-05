package googlelogin

import (
	"net/http"
	"regexp"

	"google.golang.org/api/idtoken"

	"github.com/ajpikul-com/sutils"
	"github.com/ajpikul-com/uwho"
)

type ReqByIdent interface {
	AcceptData(map[string]interface{})
}

type GoogleLogin struct {
	ClientID    string
	OriginMatch *regexp.Regexp // TODO this could be a function
}

func New(clientID string, match string) *GoogleLogin {
	return &GoogleLogin{
		ClientID:    clientID,
		OriginMatch: regexp.MustCompile(match),
	}
}

func (g *GoogleLogin) TestInterface(stateManager uwho.ReqByCoord) {
	if _, ok := stateManager.(ReqByIdent); !ok {
		panic("State manager doesn't satisfied required interface")
	}
}

func (g *GoogleLogin) VerifyCredentials(userStateCoord uwho.ReqByCoord, w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		origin := r.Header.Get("Origin")
		if g.OriginMatch.Match([]byte(origin)) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT")
			w.Header().Set("Vary", "Origin")
			w.WriteHeader(http.StatusOK)
			return
		}
		return // We set nothing explicitly
	} else if r.Method == "POST" {
		r.ParseMultipartForm(4096)
		cookieCSRFValue, err := r.Cookie("g_csrf_token")
		if err != nil {
			defaultLogger.Error(err.Error())
			return
		}
		if cookieCSRFValue.Value != r.Form["g_csrf_token"][0] {
			defaultLogger.Info("Under attack? csrf tokens didn't match")
			return
		}

		payload, err := idtoken.Validate(r.Context(), r.Form["credential"][0], "")
		if err != nil {
			defaultLogger.Error(err.Error())
			return
		}

		userState, ok := userStateCoord.(ReqByIdent)
		if !ok {
			defaultLogger.Error("Interface assertion error")
			return
		}
		userState.AcceptData(payload.Claims)
	} else if r.Method == "PUT" {
		// TODO PROCESS JAVASCRIPT HERE
	}
}

func DefaultLoginDiv(loginEndpoint string, clientID string) string {
	return `<div id="g_id_onload"
	data-client_id="` + clientID + `"
	data-context="signin"
	data-ux_mode="popup"
	data-login_uri="` + loginEndpoint + `"
	data-auto_prompt="false"
</div>

<div class="g_id_signin"
	data-type="icon"
	data-shape="circle"
	data-theme="outline"
	data-text="continue_with"
	data-size="large"
</div>`
}

func DefaultLoginPortal(loginEndpoint string, clientID string) http.Handler {
	return sutils.StringHandler{`<!DOCTYPE html>
<html>
	<head>
		<script src="https://accounts.google.com/gsi/client" async></script>
	</head>
	<body>
		` + DefaultLoginDiv(loginEndpoint, clientID) + `
	</body>
</html>`}
}
