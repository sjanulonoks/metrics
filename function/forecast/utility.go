// Copyright 2015 Square Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package forecast

import (
	"fmt"
	"math"
)

// Returns the unique integer r such that x == r (mod m) and 0 <= r < m
func mod(x int, m int) int {
	return ((x % m) + m) % m
}

// LinearRegression estimates ys as (a + b*t) and returns (a, b).
// It performs linear regression using the explicit form for minimization of least-squares error.
// When ys[i] is NaN, it is treated as a missing point. (This makes things only slightly more complicated).
func LinearRegression(ys []float64) (float64, float64) {
	xm := 0.0  // mean of xs[i]
	ym := 0.0  // mean of ys[i]
	xym := 0.0 // mean of xs[i]*ys[i]
	x2m := 0.0 // mean of xs[i]^2
	c := 0     // count
	for i := range ys {
		if math.IsNaN(ys[i]) {
			continue
		}
		c++
		xm += float64(i)
		ym += ys[i]
		xym += float64(i) * ys[i]
		x2m += float64(i) * float64(i)
	}
	xm /= float64(c)
	ym /= float64(c)
	xym /= float64(c)
	x2m /= float64(c)
	// See https://en.wikipedia.org/wiki/Simple_linear_regression#Fitting_the_regression_line for justification.
	beta := (xym - xm*ym) / (x2m - xm*xm)
	alpha := ym - beta*xm
	return alpha, beta
}

func pValueFromNormalDifferences(correct []float64, estimate []float64) ([]float64, error) {
	if len(correct) != len(estimate) {
		return nil, fmt.Errorf("p-value calculation requires two lists of equal size")
	}
	deviations := []float64{}
	for i := range correct {
		if math.IsInf(correct[i], 0) || math.IsNaN(correct[i]) {
			continue
		}
		if math.IsInf(estimate[i], 0) || math.IsNaN(estimate[i]) {
			continue
		}
		deviations = append(deviations, estimate[i]-correct[i])
	}
	meandev := 0.0
	for _, dev := range deviations {
		meandev += dev
	}
	meandev /= float64(len(deviations))
	stddev := 0.0
	for _, dev := range deviations {
		stddev += (meandev - dev) * (meandev - dev)
	}
	stddev /= float64(len(deviations)) - 1
	stddev = math.Sqrt(stddev) // Now we have a good estimate for the population standard deviation.
	pvalues := make([]float64, len(estimate))
	for i := range pvalues {
		if math.IsInf(correct[i], 0) || math.IsNaN(correct[i]) {
			pvalues[i] = math.NaN()
			continue
		}
		if math.IsInf(estimate[i], 0) || math.IsNaN(estimate[i]) {
			pvalues[i] = math.NaN()
			continue
		}
		difference := estimate[i] - correct[i]
		tvalue := (difference - meandev) / stddev
		pvalue := 2 * math.Erf(-math.Abs(tvalue))
		pvalues[i] = pvalue
	}
	return pvalues, nil
}
