package utils

func GetMapKeys(dimensionSet map[string]string) []string {
	keys := make([]string, 0, len(dimensionSet))
	for key := range dimensionSet {
		keys = append(keys, key)
	}
	return keys
}

func AreSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	aMap := make(map[string]bool)
	for _, item := range a {
		aMap[item] = true
	}
	for _, item := range b {
		if !aMap[item] {
			return false
		}
	}
	return true
}

func ToMapInterface(dimensionSet map[string]string) map[string]interface{} {
	interfaceMap := make(map[string]interface{})
	for key, value := range dimensionSet {
		interfaceMap[key] = value
	}
	return interfaceMap
}

func MergeMaps(defaults, custom interface{}) interface{} {
	mergedInterface := make(map[string]interface{})
	mergedString := make(map[string]string)

	switch d := defaults.(type) {
	case map[string]interface{}:
		for key, value := range d {
			mergedInterface[key] = value
		}
		for key, value := range custom.(map[string]interface{}) {
			mergedInterface[key] = value
		}
		return mergedInterface
	case map[string]string:
		for key, value := range d {
			mergedInterface[key] = value
		}
		for key, value := range custom.(map[string]string) {
			mergedString[key] = value
		}
		return mergedString
	}

	return nil
}

func GenerateString(char rune, length int) string {
	r := make([]rune, length)
	for i := range r {
		r[i] = char
	}
	return string(r)
}

func AreMapsEqual(map1, map2 map[string]string) bool {
	if len(map1) != len(map2) {
		return false
	}
	for key, value := range map1 {
		if v, exists := map2[key]; !exists || v != value {
			return false
		}
	}
	return true
}

func AreArrayMapsEqual(map1, map2 []map[string]string) bool {
	if len(map1) != len(map2) {
		return false
	}
	for i := range map1 {
		map1i := map1[i]
		map2i := map2[i]
		if len(map1i) != len(map2i) {
			return false
		}
		for key, value := range map1i {
			if v, exists := map2i[key]; !exists || v != value {
				return false
			}
		}
	}
	return true
}

func AreMapsStorageResolutionEqual(map1, map2 map[string]StorageResolution) bool {
	if len(map1) != len(map2) {
		return false
	}
	for key, value := range map1 {
		if v, exists := map2[key]; !exists || v != value {
			return false
		}
	}
	return true
}

func AreMapsAnyEqual(map1, map2 map[string]any) bool {
	if len(map1) != len(map2) {
		return false
	}
	for key, value := range map1 {
		if v, exists := map2[key]; !exists || v != value {
			return false
		}
	}
	return true
}

func AreFloat64SlicesEqual(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
