package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/olekukonko/tablewriter"
)

// JSONDiff represents the whole diff
type JSONDiff struct {
	Rows []JSONDiffRow
}

// JSONDiffRow represents a single non-existent subset item
type JSONDiffRow struct {
	Key      string
	Expected interface{}
	Got      interface{}
}

// PrintJSONDiff prints JSONDiff to STDERR
func PrintJSONDiff(diff *JSONDiff) {
	data := make([][]string, 0)

	for _, r := range diff.Rows {
		data = append(data, []string{r.Key, fmt.Sprintf("%s", r.Expected), fmt.Sprintf("%s", r.Got)})
	}

	table := tablewriter.NewWriter(os.Stderr)
	table.SetHeader([]string{"Key", "Expected", "Got"})
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}

// JSONIsSubset checks if a is a subset json of b
func JSONIsSubset(a, b []byte) (bool, *JSONDiff, error) {
	return jsonIsSubsetR(a, b, nil, nil)
}

func jsonIsSubsetR(a, b []byte, diff *JSONDiff, prefix interface{}) (bool, *JSONDiff, error) {
	// Initialize
	if diff == nil {
		diff = &JSONDiff{}
	}
	if diff.Rows == nil {
		diff.Rows = make([]JSONDiffRow, 0)
	}

	// Prefix for keeping around more info (path of the diffs)
	sprefix := ""
	if prefix != nil {
		sprefix = prefix.(string)
	}

	// Unmarshal both interfaces. If something fails here, we have nothing to do
	// jai: JSON A Interface
	// jbi: JSON B Interface
	var jai, jbi interface{}
	if err := json.Unmarshal(a, &jai); err != nil {
		return false, nil, err
	}
	if err := json.Unmarshal(b, &jbi); err != nil {
		return false, nil, err
	}

	// Switch JSON (map) or array of JSON (array of interface)
	// ja: JSON A (map or []interface)
	// jb: JSON B (map or []interface)
	switch ja := jai.(type) {
	case map[string]interface{}:
		// Cast B to same type as A
		// TODO: Add a check to see if this fails
		jb := jbi.(map[string]interface{})

		// Iterate all keys of ja and check if each is present
		// and equal to the same key in jb
		for k, vu := range ja {
			switch vu.(type) {
			// A primitive value such as string or number will be compared natively
			default:
				// Check if we have the key at all
				if val, ok := jb[k]; ok {
					// Check if the key matches if we have it
					if vu != val {
						diff.Rows = append(diff.Rows, JSONDiffRow{
							Key: fmt.Sprintf("%s/%s", sprefix, k), Expected: vu, Got: jb[k]})
					}
				} else {
					// We didn't find a key we wanted
					diff.Rows = append(diff.Rows, JSONDiffRow{
						Key: fmt.Sprintf("%s/%s", sprefix, k), Expected: vu, Got: "NOT FOUND"})
				}

			// Compare nested json by calling this function recursively
			case map[string]interface{}, []interface{}:
				sja, err := json.Marshal(vu)
				if err != nil {
					return false, nil, err
				}
				sjb, err := json.Marshal(jb[k])
				if err != nil {
					return false, nil, err
				}
				_, _, err = jsonIsSubsetR(sja, sjb, diff, fmt.Sprintf("%s/%s", sprefix, k))
				if err != nil {
					return false, nil, err
				}
			}
		}

	// Compare arrays
	case []interface{}:
		// Case jbi to an array as well
		// TODO: Add a check to see if this fails
		jb := jbi.([]interface{})

		// Check if length is equal first
		if len(jb) != len(ja) {
			// Length not equal so that is not good
			diff.Rows = append(diff.Rows, JSONDiffRow{
				Key: fmt.Sprintf("%s", sprefix), Expected: fmt.Sprintf("LEN=%d", len(ja)), Got: fmt.Sprintf("LEN=%d", len(jb))})
		} else {
			// Recurse for each object inside
			for i, x := range ja {
				sja, err := json.Marshal(x)
				if err != nil {
					return false, nil, err
				}
				sjb, err := json.Marshal(jb[i])
				if err != nil {
					return false, nil, err
				}
				_, _, err = jsonIsSubsetR(sja, sjb, diff, fmt.Sprintf("%s[%d]", sprefix, i))
				if err != nil {
					return false, nil, err
				}
			}
		}
	// Compare primitive types directly
	default:
		return jai == jbi, diff, nil
	}

	// No diff means all keys in A were found and equal in B
	return diff == nil || len(diff.Rows) == 0, diff, nil
}

func main() {

	// Open our jsonFile
	jsonFile, err := os.Open("kc.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened kc.json")
	// defer the closing of our jsonFile so that we can parse it later on]

	jsonFile2, err := os.Open("render.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened render.json")
	// defer the closing of our jsonFile so that we can parse it later on

	defer jsonFile.Close()

	defer jsonFile2.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	byteValue2, _ := ioutil.ReadAll(jsonFile2)

	isSubset, diff, err := JSONIsSubset(byteValue, byteValue2)
	if isSubset || err != nil {
		fmt.Println("KubeletConfig is a subset of Rendered MachineConfig")
	} else {
		fmt.Println("KubeletConfig is not a subset of Rendered MachineConfig")
		PrintJSONDiff(diff)
	}

}
