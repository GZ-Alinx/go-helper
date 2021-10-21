package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// get environment variables to interface
func EnvToInterface(i interface{}, prefix string) {
	m := make(map[string]interface{}, 0)
	Struct2StructByJson(i, &m)
	newMap := envToInterface(m, prefix)
	Struct2StructByJson(newMap, &i)
}

func envToInterface(m map[string]interface{}, prefix string) map[string]interface{} {
	newMap := make(map[string]interface{}, 0)
	// json types: string/bool/float64
	for key, item := range m {
		newKey := strings.ToUpper(SnakeCase(key))
		if prefix != "" {
			newKey = strings.ToUpper(fmt.Sprintf("%s_%s", SnakeCase(prefix), SnakeCase(key)))
		}
		switch item.(type) {
		case map[string]interface{}:
			itemM, _ := item.(map[string]interface{})
			newMap[key] = envToInterface(itemM, newKey)
			continue
		case string:
			env := strings.TrimSpace(os.Getenv(newKey))
			if env != "" {
				newMap[key] = env
				fmt.Printf("[EnvToInterface]get %s: %v\n", newKey, newMap[key])
				continue
			}
		case bool:
			env := strings.TrimSpace(os.Getenv(newKey))
			if env != "" {
				itemB, ok := item.(bool)
				b, err := strconv.ParseBool(env)
				if ok && err == nil {
					if itemB && !b {
						newMap[key] = false
						fmt.Printf("[EnvToInterface]get %s: %v\n", newKey, newMap[key])
						continue
					} else if !itemB && b {
						newMap[key] = true
						fmt.Printf("[EnvToInterface]get %s: %v\n", newKey, newMap[key])
						continue
					}
				}
			}
		case float64:
			env := strings.TrimSpace(os.Getenv(newKey))
			if env != "" {
				v, err := strconv.ParseFloat(env, 64)
				if err == nil {
					newMap[key] = v
					fmt.Printf("[EnvToInterface]get %s: %v\n", newKey, newMap[key])
					continue
				}
			}
		}
		// no difference
		newMap[key] = item
	}
	return newMap
}
