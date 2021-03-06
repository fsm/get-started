package main

import (
	"net/http"
	"os"

	"github.com/TV4/graceful"
	"github.com/fsm/alexa"
	cachestore "github.com/fsm/cache-store"
	"github.com/fsm/cli"
	"github.com/fsm/example/states"
	"github.com/fsm/fsm"
	"github.com/fsm/messenger"
	"github.com/julienschmidt/httprouter"
)

func main() {
	startCLI()
}

func startCLI() {
	cli.Start(getStateMachine(), getStore())
}

func startFacebook() {
	router := &httprouter.Router{
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
	}

	// Facebook
	router.HandlerFunc(http.MethodGet, "/facebook", messenger.SetupWebhook)
	router.HandlerFunc(http.MethodPost, "/facebook", messenger.GetMessageReceivedWebhook(getStateMachine(), getStore()))

	graceful.LogListenAndServe(
		&http.Server{
			Addr:    ":" + os.Getenv("PORT"),
			Handler: router,
		},
	)
}

func startAlexa() {
	router := &httprouter.Router{
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
	}
	router.HandlerFunc(http.MethodPost, "/alexa",
		alexa.GetWebhook(
			getStateMachine(),
			getStore(),
			func(intent alexa.Intent) string {
				switch intent.Name {
				case "EnterBankIntent":
					return intent.Slots["bank_options"].Value
				case "WithdrawIntent":
					return intent.Slots["dollars"].Value
				}
				return ""
			},
		),
	)
	graceful.LogListenAndServe(
		&http.Server{
			Addr:    ":" + os.Getenv("PORT"),
			Handler: router,
		},
	)
}

func getStateMachine() fsm.StateMachine {
	return fsm.StateMachine{
		states.GetStartState,
		states.GetEnterBankState,
		states.GetWithdrawState,
		states.GetWithdrawResultState,
		states.GetReenterBankState,
		states.GetDepositState,
		states.GetDepositResultState,
		states.GetViewBalanceState,
	}
}

func getStore() fsm.Store {
	return cachestore.New()
}
