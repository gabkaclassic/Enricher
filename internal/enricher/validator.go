package enricher

import (
	"enricher/internal/enricher/dto"
	"errors"
	"fmt"
	"net/url"
	"os"
)

func isValidURL(str string) bool {
	u, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}
	if u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}

func validateEnricherArgType(value any, requiredType dto.EnricherConfigArgType) error {

	ok := true
	switch requiredType {
	case dto.StringArg:
		_, ok = value.(string)
	case dto.IntArg:
		_, ok = value.(int)
	case dto.BoolArg:
		_, ok = value.(bool)
	case dto.FloatArg:
		_, ok = value.(float64)
	default:
		return fmt.Errorf("unknown argument type: %s", requiredType)
	}
	if !ok {
		return fmt.Errorf("invalid type for argument: %v, expected %s", value, requiredType)
	}
	return nil
}

func validateEnricherSource(source string) error {

	if len(source) > 0 && !isValidURL(source) {
		return errors.New(fmt.Sprintf("invalid source URL: %s", source))
	}
	return nil
}

func validateEnricherExecutablePath(executablePath string) error {
	if _, err := os.Stat(executablePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New(fmt.Sprintf("executable file not found: %s", executablePath))
		}
		return errors.New(fmt.Sprintf("invalid executable path: %s", executablePath))
	}
	return nil

}

func validateEnricherArgs(args []dto.EnricherConfigArg) error {
	for _, arg := range args {

		if arg.Required && arg.DefaultValue != nil {
			return errors.New(fmt.Sprintf("default value provided for required argument: %s", arg.Name))
		}

		if arg.DefaultValue == nil {
			continue
		}

		err := validateEnricherArgType(arg.DefaultValue, arg.Type)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateAllowedTypes(enricherTypes []dto.EnricherArgType) error {
	if len(enricherTypes) <= 0 {
		return errors.New(fmt.Sprintf("Enrichers allowed types list cannot be empty"))
	}
	return nil
}

func validateEnricher(enricherValue dto.Enricher) error {
	err := validateEnricherExecutablePath(enricherValue.ExecutablePath)

	if err != nil {
		return err
	}

	err = validateEnricherSource(enricherValue.Source)

	if err != nil {
		return err
	}

	err = validateEnricherArgs(enricherValue.ConfigArgs)

	if err != nil {
		return err
	}

	err = validateAllowedTypes(enricherValue.AllowedTypes)

	if err != nil {
		return err
	}

	return nil
}
