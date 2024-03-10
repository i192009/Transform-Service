package controller

import (
	"encoding/json"
	"strconv"
	"strings"
)

func CheckTargetFormatPram(parms map[string]interface{}) (map[string]interface{}, string) {
	if parms == nil {
		parms = make(map[string]interface{})
	}
	var scriptsStr strings.Builder

	if _, ok := parms["tessellate"]; !ok {
		parms["tessellate"] = map[string]float32{"maxSag": 0.2}
		tessellate := "algo.tessellate([1], maxseg, -1, -1, True, 0, 1, 0.000000, False, False, True, False)"
		tessellate = strings.Replace(tessellate, "maxseg", strconv.FormatFloat(float64(parms["tessellate"].(map[string]float32)["maxSag"]), 'f', -1, 32), -1)
		scriptsStr.WriteString(tessellate + ";")
	} else {
		var Tessellate = struct {
			MaxSag float64 `json:"maxSag"`
		}{}
		valMap, err := json.Marshal(parms["tessellate"])
		if err == nil {
			json.Unmarshal(valMap, &Tessellate)
			tessellate := "algo.tessellate([1], maxseg, -1, -1, True, 0, 1, 0.000000, False, False, True, False)"
			tessellate = strings.Replace(tessellate, "maxseg", strconv.FormatFloat(Tessellate.MaxSag, 'f', -1, 64), -1)
			scriptsStr.WriteString(tessellate + ";")
		}
	}
	if _, ok := parms["mergeByTreeLevel"]; !ok {
		parms["mergeByTreeLevel"] = map[string]int{"maxLevel": 5}
		mergeByTreeLevel := "scene.mergeByTreeLevel([1], maxLevel, 2)"
		mergeByTreeLevel = strings.Replace(mergeByTreeLevel, "maxLevel", strconv.Itoa(parms["mergeByTreeLevel"].(map[string]int)["maxLevel"]), -1)
		scriptsStr.WriteString(mergeByTreeLevel + ";")
	} else {
		var MergeByTreeLevel = struct {
			MaxLevel int `json:"maxLevel"`
		}{}
		valMap, err := json.Marshal(parms["mergeByTreeLevel"])
		if err == nil {
			json.Unmarshal(valMap, &MergeByTreeLevel)
			mergeByTreeLevel := "scene.mergeByTreeLevel([1], maxLevel, 2)"
			mergeByTreeLevel = strings.Replace(mergeByTreeLevel, "maxLevel", strconv.Itoa(MergeByTreeLevel.MaxLevel), -1)
			scriptsStr.WriteString(mergeByTreeLevel + ";")
		}
	}
	log.Debugf("scriptsStr1:%s", scriptsStr.String())
	if v, ok := parms["mergeMaterials"]; !ok || v == nil {
		parms["mergeMaterials"] = map[string]bool{"enable": true}
		mergeMaterials := "scene.mergeMaterials()"
		scriptsStr.WriteString(mergeMaterials + ";")
	} else {
		var MergeMaterials = struct {
			Enable bool `json:"enable"`
		}{}
		valMap, err := json.Marshal(parms["mergeMaterials"])
		if err == nil {
			json.Unmarshal(valMap, &MergeMaterials)
			if MergeMaterials.Enable {
				mergeMaterials := "scene.mergeMaterials()"
				scriptsStr.WriteString(mergeMaterials + ";")
			}
		}
	}
	log.Debugf("scriptsStr2:%s", scriptsStr.String())
	if v, ok := parms["identifyInstances"]; !ok || v == nil {
		parms["identifyInstances"] = map[string]interface{}{"minOccurrenceCount": 2, "enabled": true}
		identifyInstances := "scene.identifyInstances(minOccurrenceCount)"
		identifyInstances = strings.Replace(identifyInstances, "minOccurrenceCount", strconv.Itoa(parms["identifyInstances"].(map[string]interface{})["minOccurrenceCount"].(int)), -1)
		scriptsStr.WriteString(identifyInstances + ";")
	} else {
		var IdentifyInstances = struct {
			MinOccurrenceCount int  `json:"minOccurrenceCount"`
			Enabled            bool `json:"enabled"`
		}{}
		valMap, err := json.Marshal(parms["identifyInstances"])
		if err == nil {
			json.Unmarshal(valMap, &IdentifyInstances)
			if IdentifyInstances.Enabled {
				identifyInstances := "scene.identifyInstances(minOccurrenceCount)"
				identifyInstances = strings.Replace(identifyInstances, "minOccurrenceCount", strconv.Itoa(IdentifyInstances.MinOccurrenceCount), -1)
				scriptsStr.WriteString(identifyInstances + ";")
			}
		}
	}
	log.Debugf("scriptsStr3:%s", scriptsStr.String())
	if v, ok := parms["deleteEmptyOccurrences"]; !ok || v == nil {
		parms["deleteEmptyOccurrences"] = map[string]bool{"enable": true}
	} else {
		var deleteEmptyOccurrences = struct {
			Enable bool `json:"enable"`
		}{}
		valMap, err := json.Marshal(parms["deleteEmptyOccurrences"])
		if err == nil {
			json.Unmarshal(valMap, &deleteEmptyOccurrences)
			if deleteEmptyOccurrences.Enable {
				deleteEmptyOccurrences := "scene.deleteEmptyOccurrences()"
				scriptsStr.WriteString(deleteEmptyOccurrences + ";")
			}
		}
	}
	log.Debugf("scriptsStr4:%s", scriptsStr.String())
	if v, ok := parms["selectByMaxSize"]; !ok || v == nil {
		parms["selectByMaxSize"] = map[string]interface{}{"maxDiagLength": 150, "enabled": false}
		selectByMaxSize := "scene.selectByMaximumSize([1], maxDiagLength, -1.000000, False)"
		selectByMaxSize = strings.Replace(selectByMaxSize, "maxDiagLength", strconv.FormatFloat(float64(parms["selectByMaxSize"].(map[string]interface{})["maxDiagLength"].(int)), 'f', -1, 64), -1)
		scriptsStr.WriteString(selectByMaxSize + ";")
	} else {
		var SelectByMaxSize = struct {
			MaxDiagLength float32 `json:"maxDiagLength"`
			Enabled       bool    `json:"enabled"`
		}{}
		valMap, err := json.Marshal(parms["selectByMaxSize"])
		if err == nil {
			json.Unmarshal(valMap, &SelectByMaxSize)
			if SelectByMaxSize.Enabled {
				selectByMaxSize := "scene.selectByMaximumSize([1], maxDiagLength, -1.000000, False)"
				selectByMaxSize = strings.Replace(selectByMaxSize, "maxDiagLength", strconv.FormatFloat(float64(SelectByMaxSize.MaxDiagLength), 'f', -1, 64), -1)
				scriptsStr.WriteString(selectByMaxSize + ";")
			}
		}
	}
	log.Infof("scriptsStr5:%s", scriptsStr.String())
	if v, ok := parms["rake"]; !ok || v == nil {
		parms["rake"] = map[string]bool{"enable": true}
		rake := "scene.rake(1, False)"
		scriptsStr.WriteString(rake + ";")
	} else {
		var Rake = struct {
			Enable bool `json:"enable"`
		}{}
		valMap, err := json.Marshal(parms["rake"])
		if err == nil {
			json.Unmarshal(valMap, &Rake)
			if Rake.Enable {
				rake := "scene.rake(1, False)"
				scriptsStr.WriteString(rake + ";")
			}
		}
	}
	log.Debugf("scriptsStr6:%s", scriptsStr.String())
	if v, ok := parms["removeHoles"]; !ok || v == nil {
		parms["removeHoles"] = map[string]interface{}{"maxDiameter": 15, "enabled": true}
		removeHoles := "algo.removeHoles([1], True, False, False, maxDiameter, 0)"
		removeHoles = strings.Replace(removeHoles, "maxDiameter", strconv.Itoa(parms["removeHoles"].(map[string]interface{})["maxDiameter"].(int)), -1)
	} else {
		var RemoveHoles = struct {
			MaxDiameter int  `json:"maxDiameter"`
			Enabled     bool `json:"enabled"`
		}{}
		valMap, err := json.Marshal(parms["removeHoles"])
		if err == nil {
			json.Unmarshal(valMap, &RemoveHoles)
			if RemoveHoles.Enabled {
				removeHoles := "algo.removeHoles([1], True, False, False, maxDiameter, 0)"
				removeHoles = strings.Replace(removeHoles, "maxDiameter", strconv.Itoa(RemoveHoles.MaxDiameter), -1)
				scriptsStr.WriteString(removeHoles + ";")
			}
		}
	}
	log.Debugf("scriptsStr7:%s", scriptsStr.String())
	if _, ok := parms["decimate"]; !ok {
		parms["decimate"] = map[string]interface{}{"method": "ratio", "value": map[string]float32{"ratio": 50}}
		decimate := "algo.decimateTarget([1], [\"ratio\",ratio], 0, False, 5000000)"
		decimate = strings.Replace(decimate, "surfaceTolerance", strconv.Itoa(int(parms["decimate"].(map[string]interface{})["value"].(map[string]float32)["ratio"])), -1)
		scriptsStr.WriteString(decimate + ";")
	} else {
		if val, ok1 := parms["decimate"]; ok1 {
			valMap, err := json.Marshal(val)
			if err != nil {
				log.Errorf("json.Marshal err:%v", err)
			} else {
				var DecimateStruct = struct {
					Method string `json:"method"`
					Value  struct {
						Ratio             float64 `json:"ratio"`
						TriangleCount     int     `json:"triangleCount"`
						SurfacicTolerance float64 `json:"surfacicTolerance"`
					} `json:"value"`
				}{}
				json.Unmarshal(valMap, &DecimateStruct)
				var decimateTarget string
				switch DecimateStruct.Method {
				case "Target":
					decimateTarget = "algo.decimateTarget([1], [\"triangleCount\",target_count], 1, True, 5000000)"
					decimateTarget = strings.Replace(decimateTarget, "target_count", strconv.Itoa(DecimateStruct.Value.TriangleCount), -1)
				case "Quality":
					decimateTarget = "algo.decimate([1], surfaceTolerance, 0.100000, 1.000000, -1, False)"
					decimateTarget = strings.Replace(decimateTarget, "surfaceTolerance", strconv.FormatFloat(DecimateStruct.Value.SurfacicTolerance, 'f', -1, 64), -1)
				default:
					decimateTarget = "algo.decimateTarget([1], [\"ratio\",ratio_vaule], 0, False, 5000000)"
					decimateTarget = strings.Replace(decimateTarget, "ratio_vaule", strconv.FormatFloat(DecimateStruct.Value.Ratio, 'f', -1, 64), -1)

				}
				scriptsStr.WriteString(decimateTarget + ";")
			}
		}
	}
	log.Debugf("scriptsStr8:%s", scriptsStr.String())
	scriptsStr.WriteString("scene.deleteEmptyOccurrences(1);")
	return parms, scriptsStr.String()
}
