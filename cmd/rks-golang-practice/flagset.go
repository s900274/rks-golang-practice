package main

import "flag"

func Flagset() *flag.FlagSet {
    flagSet := flag.NewFlagSet("uranus", flag.ExitOnError)
    flagSet.String("config", "", "path to config file")

    return flagSet
}
