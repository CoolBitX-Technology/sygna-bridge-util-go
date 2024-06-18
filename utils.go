package bridgeutil

import (
	"encoding/json"

	"github.com/iancoleman/orderedmap"
	"github.com/samber/lo"
)

// StringToOrderedMap convert string to *orderedmapOrderedMap
func StringToOrderedMap(message string) *orderedmap.OrderedMap {
	o := orderedmap.New()
	o.UnmarshalJSON([]byte(message))
	return o
}

// OrderedMapToString convert *orderedmapOrderedMap or []*orderedmapOrderedMap to string
func OrderedMapToString(maps ...*orderedmap.OrderedMap) (string, error) {
	var message interface{}
	if len(maps) == 1 {
		message = maps[0]
	} else {
		message = maps
	}
	bMessage, err := json.Marshal(message)
	if err != nil {
		return "", err
	}
	return string(bMessage), nil
}

func castArrayToOrderedMapArray(data interface{}) []*orderedmap.OrderedMap {
	castedData := data.([]interface{})

	mapArray := make([]*orderedmap.OrderedMap, len(castedData))
	for i, v := range castedData {
		mapV := v.(orderedmap.OrderedMap)
		mapArray[i] = &mapV
	}
	return mapArray
}

func castObjectToOrderedMapObject(data interface{}) *orderedmap.OrderedMap {
	object := lo.ToPtr(data.(orderedmap.OrderedMap))
	return object
}
