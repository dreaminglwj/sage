package builder

import (
	"fmt"
	"strings"

	"xorm.io/builder"
)

type Like [2]string

var _ builder.Cond = Like{"", ""}

// WriteTo write SQL to Writer
func (like Like) WriteTo(w builder.Writer) error {
	if _, err := fmt.Fprintf(w, "%s LIKE ?", like[0]); err != nil {
		return err
	}
	w.Append("%" + WithSpecialCharacterEscapeSequences(like[1]) + "%")
	return nil
}

// And implements And with other conditions
func (like Like) And(conds ...builder.Cond) builder.Cond {
	return builder.And(like, builder.And(conds...))
}

// Or implements Or with other conditions
func (like Like) Or(conds ...builder.Cond) builder.Cond {
	return builder.Or(like, builder.Or(conds...))
}

// IsValid tests if this condition is valid
func (like Like) IsValid() bool {
	return len(like[0]) > 0 && len(like[1]) > 0
}

func WithSpecialCharacterEscapeSequences(keywords string) string {
	//"'" is supported
	var escapedCharters = []string{"\\", "\"", "_", "%"}
	var newKeywords = keywords
	for _, charter := range escapedCharters {
		if strings.Count(keywords, charter) > 0 {
			newKeywords = strings.ReplaceAll(newKeywords, charter, "\\"+charter)
		}
	}
	return newKeywords
}
