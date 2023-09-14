package utils

import (
	"fmt"
	"post/model"
	"reflect"
	"strconv"

	"gorm.io/gorm"
)

func UpdateMQTTDataToDB(data interface{}, db *gorm.DB) error {
	// Enable query logging for this operation
	db = db.Debug()

	// Update the record in the MQTTData table
	result := db.Save(data)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func FindRecordByID(id int, record *model.Post, db *gorm.DB) error {
	// Find the record by ID
	result := db.First(record, id)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}

	return nil
}

func UpdateField(post interface{}, address string, value interface{}) error {
	// Use reflection to update the field based on the address
	v := reflect.ValueOf(post).Elem()
	field := v.FieldByName(address)

	if !field.IsValid() {
		return fmt.Errorf("field not found: %s", address)
	}

	// Determine the field type and set the value accordingly
	switch field.Kind() {
	case reflect.Bool:
		// Convert the value to a boolean
		boolValue, err := strconv.ParseBool(fmt.Sprintf("%v", value))
		if err != nil {
			return err
		}
		field.SetBool(boolValue)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// Convert the value to an integer
		intValue, err := strconv.ParseFloat(fmt.Sprintf("%v", value), 64)
		if err != nil {
			return err
		}
		field.SetInt(int64(intValue))
	case reflect.Float32, reflect.Float64:
		// Convert the value to a float
		floatValue, err := strconv.ParseFloat(fmt.Sprintf("%v", value), 64)
		if err != nil {
			return err
		}
		field.SetFloat(floatValue)
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}

	return nil
}
