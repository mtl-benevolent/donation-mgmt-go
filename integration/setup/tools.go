package setup

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/oklog/ulid/v2"
)

var names = []string{
	"3d_printer", "charging_adapter", "computer", "cpu", "drone", "ebook_reader", "ethernet_cable", "external_battery", "flash_drive", "fitness_tracker",
	"gamepad", "gpu", "hard_drive", "headphones", "hdmi_cable", "heat_sink", "joystick", "keyboard", "laptop", "memory_card",
	"microphone", "modem", "monitor", "motherboard", "mouse", "network_card", "optical_drive", "power_supply", "printer", "projector",
	"ram", "robot", "router", "scanner", "smart_glasses", "smartphone", "smartwatch", "sound_card", "speaker", "ssd",
	"stand", "tablet", "touchscreen_pen", "usb_cable", "vr_headset", "webcam", "case", "cooling_fan", "dock", "graphics_card",
}

var adjectives = []string{
	"ancient", "bad", "bright", "broad", "clean", "cold", "colorful", "dark", "deep", "dirty",
	"drab", "empty", "fast", "flat", "full", "good", "hard", "heavy", "high", "large",
	"light", "long", "loud", "low", "modern", "narrow", "new", "old", "old", "poor",
	"quick", "rich", "round", "sharp", "shallow", "short", "slow", "small", "smooth", "soft",
	"square", "strong", "tall", "thick", "thin", "warm", "weak", "young", "quiet", "bright",
}

func GenerateName() string {
	// #nosec G404
	adjective := adjectives[rand.Intn(len(adjectives))]
	// #nosec G404
	name := names[rand.Intn(len(names))]
	id := ulid.Make().String()

	return fmt.Sprintf("%s %s %s", adjective, name, id)
}

func Slugify(s string, maxLen int) string {
	s = strings.ToLower(s)

	// remap spaces, remove all symbols
	s = strings.Map(func(r rune) rune {
		if r == ' ' || r == '_' || r == '.' || r == '-' {
			return '-'
		}

		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			return r
		}

		return -1
	}, s)

	// truncate string to max len
	if len(s) <= maxLen {
		return s
	}

	runeSlice := []rune(s)
	return string(runeSlice[:maxLen])
}
