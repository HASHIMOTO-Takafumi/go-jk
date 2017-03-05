package jk

import (
	"strconv"
	"strings"
)

//Morpheme is a morpheme
type Morpheme struct {
	Doukeis       Morphemes
	Surface       string //Midashi
	Pronunciation string //Yomi
	RootForm      string //Genkei
	Pos0          string //Hinshi
	Pos0ID        int    //HinshiID
	Pos1          string //Bunrui
	Pos1ID        int    //BunruiID
	CType         string //Katsuyou1
	CTypeID       int    //Katsuyou1ID
	CForm         string //Katsuyou2
	CFormID       int    //Katsuyou2ID
	Seminfo       string //Seminfo
	Rep           string //Rep
	Features      Features
}

//NewMorpheme returns a morpheme for the given line
func NewMorpheme(line string) (*Morpheme, error) {
	var items []string
	if strings.HasPrefix(line, `  \  \  `) {
		items = make([]string, 12)
		items[0] = " "
		items[1] = " "
		items[2] = " "
		rest := strings.TrimPrefix(line, `  \  \  `)
		for i, s := range strings.SplitN(rest, " ", 9) {
			items[3+i] = s
		}
	} else {
		items = strings.SplitN(line, " ", 12)
	}

	mrph := new(Morpheme)
	mrph.Surface = items[0]
	mrph.Pronunciation = items[1]
	mrph.RootForm = items[2]
	mrph.Pos0 = items[3]

	hinshiID, err := strconv.Atoi(items[4])
	if err != nil {
		return nil, err
	}
	mrph.Pos0ID = hinshiID

	mrph.Pos1 = items[5]
	bunruiID, err := strconv.Atoi(items[6])
	if err != nil {
		return nil, err
	}
	mrph.Pos1ID = bunruiID

	mrph.CType = items[7]
	katsuyo1ID, err := strconv.Atoi(items[8])
	if err != nil {
		return nil, err
	}
	mrph.CTypeID = katsuyo1ID

	mrph.CForm = items[9]
	katsuyo2ID, err := strconv.Atoi(items[10])
	if err != nil {
		return nil, err
	}
	mrph.CFormID = katsuyo2ID

	rest := items[11]
	seminfoStartPos := strings.Index(rest, "\"")
	featureStart := 0
	if seminfoStartPos == -1 {
		mrph.Seminfo = ""
		mrph.Rep = mrph.RootForm + "/" + mrph.RootForm
		featureStart = strings.Index(rest, "<")
	} else {
		seminfoCharNum := strings.Index(rest[seminfoStartPos+1:], "\"")
		mrph.Seminfo = rest[seminfoStartPos+1 : seminfoStartPos+1+seminfoCharNum]

		retName := "代表表記:"
		repStart := strings.Index(mrph.Seminfo, retName)
		if repStart != -1 {
			repEnd := strings.Index(mrph.Seminfo[repStart:], " ")
			if repEnd == -1 {
				repEnd = len(mrph.Seminfo)
			} else {
				repEnd += repStart
			}
			mrph.Rep = mrph.Seminfo[repStart+len(retName) : repEnd]
		} else {
			mrph.Rep = mrph.RootForm + "/" + mrph.RootForm
		}
		featureStart = seminfoStartPos + 1 + seminfoCharNum + 2
	}

	if featureStart >= 0 && featureStart < len(rest) {
		mrph.Features = getFeatures(rest[featureStart:], '>', 1)
	}

	return mrph, err
}
