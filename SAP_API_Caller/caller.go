package sap_api_caller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	sap_api_output_formatter "sap-api-integrations-credit-management-master-reads/SAP_API_Output_Formatter"
	"strings"
	"sync"

	"github.com/latonaio/golang-logging-library-for-sap/logger"
	"golang.org/x/xerrors"
)

type SAPAPICaller struct {
	baseURL string
	apiKey  string
	log     *logger.Logger
}

func NewSAPAPICaller(baseUrl string, l *logger.Logger) *SAPAPICaller {
	return &SAPAPICaller{
		baseURL: baseUrl,
		apiKey:  GetApiKey(),
		log:     l,
	}
}

func (c *SAPAPICaller) AsyncGetCreditManagementMaster(businessPartner string, accepter []string) {
	wg := &sync.WaitGroup{}
	wg.Add(len(accepter))
	for _, fn := range accepter {
		switch fn {
		case "BusinessPartner":
			func() {
				c.BusinessPartner(businessPartner)
				wg.Done()
			}()
		case "CreditAccount":
			func() {
				c.CreditAccount(businessPartner)
				wg.Done()
			}()
		default:
			wg.Done()
		}
	}

	wg.Wait()
}

func (c *SAPAPICaller) BusinessPartner(businessPartner string) {
	businessPartnerData, err := c.callCreditManagementMasterSrvAPIRequirementBusinessPartner("CreditMgmtBusinessPartner", businessPartner)
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(businessPartnerData)

	creditAccountData, err := c.callToCreditAccount(businessPartnerData[0].ToCreditAccount)
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(creditAccountData)

}

func (c *SAPAPICaller) callCreditManagementMasterSrvAPIRequirementBusinessPartner(api, businessPartner string) ([]sap_api_output_formatter.BusinessPartner, error) {
	url := strings.Join([]string{c.baseURL, "API_CRDTMBUSINESSPARTNER", api}, "/")
	req, _ := http.NewRequest("GET", url, nil)

	c.setHeaderAPIKeyAccept(req)
	c.getQueryWithBusinessPartner(req, businessPartner)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToBusinessPartner(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) callToCreditAccount(url string) ([]sap_api_output_formatter.ToCreditAccount, error) {
	req, _ := http.NewRequest("GET", url, nil)
	c.setHeaderAPIKeyAccept(req)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToToCreditAccount(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) CreditAccount(businessPartner string) {
	data, err := c.callCreditManagementMasterSrvAPIRequirementCreditAccount("CreditManagementAccount", businessPartner)
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(data)

}

func (c *SAPAPICaller) callCreditManagementMasterSrvAPIRequirementCreditAccount(api, businessPartner string) ([]sap_api_output_formatter.CreditAccount, error) {
	url := strings.Join([]string{c.baseURL, "API_CRDTMBUSINESSPARTNER", api}, "/")
	req, _ := http.NewRequest("GET", url, nil)

	c.setHeaderAPIKeyAccept(req)
	c.getQueryWithCreditAccount(req, businessPartner)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToCreditAccount(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) setHeaderAPIKeyAccept(req *http.Request) {
	req.Header.Set("APIKey", c.apiKey)
	req.Header.Set("Accept", "application/json")
}

func (c *SAPAPICaller) getQueryWithBusinessPartner(req *http.Request, businessPartner string) {
	params := req.URL.Query()
	params.Add("$filter", fmt.Sprintf("BusinessPartner eq '%s'", businessPartner))
	req.URL.RawQuery = params.Encode()
}

func (c *SAPAPICaller) getQueryWithCreditAccount(req *http.Request, businessPartner string) {
	params := req.URL.Query()
	params.Add("$filter", fmt.Sprintf("BusinessPartner eq '%s'", businessPartner))
	req.URL.RawQuery = params.Encode()
}
