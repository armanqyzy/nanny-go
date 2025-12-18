package validator

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	if err := validate.RegisterValidation("phone_kz", validateKazakhPhone); err != nil {
		log.Fatal("Failed to register phone_kz validator:", err)
	}
	if err := validate.RegisterValidation("pet_type", validatePetType); err != nil {
		log.Fatal("Failed to register pet_type validator:", err)
	}
	if err := validate.RegisterValidation("booking_status", validateBookingStatus); err != nil {
		log.Fatal("Failed to register booking_status validator:", err)
	}
	if err := validate.RegisterValidation("user_role", validateUserRole); err != nil {
		log.Fatal("Failed to register user_role validator:", err)
	}
}

func Validate(data interface{}) error {
	err := validate.Struct(data)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			return formatValidationErrors(validationErrors)
		}
		return err
	}
	return nil
}

func formatValidationErrors(errors validator.ValidationErrors) error {
	var messages []string

	for _, err := range errors {
		message := getErrorMessage(err)
		messages = append(messages, message)
	}

	return fmt.Errorf("%s", strings.Join(messages, "; "))
}

func getErrorMessage(err validator.FieldError) string {
	field := getFieldName(err.Field())

	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s requiered to fill", field)
	case "email":
		return fmt.Sprintf("%s must be correct email address", field)
	case "min":
		return fmt.Sprintf("%s must be atleast %s symbols", field, err.Param())
	case "max":
		return fmt.Sprintf("%s must be not larger than %s symbols", field, err.Param())
	case "gte":
		return fmt.Sprintf("%s must be not less than %s", field, err.Param())
	case "lte":
		return fmt.Sprintf("%s must be not more than %s", field, err.Param())
	case "gt":
		return fmt.Sprintf("%s must be more %s", field, err.Param())
	case "lt":
		return fmt.Sprintf("%s must be less %s", field, err.Param())
	case "phone_kz":
		return fmt.Sprintf("%s must be in this format: +7XXXXXXXXXX", field)
	case "pet_type":
		return fmt.Sprintf("%s must be one of: dog, cat, bird, rat, raptile, other", field)
	case "booking_status":
		return fmt.Sprintf("%s must be one of: pending, confirmed, cancelled, completed", field)
	case "user_role":
		return fmt.Sprintf("%s must be one of: owner, sitter, admin", field)
	default:
		return fmt.Sprintf("%s did not do validation (%s)", field, err.Tag())
	}
}

func getFieldName(field string) string {
	fieldNames := map[string]string{
		"FullName":        "Полное имя",
		"Email":           "Email",
		"Phone":           "Телефон",
		"Password":        "Пароль",
		"ExperienceYears": "Опыт работы",
		"Certificates":    "Сертификаты",
		"Preferences":     "Предпочтения",
		"Location":        "Местоположение",
		"Name":            "Имя",
		"Type":            "Тип",
		"Age":             "Возраст",
		"Notes":           "Заметки",
		"Rating":          "Рейтинг",
		"Comment":         "Комментарий",
		"OwnerID":         "ID владельца",
		"SitterID":        "ID няни",
		"PetID":           "ID питомца",
		"ServiceID":       "ID услуги",
		"BookingID":       "ID бронирования",
		"StartTime":       "Время начала",
		"EndTime":         "Время окончания",
		"PricePerHour":    "Цена за час",
		"Description":     "Описание",
	}

	if name, ok := fieldNames[field]; ok {
		return name
	}
	return field
}

func validateKazakhPhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	matched, _ := regexp.MatchString(`^\+7\d{10}$`, phone)
	return matched
}

func validatePetType(fl validator.FieldLevel) bool {
	petType := strings.ToLower(fl.Field().String())
	validTypes := []string{"dog", "cat", "bird", "rat", "raptile", "other"}

	for _, valid := range validTypes {
		if petType == valid {
			return true
		}
	}
	return false
}

func validateBookingStatus(fl validator.FieldLevel) bool {
	status := strings.ToLower(fl.Field().String())
	validStatuses := []string{"pending", "confirmed", "cancelled", "completed"}

	for _, valid := range validStatuses {
		if status == valid {
			return true
		}
	}
	return false
}

func validateUserRole(fl validator.FieldLevel) bool {
	role := strings.ToLower(fl.Field().String())
	validRoles := []string{"owner", "sitter", "admin"}

	for _, valid := range validRoles {
		if role == valid {
			return true
		}
	}
	return false
}
