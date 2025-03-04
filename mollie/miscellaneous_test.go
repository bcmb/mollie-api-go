package mollie

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/VictorAvelar/mollie-api-go/v3/testdata"
	"github.com/stretchr/testify/suite"
)

type miscellaneousServiceSuite struct{ suite.Suite }

func (ms *miscellaneousServiceSuite) SetupSuite() { setEnv() }

func (ms *miscellaneousServiceSuite) TearDownSuite() { unsetEnv() }

func (ms *miscellaneousServiceSuite) TestMiscellaneousService_ApplePaymentSession() {
	type args struct {
		ctx       context.Context
		appleSess *ApplePaymentSessionRequest
	}

	cases := []struct {
		name    string
		args    args
		wantErr bool
		err     error
		pre     func()
		handler http.HandlerFunc
	}{
		{
			"get apple payment session works as expected.",
			args{
				context.Background(),
				&ApplePaymentSessionRequest{
					Domain: "https://example.com",
				},
			},
			false,
			nil,
			noPre,
			func(w http.ResponseWriter, r *http.Request) {
				testHeader(ms.T(), r, AuthHeader, "Bearer token_X12b31ggg23")
				testMethod(ms.T(), r, "POST")

				if _, ok := r.Header[AuthHeader]; !ok {
					w.WriteHeader(http.StatusUnauthorized)
				}
				_, _ = w.Write([]byte(testdata.ListMethodsResponse))
			},
		},
		{
			"get apple payment session, an error is returned from the server",
			args{
				context.Background(),
				nil,
			},
			true,
			fmt.Errorf("500 Internal Server Error: An internal server error occurred while processing your request."),
			noPre,
			errorHandler,
		},
		{
			"get apple payment session, an error occurs when parsing json",
			args{
				context.Background(),
				nil,
			},
			true,
			fmt.Errorf("invalid character 'h' looking for beginning of object key string"),
			noPre,
			encodingHandler,
		},
		{
			"get apple payment session, invalid url when building request",
			args{
				context.Background(),
				nil,
			},
			true,
			errBadBaseURL,
			crashSrv,
			errorHandler,
		},
	}

	for _, c := range cases {
		setup()
		defer teardown()

		ms.T().Run(c.name, func(t *testing.T) {
			c.pre()
			tMux.HandleFunc("/v2/wallets/applepay/sessions", c.handler)

			res, m, err := tClient.Miscellaneous.ApplePaymentSession(c.args.ctx, c.args.appleSess)
			if c.wantErr {
				ms.NotNil(err)
				ms.EqualError(err, c.err.Error())
			} else {
				ms.Nil(err)
				ms.IsType(&ApplePaymentSession{}, m)
				ms.IsType(&http.Response{}, res.Response)
			}
		})
	}
}

func TestMiscellaneousService(t *testing.T) {
	suite.Run(t, new(miscellaneousServiceSuite))
}
