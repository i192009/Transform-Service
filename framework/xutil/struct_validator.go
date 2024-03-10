package xutil

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slices"
	"reflect"
)

// validateTemplateString is a function for JUMEAUX-213 which will validate the correct configuration syntax of string elements
func validateTemplateString(element string) bool {
	if element != "Boolean" && element != "Numerical" && element != "String" {
		println("Unknown type string ", element)
		return false
	} else {
		return true
	}
}

// validateTemplateArray is a function for JUMEAUX-213 which will validate the correct configuration syntax of array elements
func validateTemplateArray(element interface{}) bool {
	for _, element := range element.([]interface{}) {
		if reflect.TypeOf(element).Kind() == reflect.String {
			if !validateTemplateString(element.(string)) {
				return false
			}
		} else {
			println("Unknown type.")
			return false
		}
	}
	return true
}

// validateTemplateElement is a function for JUMEAUX-213 which will check the data type of the element and send to respect data type validation function
func validateTemplateElement(element interface{}) bool {
	if reflect.TypeOf(element).Kind() == reflect.String {
		if !validateTemplateString(element.(string)) {
			return false
		}
	} else if reflect.TypeOf(element).Kind() == reflect.Map {
		for _, value := range element.(map[string]interface{}) {
			if !validateTemplateElement(value) {
				return false
			}
		}
	} else if reflect.TypeOf(element).Kind() == reflect.Slice {
		if !validateTemplateArray(element) {
			return false
		}
	} else {
		return false
	}
	return true
}

// ValidateTemplate is a function for JUMEAUX-213 will loop through all elements of configuration and will validate every key and value
func ValidateTemplate(configuration map[string]any) bool {
	for key, element := range configuration {
		if !validateTemplateElement(element) {
			fmt.Printf("validate key %s error\n", key)
			return false
		}
	}
	return true
}

// validateConfigurationKey is a function for JUMEAUX-214 which will send required keys for validation
func validateConfigurationKey(key string, element interface{}, configuration map[string]any) bool {
	if key[0] == '#' {
		if reflect.TypeOf(element).Kind() == reflect.Map {
			if val, ok := configuration[key]; ok {
				for k, e := range element.(map[string]any) {
					if !validateConfigurationKey(k, e, val.(map[string]interface{})) {
						return false
					}
				}
			} else {
				return false
			}
		} else {
			if _, ok := configuration[key]; !ok {
				return false
			}
		}
	}
	return true
}

// validateConfigurationElement is a function for JUMEAUX-214 which will check compare element and configuration and then send to respect data type validation function
func validateConfigurationElement(key string, element any, currentConfig map[string]any) bool {
	if val, ok := currentConfig[key]; ok {
		if element == "String" {
			if reflect.TypeOf(val).Kind() != reflect.String {
				return false
			}
		} else if element == "Numerical" {
			if reflect.TypeOf(val).Kind() != reflect.Float64 {
				return false
			}
		} else if element == "Boolean" {
			if reflect.TypeOf(val).Kind() != reflect.Bool {
				return false
			}
		} else if reflect.TypeOf(val).Kind() == reflect.Slice {
			if len(val.(primitive.A)) != 0 {
				if reflect.TypeOf(element).Kind() == reflect.Slice {
					var configArray []string
					for _, value := range val.(primitive.A) {
						configArray = append(configArray, value.(string))
					}
					for _, e := range element.([]interface{}) {
						if reflect.TypeOf(e).Kind() == reflect.String {
							if !slices.Contains(configArray, "String") {
								return false
							}
						} else if reflect.TypeOf(e).Kind() == reflect.Float64 {
							if !slices.Contains(configArray, "Numerical") {
								return false
							}
						}
					}
				} else {
					return false
				}
			}
		} else if reflect.TypeOf(val).Kind() == reflect.Map {
			if reflect.TypeOf(element).Kind() == reflect.Map {
				for k, v := range element.(map[string]interface{}) {
					if !validateConfigurationElement(k, v, val.(map[string]any)) {
						return false
					}
				}
			} else {
				return false
			}
		} else {
			return false
		}
	} else {
		return false
	}
	return true
}

// ValidateConfigurationKeys is a function for JUMEAUX-214 will loop through all elements of configuration and will check if all required keys are present
func ValidateConfigurationKeys(currentConfig map[string]any, configuration map[string]any) bool {
	for key, element := range currentConfig {
		if !validateConfigurationKey(key, element, configuration) {
			return false
		}
	}
	return true
}

// ValidateConfiguration is a function for JUMEAUX-214 will loop through all elements of configuration and will validate every key and value
func ValidateConfiguration(currentConfig map[string]any, configuration map[string]any) bool {
	for key, element := range configuration {
		if !validateConfigurationElement(key, element, currentConfig) {
			return false
		}
	}
	return true
}
