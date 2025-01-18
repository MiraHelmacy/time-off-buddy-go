/*
Copyright © 2024 Alex Helmacy

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	SHORT_DESCRIPTION             string = "Calculate when you can take your next dream vacation"
	LONG_DESCRIPTION              string = "Calculate when you can take your next dream vacation based on your current time off and time off earned each pay period."
	EARNED_HOURS_OPTION           string = "ch"
	EARNED_MINUTES_OPTION         string = "cm"
	HOURS_PER_PAY_PERIOD_OPTION   string = "eh"
	MINUTES_PER_PAY_PERIOD_OPTION string = "em"
	TARGET_LONG_OPTION            string = "target"
	TARGET_SHORT_OPTION           string = "t"
	VERBOSE_LONG_OPTION           string = "verbose"
	VERBOSE_SHORT_OPTION          string = "v"
	INTERACTIVE_OPTION            string = "interactive"
	INTERACTIVE_SHORT_OPTION      string = "i"
)

var descriptions map[string]string = make(map[string]string)

type TimeOffBuddyConfig struct {
	EarnedHours         int
	EarnedMinutes       int
	HoursPerPayPeriod   int
	MinutesPerPayPeriod int
	TargetHours         int
	Verbose             bool
}

func (cfg TimeOffBuddyConfig) totalEarnedMinutes() (totalEarnedMinutes int) {
	totalEarnedMinutes = cfg.EarnedHours * 60
	totalEarnedMinutes += cfg.EarnedMinutes
	return
}

func (cfg TimeOffBuddyConfig) totalMinutesPerPayPeriod() (minutesPerPayPeriod int) {
	minutesPerPayPeriod = cfg.HoursPerPayPeriod * 60
	minutesPerPayPeriod += cfg.MinutesPerPayPeriod
	return
}

func (cfg TimeOffBuddyConfig) targetHoursInMinutes() (targetMinutes int) {
	targetMinutes = cfg.TargetHours * 60
	return
}

func (cfg TimeOffBuddyConfig) totalEarnedMinutesValid() bool {
	return cfg.totalEarnedMinutes() >= 0
}

func (cfg TimeOffBuddyConfig) totalMinutesPerPayPeridValid() bool {
	return cfg.totalMinutesPerPayPeriod() > 0
}

func (cfg TimeOffBuddyConfig) targetHoursInMinutesValid() bool {
	return cfg.targetHoursInMinutes() > 0
}

func (cfg TimeOffBuddyConfig) whatIsInvalid() (err error) {

	if !cfg.totalEarnedMinutesValid() {
		err = fmt.Errorf("total earned minutes invalid: %v", cfg.totalEarnedMinutes())
	} else if !cfg.totalMinutesPerPayPeridValid() {
		err = fmt.Errorf("total minutes per pay period invalid: %v", cfg.totalMinutesPerPayPeriod())
	} else if !cfg.targetHoursInMinutesValid() {
		err = fmt.Errorf("target hours invalid: %v", cfg.TargetHours)
	} else {
		err = nil
	}
	return
}

func (cfg TimeOffBuddyConfig) validateConfig() (isValid bool) {
	validEarnedMinutes := cfg.totalEarnedMinutesValid()
	validMinutesPerPayPeriod := cfg.totalMinutesPerPayPeridValid()
	validTargetMinutes := cfg.targetHoursInMinutesValid()
	isValid = validEarnedMinutes && validMinutesPerPayPeriod && validTargetMinutes
	return
}

func (cfg TimeOffBuddyConfig) data() (totalEarnedMinutes, totalMinutesPerPayPeriod, targetHoursInMinutes int) {
	totalEarnedMinutes = cfg.totalEarnedMinutes()
	totalMinutesPerPayPeriod = cfg.totalMinutesPerPayPeriod()
	targetHoursInMinutes = cfg.targetHoursInMinutes()
	return
}

func newTimeOffBuddyConfig() (cfg TimeOffBuddyConfig) {
	cfg = TimeOffBuddyConfig{}
	cfg.EarnedHours = viper.GetInt(EARNED_HOURS_OPTION)
	cfg.EarnedMinutes = viper.GetInt(EARNED_MINUTES_OPTION)
	cfg.HoursPerPayPeriod = viper.GetInt(HOURS_PER_PAY_PERIOD_OPTION)
	cfg.MinutesPerPayPeriod = viper.GetInt(MINUTES_PER_PAY_PERIOD_OPTION)
	cfg.TargetHours = viper.GetInt(TARGET_LONG_OPTION)
	cfg.Verbose = viper.GetBool(VERBOSE_LONG_OPTION)
	return
}

func isInteractive() bool {
	isInteractive := viper.GetBool(INTERACTIVE_OPTION)
	return isInteractive
}

func standardOptions() (options []string) {
	options = make([]string, 0)
	options = append(options, EARNED_HOURS_OPTION)
	options = append(options, EARNED_MINUTES_OPTION)
	options = append(options, HOURS_PER_PAY_PERIOD_OPTION)
	options = append(options, MINUTES_PER_PAY_PERIOD_OPTION)
	options = append(options, TARGET_LONG_OPTION)
	return
}

func interactiveMode() {
	for _, option := range standardOptions() {
		valueSet := false
		for !valueSet {
			currentValue := viper.GetInt(option)
			var input string
			fmt.Print("Enter an Integer for the following: "+descriptions[option]+": ", currentValue, ": ")
			fmt.Scanf("%v\n", &input)
			if input == "" {
				valueSet = true
			} else if newValue, e := strconv.Atoi(input); e == nil {
				viper.Set(option, newValue)
				valueSet = true
			}
		}
	}
	var input string
	fmt.Print("Enable Verbose Logging? y/n:")
	fmt.Scanf("%v\n", &input)
	if input == "y" {
		viper.Set(VERBOSE_LONG_OPTION, true)
	}
}
func determinePayPeriods(totalEarnedMinutes, minutesPerPayPeriod, targetHoursInMinutes int, verboseOutput bool) (payPeriods int) {
	payPeriods = 0
	for ; totalEarnedMinutes < targetHoursInMinutes; totalEarnedMinutes += minutesPerPayPeriod {
		if verboseOutput {
			totalEarnedHours := totalEarnedMinutes / 60
			remainingMinutes := totalEarnedMinutes % 60
			fmt.Printf("%v hrs %v minutes earned\n", totalEarnedHours, remainingMinutes)
		}

		payPeriods += 1
	}
	if verboseOutput {
		totalEarnedHours := totalEarnedMinutes / 60
		remainingMinutes := totalEarnedMinutes % 60
		fmt.Printf("%v hrs %v minutes earned\n", totalEarnedHours, remainingMinutes)
	}
	return
}
func executeTimeOffBuddy() (payPeriods int, err error) {
	if isInteractive() {
		interactiveMode()
	}
	if cfg := newTimeOffBuddyConfig(); cfg.validateConfig() {
		verboseOutput := cfg.Verbose
		totalEarnedMinutes, minutesPerPayPeriod, targetHoursInMinutes := cfg.data()
		payPeriods, err = determinePayPeriods(totalEarnedMinutes, minutesPerPayPeriod, targetHoursInMinutes, verboseOutput), nil
	} else if cfg.totalMinutesPerPayPeriod() <= 0 {
		payPeriods, err = math.MinInt, nil
	} else {
		payPeriods, err = -1, cfg.whatIsInvalid()
	}
	return
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tobuddy",
	Short: SHORT_DESCRIPTION,
	Long:  LONG_DESCRIPTION,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		if payPeriods, err := executeTimeOffBuddy(); err != nil {
			return fmt.Errorf("failed to calculate pay periods: %v", err)
		} else if payPeriods == math.MinInt {
			fmt.Println("No Time Off Earned")
		} else {
			fmt.Printf("%v pay periods\n", payPeriods)
		}
		return nil
	},
	Version: "2.1.0",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().Int(EARNED_HOURS_OPTION, 0, "Number of Hours accrued.")
	descriptions[EARNED_HOURS_OPTION] = rootCmd.Flag(EARNED_HOURS_OPTION).Usage

	rootCmd.Flags().Int(EARNED_MINUTES_OPTION, 0, "Number of Minutes accrued.")
	descriptions[EARNED_MINUTES_OPTION] = rootCmd.Flag(EARNED_MINUTES_OPTION).Usage

	rootCmd.Flags().Int(HOURS_PER_PAY_PERIOD_OPTION, 0, "Number of hours earned per pay period")
	descriptions[HOURS_PER_PAY_PERIOD_OPTION] = rootCmd.Flag(HOURS_PER_PAY_PERIOD_OPTION).Usage

	rootCmd.Flags().Int(MINUTES_PER_PAY_PERIOD_OPTION, 0, "Number of minutes earned per pay period")
	descriptions[MINUTES_PER_PAY_PERIOD_OPTION] = rootCmd.Flag(MINUTES_PER_PAY_PERIOD_OPTION).Usage

	rootCmd.Flags().IntP(TARGET_LONG_OPTION, TARGET_SHORT_OPTION, 40, "Total time off time required in hours.")
	descriptions[TARGET_LONG_OPTION] = rootCmd.Flag(TARGET_LONG_OPTION).Usage

	rootCmd.Flags().BoolP(VERBOSE_LONG_OPTION, VERBOSE_SHORT_OPTION, false, "Print verbose output.")
	rootCmd.Flags().BoolP(INTERACTIVE_OPTION, INTERACTIVE_SHORT_OPTION, false, "Start tobuddy in interactive mode.")
}

func initConfig() {
	viper.BindPFlags(rootCmd.LocalFlags())
}
