package main

import (
	"context"
	"fmt"
	"hw-async/domain"
	"hw-async/generator"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

var tickers = []string{"SBER", "AAPL"}

type TickerByTime struct {
	Ticker string
	TS     time.Time
}

func generate1mCandle(ctx context.Context, input <-chan domain.Price) <-chan domain.Candle {
	output := make(chan domain.Candle)

	go func() {
		defer close(output)
		candlesByTime := make(map[string]domain.Candle)
		for {
			select {
			case <-ctx.Done():
				for candle := range candlesByTime {
					output <- candlesByTime[candle]
				}
				return
			case price := <-input:
				fmt.Printf("prices: %+v\n", price)
				myTime, _ := domain.PeriodTS(domain.CandlePeriod1m, price.TS)
				ticker := price.Ticker
				candle1m := domain.Candle{
					Ticker: price.Ticker,
					Period: domain.CandlePeriod1m,
					Open:   price.Value,
					High:   price.Value,
					Low:    price.Value,
					Close:  price.Value,
					TS:     myTime,
				}

				if candle, found := candlesByTime[ticker]; found {
					if candle.TS == myTime {
						candle1m.Open = candle.Open
						candle1m.TS = candle.TS
						if candle1m.High < candle.High {
							candle1m.High = candle.High
						}
						if candle1m.Low > candle.Low {
							candle1m.Low = candle.Low
						}
					} else {
						output <- candle
					}
				}
				candlesByTime[ticker] = candle1m
			}
		}
	}()

	return output
}

func generate2mCandle(ctx context.Context, input <-chan domain.Candle) <-chan domain.Candle {
	output := make(chan domain.Candle)

	go func() {
		defer close(output)
		candlesByTime := make(map[string]domain.Candle)

		for {
			select {
			case <-ctx.Done():
				for candle := range candlesByTime {
					output <- candlesByTime[candle]
				}
				return
			case candle1m := <-input:
				fmt.Printf("candle1m: %+v\n", candle1m)
				myTime, _ := domain.PeriodTS(domain.CandlePeriod2m, candle1m.TS)
				ticker := candle1m.Ticker
				candle2m := candle1m
				candle2m.Period = domain.CandlePeriod2m
				candle2m.TS = myTime

				if candle, found := candlesByTime[ticker]; found {
					if candle.TS == myTime {
						candle2m.Open = candle.Open
						if candle2m.High < candle.High {
							candle2m.High = candle.High
						}
						if candle2m.Low > candle.Low {
							candle2m.Low = candle.Low
						}
					} else {
						output <- candle
					}
				}
				candlesByTime[ticker] = candle2m
			}
		}
	}()

	return output
}

func main() {
	logger := log.New()
	ctx, cancel := context.WithCancel(context.Background())

	pg := generator.NewPricesGenerator(generator.Config{
		Factor:  10,
		Delay:   time.Millisecond * 500,
		Tickers: tickers,
	})

	logger.Info("start prices generator...")
	prices := pg.Prices(ctx)
	candleses1m := generate1mCandle(ctx, prices)
	candleses2m := generate2mCandle(ctx, candleses1m)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)
	go func() {
		<-sigCh
		cancel()
	}()

	for {
		select {
		case <-ctx.Done():
			{
				for candle := range candleses1m {
					fmt.Printf("candle1m: %+v\n", candle)
				}
				for candle := range candleses2m {
					fmt.Printf("candle2m: %+v\n", candle)
				}
				return
			}
		case c2m := <-candleses2m:
			{
				fmt.Printf("candle2m: %+v\n", c2m)
			}
		}
	}

}
