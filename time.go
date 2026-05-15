// Copyright ©2026 xiayoudi. All rights reserved.
// Author: xiayoudi
// Email: ur@xiaud.com

package ud

import "time"

var (
	CNLoc     = time.FixedZone("CST", 8*3600)
	UnixEpoch = time.Unix(0, 0).UTC()
)

func Now() time.Time {
	return time.Now().In(CNLoc)
}

func ToCN(t time.Time) time.Time {
	return t.In(CNLoc)
}

func ToUTC(t time.Time) time.Time {
	return t.UTC()
}
