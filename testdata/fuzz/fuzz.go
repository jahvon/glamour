package fuzzing

import "github.com/jahvon/glamour"

func Fuzz(data []byte) int {
	_, err := glamour.RenderBytes(data, glamour.DarkStyle)
	if err != nil {
		return 0
	}
	return 1
}
