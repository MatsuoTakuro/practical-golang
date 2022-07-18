package ch13

import (
	"fmt"
	"testing"
)

func TestMain(m *testing.M) {
	setup()
	defer teardown()
	m.Run()
}

func setup() {
	// setup code
	fmt.Println("1. before the whole of testing")
}

func teardown() {
	// tear down code
	fmt.Println("6. after the whole of testing")
}

func Calc(a, b int, operator string) (int, error) {
	switch operator {
	case "+":
		return a + b, nil

	case "-":
		return a - b, nil

	case "*":
		return a * b, nil

	case "/":
		if b == 0 {
			return 0, fmt.Errorf("division by 0	is undefined")
		}
		return a / b, nil
	}
	return 0, fmt.Errorf("unexpected operator %v is specified", operator)
}

func TestCalc(t *testing.T) {
	type args struct {
		a        int
		b        int
		operator string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "addition",
			args: args{
				a:        10,
				b:        2,
				operator: "+",
			},
			want:    12,
			wantErr: false,
		},
		{
			name: "invaild operator specified",
			args: args{
				a:        10,
				b:        2,
				operator: "?",
			},
			want:    0,
			wantErr: true,
		},
	}
	fmt.Println("2. before excuting test func")
	defer func() {
		fmt.Println("5. after excuting test func")
	}()

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fmt.Printf("3. before excuting test case #%d\n", i+1)
			defer func() {
				fmt.Printf("4. after excuting test case #%d\n", i+1)
			}()

			got, err := Calc(tt.args.a, tt.args.b, tt.args.operator)
			// with regard to errors, an unexpected situation occurred with or without errors.
			if (err != nil) != tt.wantErr {
				t.Errorf("Calc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Calc() = %v, want %v", got, tt.want)
			}
		})
	}
}
