package processor

import (
	"encoding/json"
	"log/slog"

	"github.com/WWoi/web-parcer/internal/models"
)

type Processor struct {
	inputChan  <-chan []byte
	outputChan chan<- models.UniversalTrade
}

func New(inChan chan []byte, outChan chan models.UniversalTrade) *Processor {
	return &Processor{
		inputChan:  inChan,
		outputChan: outChan,
	}
}

func (p *Processor) Start() {
	for range 10 {
		go p.worker()
	}
}

func (p *Processor) worker() {
	for rawMsg := range p.inputChan {

		trade, err := p.parse(rawMsg)
		if err != nil {
			continue
		}

		p.outputChan <- trade
	}
}

func (p *Processor) parse(rawMsg []byte) (models.UniversalTrade, error) {
	var baseEvent struct {
		EventType string `json:"e"`
	}

	if err := json.Unmarshal(rawMsg, &baseEvent); err != nil {
		slog.Error("❌Could not parse from JSON", slog.String("error", err.Error()))
		return models.UniversalTrade{}, err
	}

	var unTrade models.UniversalTrade
	var err error

	switch baseEvent.EventType {
	case "aggTrade":
		var aggTrade models.AggTrade
		if err := json.Unmarshal(rawMsg, &aggTrade); err != nil {
			slog.Error("❌Could not parse from JSON", slog.String("error", err.Error()))
			return models.UniversalTrade{}, err
		}

		unTrade, err = convertAggTradeToUniversalTrade(aggTrade)
		if err != nil {
			return models.UniversalTrade{}, err
		}

	case "miniTicker":
		var miniTicker models.MiniTicker
		if err := json.Unmarshal(rawMsg, &miniTicker); err != nil {
			slog.Error("❌Could not parse JSON", slog.String("error", err.Error()))
			return models.UniversalTrade{}, err
		}

		unTrade, err = convertMiniTickerToUniversalTrade(miniTicker)
		if err != nil {
			return models.UniversalTrade{}, err
		}
	}

	return unTrade, nil
}
