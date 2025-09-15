package util

import (
	"fmt"
)

func GetRandomProfilePic(index int) string {
	return fmt.Sprintf("https://randomuser.me/api/portraits/lego/%d.jpg", index)
}
