package interactivity

import "crypto/rand"

type InteractivityTemplate struct{}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomString(length int) string {
	b := make([]byte, length)
	_, _ = rand.Read(b)

	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}

	return string(b)
}

// func GetInteractivity(ID string) Interactivity {
// 	i, ok := interactives[ID]
// 	if !ok {
// 		return nil
// 	}
//
// 	return i
// }
//
// func GetInteractivityByCustomID(customID string) (result []Interactivity) {
// 	for _, interactive := range interactives {
// 		if interactive.CustomID() == customID {
// 			result = append(result, interactive)
// 		}
// 	}
//
// 	return result
// }
//
// func CancelInteractivity(ID string) {
// 	i, ok := interactives[ID]
// 	if !ok {
// 		return
// 	}
//
// 	i.Cancel()
// }
