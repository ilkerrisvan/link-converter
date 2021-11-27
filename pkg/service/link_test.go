package service

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"trendyolcase/pkg/repository/link"
)

/*
Initialize db for use in tests.
*/
func dbInıt() *sql.DB {
	var DB *sql.DB
	connectionString := fmt.Sprintf(os.Getenv("CONNECTION_STRING"))
	DB, _ = sql.Open("postgres", connectionString)
	return DB
}

func TestCreateWebURL(t *testing.T) {
	assert := assert.New(t)
	db := dbInıt()
	converterRepositoryTest := link.NewRepository(db)
	c := NewConverterService(converterRepositoryTest)

	deepLinks := []struct {
		testDeepLinkRequest string
		expectedURL         string
	}{

		{"ty://?Page= TestWithSpaceChar", ""},
		{"ty://?Page=Favorites", "https://www.trendyol.com"},
		{"ty://?Page=Siparişlerim", "https://www.trendyol.com"},
		{"ty://?Page=", "https://www.trendyol.com"},
		{"ty://?Page?", "https://www.trendyol.com"},
		{"ty://?Page?orders", "https://www.trendyol.com"},
		{"ty://?Page=orders&Search=%C3%BCt%C3%BC", "https://www.trendyol.com"},
		{"ty://?Page?orders&Search=%C3%BCt%C3%BC", "https://www.trendyol.com"},
		{"ty://?Page?orders?Search=%C3%BCt%C3%BC", "https://www.trendyol.com"},
		{"ty://?Page?ordersSearch=%C3%BCt%C3%BC", "https://www.trendyol.com"},
		{"ty://?Pageorders", "https://www.trendyol.com"},

		{"ty://?Page=Search&Query= %C3%BCt%C3%BC", ""},
		{"ty://?Page=Search&Query=elbise", "https://www.trendyol.com/sr?q=elbise"},
		{"ty://?Page=Search&Query=süpürge", "https://www.trendyol.com/sr?q=s%C3%BCp%C3%BCrge"},
		{"ty://?Page=Search&Query=%C3%BCt%C3%BC", "https://www.trendyol.com/sr?q=%C3%BCt%C3%BC"},
		{"ty://?Page=Search&Query=%C3%BCt%C3%BC&testdata", "https://www.trendyol.com/sr?q=%C3%BCt%C3%BC"},
		{"ty://=Page=Search&Query=%C3%BCt%C3%BC", "https://www.trendyol.com"},
		{"ty://?Page=Search&Query=", "https://www.trendyol.com"},
		{"ty//?Search=%C3%BCt%C3%BC", "https://www.trendyol.com"},
		{"ty//?Query=%C3%BCt%C3%BC", "https://www.trendyol.com"},
		{"ty//&Query=%C3%BCt%C3%BC", "https://www.trendyol.com"},
		{"ty://?Page=Search?Query=test", "https://www.trendyol.com"},
		{"ty://?Page=&Query=test", "https://www.trendyol.com"},
		{"ty://?Page?&Query=test", "https://www.trendyol.com"},
		{"ty://?Page?Search&Query=test", "https://www.trendyol.com"},
		{"ty://?Page&Search&Query=test", "https://www.trendyol.com"},
		{"ty://&Page&Search&Query=test", "https://www.trendyol.com"},
		{"ty://?Page&Search?Query=test", "https://www.trendyol.com"},
		{"ty://?Page?Search&Query=test", "https://www.trendyol.com"},
		{"ty://?Page=Query=%C3%BCt%C3%BC", "https://www.trendyol.com"},

		{"ty://?Page=Product&ContentId=1925865& MerchantId=105064", ""},
		{"ty://?Page=Product&ContentId=1&CampaignId=12&MerchantId=123", "https://www.trendyol.com/brand/name-p-1?boutiqueId=12&merchantId=123"},
		{"ty://?Page=Product&ContentId=12", "https://www.trendyol.com/brand/name-p-12"},
		{"ty://?Page=Product&ContentId=123&CampaignId=1234", "https://www.trendyol.com/brand/name-p-123?boutiqueId=1234"},
		{"ty://?Page=Product&ContentId=1234&MerchantId=12345", "https://www.trendyol.com/brand/name-p-1234?merchantId=12345"},
		{"ty://?Page=Product&ContentId=12345&campaignId=123", "https://www.trendyol.com/brand/name-p-12345"},
		{"ty://?Page=Product&ContentId=678&merchantId=123", "https://www.trendyol.com/brand/name-p-678"},
		{"ty://?Page=Product&ContentId=&CampaignId=000", "https://www.trendyol.com"},
		{"ty://?Page=Product&ContentId=", "https://www.trendyol.com"},
		{"ty://?Page=Product&ContentId=678&MerchantId=", "https://www.trendyol.com"},
	}
	for _, testDeepLinks := range deepLinks {
		actualWebURL, _:= c.CreateWebURL(testDeepLinks.testDeepLinkRequest)
		assert.Equal(testDeepLinks.expectedURL, actualWebURL, "Should be %s", testDeepLinks.expectedURL)
	}
}
func TestConvertSearchPageToURL(t *testing.T) {

	assert := assert.New(t)

	webURLs := []struct {
		testDeepLinkRequest string
		expectedURL         string
	}{
		{"ty://?Page=Search&Query=elbisetestdata", "https://www.trendyol.com/sr?q=elbisetestdata"},
		{"ty://?Page=Search&Query=süpürgetestdata", "https://www.trendyol.com/sr?q=s%C3%BCp%C3%BCrgetestdata"},
		{"ty://?Page=Search&Query=%C3%BCt%C3%BCtestdata", "https://www.trendyol.com/sr?q=%C3%BCt%C3%BCtestdata"},
		{"ty://?Page=Search&Query=testdata99", "https://www.trendyol.com/sr?q=testdata99"},
	}
	for _, testDeeplinks := range webURLs {
		testDeepLinkRequest := ConvertSearchPageToURL(testDeeplinks.testDeepLinkRequest)
		assert.Equal(testDeeplinks.expectedURL, testDeepLinkRequest, "Should be %s", testDeeplinks.expectedURL)
	}
}
func TestConvertProductDetailPageToURL(t *testing.T) {
	assert := assert.New(t)

	// BU DATALAR DATABASEDE OLMAMALI
	// testleri arttır.
	deepLinks := []struct {
		testDeepLinkRequest string
		expectedURL         string
	}{
		{"ty://?Page=Product&ContentId=1&CampaignId=1&MerchantId=1", "https://www.trendyol.com/brand/name-p-1?boutiqueId=1&merchantId=1"},
		{"ty://?Page=Product&ContentId=1&MerchantId=2", "https://www.trendyol.com/brand/name-p-1?merchantId=2"},
		{"ty://?Page=Product&ContentId=1&CampaignId=2", "https://www.trendyol.com/brand/name-p-1?boutiqueId=2"},
		{"ty://?Page=Product&ContentId=123", "https://www.trendyol.com/brand/name-p-123"},
		{"ty://?Page=Product&ContentId=222&merchantId=333", "https://www.trendyol.com/brand/name-p-222"},
		{"ty://?Page=Product&ContentId=666&campaignId=777", "https://www.trendyol.com/brand/name-p-666"},
		{"ty://?Page=Product&ContentId=1&campaignId=12&merchantId=123", "https://www.trendyol.com/brand/name-p-1"},
		{"ty://?Page=Product&ContentId=123&CampaignId=", "https://www.trendyol.com"},
		{"ty://?Page=Product&ContentId=345MerchantId=567", "https://www.trendyol.com"},
		{"ty://?Page=Product&ContentId=100?MerchantId=111", "https://www.trendyol.com"},
		{"ty://?Page=Product&ContentId=444?CampaignId=555", "https://www.trendyol.com"},
	}
	for _, testDeepLinks := range deepLinks {
		testDeepLinkRequest := ConvertProductDetailPageToURL(testDeepLinks.testDeepLinkRequest)
		assert.Equal(testDeepLinks.expectedURL, testDeepLinkRequest, "Should be %s", testDeepLinks.expectedURL)

	}
}

func TestCreateDeepLink(t *testing.T) {
	assert := assert.New(t)
	db := dbInıt()

	converterRepository := link.NewRepository(db)
	c := NewConverterService(converterRepository)

	webURLs := []struct {
		testWebURLRequest string
		expectedDeepLink  string
	}{
		{"https://www.trendyol.com/ WithSpace", ""},
		{"https://www.trendyol.com/Hesabim/Favoriler", "ty://?Page=Home"},
		{"https://www.trendyol.com", "ty://?Page=Home"},
		{"https://www.trendyol.com/sr?q=%C3%BCt%C3%BC", "ty://?Page=Search&Query=%C3%BCt%C3%BC"},
		{"https://www.trendyol.com/sr?q=elbise", "ty://?Page=Search&Query=elbise"},
		{"https://www.trendyol.com/sr?=elbise", "ty://?Page=Home"},
		{"https://www.trendyol.com/sr?q=", "ty://?Page=Home"},
		{"https://www.trendyol.com/test/saat-p-0?boutiqueId=0&merchantId=0", "ty://?Page=Product&ContentId=0&CampaignId=0&MerchantId=0"},
		{"https://www.trendyol.com/test/saat-p-1", "ty://?Page=Product&ContentId=1"},
		{"https://www.trendyol.com/test/saat", "ty://?Page=Home"},
		{"https://www.trendyol.com/test/erkek-kol-saati-p-22?boutiqueId=33", "ty://?Page=Product&ContentId=22&CampaignId=33"},
		{"https://www.trendyol.com/test/erkek-kol-saati-p-4444?merchantId=5555", "ty://?Page=Product&ContentId=4444&MerchantId=5555"},
	}
	for _, webURLTest := range webURLs {
		actualDeepLink, _ := c.CreateDeepLink(webURLTest.testWebURLRequest)
		assert.Equal(webURLTest.expectedDeepLink, actualDeepLink, "%s olmalı", webURLTest.expectedDeepLink)

	}
}
func TestConvertSearchPageToDeepLink(t *testing.T) {
	assert := assert.New(t)

	webURLs := []struct {
		testWebURL  string
		expectedURL string
		err         error
	}{
		{"https://www.trendyol.com/sr?q=elbisetestdata", "ty://?Page=Search&Query=elbisetestdata", nil},
		{"https://www.trendyol.com/sr?q=ütütestdata", "ty://?Page=Search&Query=%C3%BCt%C3%BCtestdata", nil},
		{"https://www.trendyol.com/sr?q=süpürge", "ty://?Page=Search&Query=s%C3%BCp%C3%BCrge", nil},
		{"https://www.trendyol.com/sr?q=", "ty://?Page=Home", errors.New("Bad request.")},
		{"https://www.trendyol.com/sr?q=elbise&merchantId=123", "ty://?Page=Home", errors.New("Bad request.")},
	}
	for _, testData := range webURLs {
		actualWebURL := ConvertSearchPageToDeepLink(testData.testWebURL)
		assert.Equal(testData.expectedURL, actualWebURL, "Should be %s", testData.expectedURL)

	}

}
func TestConvertProductDetailPageToDeepLink(t *testing.T) {
	assert := assert.New(t)

	webURLs := []struct {
		testWebURL  string
		expectedURL string
	}{
		{"https://www.trendyol.com/testbrand/saat-p-00000?boutiqueId=00000&merchantId=00000", "ty://?Page=Product&ContentId=00000&CampaignId=00000&MerchantId=00000"},
		{"https://www.trendyol.com/testbrand/saat-p-1", "ty://?Page=Product&ContentId=1"},
		{"https://www.trendyol.com/casio/erkek-kol-saati-p-1925865?boutiqueId=439892", "ty://?Page=Product&ContentId=1925865&CampaignId=439892"},
		{"https://www.trendyol.com/casio/erkek-kol-saati-p-1925865?merchantId=105064", "ty://?Page=Product&ContentId=1925865&MerchantId=105064"},
		{"https://www.trendyol.com/testbrand/erkek-kol-saati-p-4444?MerchantId=5555", "ty://?Page=Home"},
	}
	for _, v := range webURLs {
		actualWebURL := ConvertProductDetailPageToDeepLink(v.testWebURL)
		assert.Equal(v.expectedURL, actualWebURL, "Should be %s", v.expectedURL)
	}

}
