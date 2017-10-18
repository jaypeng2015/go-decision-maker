package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	sparta "github.com/mweagle/Sparta"
	"github.com/mweagle/Sparta/explore"
)

func TestRoll(t *testing.T) {
	// 1. Create the function(s) we want to test
	var lambdaFunctions []*sparta.LambdaAWSInfo
	lambdaFn := sparta.HandleAWSLambda(sparta.LambdaName(makeDecision),
		http.HandlerFunc(makeDecision),
		sparta.IAMRoleDefinition{})
	lambdaFunctions = append(lambdaFunctions, lambdaFn)

	// 2. Mock event specific data to send to the lambda function
	eventData := url.Values{}
	eventData.Set("user_id", "U12345")
	eventData.Set("text", "coin")

	// 3. Make the request and confirm
	// Make the request and confirm
	logger, _ := sparta.NewLogger("warning")
	ts := httptest.NewServer(sparta.NewServeMuxLambda(lambdaFunctions, logger))
	defer ts.Close()
	whitelistParamValues := map[string]string{
		"method.request.header.Content-type": "application/x-www-form-urlencoded",
	}
	resp, err := explore.NewAPIGatewayRequest(lambdaFn.URLPath(),
		"POST",
		whitelistParamValues,
		eventData.Encode(),
		ts.URL)

	if err != nil {
		t.Fatal(err.Error())
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	t.Log("Status: ", resp.Status)
	t.Log("Headers: ", resp.Header)
	t.Log("Body: ", string(body))
}
