package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Sirupsen/logrus"
	sparta "github.com/mweagle/Sparta"
	cloudformation "github.com/mweagle/go-cloudformation"
)

// slackLambdaJSONEvent provides a pass through mapping
// of all whitelisted Parameters.  The transformation is defined
// by the resources/gateway/inputmapping_json.vtl template.
type slackLambdaJSONEvent struct {
	// HTTPMethod
	Method string `json:"method"`
	// Body, if available.  This is going to be an interface s.t. we can support
	// testing through APIGateway, which by default sends 'application/json'
	Body interface{} `json:"body"`
	// Whitelisted HTTP headers
	Headers map[string]string `json:"headers"`
	// Whitelisted HTTP query params
	QueryParams map[string]string `json:"queryParams"`
	// Whitelisted path parameters
	PathParams map[string]string `json:"pathParams"`
	// Context information - http://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-mapping-template-reference.html#context-variable-reference
	Context sparta.APIGatewayContext `json:"context"`
}

func makeDecision(w http.ResponseWriter, r *http.Request) {
	// logger, _ := r.Context().Value(sparta.ContextKeyLogger).(*logrus.Logger)
	logger, _ := sparta.NewLogger("info")

	// 1. Unmarshal the primary event
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var lambdaEvent slackLambdaJSONEvent
	err := decoder.Decode(&lambdaEvent)
	if err != nil {
		logger.Error("Failed to unmarshal event data: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// 2. Conditionally unmarshal to get the Slack text.  See
	// https://api.slack.com/slash-commands
	// for the value name list
	requestParams := url.Values{}
	if bodyData, ok := lambdaEvent.Body.(string); ok {
		requestParams, err = url.ParseQuery(bodyData)
		if err != nil {
			logger.Error("Failed to parse query: ", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		logger.WithFields(logrus.Fields{
			"Values": requestParams,
		}).Info("Slack slashcommand values")
	} else {
		logger.Info("Event body empty")
	}

	// 3. Create the response
	// Slack formatting:
	// https://api.slack.com/docs/formatting
	var userID string
	if requestParams["user_id"] != nil && len(requestParams["user_id"]) > 0 {
		userID = "<@" + requestParams["user_id"][0] + ">"
	} else {
		userID = "You"
	}
	var text string
	if requestParams["text"] != nil && len(requestParams["user_id"]) > 0 && requestParams["text"][0] == "coin" {
		text = "coin"
	} else {
		text = "dice"
	}

	var responseText string
	if text == "dice" {
		responseText = userID + " just rolled :dice_" + strconv.Itoa(rand.Intn(6)+1) + ":"
	} else {
		responseText = userID + " just rolled :coin_" + strconv.Itoa(rand.Intn(2)+1) + ":"
	}

	// 4. Setup the response object:
	// https://api.slack.com/slash-commands, "Responding to a command"
	responseData := sparta.ArbitraryJSONObject{
		"response_type": "in_channel",
		"text":          responseText,
	}
	// 5. Send it off
	responseBody, err := json.Marshal(responseData)
	if err != nil {
		logger.Error("Failed to marshal response: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
}

func spartaLambdaFunctions(api *sparta.API) []*sparta.LambdaAWSInfo {
	var lambdaFunctions []*sparta.LambdaAWSInfo
	roleDefinition := sparta.IAMRoleDefinition{}
	roleDefinition.Privileges = append(roleDefinition.Privileges,
		sparta.IAMRolePrivilege{
			Actions:  []string{"ec2:CreateNetworkInterface", "ec2:DescribeNetworkInterfaces", "ec2:DeleteNetworkInterface"},
			Resource: "*",
		})
	roleDefinition.Privileges = append(roleDefinition.Privileges,
		sparta.IAMRolePrivilege{
			Actions:  []string{"xray:PutTraceSegments", "xray:PutTelemetryRecords"},
			Resource: "*",
		})
	lambdaFn := sparta.HandleAWSLambda(sparta.LambdaName(makeDecision),
		http.HandlerFunc(makeDecision),
		roleDefinition)

	vpcConfig := cloudformation.LambdaFunctionVPCConfig{}
	vpcConfig.SecurityGroupIDs = &cloudformation.StringListExpr{
		Literal: []*cloudformation.StringExpr{&cloudformation.StringExpr{Literal: "sg-cc2f31a0"}},
	}
	vpcConfig.SubnetIDs = &cloudformation.StringListExpr{
		Literal: []*cloudformation.StringExpr{&cloudformation.StringExpr{Literal: "subnet-7ec95317"}},
	}
	lambdaFn.Options.VpcConfig = &vpcConfig

	if nil != api {
		apiGatewayResource, _ := api.NewResource("/roll", lambdaFn)
		_, err := apiGatewayResource.NewMethod("POST", http.StatusOK)
		if nil != err {
			panic("Failed to create /roll resource")
		}
	}
	return append(lambdaFunctions, lambdaFn)
}

////////////////////////////////////////////////////////////////////////////////
// Main
func main() {
	// Register the function with the API Gateway
	apiStage := sparta.NewStage("dev")
	apiGateway := sparta.NewAPIGateway("GoDecisionMaker", apiStage)

	// Deploy it
	sparta.Main("GoDecisionMaker",
		"The Decision Maker in Go.",
		spartaLambdaFunctions(apiGateway),
		apiGateway,
		nil)
}
