package redis

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Exists to check if key existed
func (s *Redis) Exists(key string) bool {
	result, err := s.rd.Exists(s.ctx, key).Result()
	if err != nil || result == 0 {
		return false
	}

	return true
}

// Set to set key value
func (s *Redis) Set(key string, input interface{}, expiration time.Duration) error {
	var value string
	switch v := input.(type) {
	case int64:
		// If input is int64, convert it to string
		value = strconv.FormatInt(v, 10)
	case string:
		// If input is a string, assume it's already in JSON format
		value = v
	case *string:
		// If input is a pointer to a string, dereference the pointer and use it
		value = *v
	default:
		// For other types, marshal them to JSON
		data, err := json.Marshal(input)
		if err != nil {
			return fmt.Errorf("failed to marshal input to JSON: %v", err)
		}
		value = string(data)
	}

	return s.rd.Set(s.ctx, key, value, expiration).Err()
}

// Get to get key value
func (s *Redis) Get(key string, output interface{}) error {
	result, err := s.rd.Get(s.ctx, key).Result()
	if err != nil {
		return err
	}

	switch v := output.(type) {
	case *int64:
		// If output is an int64 pointer, convert result to int64 and assign it
		num, err := strconv.ParseInt(result, 10, 64)
		if err != nil {
			return err
		}
		*v = num
	case string:
		v = result
	case *string:
		// If output is a string pointer, assign the result to it
		*v = result
	default:
		// For other types, attempt to unmarshal JSON into it
		err := json.Unmarshal([]byte(result), output)
		if err != nil {
			return err
		}
	}

	return nil
}

// Del to delete key
func (s *Redis) Del(key string) error {
	return s.rd.Del(s.ctx, key).Err()
}

// SetExp to increase expiration time of key
func (s *Redis) SetExp(key string, expiration time.Duration) error {
	return s.rd.Expire(s.ctx, key, expiration).Err()
}

// SetNX to set key value if successfully returns bool value
// if key existed returns false
func (s *Redis) SetNX(key string, input interface{}, expiration time.Duration) (bool, error) {
	// parse input to JSON
	value, err := json.Marshal(input)
	if err != nil {
		return false, err
	}

	return s.rd.SetNX(s.ctx, key, value, expiration).Result()
}

// HGet to get value from key as hashmap
func (s *Redis) HGet(key, field string) (string, error) {
	return s.rd.HGet(s.ctx, key, field).Result()
}

// HSet to set value to key as hashmap
func (s *Redis) HSet(key string, values interface{}, expiration time.Duration) error {
	s.rd.HSet(s.ctx, key, values)
	return s.rd.Expire(s.ctx, key, expiration).Err()
}

// HSetArray to set array value to key as hashmap
func (s *Redis) HSetArray(hashKey, field string, arr []string) error {
	// Serialize the array to JSON
	arrJSON, err := json.Marshal(arr)
	if err != nil {
		return fmt.Errorf("error marshalling array to JSON: %v", err)
	}

	// Use HSET to store the serialized array as a string value
	err = s.rd.HSet(s.ctx, hashKey, field, string(arrJSON)).Err()
	if err != nil {
		return fmt.Errorf("error setting array in hash: %v", err)
	}

	return nil
}

// HDel to delete field from key as hashmap
func (s *Redis) HDel(key string, fields ...string) error {
	return s.rd.HDel(s.ctx, key, fields...).Err()
}

// HGetAll to get key value as hashmap
// The output must be a struct
func (s *Redis) HGetAll(key string, output interface{}) error {
	// ! remember to input "redis" tag to the struct fields
	return s.rd.HGetAll(s.ctx, key).Scan(output)
}
