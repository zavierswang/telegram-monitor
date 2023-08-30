package cst

const (
	AppName           = "telegram-monitor"
	BaseName          = "telegram"
	DateTimeFormatter = "2006-01-02 15:04:05"
	TimeFormatter     = "02 15:04:05"
	PerCountEnergy    = 32000
	LicenseApi        = "http://127.0.0.1/api/license"
	UserAgent         = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"
)

const (
	OrderStatus = iota
	OrderStatusSuccess
	OrderStatusRunning
	OrderStatusReceived
	OrderStatusApiSuccess
	OrderStatusApiFailure
	OrderStatusFailure
	OrderStatusNotSufficientFunds
	OrderStatusCancel
)

const (
	OkxMarketTradesApi  = "https://www.okx.com/priapi/v5/market/trades"
	OkxTradingOrdersApi = "https://www.okx.com/v3/c2c/tradingOrders/books"
)
