package tag

import (
	"encoding/json"
	"inventory_money_tracking_software/cmd/api"
	cmn "inventory_money_tracking_software/cmd/common"
	mdls "inventory_money_tracking_software/cmd/models"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func findTagByIdHandler(w http.ResponseWriter, r *http.Request) {

	var allowedMethods = []string{http.MethodGet, http.MethodOptions}
	allowedHeaders := []string{"Content-Type"}
	if !api.ValidateHttpRequestSecurityRules(w, r,
		allowedMethods, allowedHeaders) {
		api.HandleBadRequest(w)
		return
	}
	path := strings.Trim(r.URL.Path, "")
	idString := api.ExtractSuffixId(path)
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "malformed request")
		return
	}
	tag, errCode := ReadTagById(idInt)
	if errCode != cmn.ErrorCodes.NoError {
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "Something went wrong, try again please.")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jsonResp, err := json.Marshal(tag)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "Something went wrong, try again please.")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
	return
}
func createTagHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var allowedMethods = []string{http.MethodPost, http.MethodOptions}
	allowedHeaders := []string{"Content-Type"}
	if !api.ValidateHttpRequestSecurityRules(w, r,
		allowedMethods, allowedHeaders) {
		api.HandleBadRequest(w)
		return
	}
	if r.Method == http.MethodPost {
		w.Header().Set("Content-Type", "text/plain")
		tag := mdls.Tag{}
		err := json.NewDecoder(r.Body).Decode(&tag)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			api.HandleRequestError(w, cmn.ErrorCodes.CreateFailed, "")
			return
		}
		tId, errorCode := CreateTag(tag)
		if errorCode != cmn.ErrorCodes.NoError {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			api.HandleRequestError(w, errorCode, "")
			return
		}
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, strconv.Itoa(tId))
	}
	return
}
func RegisterHandlers() {
	mux := api.GetHttpMux()
	mux.HandleFunc(tagRoutes.create, createTagHandler)
	mux.HandleFunc(tagRoutes.read, findTagByIdHandler)
}
