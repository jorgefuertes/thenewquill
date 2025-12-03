package io

// import (
// 	"bytes"
// 	"compress/gzip"
// 	"fmt"
// 	"os"
// 	"reflect"
// 	"strings"

// 	"github.com/fxamacker/cbor/v2"
// 	"github.com/jorgefuertes/thenewquill/internal/adventure/database"
// 	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
// 	"github.com/jorgefuertes/thenewquill/internal/util"
// 	"github.com/jorgefuertes/thenewquill/pkg/log"
// 	"github.com/jorgefuertes/thenewquill/pkg/tms"
// )

// const (
// 	paramTitle       = "title"
// 	paramDescription = "description"
// 	paramAuthor      = "author"
// 	paramVersion     = "version"
// 	paramDate        = "date"
// 	paramLanguage    = "language"

// 	commentChar   = "#"
// 	dataSeparator = "---"
// )

// var headerSectionSeparator = strings.Repeat(commentChar, 80)

// func Export(db *db.DB, path string) (int, error) {
// 	headers := composeExportHeaders(db)

// 	cbData, err := cbor.Marshal(db)
// 	if err != nil {
// 		return 0, err
// 	}

// 	f, err := os.Create(path)
// 	if err != nil {
// 		return 0, err
// 	}

// 	f.WriteString(headers)
// 	f.WriteString(dataSeparator + "\n")

// 	plainBuffer := bytes.NewBuffer(nil)
// 	gz := gzip.NewWriter(plainBuffer)

// 	n, err := gz.Write(cbData)
// 	if err != nil {
// 		return 0, err
// 	}

// 	if err := gz.Close(); err != nil {
// 		log.Warning("error closing file %q: %s", path, err)
// 	}

// 	key := tms.GenerateKey(headers)
// 	encrypted, err := tms.Encrypt(key, plainBuffer.Bytes())
// 	if err != nil {
// 		return 0, err
// 	}

// 	if _, err := f.Write(encrypted); err != nil {
// 		return 0, err
// 	}

// 	if err := f.Close(); err != nil {
// 		log.Warning("error closing file %q: %s", path, err)
// 	}

// 	return n, nil
// }

// func composeExportHeaders(db *db.DB) string {
// 	cfgParams := getConfigParams(db)
// 	var headers []string

// 	title := strToHeaderLines(cfgParams[paramTitle])
// 	author := strToHeaderLines(cfgParams[paramAuthor])
// 	description := strToHeaderLines(cfgParams[paramDescription])
// 	version := strToHeaderLines(fmt.Sprintf("Version.: %s %s", cfgParams[paramVersion], cfgParams[paramLanguage]))
// 	ts := strToHeaderLines("Compiled: " + cfgParams[paramDate])

// 	headers = append(headers, headerSectionSeparator)
// 	headers = append(headers, "# "+strings.Repeat(" ", 76)+" #")
// 	headers = append(headers, title...)
// 	headers = append(headers, author...)
// 	headers = append(headers, description...)
// 	headers = append(headers, "# "+strings.Repeat(" ", 76)+" #")
// 	headers = append(headers, version...)
// 	headers = append(headers, ts...)
// 	headers = append(headers, "# "+strings.Repeat(" ", 76)+" #")
// 	headers = append(headers, headerSectionSeparator)

// 	return strings.Join(headers, "\n") + "\n"
// }

// func strToHeaderLines(s string) []string {
// 	lines := util.SplitIntoLines(s, 76)
// 	for i, line := range lines {
// 		lines[i] = fmt.Sprintf("# %-76s #", line)
// 	}

// 	return lines
// }

// func getConfigParams(db *db.DB) map[string]string {
// 	cfgParams := map[string]string{
// 		paramTitle:       "",
// 		paramDescription: "",
// 		paramAuthor:      "",
// 		paramVersion:     "",
// 		paramDate:        "",
// 		paramLanguage:    "",
// 	}

// 	for _, r := range db.Data {
// 		if kind.KindOf(r) == kind.Param {
// 			l := db.GetLabelName(r.GetID())

// 			v := reflect.ValueOf(r).FieldByName("V")
// 			if v.IsValid() && v.Kind() == reflect.String {
// 				cfgParams[l] = v.String()
// 			} else {
// 				log.Warning("cannot get config param %q: %+v", l, r)
// 			}
// 		}
// 	}

// 	return cfgParams
// }
