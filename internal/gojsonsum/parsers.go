package gojsonsum

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

func ParseDefinitionComments(r io.Reader) ([][]string, error) {
	var out [][]string
	var err error

	// read each line of the file looking for the gosumtype: <type> comment/directive
	scanner := bufio.NewScanner(r)

	var (
		withinDirective = false
		partial         = []string{}
	)

	for scanner.Scan() {
		line := scanner.Text()
		lineTrimmedSpace := strings.TrimSpace(line)

		if withinDirective {
			if strings.HasPrefix(lineTrimmedSpace, "type") || !strings.HasPrefix(lineTrimmedSpace, "//") {
				// insert line and assume that the next line is the
				// not related. This only supports single line type
				// definitions
				//
				// eg. `type AutomationType string`

				if strings.HasPrefix(line, "\t") {
					// type definition is using the
					// type (
					//   TypeA string
					// )
					// format, so we'll cheat and prepend type to the line to
					// normalize it to the single line format
					lineTrimmedSpace = "type " + lineTrimmedSpace
				}

				partial = append(partial, lineTrimmedSpace)
				out = append(out, partial)

				// reset
				partial = []string{}
				withinDirective = false
				continue
			}

			partial = append(partial, lineTrimmedSpace)
			continue
		}

		if strings.HasPrefix(lineTrimmedSpace, "// gosumtype:") {
			partial = append(partial, lineTrimmedSpace)
			withinDirective = true
		}
	}

	err = scanner.Err()
	if err != nil {
		return nil, err
	}

	return out, nil
}

var typeMapRe = regexp.MustCompile(`['"]?([\w-]+)['"]?\s*:\s*['"]?([\w-]+)['"]?`)

func ParseSumTypeDef(input []string) SumTypeDef {
	result := SumTypeDef{
		TypeMap: make(map[string]string),
		Opts: SumTypeOptions{
			Tag:  "type",
			Name: "Type",
		},
	}

	for _, line := range input {
		// trim '//' and whitespace
		trimmed := strings.TrimPrefix(line, "//")
		trimmed = strings.TrimSpace(trimmed)

		if trimmed == "" {
			continue
		}

		if strings.HasPrefix(trimmed, "gosumtype:") {
			value := strings.TrimPrefix(trimmed, "gosumtype:")
			value = strings.TrimSpace(value)

			result.TypeName = value
			continue
		}

		if strings.HasPrefix(trimmed, "opt:") {
			split := strings.Split(trimmed, ":")
			if len(split) != 3 {
				log.Warn().Msgf("invalid opt directive: %s", trimmed)
				continue
			}

			key := strings.TrimSpace(split[1])

			switch strings.ToLower(key) {
			case "tag":
				result.Opts.Tag = strings.TrimSpace(split[2])
			case "name":
				result.Opts.Name = strings.TrimSpace(split[2])
			}

			continue
		}

		if strings.HasPrefix(trimmed, "type") {
			// assume it's a type definition
			result.Discriminator = strings.Split(trimmed, " ")[1]
			continue
		}

		// else assume it's a type map definition
		typeMapMatches := typeMapRe.FindAllStringSubmatch(trimmed, -1)
		for _, match := range typeMapMatches {
			if len(match) > 2 {
				result.TypeMap[match[1]] = match[2]
			}
		}
	}

	return result
}
