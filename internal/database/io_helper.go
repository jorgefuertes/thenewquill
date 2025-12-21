package database

import (
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
)

func recordFromLine(line string) (uint32, Record, error) {
	matches := recordRg.FindAllStringSubmatch(line, 1)

	if len(matches[0]) < 5 {
		return 0, Record{}, fmt.Errorf("invalid record line: %s", line)
	}

	i, err := strconv.Atoi(matches[0][1])
	if err != nil {
		return 0, Record{}, err
	}

	labelID, err := strconv.Atoi(matches[0][2])
	if err != nil {
		return 0, Record{}, err
	}

	k, err := strconv.Atoi(matches[0][3])
	if err != nil {
		return 0, Record{}, err
	}

	data, err := base64.StdEncoding.DecodeString(matches[0][4])
	if err != nil {
		return 0, Record{}, err
	}

	r := Record{
		LabelID: uint32(labelID),
		Kind:    kind.KindFromByte(byte(k)),
		Data:    []byte(data),
	}

	return uint32(i), r, nil
}

func labelFromLine(line string) (uint32, string, error) {
	matches := labelRg.FindAllStringSubmatch(line, 1)

	if len(matches[0]) < 3 {
		return 0, "", fmt.Errorf("invalid label line: %s", line)
	}

	i, err := strconv.Atoi(matches[0][1])
	if err != nil {
		return 0, "", err
	}

	return uint32(i), matches[0][2], nil
}
