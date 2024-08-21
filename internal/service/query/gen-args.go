package query

import (
	"fmt"
	"strings"
)

func genInsertFArgs(fields []string) (string, string) {
	fArgs := make([]string, len(fields))
	for i := range fArgs {
		fArgs[i] = "?"
	}
	return strings.Join(fields, ","), strings.Join(fArgs, ",")
}

func genUpdateFArgs(fields []string) string {
	var args []string
	for _, f := range fields {
		args = append(args, fmt.Sprintf("%s=?", f))
	}
	return strings.Join(args, ",")
}

func genUpdateOnConflictFArgs(fields []string) string {
	var args []string
	for _, f := range fields {
		args = append(args, fmt.Sprintf("%s=excluded.%s", f, f))
	}
	return strings.Join(args, ",")
}

func genUpdateFArgsToNull(fields []string) string {
	var args []string
	for _, f := range fields {
		args = append(args, fmt.Sprintf("%s=null", f))
	}
	return strings.Join(args, ",")
}
