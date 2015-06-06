package ergo

import (
	"strings"
)

func firstPathParam(chunk string) (int, int) {
	openedFrom := strings.Index(chunk, "{")
	openedTo := strings.Index(chunk, "}")
	if openedFrom < 0 || openedTo < 0 {
		return 0, 0
	}
	if openedFrom >= openedTo {
		return 0, 0
	}
	if strings.Index(chunk[openedFrom+1:openedTo], "{") != -1 {
		return 0, 0
	}
	return openedFrom, openedTo
}

func containsParam(chunk string) bool {
	from, to := firstPathParam(chunk)
	if from == to {
		return false
	}
	return true
}

func matchParams(rpchunk string, pchunk string) (bool, string) {
	rplen := len(rpchunk)
	plen := len(pchunk)

	if rplen == 0 || plen == 0 {
		return false, ""
	}

	pFrom, pTo := firstPathParam(rpchunk) // no more params
	if pFrom == pTo {
		if rpchunk != pchunk {
			return false, ""
		}
		return true, ""
	}

	if pFrom > 0 { // param is not at the beginning
		if pchunk[:pFrom] != rpchunk[:pFrom] {
			return false, ""
		}

		rpchunk = rpchunk[pFrom:]
		pchunk = pchunk[pFrom:]
		rplen = len(rpchunk)
		plen = len(pchunk)
		pTo = pTo - pFrom
		pFrom = 0
	}

	var params string
	params += rpchunk[1:pTo] + ":"

	// param is at the beginning
	if rplen == pTo+1 { // if param is only thing left
		params += pchunk + ";"
		return true, params
	}

	var npFrom, npTo int
	var nc string

	if pTo+1 < rplen { // has more path
		npFrom, npTo = firstPathParam(rpchunk[pTo+1:]) // next path param
	}

	if npFrom == npTo { // no more path params
		nc = rpchunk[pTo+1:]
	} else { // more path params
		npFrom += pTo
		npTo += pTo
		nc = rpchunk[pTo+1 : npFrom+1] // from end of first param to start of next one
	}

	clen := len(nc)
	ci := strings.Index(pchunk, nc)

	if clen == 0 {
		params += pchunk + ";"
		return true, params
	}
	if clen > 0 && ci < 0 { // chunk isnt found
		return false, ""
	}
	params += pchunk[:ci] + ";"

	nm, np := matchParams(rpchunk[pTo+clen:], pchunk[ci+clen-1:])
	if !nm {
		return false, ""
	}
	return true, params + np
}

func match(rpath string, path string) (matches bool, remaining string, params string) {
	fp := strings.Index(rpath, "{")
	if fp > 0 { // param not at start
		pre := rpath[:fp]
		if !strings.HasPrefix(path, pre) {
			remaining = path
			return
		}
		rpath = rpath[fp:]
		path = path[fp:]
	} else if fp == -1 { // no pathparam
		if strings.HasPrefix(path, rpath+"/") {
			matches = true
			remaining = strings.TrimPrefix(path, rpath+"/")
			return
		} else if path == rpath {
			matches = true
			remaining = ""
			return
		}
		matches = false
		remaining = path
		return
	}

	rsi := strings.Index(rpath, "/")
	si := strings.Index(path, "/")

	var rc, c string
	repeat := true
	for repeat {

		if rsi > -1 {
			rc = rpath[:rsi]
		} else {
			rc = rpath
			repeat = false
		}

		if si > -1 {
			c = path[:si]
		} else {
			c = path
		}
		if !containsParam(rc) {
			if rc != c {
				matches = false
				remaining = path
				params = ""
				return
			}
		} else {
			nm, np := matchParams(rc, c)
			if !nm {
				matches = false
				remaining = path
				params = ""
				return
			} else {
				params = params + np
			}
		}
		rpath = rpath[rsi+1:]
		path = path[si+1:]

		if len(path) == 0 && len(rpath) > 0 {
			matches = false
			remaining = path
			params = ""
			return
		}

		rsi = strings.Index(rpath, "/")
		si = strings.Index(path, "/")
	}
	remaining = path
	if si == -1 {
		remaining = ""
	}
	matches = true
	return
}
