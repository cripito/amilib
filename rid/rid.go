package rid

import (
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"golang.org/x/exp/rand"
)

func New(kind string) (string, error) {

	entropy := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	ms := ulid.Timestamp(time.Now())
	ulIDvalue, err := ulid.New(ms, entropy)
	if err != nil {
		return "", err
	}

	return kind + "." + strings.ToLower(ulIDvalue.String()), nil
}
