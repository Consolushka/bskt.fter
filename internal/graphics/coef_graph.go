package graphics

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"math"
	"os"
	"strconv"
)

func DrawCoefGraph() {
	line := charts.NewLine()

	// Увеличиваем размер слайсов для большего количества точек
	points := int(55/0.01) + 1
	xAxis := make([]string, points)
	yAxis := make([]opts.LineData, points)

	// Заполняем данные с шагом 0.1
	for i := 0; i <= int(38/0.01); i++ {
		x := float64(i) * 0.01
		xAxis[i] = strconv.FormatFloat(x, 'f', 1, 64)
		y := math.Pow(x, 3) / math.Pow(38, 3)
		yAxis[i] = opts.LineData{Value: y}
	}

	for i := int(38/0.01) + 1; i <= int(48/0.01); i++ {
		x := float64(i) * 0.01
		xAxis[i] = strconv.FormatFloat(x, 'f', 1, 64)
		y := 1 + math.Pow(38-x, 3)/math.Pow(38, 3)
		yAxis[i] = opts.LineData{Value: y}
	}

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Кубическая функция",
			Subtitle: "y = x³/38³",
		}),
		charts.WithTooltipOpts(opts.Tooltip{}),
		charts.WithLegendOpts(opts.Legend{}),
	)

	line.SetXAxis(xAxis).
		AddSeries("f(x)", yAxis).
		SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true)}),
			charts.WithLabelOpts(opts.Label{}),
		)

	f, _ := os.Create("cubic_styled.html")
	line.Render(f)
}
