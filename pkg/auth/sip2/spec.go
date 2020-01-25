package sip2

import (
	"log"
	"strconv"
	"strings"
)

type FixedFieldSpec struct {
	length int
	label  string
}

type FixedField struct {
	spec  FixedFieldSpec
	value string
}

type VariableFieldSpec struct {
	id    string
	label string
}

type VariableField struct {
	spec  VariableFieldSpec
	value string
}

type MessageSpec struct {
	id     int
	label  string
	fields []FixedFieldSpec
}

type Message struct {
	fields       []VariableField
	fixed_fields []FixedField
	msg_txt      string
}

type Parser struct {
	getMessageSpecByCode   func(int) MessageSpec
	getVariableFieldByCode func(string) VariableFieldSpec
	getFixedFieldValue     func(*Message, string) (string, bool)
	getFieldValue          func(*Message, string) (string, bool)
	parseMessage           func(string) *Message
}

const (
	YES                  string = "Y"
	TRUE                 string = "1"
	Ok                   string = "ok"
	Language             string = "language"
	PatronStatus         string = "patron status"
	Date                 string = "transaction date"
	PatronIdentifier     string = "patron identifier"
	PatronPassword       string = "patron password"
	PersonalName         string = "personal name"
	ScreenMessage        string = "screen message"
	InstitutionId        string = "institution id"
	ValidPatron          string = "valid patron"
	ValidPatronPassword  string = "valid patron password"
	LoginResponse        string = "Login Response"
	PatronStatusResponse string = "Patron Status Response"
)

func getParser() *Parser {

	LanguageSpec := FixedFieldSpec{3, Language}
	PatronStatusSpec := FixedFieldSpec{14, PatronStatus}
	DateSpec := FixedFieldSpec{18, Date}
	OkSpec := FixedFieldSpec{1, Ok}

	msgByCodeMap := map[int]MessageSpec{
		94: MessageSpec{94, LoginResponse, []FixedFieldSpec{OkSpec}},
		24: MessageSpec{24, PatronStatusResponse, []FixedFieldSpec{PatronStatusSpec, LanguageSpec, DateSpec}},
	}

	variableFieldByCodeMap := map[string]VariableFieldSpec{
		"AA": VariableFieldSpec{"AA", PatronIdentifier},
		"AD": VariableFieldSpec{"AD", PatronPassword},
		"AE": VariableFieldSpec{"AE", PersonalName},
		"AF": VariableFieldSpec{"AF", ScreenMessage},
		"AO": VariableFieldSpec{"AO", InstitutionId},
		"BL": VariableFieldSpec{"BL", ValidPatron},
		"CQ": VariableFieldSpec{"CQ", ValidPatronPassword},
	}

	parser := new(Parser)
	parser.getMessageSpecByCode = func(code int) MessageSpec {
		return msgByCodeMap[code]
	}
	parser.getVariableFieldByCode = func(code string) VariableFieldSpec {
		return variableFieldByCodeMap[code]
	}
	parser.getFixedFieldValue = func(msg *Message, field string) (string, bool) {
		for _, v := range msg.fixed_fields {
			if v.spec.label == field {
				return v.value, true
			}
		}
		return "", false
	}
	parser.getFieldValue = func(msg *Message, field string) (string, bool) {
		for _, v := range msg.fields {
			if v.spec.label == field {
				return v.value, true
			}
		}
		return "", false
	}

	parser.parseMessage = func(msg string) *Message {
		txt := msg[:len(msg)-len(terminator)]
		code, err := strconv.Atoi(txt[:2])
		if nil != err {
			log.Println("Error parsing integer: %s", txt[:2])
		}
		spec := parser.getMessageSpecByCode(code)
		txt = txt[2:]

		message := new(Message)
		for _, sp := range spec.fields {
			value := txt[:sp.length]
			txt = txt[sp.length:]
			message.fixed_fields = append(message.fixed_fields, FixedField{sp, value})
		}
		if len(txt) == 0 {
			return message
		}
		for _, part := range strings.Split(txt, "|") {
			if len(part) > 0 {
				part_spec := parser.getVariableFieldByCode(part[:2])
				value := part[2:]
				message.fields = append(message.fields, VariableField{part_spec, value})
			}
		}
		return message
	}
	return parser
}
