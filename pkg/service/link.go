package service

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"log"
	"net/url"
	"strings"
	"trendyolcase/pkg/repository/link"
)

type ConverterService struct {
	ConverterRepository *link.Repository
}

func NewConverterService(l *link.Repository) ConverterService {
	return ConverterService{ConverterRepository: l}
}

func (l *ConverterService) GetDeepLinkIfWebURLExist(webURL string) (string, error) {
	return l.ConverterRepository.GetDeepLinkIfWebURLExist(webURL)
}

func (l *ConverterService) GetWebURLIfDeepLinkExist(deepLink string) (string, error) {
	return l.ConverterRepository.GetWebURLIfDeepLinkExist(deepLink)
}

func (l *ConverterService) Insert(webURL string, deepLink string) bool {
	return l.ConverterRepository.Insert(webURL, deepLink)
}

func (l *ConverterService) InsertLog(logInformation string) bool {
	return l.ConverterRepository.InsertLog(logInformation)
}

/*
Deeplink doesn't exist in db so create a deeplink.
First check is the request URL about product detail page, search page or other page.
After that return response created deeplink.
If request is not url, bad request is returned. In all other faulty cases, the other page is returned.
*/

func (l *ConverterService) CreateDeepLink(requestLink string) (string, error) {

	const (
		trendyolHomePage     = "https://www.trendyol.com"
		trendyolSearchPage   = "https://www.trendyol.com/sr?q="
		deepLinkHomePage     = "ty://?Page=Home"
		productPageSeperator = "-p-"
	)

	if !(govalidator.IsURL(requestLink)) {
		return "", errors.New("There is an error in the requested data. Check the data. Tag should be 'weburl' and links doesn't contain space")
	}

	if strings.HasPrefix(requestLink, trendyolHomePage) && strings.Contains(requestLink, productPageSeperator) {
		responseDeepLink := ConvertProductDetailPageToDeepLink(requestLink)
		return responseDeepLink, nil
	} else if strings.HasPrefix(requestLink, trendyolSearchPage) {
		responseDeepLink := ConvertSearchPageToDeepLink(requestLink)
		return responseDeepLink, nil
	} else {
		responseDeepLink := ConvertOtherPageToDeepLink(deepLinkHomePage)
		return responseDeepLink, nil
	}
}

/*
Converts to product page URL to product page deeplink.
*/

func ConvertProductDetailPageToDeepLink(requestLink string) string {
	const (
		productPageDeepLinkBase = "ty://?Page=Product&ContentId="
		deepLinkHomePage        = "ty://?Page=Home"
		boutiqueIDSeperator     = "boutiqueId"
		merchantIDSeparator     = "merchantId"
	)
	query := strings.Split(requestLink, "-p-")[1]

	u, _ := url.Parse(requestLink)
	q, _ := url.ParseQuery(u.RawQuery)

	boutiqueID := q.Get(boutiqueIDSeperator)
	merchantID := q.Get(merchantIDSeparator)

	boutiqueIDFlag := q.Has(boutiqueIDSeperator)
	merchantIDFlag := q.Has(merchantIDSeparator)
	badRequestWithCampaignID := boutiqueIDFlag && boutiqueID == ""
	badRequestWithMerchantID := merchantIDFlag && merchantID == ""

	if badRequestWithCampaignID || badRequestWithMerchantID {
		return deepLinkHomePage
	}

	idx := strings.Index(query, "?")

	if idx != -1 {
		contentID := query[:idx]
		responseDeepLink := productPageDeepLinkBase + contentID
		if contentID == "" {
			return deepLinkHomePage
		}
		if boutiqueID != "" && merchantID != "" {
			responseDeepLink := responseDeepLink + "&CampaignId=" + boutiqueID + "&MerchantId=" + merchantID
			return responseDeepLink
		} else if merchantID != "" {
			responseDeepLink := responseDeepLink + "&MerchantId=" + merchantID
			return responseDeepLink
		} else if boutiqueID != "" {
			responseDeepLink := responseDeepLink + "&CampaignId=" + boutiqueID
			return responseDeepLink
		} else {
			responseDeepLink := deepLinkHomePage
			return responseDeepLink
		}
	} else {
		responseDeepLink := productPageDeepLinkBase + query
		return responseDeepLink
	}
}

/*
Converts to search page URL to search page deeplink. If query is empty returns homepage.
*/

func ConvertSearchPageToDeepLink(requestLink string) string {
	const (
		deepLinkBaseSearch     = "ty://?Page=Search&Query="
		querySeparator         = "q="
		specificTurkishLetters = "çÇğĞıİöÖşŞüÜ"
		deepLinkHomePage       = "ty://?Page=Home"
	)
	query := strings.Split(requestLink, querySeparator)[1]

	if query == "" || strings.ContainsAny(query, "?&/") {
		return deepLinkHomePage
	}

	if strings.ContainsAny(query, specificTurkishLetters) {
		modifiedQuery := url.QueryEscape(query)
		responseDeepLink := deepLinkBaseSearch + modifiedQuery
		return responseDeepLink
	}
	responseDeepLink := deepLinkBaseSearch + query
	return responseDeepLink
}

/*
Converts to all other urls like 'https://www.trendyol.com/Hesabim/Favoriler' to 'ty://?Page=Home'.
*/

func ConvertOtherPageToDeepLink(deepLinkHomePage string) string {
	return deepLinkHomePage
}

/*
webURL doesn't exist in db so create a webURL.
First check is the deeplink about product detail page, search page or other page.
After that return response created webURL.
*/

func (l *ConverterService) CreateWebURL(requestLink string) (string, error) {
	const (
		trendyolHomePageURL     = "https://www.trendyol.com"
		searchPageDeepLinkBase  = "ty://?Page=Search&Query="
		productPageDeepLinkBase = "ty://?Page=Product&ContentId="
	)

	if requestLink == "" || strings.Contains(requestLink, " ") {
		return "", errors.New("There is an error in the requested data. Check the data. Tag should be 'deeplink' and links doesn't contain space")
	}

	if strings.HasPrefix(requestLink, searchPageDeepLinkBase) {
		responseWebURL := ConvertSearchPageToURL(requestLink)
		return responseWebURL, nil
	} else if strings.HasPrefix(requestLink, productPageDeepLinkBase) {
		responseWebURL := ConvertProductDetailPageToURL(requestLink)
		return responseWebURL, nil
	} else {
		responseWebURL := ConvertOtherPageToURL(trendyolHomePageURL)
		return responseWebURL, nil
	}
}

/*
Converts to all other deeplinks like 'ty://?Page=Favorites' to 'www.trendyol.com'.
*/

func ConvertOtherPageToURL(trendyolHomePage string) string {
	return trendyolHomePage
}

/*
Converts deeplink search page to search page URL.
*/

func ConvertSearchPageToURL(requestLink string) string {

	const (
		baseSearchWebURL = "https://www.trendyol.com/sr?q="
		querySeparator   = "Query"
	)
	u, _ := url.Parse(requestLink)
	q, _ := url.ParseQuery(u.RawQuery)

	query := q.Get(querySeparator)
	query = url.QueryEscape(query)

	if !(strings.Contains(requestLink, "ty://?Page=Search&Query")) || query == "" || strings.Contains(query, "?") {
		return "https://www.trendyol.com"
	}
	responseWebURL := baseSearchWebURL + query
	return responseWebURL
}

/*
Converts deeplink product page to product page URL.
*/

func ConvertProductDetailPageToURL(requestLink string) string {
	const (
		baseWebURL          = "https://www.trendyol.com/brand/name-p-"
		campaignIDSeperator = "CampaignId"
		merchantIDSeparator = "MerchantId"
		contentIDSeparator  = "ContentId"
	)

	u, err := url.Parse(requestLink)
	if err != nil {
		log.Printf("%s", err)
	}
	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		log.Printf("%s", err)
	}

	campaignID := q.Get(campaignIDSeperator)
	merchantID := q.Get(merchantIDSeparator)
	contentID := q.Get(contentIDSeparator)

	responseWebURL := baseWebURL + contentID

	campaignIDFlag := q.Has(campaignIDSeperator)
	merchantIDFlag := q.Has(merchantIDSeparator)
	contentIDFlag := q.Has(contentIDSeparator)

	badRequestWithCampaignID := campaignIDFlag && campaignID == ""
	badRequestWithMerchantID := merchantIDFlag && merchantID == ""
	badRequestWithContentID := contentIDFlag && contentID == ""

	if badRequestWithCampaignID || badRequestWithMerchantID || badRequestWithContentID {
		return "https://www.trendyol.com"
	}

	if campaignIDFlag && merchantIDFlag {
		responseDeepLink := responseWebURL + "?boutiqueId=" + campaignID + "&merchantId=" + merchantID
		return responseDeepLink
	} else if merchantIDFlag {
		responseDeepLink := responseWebURL + "?merchantId=" + merchantID
		return responseDeepLink
	} else if campaignIDFlag {
		responseDeepLink := responseWebURL + "?boutiqueId=" + campaignID
		return responseDeepLink
	} else if strings.ContainsAny(contentID, "&=/") {
		return "https://www.trendyol.com"
	} else {
		return responseWebURL
	}
}
