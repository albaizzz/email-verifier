package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	emailverifier "github.com/AfterShip/email-verifier"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
)

var (
	verifier = emailverifier.NewVerifier()
)

func GetEmailVerification(w http.ResponseWriter, r *http.Request) {
	email := mux.Vars(r)["email"]
	ret, err := verifier.Verify(email)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if !ret.Syntax.Valid {
		_, _ = fmt.Fprint(w, "email address syntax is invalid")
		return
	}

	bytes, err := json.Marshal(ret)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	_, _ = fmt.Fprint(w, string(bytes))

}

var gorillaLambda *gorillamux.GorillaMuxAdapter

func init() {

	r := mux.NewRouter()

	r.HandleFunc("/v1/{email}/verification", GetEmailVerification).Methods("GET")
	// func(w http.ResponseWriter, r *http.Request) {
	// 	json.NewEncoder(w).Encode("")
	// })

	gorillaLambda = gorillamux.New(r)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	r, err := gorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&req))
	return *r.Version1(), err
}
func main() {
	lambda.Start(Handler)
}
