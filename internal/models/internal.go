package models

import "time"

type UniversalTrade struct {
	// ============================================
	// ОБЩИЕ ПОЛЯ (есть у всех типов)
	// ============================================
	Symbol    string    `json:"symbol"`     // "BTCUSDT"
	Timestamp time.Time `json:"timestamp"`  // Когда произошло
	EventType string    `json:"event_type"` // "aggTrade", "miniTicker", "bookTicker"

	// ============================================
	// ОСНОВНАЯ ЦЕНА (всегда заполнено)
	// ============================================
	Price float64 `json:"price"` // Текущая/последняя цена

	// ============================================
	// ДЛЯ aggTrade
	// ============================================
	Quantity float64 `json:"quantity,omitempty"` // Объем сделки
	//TradeID      int64   `json:"trade_id,omitempty"`       // ID сделки
	IsBuyerMaker bool `json:"is_buyer_maker,omitempty"` // Направление

	// ============================================
	// ДЛЯ miniTicker (24ч статистика)
	// ============================================
	OpenPrice   float64 `json:"open_price,omitempty"`   // Цена открытия 24ч
	HighPrice   float64 `json:"high_price,omitempty"`   // Максимум 24ч
	LowPrice    float64 `json:"low_price,omitempty"`    // Минимум 24ч
	Volume      float64 `json:"volume,omitempty"`       // Объем 24ч
	QuoteVolume float64 `json:"quote_volume,omitempty"` // Объем в quote asset

	// ============================================
	// ДЛЯ bookTicker (order book)
	// ============================================
	BidPrice float64 `json:"bid_price,omitempty"` // Лучший bid
	AskPrice float64 `json:"ask_price,omitempty"` // Лучший ask
	Spread   float64 `json:"spread,omitempty"`    // ask - bid (вычисляется)
}
