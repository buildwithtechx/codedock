package compose

import (
	"math"
	"strconv"
	"strings"
)

var composeDurationUnits = []struct {
	suffix string
	ns     int64
}{
	{"ns", 1},
	{"us", 1_000},
	{"µs", 1_000},
	{"μs", 1_000},
	{"ms", 1_000_000},
	{"s", 1_000_000_000},
	{"m", 60 * 1_000_000_000},
	{"h", 3_600 * 1_000_000_000},
	{"d", 86_400 * 1_000_000_000},
	{"w", 604_800 * 1_000_000_000},
}

var dockerDurationLimitNs = math.Ldexp(1, 63)

func ParseComposeDuration(value string) (int64, bool) {
	str := strings.TrimSpace(value)
	if str == "" {
		return 0, false
	}
	if str[0] == '+' {
		str = str[1:]
		if str == "" {
			return 0, false
		}
	}

	if number, rest, ok := scanComposeNumber(str); ok && rest == "" {
		seconds, err := strconv.ParseFloat(number, 64)
		if err != nil {
			return 0, false
		}
		return finishComposeDuration(seconds * 1_000_000_000)
	}

	total := 0.0
	groups := 0
	for work := str; work != ""; {
		number, rest, ok := scanComposeNumber(work)
		if !ok {
			return 0, false
		}
		unitNs, rest2, ok := scanComposeUnit(rest)
		if !ok {
			return 0, false
		}
		scaled, err := strconv.ParseFloat(number, 64)
		if err != nil {
			return 0, false
		}
		total += scaled * float64(unitNs)
		groups++
		work = rest2
	}
	if groups == 0 {
		return 0, false
	}
	return finishComposeDuration(total)
}

func finishComposeDuration(ns float64) (int64, bool) {
	if math.IsNaN(ns) || math.IsInf(ns, 0) || ns >= dockerDurationLimitNs {
		return 0, false
	}
	rounded := math.Round(ns)
	if rounded < 0 || rounded >= dockerDurationLimitNs {
		return 0, false
	}
	return int64(rounded), true
}

func scanComposeNumber(s string) (number, rest string, ok bool) {
	i := 0
	for i < len(s) && isComposeDigit(s[i]) {
		i++
	}
	intDigits := i
	if i < len(s) && s[i] == '.' {
		i++
		fracStart := i
		for i < len(s) && isComposeDigit(s[i]) {
			i++
		}
		if intDigits == 0 && i == fracStart {
			return "", "", false
		}
		return s[:i], s[i:], true
	}
	if intDigits == 0 {
		return "", "", false
	}
	return s[:i], s[i:], true
}

func scanComposeUnit(s string) (ns int64, rest string, ok bool) {
	for _, unit := range composeDurationUnits {
		if strings.HasPrefix(s, unit.suffix) {
			return unit.ns, s[len(unit.suffix):], true
		}
	}
	return 0, "", false
}

func isComposeDigit(b byte) bool { return b >= '0' && b <= '9' }
