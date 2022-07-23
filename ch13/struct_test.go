package ch13

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

type User struct {
	UserID   string
	UserName string
	Language []string
}

func TestTom(t *testing.T) {
	tom := User{
		UserID:   "0001",
		UserName: "Tom",
		Language: []string{"Java", "Go"},
	}

	tom2 := User{
		UserID:   "0001",
		UserName: "Tom",
		Language: []string{"Java", "Go"},
	}

	if !reflect.DeepEqual(tom, tom2) {
		t.Errorf("User tom is mismatch by DeepEqual, tom=%v, tom2=%v", tom, tom2)
	}

	if diff := cmp.Diff(tom, tom2); diff != "" {
		t.Errorf("User value is mismatch by go-cmp.Diff, (-tom +tom2):\n%s", diff)
	}
}

func TestByTestify(t *testing.T) {
	result, err := Add(1, 2)
	assert.Nil(t, err)
	assert.Equal(t, 3, result)
}

func TestErrorMessage(t *testing.T) {
	assert.Equal(t, 2, 3)
	type Person struct {
		Name string
		Age  int
	}
	assert.Equal(t, Person{"織田信長", 49}, Person{"徳川家康", 73})
}

// go test -v -run TestX github.com/MatsuoTakuro/practical-golang/ch13
func TestX(t *testing.T) {
	type X struct {
		numUnExport int
		NumExport   int
	}

	num1 := X{
		numUnExport: 100,
		NumExport:   -1,
	}

	num2 := X{
		numUnExport: 999,
		NumExport:   -2,
	}

	opt := cmpopts.IgnoreUnexported(X{})

	if diff := cmp.Diff(num1, num2, opt); diff != "" {
		t.Errorf("X value is mismatch with IgnoreUnexported, (-num +num2):\n%s", diff)
	}

	opt = cmp.AllowUnexported(X{})

	if diff := cmp.Diff(num1, num2, opt); diff != "" {
		t.Errorf("X value is mismatch with AllowUnexported, (-num +num2):\n%s", diff)
	}
}

func TestX2(t *testing.T) {
	type X struct {
		NumExport int
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	num1 := X{
		NumExport: -1,
		CreatedAt: time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Now(),
	}

	num2 := X{
		NumExport: -1,
		CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Now().Add(24 * time.Hour),
	}

	opt := cmpopts.IgnoreFields(X{}, "CreatedAt", "UpdatedAt")

	if diff := cmp.Diff(num1, num2, opt); diff != "" {
		t.Errorf("X value is mismatch with AllowUnexported, (-num +num2):\n%s", diff)
	}
}

func TestX3(t *testing.T) {
	type X struct {
		NumExport   int
		numUnExport int
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}

	num1 := X{
		NumExport:   100,
		numUnExport: -1,
		CreatedAt:   time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Now(),
	}

	num2 := X{
		NumExport:   100,
		numUnExport: -111,
		CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Now().Add(24 * time.Hour),
	}

	opts := []cmp.Option{
		cmpopts.IgnoreUnexported(X{}),
		cmpopts.IgnoreFields(X{}, "CreatedAt", "UpdatedAt"),
	}

	if diff := cmp.Diff(num1, num2, opts...); diff != "" {
		t.Errorf("X value is mismatch with AllowUnexported, (-num +num2):\n%s", diff)
	}
}
