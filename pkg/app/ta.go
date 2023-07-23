package app

import (
	"fmt"

	"github.com/sdcoffey/techan"
)

func getBolBands(series *techan.TimeSeries) (techan.Indicator, techan.Indicator, error) {
	/*
		Calculates the Bollinger Bands indicator for a given time series.
		The upper band and lower band indicators are returned.

		Parameters:
		- series: a pointer to a techan.TimeSeries object for which to calculate the Bollinger Bands.

		Returns:
		- A techan.Indicator representing the upper band of the Bollinger Bands.
		- A techan.Indicator representing the lower band of the Bollinger Bands.
	*/
	if series == nil || series == techan.NewTimeSeries() {
		return nil, nil, fmt.Errorf("input series is nil or empty")
	}

	// Get the Closing price for the series
	closePrices := techan.NewClosePriceIndicator(series)

	// Generate the Upper Bollinger Band
	bolHigh := techan.NewBollingerUpperBandIndicator(closePrices, 21, 2.0)

	// Generate the Lower Bollinger Band
	bolLow := techan.NewBollingerLowerBandIndicator(closePrices, 21, 2.0)

	return bolHigh, bolLow, nil
}

func Strategy1(series *techan.TimeSeries) (techan.RuleStrategy, error) {
	/*
		Returns a RuleStrategy that uses Bollinger Bands and cross
		indicator rules to generate buy and sell signals.

		Parameters:
		- series: A pointer to a TimeSeries containing the data to be analyzed.

		Returns:
		- A RuleStrategy struct containing the entry and exit rules for the strategy.
		- An error if there was an issue generating the strategy.
	*/
	if series == nil || len(series.Candles) == 0 {
		return techan.RuleStrategy{}, fmt.Errorf("empty or nil time series")
	}

	// Get the High, Low, and Close prices from the series
	highPrices := techan.NewHighPriceIndicator(series)
	lowPrices := techan.NewLowPriceIndicator(series)
	closePrices := techan.NewClosePriceIndicator(series)

	// Get the Upper and Lower Bollinger Bands
	bolHigh, bolLow, err := getBolBands(series)
	if err != nil {
		return techan.RuleStrategy{}, fmt.Errorf("failed to get Bollinger Bands %v\n", err)
	}

	// Calculate the RSI indicator
	rsi := techan.NewRelativeStrengthIndexIndicator(closePrices, 14)
	rsiOversold := techan.NewConstantIndicator(30)
	rsiUndersold := techan.NewConstantIndicator(70)

	// Create the entry rule with Bollinger Bands and RSI filter conditions
	entryRule := techan.And(
		techan.And(
			techan.NewCrossDownIndicatorRule(lowPrices, bolLow), // Price below lower Bollinger Band
			techan.NewCrossUpIndicatorRule(rsi, rsiOversold),    // RSI crosses above 30 (oversold condition)
		),
		techan.PositionNewRule{}, // No existing position
	)

	// Create the exit rule with Bollinger Bands and RSI filter conditions
	exitRule := techan.And(
		techan.And(
			techan.NewCrossUpIndicatorRule(highPrices, bolHigh), // Price above upper Bollinger Band
			techan.NewCrossDownIndicatorRule(rsi, rsiUndersold), // RSI crosses below 70 (overbought condition)
		),
		techan.PositionOpenRule{}, // Existing position open
	)

	// Create the strategy with the above entry and exit rules
	strategy := techan.RuleStrategy{
		UnstablePeriod: 0,
		EntryRule:      entryRule,
		ExitRule:       exitRule,
	}
	return strategy, nil
}
