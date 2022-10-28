package formulas

import "fmt"

func IED(base float64, ieds ...float64) (float64, error) {
	// 1 - (1 - x1) * (1 - x2) ..
	if !(base > 0 && base < 100) {
		return 0, fmt.Errorf("base should be in range (0, 100)")
	}
	i := 1.0 - base/100.0
	for _, ied := range ieds {
		if !(ied > 0 && ied < 100) {
			return 0, fmt.Errorf("ied should be in range (0, 100)")
		}
		c := 1 - ied/100.0
		i = i * c
	}

	res := 1.0 - i
	return res * 100.0, nil
}
