package handlers

import (
	"fmt"
	"net/http"

	"bitbucket.org/hebertthome/traning-oauth-go/context"
	"bitbucket.org/hebertthome/traning-oauth-go/logger"
)

func API(ctx *context.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	// Build Response
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("{ \"AtivarPrePagoEnriquecimentoResponse\": { \"MessageHeader\": { \"Response\": { \"Return\": { \"Type\": \"N\", \"Description\": \"Solicitação de ativação encaminhada ao Siebel.\", \"NativeReturn\": { \"Type\": \"S\", \"Description\": \"-\", \"AppId\": \"AN0448\", \"Code\": \"Y\" }, \"Code\": \"00094\" } }, \"CorrelationId\": \"\", \"Timestamp\": \"2021-03-10T10:01:26.736-03:00\", \"BusinessId\": \"\", \"Credentials\": { \"AppToken\": \"\", \"AppId\": \"FWSOA\", \"UserId\": \"\", \"UserToken\": \"\" }, \"TransactionId\": \"25accd21-1c7a-46a8-a310-5f6187de98db\" } } }")))

	ctx.Logger.Info("API",
		logger.String("Request", "Success"),
	)

	return http.StatusOK, nil
}
