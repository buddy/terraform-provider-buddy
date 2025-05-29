package util

import (
	"github.com/buddy/api-go-sdk/buddy"
	"regexp"
)

func FilterTargetByName(targets []*buddy.Target, name string) *buddy.Target {
	for _, target := range targets {
		if target.Name == name {
			return target
		}
	}
	return nil
}

func FilterTargetListByNameRegex(targets []*buddy.Target, nameRegex string) []*buddy.Target {
	re, err := regexp.Compile(nameRegex)
	if err != nil {
		return nil
	}
	var result []*buddy.Target
	for _, target := range targets {
		if re.MatchString(target.Name) {
			result = append(result, target)
		}
	}
	return result
}