package utils

import (
	"github.com/xuri/excelize/v2"
)

type ExcelizeStyle struct {
	Color               string
	HorizontalAlignment string
	VerticalAlignment   string
	Border              string
	FontColor           string
	Bold                bool
	Italic              bool
	TextRotation        int
}

func (rs ExcelizeStyle) GenerateStyle() (s *excelize.Style) {

	var border []excelize.Border
	alignment := excelize.Alignment{}
	font := excelize.Font{}

	color := make([]string, 1)
	if rs.Color != "" {
		color[0] = rs.Color
	} else {
		color[0] = "#FFFFFF"
	}
	if rs.HorizontalAlignment != "" {
		alignment.Horizontal = rs.HorizontalAlignment
	}
	if rs.VerticalAlignment != "" {
		alignment.Vertical = rs.VerticalAlignment
	}
	if rs.TextRotation > 0 {
		alignment.TextRotation = rs.TextRotation
	}
	if rs.Border == "all" {
		topBorder := excelize.Border{
			Type:  "top",
			Color: "#000000",
			Style: 1,
		}
		btmBorder := excelize.Border{
			Type:  "bottom",
			Color: "#000000",
			Style: 1,
		}
		leftBorder := excelize.Border{
			Type:  "left",
			Color: "#000000",
			Style: 1,
		}
		rightBorder := excelize.Border{
			Type:  "right",
			Color: "#000000",
			Style: 1,
		}
		border = append(border, topBorder)
		border = append(border, btmBorder)
		border = append(border, leftBorder)
		border = append(border, rightBorder)
	}

	if rs.FontColor != "" {
		font.Color = rs.FontColor
	}

	if rs.Bold {
		font.Bold = true
	}

	if rs.Italic {
		font.Italic = true
	}

	s = &excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   color,
			Pattern: 1,
		},
		Alignment: &alignment,
		Font:      &font,
		Border:    border,
	}
	return s
}

func GetExcelColumns() (columns []string) {
	columns = append(columns, "A")
	columns = append(columns, "B")
	columns = append(columns, "C")
	columns = append(columns, "D")
	columns = append(columns, "E")
	columns = append(columns, "F")
	columns = append(columns, "G")
	columns = append(columns, "H")
	columns = append(columns, "I")
	columns = append(columns, "J")
	columns = append(columns, "K")
	columns = append(columns, "L")
	columns = append(columns, "M")
	columns = append(columns, "N")
	columns = append(columns, "O")
	columns = append(columns, "P")
	columns = append(columns, "Q")
	columns = append(columns, "R")
	columns = append(columns, "S")
	columns = append(columns, "T")
	columns = append(columns, "U")
	columns = append(columns, "V")
	columns = append(columns, "W")
	columns = append(columns, "X")
	columns = append(columns, "Y")
	columns = append(columns, "Z")
	columns = append(columns, "AA")
	columns = append(columns, "AB")
	columns = append(columns, "AC")
	columns = append(columns, "AD")
	columns = append(columns, "AE")
	columns = append(columns, "AF")
	columns = append(columns, "AG")
	columns = append(columns, "AH")
	columns = append(columns, "AI")
	columns = append(columns, "AJ")
	columns = append(columns, "AK")
	columns = append(columns, "AL")
	columns = append(columns, "AM")
	columns = append(columns, "AN")
	columns = append(columns, "AO")
	columns = append(columns, "AP")
	columns = append(columns, "AQ")
	columns = append(columns, "AR")
	columns = append(columns, "AS")
	columns = append(columns, "AT")
	columns = append(columns, "AU")
	columns = append(columns, "AV")
	columns = append(columns, "AW")
	columns = append(columns, "AX")
	columns = append(columns, "AY")
	columns = append(columns, "AZ")
	columns = append(columns, "BA")
	columns = append(columns, "BB")
	columns = append(columns, "BC")
	columns = append(columns, "BD")
	columns = append(columns, "BE")
	columns = append(columns, "BF")
	columns = append(columns, "BG")
	columns = append(columns, "BH")
	columns = append(columns, "BI")
	columns = append(columns, "BJ")
	columns = append(columns, "BK")
	columns = append(columns, "BL")
	columns = append(columns, "BM")
	columns = append(columns, "BN")
	columns = append(columns, "BO")
	columns = append(columns, "BP")
	columns = append(columns, "BQ")
	columns = append(columns, "BR")
	columns = append(columns, "BS")
	columns = append(columns, "BT")
	columns = append(columns, "BU")
	columns = append(columns, "BV")
	columns = append(columns, "BW")
	columns = append(columns, "BX")
	columns = append(columns, "BY")
	columns = append(columns, "BZ")

	columns = append(columns, "CA")
	columns = append(columns, "CB")
	columns = append(columns, "CC")
	columns = append(columns, "CD")
	columns = append(columns, "CE")
	columns = append(columns, "CF")
	columns = append(columns, "CG")
	columns = append(columns, "CH")
	columns = append(columns, "CI")
	columns = append(columns, "CJ")
	columns = append(columns, "CK")
	columns = append(columns, "CL")
	columns = append(columns, "CM")
	columns = append(columns, "CN")
	columns = append(columns, "CO")
	columns = append(columns, "CP")
	columns = append(columns, "CQ")
	columns = append(columns, "CR")
	columns = append(columns, "CS")
	columns = append(columns, "CT")
	columns = append(columns, "CU")
	columns = append(columns, "CV")
	columns = append(columns, "CW")
	columns = append(columns, "CX")
	columns = append(columns, "CY")
	columns = append(columns, "CZ")

	columns = append(columns, "DA")
	columns = append(columns, "DB")
	columns = append(columns, "DC")
	columns = append(columns, "DD")
	columns = append(columns, "DE")
	columns = append(columns, "DF")
	columns = append(columns, "DG")
	columns = append(columns, "DH")
	columns = append(columns, "DI")
	columns = append(columns, "DJ")
	columns = append(columns, "DK")
	columns = append(columns, "DL")
	columns = append(columns, "DM")
	columns = append(columns, "DN")
	columns = append(columns, "DO")
	columns = append(columns, "DP")
	columns = append(columns, "DQ")
	columns = append(columns, "DR")
	columns = append(columns, "DS")
	columns = append(columns, "DT")
	columns = append(columns, "DU")
	columns = append(columns, "DV")
	columns = append(columns, "DW")
	columns = append(columns, "DX")
	columns = append(columns, "DY")
	columns = append(columns, "DZ")

	columns = append(columns, "EA")
	columns = append(columns, "EB")
	columns = append(columns, "EC")
	columns = append(columns, "ED")
	columns = append(columns, "EE")
	columns = append(columns, "EF")
	columns = append(columns, "EG")
	columns = append(columns, "EH")
	columns = append(columns, "EI")
	columns = append(columns, "EJ")
	columns = append(columns, "EK")
	columns = append(columns, "EL")
	columns = append(columns, "EM")
	columns = append(columns, "EN")
	columns = append(columns, "EO")
	columns = append(columns, "EP")
	columns = append(columns, "EQ")
	columns = append(columns, "ER")
	columns = append(columns, "ES")
	columns = append(columns, "ET")
	columns = append(columns, "EU")
	columns = append(columns, "EV")
	columns = append(columns, "EW")
	columns = append(columns, "EX")
	columns = append(columns, "EY")
	columns = append(columns, "EZ")

	columns = append(columns, "FA")
	columns = append(columns, "FB")
	columns = append(columns, "FC")
	columns = append(columns, "FD")
	columns = append(columns, "FE")
	columns = append(columns, "FF")
	columns = append(columns, "FG")
	columns = append(columns, "FH")
	columns = append(columns, "FI")
	columns = append(columns, "FJ")
	columns = append(columns, "FK")
	columns = append(columns, "FL")
	columns = append(columns, "FM")
	columns = append(columns, "FN")
	columns = append(columns, "FO")
	columns = append(columns, "FP")
	columns = append(columns, "FQ")
	columns = append(columns, "FR")
	columns = append(columns, "FS")
	columns = append(columns, "FT")
	columns = append(columns, "FU")
	columns = append(columns, "FV")
	columns = append(columns, "FW")
	columns = append(columns, "FX")
	columns = append(columns, "FY")
	columns = append(columns, "FZ")

	columns = append(columns, "GA")
	columns = append(columns, "GB")
	columns = append(columns, "GC")
	columns = append(columns, "GD")
	columns = append(columns, "GE")
	columns = append(columns, "GF")
	columns = append(columns, "GG")
	columns = append(columns, "GH")
	columns = append(columns, "GI")
	columns = append(columns, "GJ")
	columns = append(columns, "GK")
	columns = append(columns, "GL")
	columns = append(columns, "GM")
	columns = append(columns, "GN")
	columns = append(columns, "GO")
	columns = append(columns, "GP")
	columns = append(columns, "GQ")
	columns = append(columns, "GR")
	columns = append(columns, "GS")
	columns = append(columns, "GT")
	columns = append(columns, "GU")
	columns = append(columns, "GV")
	columns = append(columns, "GW")
	columns = append(columns, "GX")
	columns = append(columns, "GY")
	columns = append(columns, "GZ")

	columns = append(columns, "HA")
	columns = append(columns, "HB")
	columns = append(columns, "HC")
	columns = append(columns, "HD")
	columns = append(columns, "HE")
	columns = append(columns, "HF")
	columns = append(columns, "HG")
	columns = append(columns, "HH")
	columns = append(columns, "HI")
	columns = append(columns, "HJ")
	columns = append(columns, "HK")
	columns = append(columns, "HL")
	columns = append(columns, "HM")
	columns = append(columns, "HN")
	columns = append(columns, "HO")
	columns = append(columns, "HP")
	columns = append(columns, "HQ")
	columns = append(columns, "HR")
	columns = append(columns, "HS")
	columns = append(columns, "HT")
	columns = append(columns, "HU")
	columns = append(columns, "HV")
	columns = append(columns, "HW")
	columns = append(columns, "HX")
	columns = append(columns, "HY")
	columns = append(columns, "HZ")

	columns = append(columns, "IA")
	columns = append(columns, "IB")
	columns = append(columns, "IC")
	columns = append(columns, "ID")
	columns = append(columns, "IE")
	columns = append(columns, "IF")
	columns = append(columns, "IG")
	columns = append(columns, "IH")
	columns = append(columns, "II")
	columns = append(columns, "IJ")
	columns = append(columns, "IK")
	columns = append(columns, "IL")
	columns = append(columns, "IM")
	columns = append(columns, "IN")
	columns = append(columns, "IO")
	columns = append(columns, "IP")
	columns = append(columns, "IQ")
	columns = append(columns, "IR")
	columns = append(columns, "IS")
	columns = append(columns, "IT")
	columns = append(columns, "IU")
	columns = append(columns, "IV")
	columns = append(columns, "IW")
	columns = append(columns, "IX")
	columns = append(columns, "IY")
	columns = append(columns, "IZ")

	columns = append(columns, "JA")
	columns = append(columns, "JB")
	columns = append(columns, "JC")
	columns = append(columns, "JD")
	columns = append(columns, "JE")
	columns = append(columns, "JF")
	columns = append(columns, "JG")
	columns = append(columns, "JH")
	columns = append(columns, "JI")
	columns = append(columns, "JJ")
	columns = append(columns, "JK")
	columns = append(columns, "JL")
	columns = append(columns, "JM")
	columns = append(columns, "JN")
	columns = append(columns, "JO")
	columns = append(columns, "JP")
	columns = append(columns, "JQ")
	columns = append(columns, "JR")
	columns = append(columns, "JS")
	columns = append(columns, "JT")
	columns = append(columns, "JU")
	columns = append(columns, "JV")
	columns = append(columns, "JW")
	columns = append(columns, "JX")
	columns = append(columns, "JY")
	columns = append(columns, "JZ")

	columns = append(columns, "KA")
	columns = append(columns, "KB")
	columns = append(columns, "KC")
	columns = append(columns, "KD")
	columns = append(columns, "KE")
	columns = append(columns, "KF")
	columns = append(columns, "KG")
	columns = append(columns, "KH")
	columns = append(columns, "KI")
	columns = append(columns, "KJ")
	columns = append(columns, "KK")
	columns = append(columns, "KL")
	columns = append(columns, "KM")
	columns = append(columns, "KN")
	columns = append(columns, "KO")
	columns = append(columns, "KP")
	columns = append(columns, "KQ")
	columns = append(columns, "KR")
	columns = append(columns, "KS")
	columns = append(columns, "KT")
	columns = append(columns, "KU")
	columns = append(columns, "KV")
	columns = append(columns, "KW")
	columns = append(columns, "KX")
	columns = append(columns, "KY")
	columns = append(columns, "KZ")

	columns = append(columns, "LA")
	columns = append(columns, "LB")
	columns = append(columns, "LC")
	columns = append(columns, "LD")
	columns = append(columns, "LE")
	columns = append(columns, "LF")
	columns = append(columns, "LG")
	columns = append(columns, "LH")
	columns = append(columns, "LI")
	columns = append(columns, "LJ")
	columns = append(columns, "LK")
	columns = append(columns, "LL")
	columns = append(columns, "LM")
	columns = append(columns, "LN")
	columns = append(columns, "LO")
	columns = append(columns, "LP")
	columns = append(columns, "LQ")
	columns = append(columns, "LR")
	columns = append(columns, "LS")
	columns = append(columns, "LT")
	columns = append(columns, "LU")
	columns = append(columns, "LV")
	columns = append(columns, "LW")
	columns = append(columns, "LX")
	columns = append(columns, "LY")
	columns = append(columns, "LZ")

	columns = append(columns, "MA")
	columns = append(columns, "MB")
	columns = append(columns, "MC")
	columns = append(columns, "MD")
	columns = append(columns, "ME")
	columns = append(columns, "MF")
	columns = append(columns, "MG")
	columns = append(columns, "MH")
	columns = append(columns, "MI")
	columns = append(columns, "MJ")
	columns = append(columns, "MK")
	columns = append(columns, "ML")
	columns = append(columns, "MM")
	columns = append(columns, "MN")
	columns = append(columns, "MO")
	columns = append(columns, "MP")
	columns = append(columns, "MQ")
	columns = append(columns, "MR")
	columns = append(columns, "MS")
	columns = append(columns, "MT")
	columns = append(columns, "MU")
	columns = append(columns, "MV")
	columns = append(columns, "MW")
	columns = append(columns, "MX")
	columns = append(columns, "MY")
	columns = append(columns, "MZ")

	columns = append(columns, "NA")
	columns = append(columns, "NB")
	columns = append(columns, "NC")
	columns = append(columns, "ND")
	columns = append(columns, "NE")
	columns = append(columns, "NF")
	columns = append(columns, "NG")
	columns = append(columns, "NH")
	columns = append(columns, "NI")
	columns = append(columns, "NJ")
	columns = append(columns, "NK")
	columns = append(columns, "NL")
	columns = append(columns, "NM")
	columns = append(columns, "NN")
	columns = append(columns, "NO")
	columns = append(columns, "NP")
	columns = append(columns, "NQ")
	columns = append(columns, "NR")
	columns = append(columns, "NS")
	columns = append(columns, "NT")
	columns = append(columns, "NU")
	columns = append(columns, "NV")
	columns = append(columns, "NW")
	columns = append(columns, "NX")
	columns = append(columns, "NY")
	columns = append(columns, "NZ")
	return columns
}
