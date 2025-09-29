// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package log4j

import (
	"bytes"
	"fmt"
	"strings"
)

/*type formatCacheType struct {
	LastUpdateSeconds    int64
	shortTime, shortDate string
	longTime, longDate   string
}

var formatCache = &formatCacheType{}*/

func formatLogRecord(format string, rec *logRecord) string {
	if rec == nil {
		return "\n"
	}
	if len(format) == 0 {
		return ""
	}

	out := bytes.NewBuffer(make([]byte, 0, 64))
	/*secs := rec.Created.UnixNano() / 1e9
	  cache := *formatCache
	  if cache.LastUpdateSeconds != secs {
	      month, day, year := rec.Created.Month(), rec.Created.Day(), rec.Created.Year()
	      hour, minute, second := rec.Created.Hour(), rec.Created.Minute(), rec.Created.Second()
	      zone, _ := rec.Created.Zone()
	      ms := rec.Created.UnixNano() / 1e6 % 1000
	      updated := &formatCacheType{
	          LastUpdateSeconds: secs,
	          shortTime:         fmt.Sprintf("%02d:%02d", hour, minute),
	          shortDate:         fmt.Sprintf("%02d/%02d/%02d", month, day, year%100),
	          longTime:          fmt.Sprintf("%02d:%02d:%02d.%03d %s", hour, minute, second, ms, zone),
	          longDate:          fmt.Sprintf("%04d/%02d/%02d", year, month, day),
	      }
	      cache = *updated
	      formatCache = updated
	  }*/

	// Split the string into pieces by % signs
	pieces := bytes.Split([]byte(format), []byte{'%'})

	// Iterate over the pieces, replacing known formats
	for i, piece := range pieces {
		if i > 0 && len(piece) > 0 {
			switch piece[0] {
			case 'T':
				zone, _ := rec.Created.Zone()
				out.WriteString(rec.Created.Format("15:04:05"))
				out.WriteString(fmt.Sprintf(".%03d", rec.Created.UnixNano()/1e6%1000))
				out.WriteByte(' ')
				out.WriteString(zone)
			case 't':
				out.WriteString(rec.Created.Format("15:04"))
			case 'D':
				out.WriteString(rec.Created.Format("2006/01/02"))
			case 'd':
				out.WriteString(rec.Created.Format("01/02/"))
				out.WriteString(rec.Created.Format("2006")[2:])
			case 'L':
				out.WriteString(levelStrings[rec.Level])
			case 'S':
				out.WriteString(rec.Source)
			case 's':
				sources := strings.Split(rec.Source, "/")
				out.WriteString(sources[len(sources)-1])
			case 'f':
				sources := strings.Split(rec.Source, "/")
				names := strings.Split(sources[len(sources)-1], ".")
				out.WriteString(names[len(names)-1])
			case 'M':
				out.WriteString(rec.Message)
			case 'B':
				out.WriteByte('\n')
			}
			if len(piece) > 1 {
				out.Write(piece[1:])
			}
		} else if len(piece) > 0 {
			out.Write(piece)
		}
	}
	out.WriteByte('\n')

	return out.String()
}
