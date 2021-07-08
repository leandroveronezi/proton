package main

import (
	"os"
	"strings"
)

const lineBreak = "\r\n"
const tabulation = "\t"

func main() {

	var protocol TypeBrowserProtocol

	if len(os.Args) < 2 {
		Log("Arquivo não informado")
		return
	}

	files := make(map[string]goFile)

	arquivos := os.Args[1:]

	for _, arq := range arquivos {

		err := protocol.Load(arq)

		if err != nil {
			Log(err)
			return
		}

		major := protocol.FValue("version", "major").DString("0")
		minor := protocol.FValue("version", "minor").DString("0")

		Log("Versão:", minor+"."+major)

		err, domains := protocol.FValue("domains").Array()

		if err != nil {
			Log(err)
			return
		}

		for _, d := range domains {

			err, aux := getMap(d)

			if err != nil {
				Log(err)
				return
			}

			domain := TypeBrowserProtocol(aux)

			file, ok := files[domain.FValue("domain").DString("")]

			if !ok {
				file = NewFile()
				file.Name = domain.FValue("domain").DString("")
				file.Description = domain.FValue("description").DString("")
				file.Minor = minor
				file.Major = major
			}

			if domain.Exist("types") {

				err, types := domain.FValue("types").Array()

				if err != nil {
					Log(err)
				}

				file.Add(generateType(types))

			}

			if domain.Exist("commands") {

				err, commands := domain.FValue("commands").Array()

				if err != nil {
					Log(err)
				}

				file.Add(generateFunc(file.Name, commands))

			}

			files[file.Name] = file

		}

	}

	for _, f := range files {
		f.Save()
	}

}

func jsonTypeToGOType(tp string) string {

	switch tp {
	case "string":
		return "string"

	case "integer":
		return "int"

	case "number":
		return "int"

	case "boolean":
		return "bool"

	case "any":
		return "interface{}"

	}

	if strings.Index(tp, ".") > -1 {

		s := strings.Split(tp, ".")

		tp = s[len(s)-1]

	}

	return tp

}

func generateFunc(domain string, types []interface{}) (result string) {

	for _, d := range types {

		err, aux := getMap(d)

		if err != nil {
			Log(err)
			return ""
		}

		tp := TypeBrowserProtocol(aux)

		err, parameters := tp.FValue("parameters").Array()

		if err != nil {
			//Log(err)
			//return ""
		}

		if len(parameters) > 0 {

			result += "type " + ToCamel(domain+"_"+tp.FValue("name").DString("")+"_parameters") + " struct {" + lineBreak

			for _, v := range parameters {

				err, aux := getMap(v)

				if err != nil {
					Log(err)
					return ""
				}

				auxField := TypeBrowserProtocol(aux)

				result += tabulation + ToCamel(auxField.FValue("name").DString("")) + " "

				if auxField.FValue("optional").DBool(false) {
					result += "*"
				}

				if auxField.FValue("type").DString("") == "array" {

					if auxField.Exist("items", "$ref") {
						result += "[]" + jsonTypeToGOType(auxField.FValue("items", "$ref").DString("")) + " "
					} else {
						result += "[]" + jsonTypeToGOType(auxField.FValue("items", "type").DString("")) + " "
					}

				} else {

					if auxField.Exist("$ref") {
						result += jsonTypeToGOType(auxField.FValue("$ref").DString("")) + " "
					} else {
						result += jsonTypeToGOType(auxField.FValue("type").DString("")) + " "
					}

				}

				result += "`json:\"" + auxField.FValue("name").DString("") + "\"` "

				if auxField.Exist("description") {
					result += "/* " + auxField.FValue("description").DString("") + " */"
				}

				result += lineBreak

			}

			result += "}" + lineBreak + lineBreak

		}

		err, returns := tp.FValue("returns").Array()

		if err != nil {
			//Log(err)
			//return ""
		}

		if len(returns) > 0 {

			result += "type " + ToCamel(domain+"_"+tp.FValue("name").DString("")+"_returns") + " struct {" + lineBreak

			for _, v := range returns {

				err, aux := getMap(v)

				if err != nil {
					Log(err)
					return ""
				}

				auxField := TypeBrowserProtocol(aux)

				result += tabulation + ToCamel(auxField.FValue("name").DString("")) + " "

				if auxField.FValue("optional").DBool(false) {
					result += "*"
				}

				if auxField.FValue("type").DString("") == "array" {

					if auxField.Exist("items", "$ref") {
						result += "[]" + jsonTypeToGOType(auxField.FValue("items", "$ref").DString("")) + " "
					} else {
						result += "[]" + jsonTypeToGOType(auxField.FValue("items", "type").DString("")) + " "
					}

				} else {

					if auxField.Exist("$ref") {
						result += jsonTypeToGOType(auxField.FValue("$ref").DString("")) + " "
					} else {
						result += jsonTypeToGOType(auxField.FValue("type").DString("")) + " "
					}

				}

				result += "`json:\"" + auxField.FValue("name").DString("") + "\"` "

				if auxField.Exist("description") {
					result += "/* " + auxField.FValue("description").DString("") + " */"
				}

				result += lineBreak

			}

			result += "}" + lineBreak + lineBreak

		}

		typeDescription := tp.FValue("description").DString("")

		result += "/*" + ToCamel(domain+"_"+tp.FValue("name").DString("")) + " "

		if len(strings.Trim(typeDescription, " ")) > 0 {
			result += typeDescription
		}

		result += "*/"

		result += lineBreak

		result += "func (_this *Browser) " + ToCamel(domain+"_"+tp.FValue("name").DString("")) + "("

		if len(parameters) > 0 {
			result += "Parameters " + ToCamel(domain+"_"+tp.FValue("name").DString("")+"_parameters")
		}

		result += ")"

		if len(returns) > 0 {
			result += "(" + ToCamel(domain+"_"+tp.FValue("name").DString("")+"_returns") + ",error) {" + lineBreak

			result += tabulation + `result, err := _this.send("` + domain + "." + tp.FValue("name").DString("")

			if len(parameters) > 0 {
				result += `", structToMap(Parameters))` + lineBreak
			} else {
				result += `", h{})` + lineBreak
			}

			result += tabulation + `data := ` + ToCamel(domain+"_"+tp.FValue("name").DString("")+"_returns") + `{}` + lineBreak

			result += tabulation + `if err != nil {return data, err}` + lineBreak

			result += tabulation + `err = json.Unmarshal(result, &data)` + lineBreak

			result += tabulation + `return data, err` + lineBreak

		} else {
			result += " error {" + lineBreak

			result += tabulation + `_, err := _this.send("` + domain + "." + tp.FValue("name").DString("")

			if len(parameters) > 0 {
				result += `", structToMap(Parameters))` + lineBreak
			} else {
				result += `", h{})` + lineBreak
			}

			result += `return err` + lineBreak

		}

		result += "}" + lineBreak

		result += lineBreak + lineBreak

	}

	return result
}

func generateType(types []interface{}) (result string) {

	for _, d := range types {

		err, aux := getMap(d)

		if err != nil {
			Log(err)
			return ""
		}

		tp := TypeBrowserProtocol(aux)

		typeNome := tp.FValue("id").DString("")

		typeTipo := tp.FValue("type").DString("")

		typeDescription := tp.FValue("description").DString("")

		if typeTipo == "object" {

			if len(strings.Trim(typeDescription, " ")) > 0 {
				result += "/* " + typeDescription + " */" + lineBreak
			}

			result += "type " + typeNome + " struct {" + lineBreak

			if tp.Exist("properties") {

				err, prop := tp.FValue("properties").Array()

				if err != nil {
					Log(err)
					return ""
				}

				for _, v := range prop {

					err, aux := getMap(v)

					if err != nil {
						Log(err)
						return ""
					}

					auxField := TypeBrowserProtocol(aux)

					result += tabulation + ToCamel(auxField.FValue("name").DString("")) + " "

					if auxField.FValue("optional").DBool(false) {
						result += "*"
					}

					if auxField.FValue("type").DString("") == "array" {

						if auxField.Exist("items", "$ref") {
							result += "[]" + jsonTypeToGOType(auxField.FValue("items", "$ref").DString("")) + " "
						} else {
							result += "[]" + jsonTypeToGOType(auxField.FValue("items", "type").DString("")) + " "
						}

					} else {

						if auxField.Exist("$ref") {
							result += jsonTypeToGOType(auxField.FValue("$ref").DString("")) + " "
						} else {
							result += jsonTypeToGOType(auxField.FValue("type").DString("")) + " "
						}

					}

					result += "`json:\"" + auxField.FValue("name").DString("") + "\"` "

					if auxField.Exist("description") {
						result += "/* " + auxField.FValue("description").DString("") + " */"
					}

					result += lineBreak

				}

			}

			result += "}" + lineBreak + lineBreak

		} else if tp.Exist("enum") {

			result += "type " + typeNome + " " + jsonTypeToGOType(typeTipo)

			if len(strings.Trim(typeDescription, " ")) > 0 {
				result += " /* " + typeDescription + " */"
			}

			result += lineBreak + lineBreak

			result += "const ("
			result += lineBreak

			err, valores := tp.FValue("enum").Array()

			if err != nil {
				Log(err)
				return ""
			}

			for _, v := range valores {

				result += tabulation + ToCamel(v.(string)) + "_" + typeNome + " " + typeNome + ` = "` + v.(string) + `"`
				result += lineBreak

			}

			result += ")"
			result += lineBreak + lineBreak

			result += "func (_this " + typeNome + ") Pointer() *" + typeNome + " { return &_this }"
			result += lineBreak + lineBreak

		} else {

			result += "type " + typeNome + " "

			if typeTipo == "array" {

				auxTipo := ""

				if tp.Exist("items", "$ref") {
					auxTipo = tp.FValue("items", "$ref").DString("")
				} else {
					auxTipo = tp.FValue("items", "type").DString("")
				}

				result += "[]" + jsonTypeToGOType(auxTipo)
			} else {
				result += jsonTypeToGOType(typeTipo)
			}

			if len(strings.Trim(typeDescription, " ")) > 0 {
				result += " /* " + typeDescription + " */"
			}

			result += lineBreak + lineBreak

		}

	}

	return result

}
