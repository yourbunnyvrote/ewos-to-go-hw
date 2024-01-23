package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"hw-async/domain"
	"hw-async/generator"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

func convertPriceToCandle(input <-chan domain.Price) <-chan domain.Candle {
	output := make(chan domain.Candle)

	go func() {
		defer close(output)

		for price := range input {
			candle := domain.Candle{
				Ticker: price.Ticker,
				Period: "",
				Open:   price.Value,
				High:   price.Value,
				Low:    price.Value,
				Close:  price.Value,
				TS:     price.TS,
			}
			output <- candle
		}
	}()

	return output
}

func createLogFile(period domain.CandlePeriod) (*os.File, *csv.Writer, error) {
	filename := fmt.Sprintf("candles_%s_log.csv", period)
	file, err := os.Create(filename)

	if err != nil {
		return nil, nil, err
	}

	writer := csv.NewWriter(file)

	return file, writer, nil
}

func writeCandleToFile(candle domain.Candle, writer *csv.Writer, period domain.CandlePeriod) {
	row := []string{
		candle.Ticker,
		candle.TS.Format("2006-01-02T15:04:05-07:00"),
		strconv.FormatFloat(candle.Open, 'f', 6, 64),
		strconv.FormatFloat(candle.High, 'f', 6, 64),
		strconv.FormatFloat(candle.Low, 'f', 6, 64),
		strconv.FormatFloat(candle.Close, 'f', 6, 64),
	}

	if err := writer.Write(row); err != nil {
		log.Printf("Ошибка записи свечи в файл: %v\n", err)
		return
	}

	log.Printf("candle%s: %+v\n", period, candle)
}

func processCandles(input <-chan domain.Candle, period domain.CandlePeriod, writer *csv.Writer, output chan<- domain.Candle) {
	defer close(output)

	candlesByTicker := make(map[string]domain.Candle)

	for candleLowTime := range input {
		myTime, err := domain.PeriodTS(period, candleLowTime.TS)
		if err != nil {
			fmt.Println(err)
			break
		}

		ticker := candleLowTime.Ticker
		candleHighTime := candleLowTime
		candleHighTime.Period = period
		candleHighTime.TS = myTime

		candle, found := candlesByTicker[ticker]
		if !found {
			candlesByTicker[ticker] = candleHighTime
			continue
		}

		if candle.TS != myTime {
			output <- candle
			writeCandleToFile(candle, writer, period)

			candlesByTicker[ticker] = candleHighTime

			continue
		}

		candleHighTime.Open = candle.Open
		if candleHighTime.High < candle.High {
			candleHighTime.High = candle.High
		}

		if candleHighTime.Low > candle.Low {
			candleHighTime.Low = candle.Low
		}

		candlesByTicker[ticker] = candleHighTime
	}

	for _, candle := range candlesByTicker {
		output <- candle
		writeCandleToFile(candle, writer, period)
	}
}

func generateCandle(input <-chan domain.Candle, period domain.CandlePeriod) <-chan domain.Candle {
	output := make(chan domain.Candle)

	file, writer, err := createLogFile(period)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		close(output)

		return output
	}

	go func() {
		defer file.Close()
		defer writer.Flush()

		processCandles(input, period, writer, output)
	}()

	return output
}

func main() {
	logger := log.New()
	ctx, cancel := context.WithCancel(context.Background())

	var tickers = []string{"SBER", "AAPL", "TSLA", "NVDA"}

	const delayMultiplier = 10

	factor := 40.0
	delay := delayMultiplier * time.Millisecond

	pg := generator.NewPricesGenerator(generator.Config{
		Factor:  factor,
		Delay:   delay,
		Tickers: tickers,
	})

	logger.Info("start prices generator...")

	prices := pg.Prices(ctx)
	candleses := convertPriceToCandle(prices)
	candleses1m := generateCandle(candleses, domain.CandlePeriod1m)
	candleses2m := generateCandle(candleses1m, domain.CandlePeriod2m)
	candleses10m := generateCandle(candleses2m, domain.CandlePeriod10m)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)

	go func() {
		<-sigCh
		cancel()
	}()

	for {
		if _, ok := <-candleses10m; !ok {
			break
		}
	}
}
