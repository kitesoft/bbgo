package bbgo

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/c9s/bbgo/pkg/fixedpoint"
	"github.com/c9s/bbgo/pkg/types"
)

func TestPosition(t *testing.T) {

	var testcases = []struct {
		name                string
		trades              []types.Trade
		expectedAverageCost fixedpoint.Value
		expectedBase        fixedpoint.Value
		expectedQuote       fixedpoint.Value
		expectedProfit      fixedpoint.Value
	}{
		{
			name: "long",
			trades: []types.Trade{
				{
					Side:          types.SideTypeBuy,
					Price:         1000.0,
					Quantity:      0.01,
					QuoteQuantity: 1000.0 * 0.01,
				},
				{
					Side:          types.SideTypeBuy,
					Price:         2000.0,
					Quantity:      0.03,
					QuoteQuantity: 2000.0 * 0.03,
				},
			},
			expectedAverageCost: fixedpoint.NewFromFloat((1000.0*0.01 + 2000.0*0.03) / 0.04),
			expectedBase:        fixedpoint.NewFromFloat(0.01 + 0.03),
			expectedQuote:       fixedpoint.NewFromFloat(0 - 1000.0*0.01 - 2000.0*0.03),
			expectedProfit:      fixedpoint.NewFromFloat(0.0),
		},

		{
			name: "long and sell",
			trades: []types.Trade{
				{
					Side:          types.SideTypeBuy,
					Price:         1000.0,
					Quantity:      0.01,
					QuoteQuantity: 1000.0 * 0.01,
				},
				{
					Side:          types.SideTypeBuy,
					Price:         2000.0,
					Quantity:      0.03,
					QuoteQuantity: 2000.0 * 0.03,
				},
				{
					Side:          types.SideTypeSell,
					Price:         3000.0,
					Quantity:      0.01,
					QuoteQuantity: 3000.0 * 0.01,
				},
			},
			expectedAverageCost: fixedpoint.NewFromFloat((1000.0*0.01 + 2000.0*0.03) / 0.04),
			expectedBase:        fixedpoint.NewFromFloat(0.03),
			expectedQuote:       fixedpoint.NewFromFloat(0 - 1000.0*0.01 - 2000.0*0.03 + 3000.0*0.01),
			expectedProfit:      fixedpoint.NewFromFloat((3000.0 - (1000.0*0.01+2000.0*0.03)/0.04) * 0.01),
		},

		{
			name: "long, sell to short",
			trades: []types.Trade{
				{
					Side:          types.SideTypeBuy,
					Price:         1000.0,
					Quantity:      0.01,
					QuoteQuantity: 1000.0 * 0.01,
				},
				{
					Side:          types.SideTypeBuy,
					Price:         2000.0,
					Quantity:      0.03,
					QuoteQuantity: 2000.0 * 0.03,
				},
				{
					Side:          types.SideTypeSell,
					Price:         3000.0,
					Quantity:      0.10,
					QuoteQuantity: 3000.0 * 0.10,
				},
			},

			expectedAverageCost: fixedpoint.NewFromFloat(3000.0),
			expectedBase:        fixedpoint.NewFromFloat(-0.06),
			expectedQuote:       fixedpoint.NewFromFloat(-1000.0*0.01 - 2000.0*0.03 + 3000.0*0.1),
			expectedProfit:      fixedpoint.NewFromFloat((3000.0 - (1000.0*0.01+2000.0*0.03)/0.04) * 0.04),
		},

		{
			name: "short",
			trades: []types.Trade{
				{
					Side:          types.SideTypeSell,
					Price:         2000.0,
					Quantity:      0.01,
					QuoteQuantity: 2000.0 * 0.01,
				},
				{
					Side:          types.SideTypeSell,
					Price:         3000.0,
					Quantity:      0.03,
					QuoteQuantity: 3000.0 * 0.03,
				},
			},

			expectedAverageCost: fixedpoint.NewFromFloat((2000.0*0.01 + 3000.0*0.03) / (0.01 + 0.03)),
			expectedBase:        fixedpoint.NewFromFloat(0 - 0.01 - 0.03),
			expectedQuote:       fixedpoint.NewFromFloat(2000.0*0.01 + 3000.0*0.03),
			expectedProfit:      fixedpoint.NewFromFloat(0.0),
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			pos := Position{}
			profitAmount, profit := pos.AddTrades(testcase.trades)

			assert.Equal(t, testcase.expectedQuote, pos.Quote, "expectedQuote")
			assert.Equal(t, testcase.expectedBase, pos.Base, "expectedBase")
			assert.Equal(t, testcase.expectedAverageCost, pos.AverageCost, "expectedAverageCost")
			if profit {
				assert.Equal(t, testcase.expectedProfit, profitAmount, "expectedProfit")
			}
		})
	}
}