package alerting

import (
	"github.com/back-end-labs/ruok/pkg/alerting/httpmsg"
	"github.com/back-end-labs/ruok/pkg/alerting/models"

	"github.com/rs/zerolog/log"
)

const (
	STATUS_OK = iota
	STATUS_FN_NOT_REGISTERED
	STATUS_ERR_WHILE_SENDING
)

type AlertManager struct {
	alertStrategies map[string]models.AlertFunc
}

var RegisteredFn = models.PluginList{
	httpmsg.Plugin,
}

func (am *AlertManager) SendAlert(i models.AlertInput) (string, int) {
	sendAlert, ok := am.alertStrategies[i.AlertStrategy]
	if !ok {
		return "", STATUS_FN_NOT_REGISTERED
	}
	result, err := sendAlert(i)
	if err != nil {
		log.Error().Err(err).Msg("couldn't send message")
		return result, STATUS_ERR_WHILE_SENDING
	}
	return result, STATUS_OK

}

func CreateAlertManager(availableChannels []string, registeredFn models.PluginList) *AlertManager {
	chanMapping := map[string]models.AlertFunc{}
	for _, plugin := range registeredFn {
		chanKey, alertFn := plugin()

		for _, chanel := range availableChannels {
			if chanKey == chanel {
				chanMapping[chanKey] = alertFn
				continue
			}
		}
	}
	return &AlertManager{
		alertStrategies: chanMapping,
	}
}
