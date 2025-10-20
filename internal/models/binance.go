package models

type AggTrade struct {
    EventType        string `json:"e"`  // "aggTrade"
    EventTime        int64  `json:"E"`  // Время когда сервер отправил
    Symbol           string `json:"s"`  // Торговая пара
    AggregateTradeID int64  `json:"a"`  // Уникальный ID
    Price            string `json:"p"`  // Цена сделки
    Quantity         string `json:"q"`  // Объем сделки
    FirstTradeID     int64  `json:"f"`  // ID первой микросделки
    LastTradeID      int64  `json:"l"`  // ID последней микросделки
    TradeTime        int64  `json:"T"`  // Время самой сделки
    IsBuyer          bool   `json:"m"`  // Направление
    Ignore           bool   `json:"M"`  // Игнорировать (всегда true)
}

type MiniTicker struct {
    EventType     string `json:"e"`  // "24hrMiniTicker"
    EventTime     int64  `json:"E"`  // Время отправки
    Symbol        string `json:"s"`  // Торговая пара
    ClosePrice    string `json:"c"`  // Текущая цена
    OpenPrice     string `json:"o"`  // Цена 24ч назад
    HighPrice     string `json:"h"`  // Максимум за 24ч
    LowPrice      string `json:"l"`  // Минимум за 24ч
    TotalBaseVol  string `json:"v"`  // Объем (base asset)
    TotalQuoteVol string `json:"q"`  // Объем (quote asset)
}

type BookTicker struct {
	OrderBookUpdateID int64  `json:"u"`
	Symbol            string `json:"s"`
	BestBidPrice      string `json:"b"`
	BestBidQty        string `json:"B"`
	BestAskPrice      string `json:"a"`
	BestAskQty        string `json:"A"`
}
