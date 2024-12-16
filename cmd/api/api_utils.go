package api

import (
	"fmt"
	cmn "inventory_money_tracking_software/cmd/common"
	mdls "inventory_money_tracking_software/cmd/models"
	"net/http"
	"os"
	"strings"
)

const StdErrMsg = "Something went wrong, try again please."

func logReq(r *http.Request) {
	cmn.Log(fmt.Sprintf("http %s request on %s", r.Method, r.URL), cmn.LogLevels.Info)
	return
}

// we only allow certain methods on each endpoint. For example if we expect GET request
// and we receive a POST, then this function will return false. This is redundent with
// the CORS header "allow-methods" applied in the func configureHttpAccessControlHeaders.
// We do redundant http method control
// in case the client is not a browser that follows the CORS policy.
func ensureExpectedMethods(reqMethod string, w http.ResponseWriter, expectedHttpMethods ...string) bool {
	for _, expectedMethod := range expectedHttpMethods {
		if reqMethod == expectedMethod {
			return true
		}
	}
	http.Error(w, "method is not supported", http.StatusNotFound)
	return false
}

// return the frontend server address.
func getFrontendOrigin() string {
	return os.Getenv("FRONTEND")
}

func ValidateHttpRequestSecurityRules(w http.ResponseWriter, r *http.Request, allowedHttpMethods, allowedHttpHeaders []string) bool {
	configureHttpAccessControlHeaders(w, allowedHttpMethods, allowedHttpHeaders)
	return ensureExpectedMethods(r.Method, w, allowedHttpMethods...)
}

/*
The browser trusts this server to tell it if the frontend scripts sent from the frontend server are allowed to
read content that this backend server is sending, so we must allow the frontend origin to read this server content(CORS).
Not only that we must also define the allowed http methods and allowed headers. This AccessControl headers
are not usually required for GET requests. If we don't allow a method, header, or origin then the browser will not
allow the javascript to read the contnet that we send it following the CORS policy
*/
func configureHttpAccessControlHeaders(w http.ResponseWriter, allowedMethods []string, allowedHeaders []string) {
	// allow CORS
	w.Header().Set("Access-Control-Allow-Origin", getFrontendOrigin())
	//allow http methods
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods[:], ","))
	// allow http headers
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders[:], ","))
	return //force return out of the function right here, just in case there is code injection or memory manupulation to append code
}

func SanitizeRequest(req *http.Request) *http.Request {
	/* make sure the request body doesn't contain malicous code injection, or illegal chars */
	// todo
	// create a black list of unallowed chars like && ||
	return req
}

// handles hte case when the client made a bad or malformed request
func HandleBadRequest(w http.ResponseWriter) {
	http.Error(w, "Bad Request", http.StatusBadRequest)
	// closeRequestConnectionForcefully(w)
}
func HandleRequestError(w http.ResponseWriter, errCode cmn.ErrorCode, description string) {
	err := mdls.ApiError{Code: errCode, Description: description}
	errStr := fmt.Sprintf("code:%s, description:%s", err.Code, err.Description)
	http.Error(w, errStr, http.StatusNotFound)
}

// this will close active connection without exiting the server and without extra logging.
// when the server closes the connection all of a sudden the client might get errors
// so you shouldn't panic like this unless you really need to. for now this is commented.
// the prefered way to ending requests is to just return from the function
func closeRequestConnectionForcefully(w http.ResponseWriter) {
	cmn.Log("ending http connection", cmn.LogLevels.Info)
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	panic(http.ErrAbortHandler)
}

// this function should run on every api request. It has the logging, monitoring, and authentication
func BeforeEveryRequest(r *http.Request) {
	logReq(r)
}

// takes in something like /accounts/1 and returns 1
// it grabs the suffix id in the path
func ExtractSuffixId(url string) string {
	path := strings.Split(url, "/")
	if len(path) < 2 {
		return "" // nothing was found...
	}
	return path[len(path)-1]
}
