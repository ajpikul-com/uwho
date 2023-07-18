package uwho

import (
	"net/http"
	"net/url"
)

// TODO: Probably need a function to collect errors and inform user/other people

type coordinator struct {
	identifiers     []Identifier
	sessionManager  SessionManager
	desiredResource http.Handler
	loginResult     http.Handler
	accessDenied    http.Handler
	logoutResult    http.Handler
	expiredResult   http.Handler
	loginEndpoint   *url.URL
	logoutEndpoint  *url.URL
	hooks           struct {
		loggedOut   []Hook
		loggedIn    []Hook
		authorized  []Hook
		aboutToLoad []Hook
	}
	applicationUserinfo AppUserinfo
}

func (c *coordinator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userStatus := NewUserStatus()
	defaultLogger.Info("Serving HTTP from " + r.URL.Path)

	// Check to see if user loggedout
	if c.checkLogout(w, r) {
		defaultLogger.Info(r.URL.Path + ": We're about to logout")
		c.sessionManager.EndSession(w, r)
		c.applicationUserinfo.LogOut()
		userStatus.ReconcileStatus(LOGGEDOUT)
		c.CallHooks(c.hooks.loggedOut, w, r)
		c.logoutResult.ServeHTTP(w, r)
		return
	}

	// Try to read the session
	data, opinion := c.sessionManager.ReadSession(w, r)
	defaultLogger.Info(r.URL.Path + ": Just read session")
	defaultLogger.Info("Session data: " + data)
	userStatus.ReconcileStatus(opinion)
	defaultLogger.Info(userStatus.StatusStr())

	// Found a session
	if userStatus.IsStatus(KNOWN) {
		defaultLogger.Info(r.URL.Path + ": KNOWN, attempting to read data and authorize user")
		// Store cookie data in user structure
		c.applicationUserinfo.SessionDestring(data)

		// Use stored data to try and authorize user
		userStatus.ReconcileStatus(c.applicationUserinfo.AuthorizeUser(w, r))
		defaultLogger.Info(userStatus.StatusStr())
	}

	// User is authorized
	if userStatus.IsStatus(AUTHORIZED) {
		defaultLogger.Info(r.URL.Path + ": We are freshly authorized")
		defaultLogger.Info(userStatus.StatusStr())
		c.CallHooks(c.hooks.authorized, w, r)
		c.CallHooks(c.hooks.aboutToLoad, w, r)
		// If we want to login again (ie multiple logins), should we hijack here?
		c.desiredResource.ServeHTTP(w, r)
		return
	}

	// See if we're trying to login
	if c.checkLogin(w, r) {
		defaultLogger.Info(r.URL.Path + ": checkLogin returned true")
		userStatus.ReconcileStatus(KNOWN)
		if c.applicationUserinfo.InitSession() == ErrSessionExists {
			defaultLogger.Info("Starting second session? Not possible right now.")
		} else {
			c.sessionManager.MarkSession(c.applicationUserinfo.SessionString(), w, r)
			c.CallHooks(c.hooks.loggedIn, w, r)
			c.CallHooks(c.hooks.aboutToLoad, w, r)
		}
		c.loginResult.ServeHTTP(w, r)
		return
	}

	// See if we're logged out or expired
	if userStatus.IsStatus(LOGGEDOUT) || userStatus.IsStatus(EXPIRED) {
		defaultLogger.Info(r.URL.Path + ": We logged out..")
		c.sessionManager.EndSession(w, r)
		c.applicationUserinfo.LogOut()
		defaultLogger.Info(userStatus.StatusStr())
		c.CallHooks(c.hooks.loggedOut, w, r)
		c.CallHooks(c.hooks.aboutToLoad, w, r)
		if userStatus.IsStatus(EXPIRED) {
			c.expiredResult.ServeHTTP(w, r)
		} else {
			c.logoutResult.ServeHTTP(w, r)
		}
		return
	}

	// Not authorized, known or unkown
	defaultLogger.Info(r.URL.Path + " but " + userStatus.StatusStr() + " so DENIED")
	c.CallHooks(c.hooks.aboutToLoad, w, r)
	c.accessDenied.ServeHTTP(w, r)
}

func (c *coordinator) checkLogin(w http.ResponseWriter, r *http.Request) bool {
	loggedIn := false
	if r.URL.Path == c.loginEndpoint.Path {
		defaultLogger.Info("Equal paths")
		defaultLogger.Info(r.URL.Path)
		defaultLogger.Info(c.loginEndpoint.Path)
		for _, identifier := range c.identifiers {
			opinion := identifier.VerifyCredentials(w, r)
			if opinion == KNOWN {
				loggedIn = true
				defaultLogger.Info("Found a user.")
			} else if opinion == SPOKEN {
				// TODO: Identifier trying to hijack whole process
			} else if opinion != UNKNOWN {
				defaultLogger.Error("An identifier is returning a strange user status: " + string(opinion))
			}
		}
	}
	return loggedIn
}

func (c *coordinator) checkLogout(w http.ResponseWriter, r *http.Request) bool {
	if r.URL.Path == c.logoutEndpoint.Path { // I want to do URL comparisons TODO
		return true
	}
	return false
}