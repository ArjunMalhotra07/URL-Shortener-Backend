package tier

import "time"

type Tier string

const (
	Free     Tier = "free"
	Pro      Tier = "pro"
	Business Tier = "business"
)

type Limits struct {
	URLsPerMonth       int
	AnalyticsDays      int  // 0 means unlimited
	APIAccess          bool
	BulkImport         bool
	MaxAnalyticsExport int // 0 means unlimited
}

var tierLimits = map[Tier]Limits{
	Free: {
		URLsPerMonth:  100,
		AnalyticsDays: 7,
		APIAccess:     false,
		BulkImport:    false,
	},
	Pro: {
		URLsPerMonth:  500,
		AnalyticsDays: 90,
		APIAccess:     true,
		BulkImport:    false,
	},
	Business: {
		URLsPerMonth:  2000,
		AnalyticsDays: 0, // unlimited
		APIAccess:     true,
		BulkImport:    true,
	},
}

// GetLimits returns the limits for a given tier
func GetLimits(t Tier) Limits {
	if limits, ok := tierLimits[t]; ok {
		return limits
	}
	return tierLimits[Free] // default to free
}

// GetLimitsFromString returns limits from tier string
func GetLimitsFromString(tierStr string) Limits {
	return GetLimits(Tier(tierStr))
}

// GetAnalyticsStartDate returns the earliest allowed start date for analytics
// based on tier. Returns zero time if unlimited.
func GetAnalyticsStartDate(t Tier) time.Time {
	limits := GetLimits(t)
	if limits.AnalyticsDays == 0 {
		return time.Time{} // unlimited
	}
	return time.Now().AddDate(0, 0, -limits.AnalyticsDays)
}

// CanCreateURL checks if user can create more URLs this month
func CanCreateURL(t Tier, currentMonthCount int64) bool {
	limits := GetLimits(t)
	return int(currentMonthCount) < limits.URLsPerMonth
}

// GetRemainingURLs returns how many URLs user can still create this month
func GetRemainingURLs(t Tier, currentMonthCount int64) int {
	limits := GetLimits(t)
	remaining := limits.URLsPerMonth - int(currentMonthCount)
	if remaining < 0 {
		return 0
	}
	return remaining
}

// HasAPIAccess checks if tier has API access
func HasAPIAccess(t Tier) bool {
	return GetLimits(t).APIAccess
}

// HasBulkImport checks if tier has bulk import feature
func HasBulkImport(t Tier) bool {
	return GetLimits(t).BulkImport
}

// IsValidTier checks if a string is a valid tier
func IsValidTier(tierStr string) bool {
	switch Tier(tierStr) {
	case Free, Pro, Business:
		return true
	default:
		return false
	}
}

// PricePerMonth returns price in USD for each tier
func PricePerMonth(t Tier) int {
	switch t {
	case Pro:
		return 10
	case Business:
		return 29
	default:
		return 0
	}
}
