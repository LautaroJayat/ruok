package models

type AlertInput struct {
	AlertStrategy  string
	Url            string
	Method         string
	Payload        string
	ExpectedStatus int
	ExpectedMsg    string
	Headers        map[string]string
}

type AlertFunc func(AlertInput) (string, error)

type AlertPlugin func() (string, AlertFunc)

type PluginList []AlertPlugin
