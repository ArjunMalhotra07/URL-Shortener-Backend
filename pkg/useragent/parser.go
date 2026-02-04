package useragent

import (
	"strings"

	"github.com/mssola/useragent"
)

type ParsedUA struct {
	Browser    string
	OS         string
	DeviceType string
	IsBot      bool
}

func Parse(userAgentStr string) *ParsedUA {
	ua := useragent.New(userAgentStr)

	browser, _ := ua.Browser()
	os := ua.OS()
	isBot := ua.Bot()

	deviceType := determineDeviceType(ua, userAgentStr)

	return &ParsedUA{
		Browser:    browser,
		OS:         os,
		DeviceType: deviceType,
		IsBot:      isBot,
	}
}

func determineDeviceType(ua *useragent.UserAgent, userAgentStr string) string {
	if ua.Bot() {
		return "bot"
	}

	uaLower := strings.ToLower(userAgentStr)

	// Check for tablets first (before mobile, as tablets often contain "mobile")
	if strings.Contains(uaLower, "ipad") ||
		strings.Contains(uaLower, "tablet") ||
		(strings.Contains(uaLower, "android") && !strings.Contains(uaLower, "mobile")) {
		return "tablet"
	}

	// Check for mobile devices
	if ua.Mobile() ||
		strings.Contains(uaLower, "iphone") ||
		strings.Contains(uaLower, "ipod") ||
		(strings.Contains(uaLower, "android") && strings.Contains(uaLower, "mobile")) ||
		strings.Contains(uaLower, "windows phone") ||
		strings.Contains(uaLower, "blackberry") {
		return "mobile"
	}

	return "desktop"
}
