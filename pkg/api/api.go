package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"trendyolcase/pkg/model"
	"trendyolcase/pkg/service"
)

type ConverterAPI struct {
	ConverterService service.ConverterService
}

func NewConverterAPI(c service.ConverterService) ConverterAPI {
	return ConverterAPI{ConverterService: c}
}

/* The URL is taken from the incoming request and if there is already
a deeplink for this URL,it is returned as a response,
otherwise a deeplink is created for this URL. */

func (c ConverterAPI) GenerateDeepLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		link := model.Link{}

		body, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(body, &link)
		if err != nil {
			message := "There is an error in the requested data. Check the data. Data should be JSON."
			_ = c.ConverterService.InsertLog(message)
			RespondError(w, http.StatusBadRequest, message)
			return
		}
		requestLink := link.WebUrl

		link.Deeplink, _ = c.ConverterService.GetDeepLinkIfWebURLExist(requestLink)

		if link.Deeplink != "" {
			logMessage := "WebURL= " + requestLink + " exists in db. Response= " + link.Deeplink + " successfully returned with data from db."
			_ = c.ConverterService.InsertLog(logMessage)
			RespondDeepLinkWithJSON(w, http.StatusOK, link.Deeplink)
			return
		} else {
			link.Deeplink, err = c.ConverterService.CreateDeepLink(requestLink)
			if err != nil {
				message := err.Error()
				_ = c.ConverterService.InsertLog(message)
				RespondError(w, http.StatusBadRequest, message)
				return
			}
			logMessage := "Response=" + link.Deeplink + "successfully created and returned as response. Saved to DB."
			_ = c.ConverterService.InsertLog(logMessage)
			RespondDeepLinkWithJSON(w, http.StatusOK, link.Deeplink)
			c.ConverterService.Insert(requestLink, link.Deeplink)
			return
		}
	}
}

/* The deeplink is taken from the incoming request and if there is already
a URL for this deeplink,it is returned as a response,
otherwise, a URL is created for this deeplink. */

func (c ConverterAPI) GenerateWebURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		link := model.Link{}
		body, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(body, &link)
		if err != nil {
			message := "There is an error in the requested data. Check the data. Data should be JSON."
			_ = c.ConverterService.InsertLog(message)
			RespondError(w, http.StatusBadRequest, message)
			return
		}
		requestLink := link.Deeplink

		// Check if there is a URL for the deeplink received from the request
		link.WebUrl, _ = c.ConverterService.GetWebURLIfDeepLinkExist(requestLink)
		if link.WebUrl != "" {
			logMessage := "Deeplink= " + requestLink + " exists in db. Response= " + link.WebUrl + " successfully returned with data from db."
			_ = c.ConverterService.InsertLog(logMessage)
			RespondWebURLWithJSON(w, http.StatusOK, link.WebUrl)
			return
		} else {
			link.WebUrl, err = c.ConverterService.CreateWebURL(requestLink)
			if err != nil {
				message := err.Error()
				_ = c.ConverterService.InsertLog(message)
				RespondError(w, http.StatusBadRequest, message)
				return
			}
			logMessage := "Response= " + link.WebUrl + " successfully created and returned as response. Saved to DB."
			_ = c.ConverterService.InsertLog(logMessage)
			RespondWebURLWithJSON(w, http.StatusOK, link.WebUrl)
			c.ConverterService.Insert(link.WebUrl, requestLink)
			return
		}
	}
}
