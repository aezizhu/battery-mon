package battery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type BatteryInfo struct {
	Percent    int
	Status     string
	Remaining  string
	IsCharging bool
	Source     string
}

type AdvancedInfo struct {
	CycleCount  int
	Condition   string
	MaxCapacity string
	Wattage     string
	ChargerName string
	Serial      string
}

// SPPowerDataType structure for JSON unmarshalling
type spPowerData struct {
	SPPowerDataType []map[string]interface{} `json:"SPPowerDataType"`
}

func GetBasicInfo() (BatteryInfo, error) {
	info := BatteryInfo{
		Status:    "Unknown",
		Remaining: "Unknown",
		Source:    "Unknown",
	}

	cmd := exec.Command("pmset", "-g", "batt")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return info, err
	}

	output := out.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 0 {
		return info, fmt.Errorf("no output from pmset")
	}

	// Line 1: Source
	// "Now drawing from 'AC Power'"
	reSource := regexp.MustCompile(`Now drawing from '([^']+)'`)
	if matches := reSource.FindStringSubmatch(lines[0]); len(matches) > 1 {
		info.Source = matches[1]
	}

	if len(lines) > 1 {
		detailLine := lines[1]
		
		// Percentage
		rePct := regexp.MustCompile(`(\d+)%`)
		if matches := rePct.FindStringSubmatch(detailLine); len(matches) > 1 {
			info.Percent, _ = strconv.Atoi(matches[1])
		}

		// Status
		reStatus := regexp.MustCompile(`;\s(.*?);\s`)
		if matches := reStatus.FindStringSubmatch(detailLine); len(matches) > 1 {
			info.Status = matches[1]
			if strings.Contains(info.Status, "charging") || strings.Contains(info.Status, "finishing charge") {
				info.IsCharging = true
			}
		} else if strings.Contains(detailLine, "charged") {
			info.Status = "charged"
		}

		// Time Remaining
		reTime := regexp.MustCompile(`(\d+:\d+)\sremaining`)
		if matches := reTime.FindStringSubmatch(detailLine); len(matches) > 1 {
			info.Remaining = matches[1]
		} else if strings.Contains(detailLine, "no estimate") {
			info.Remaining = "Calculating..."
		}
	}

	return info, nil
}

func GetAdvancedInfo() (AdvancedInfo, error) {
	info := AdvancedInfo{
		Condition:   "N/A",
		MaxCapacity: "N/A",
		Wattage:     "N/A",
		ChargerName: "N/A",
		Serial:      "N/A",
	}

	cmd := exec.Command("system_profiler", "SPPowerDataType", "-json")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return info, err
	}

	var data spPowerData
	if err := json.Unmarshal(out.Bytes(), &data); err != nil {
		return info, err
	}

	for _, item := range data.SPPowerDataType {
		// Battery Info
		if battInfo, ok := item["spbattery_information"].(map[string]interface{}); ok {
			if healthInfo, ok := battInfo["sppower_battery_health_info"].(map[string]interface{}); ok {
				if cycles, ok := healthInfo["sppower_battery_cycle_count"].(float64); ok {
					info.CycleCount = int(cycles)
				}
				if cond, ok := healthInfo["sppower_battery_health"].(string); ok {
					info.Condition = cond
				}
				if maxCap, ok := healthInfo["sppower_battery_health_maximum_capacity"].(string); ok {
					info.MaxCapacity = maxCap
				}
			}
			if modelInfo, ok := battInfo["sppower_battery_model_info"].(map[string]interface{}); ok {
				if serial, ok := modelInfo["sppower_battery_serial_number"].(string); ok {
					info.Serial = serial
				}
			}
		}

		// Charger Info
		if chargerInfo, ok := item["sppower_ac_charger_information"].(map[string]interface{}); ok {
			if watts, ok := chargerInfo["sppower_ac_charger_watts"].(string); ok {
				info.Wattage = watts
			}
			if name, ok := chargerInfo["sppower_ac_charger_name"].(string); ok {
				info.ChargerName = name
			}
		}
	}

	return info, nil
}
