package processor

import (
	"strconv"
	"time"

	"github.com/WWoi/web-parcer/internal/models"
)

func convertAggTradeToUniversalTrade(model models.AggTrade) (models.UniversalTrade, error) {
	price, err := strconv.ParseFloat(model.Price, 64)
	if err != nil {
		return models.UniversalTrade{}, err
	}
	quantity, err := strconv.ParseFloat(model.Quantity, 64)
	if err != nil {
		return models.UniversalTrade{}, err
	}

	return models.UniversalTrade{
		Symbol:    model.Symbol,
		Timestamp: time.UnixMilli(model.EventTime),
		EventType: model.EventType,

		Price:        price,
		Quantity:     quantity,
		IsBuyerMaker: model.IsBuyer,
	}, nil
}

func convertMiniTickerToUniversalTrade(model models.MiniTicker) (models.UniversalTrade, error) {
	cPrice, err := strconv.ParseFloat(model.ClosePrice, 64)
	if err != nil {
		return models.UniversalTrade{}, err
	}
	oPrice, err := strconv.ParseFloat(model.OpenPrice, 64)
	if err != nil {
		return models.UniversalTrade{}, err
	}
	hPrice, err := strconv.ParseFloat(model.HighPrice, 64)
	if err != nil {
		return models.UniversalTrade{}, err
	}
	lPrice, err := strconv.ParseFloat(model.LowPrice, 64)
	if err != nil {
		return models.UniversalTrade{}, err
	}
	volume, err := strconv.ParseFloat(model.TotalBaseVol, 64)
	if err != nil {
		return models.UniversalTrade{}, err
	}
	quoteVolume, err := strconv.ParseFloat(model.TotalQuoteVol, 64)
	if err != nil {
		return models.UniversalTrade{}, err
	}

	return models.UniversalTrade{
		Symbol:    model.Symbol,
		Timestamp: time.UnixMilli(model.EventTime),
		EventType: model.EventType,
		Price:     cPrice,

		OpenPrice:   oPrice,
		HighPrice:   hPrice,
		LowPrice:    lPrice,
		Volume:      volume,
		QuoteVolume: quoteVolume,
	}, nil
}
